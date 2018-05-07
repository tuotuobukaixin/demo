package algorithm

import (
	"container/list"
	"fmt"
	"sort"
	"time"
)

const (
	nMax       = 30
	nnMax      = 900
	nResource  = 10
	cnResource = 2
)

var nApps int
var nClusters int
var maxCluster int

type nodeType struct {
	cLoad        [nMax][nResource]float32
	uLoad        [nMax][nResource]float32
	rules        [nMax][nMax]float32
	cnt          int
	maxLoadRatio float32
	x            [nMax][nMax]int
}

type phyResourceType struct {
	name        int
	capacity    float32
	demand      float32
	used        float32
	utilizaiton float32
	i           int
}

type appResourceType struct {
	name        int
	demand      float32
	used        float32
	utilization float32
	i           int
}

type clusterType struct {
	resourcePhy [cnResource]phyResourceType
	name        string
}

type appType struct {
	resourceItem [cnResource]appResourceType
	name         string
	clusterid    string
}

type dynamicLoadBlanceAlgorithom struct{}

var minn nodeType
var gCluster [nMax]clusterType
var gApplication [nMax]appType

//Queue is a ring queue
var Queue *list.List

func (p dynamicLoadBlanceAlgorithom) Calculate(param interface{}) interface{} {
	fmt.Println("dynamicLoadBlanceAlgorithom Calculate...")

	appLoad, clusterLoad := get5MinAverageLoad(param.(DynamicSchedulingParam))

	fmt.Println("current appload=", appLoad)
	fmt.Println("current clusterLoad=", clusterLoad)
	//1.Compose a map[AppName]ClusterName from apps string in clusters
	deployedApps := map[string]string{}
	for cltName, cluster := range clusterLoad {
		for appName := range cluster.AppList {
			deployedApps[appName] = cltName
		}
	}

	parseInput(appLoad, clusterLoad)
	nApps = len(appLoad)
	nClusters = len(clusterLoad)

	Queue = list.New()
	minn.inital()
	Lbs()

	recommendation := Recommendation{
		ID:           time.Now().String(),
		AppToCluster: map[string]string{},
		ClusterOffID: list.New(),
	}
	for i := 0; i < len(appLoad); i++ {
		for j := 0; j < len(clusterLoad); j++ {
			if minn.x[i][j] == 1 {
				if deployedApps[gApplication[i].name] != gCluster[j].name {
					recommendation.AppToCluster[appLoad[gApplication[i].name].UUID] = clusterLoad[gCluster[j].name].UUID
					fmt.Println("App", gApplication[i].name, "xxx------->Cluster", gCluster[j].name)
				}
			}
		}
	}

	if minn.maxLoadRatio <= 0 {
		return nil
	}
	return recommendation
}

func (p *nodeType) inital() {
	p.cnt = 0
	p.maxLoadRatio = -1
	for m := 0; m < nMax; m++ {
		for r := 0; r < nResource; r++ {
			p.cLoad[m][r] = 0
			p.uLoad[m][r] = 0
		}
	}
	for n := 0; n < nMax; n++ {
		for m := 0; m < nMax; m++ {
			p.x[n][m] = 0
		}
	}
}

func parseInput(apps map[string]Application, clusters map[string]Cluster) {

	clusterList := []string{}
	uuidtoName := map[string]string{}
	for _, clt := range clusters {
		clusterList = append(clusterList, clt.UUID)
		uuidtoName[clt.UUID] = clt.Name
	}
	sort.Sort(sort.StringSlice(clusterList))

	appList := []string{}
	for _, app1 := range apps {
		uuidtoName[app1.UUID] = app1.Name
		appList = append(appList, app1.UUID)
	}
	sort.Sort(sort.StringSlice(appList))

	for i := 0; i < len(appList); i++ {
		app1 := apps[uuidtoName[appList[i]]]
		j := 0
		for resType, res := range app1.ResourceList {
			gApplication[i].resourceItem[resType].name = resType
			gApplication[i].resourceItem[resType].demand = res.Total
			gApplication[i].resourceItem[resType].used = res.Used
			gApplication[i].resourceItem[resType].utilization = res.Utilization
			j++
		}
		gApplication[i].name = app1.Name
	}

	for i := 0; i < len(clusterList); i++ {
		clt := clusters[uuidtoName[clusterList[i]]]
		j := 0
		for resType, res := range clt.ResourceList {
			gCluster[i].resourcePhy[resType].name = resType
			gCluster[i].resourcePhy[resType].capacity = res.Total
			j++
		}
		gCluster[i].name = clt.Name
	}
}

func insert(flag bool, i int, cL [nMax][nResource]float32, uL [nMax][nResource]float32, temp nodeType) {

	for m := 0; m < nClusters; m++ {
		for r := 0; r < cnResource; r++ {
			if cL[m][r] > gCluster[m].resourcePhy[r].capacity || (minn.maxLoadRatio != -1 && uL[m][r]/gCluster[m].resourcePhy[r].capacity >= minn.maxLoadRatio) {
				return
			}
		}
	}

	for n := 0; n <= i/nClusters; n++ {
		aL := 0
		for m := 0; m < nClusters; m++ {
			aL += temp.x[n][m]
			if aL > 1 {
				return
			}
		}
	}

	if flag == false && i%nClusters == nClusters-1 {
		aL := 0
		for m := 0; m < nClusters; m++ {
			aL += temp.x[i/nClusters][m]
		}
		if aL == 0 {
			return
		}
	}

	var cApp nodeType
	for m := 0; m < nClusters; m++ {
		for r := 0; r < cnResource; r++ {
			cApp.cLoad[m][r] = cL[m][r]
			cApp.uLoad[m][r] = uL[m][r]
		}
	}

	for n := 0; n < nApps; n++ {
		for m := 0; m < nClusters; m++ {
			cApp.x[n][m] = temp.x[n][m]
		}
	}
	cApp.cnt = i + 1
	Queue.PushBack(cApp)
}

//UmaxLoad is the get the max usage of Nodetype
// @Title UmaxLoad
// @Description
func UmaxLoad(temp nodeType) float32 {
	var umax float32 = -1
	var uResourceResult float32 = -1
	for ml := 0; ml < nClusters; ml++ {
		for rl := 0; rl < cnResource; rl++ {
			uResourceResult = temp.uLoad[ml][rl] / gCluster[ml].resourcePhy[rl].capacity
			if uResourceResult > umax {
				umax = uResourceResult
			}
		}
	}

	return umax
}

//Lbs is the interface of loadbalance
// @Title Lbs
// @Description
func Lbs() {
	var node nodeType
	for m := 0; m < nClusters; m++ {
		for r := 0; r < cnResource; r++ {
			node.cLoad[m][r] = 0
			node.uLoad[m][r] = 0
		}
	}
	for n := 0; n < nApps; n++ {
		for m := 0; m < nClusters; m++ {
			node.x[n][m] = 0
		}
	}
	node.maxLoadRatio = -1
	node.cnt = 0
	Queue.PushBack(node)

	for Queue.Len() > 0 {
		temp := Queue.Back().Value.(nodeType)
		Queue.Remove(Queue.Back())

		if temp.cnt >= nClusters*nApps {
			temp.maxLoadRatio = UmaxLoad(temp)

			if minn.maxLoadRatio == -1 || minn.maxLoadRatio > temp.maxLoadRatio {
				minn = temp
			}
			continue
		}

		var CLoad [nMax][nResource]float32
		var ULoad [nMax][nResource]float32

		for flag := 0; flag < 2; flag++ {
			i := temp.cnt

			for m := 0; m < nClusters; m++ {
				for r := 0; r < cnResource; r++ {
					CLoad[m][r] = temp.cLoad[m][r]
					ULoad[m][r] = temp.uLoad[m][r]
				}
			}

			if flag == 1 {
				for r := 0; r < cnResource; r++ {
					CLoad[i%nClusters][r] += gApplication[i/nClusters].resourceItem[r].demand
					ULoad[i%nClusters][r] += gApplication[i/nClusters].resourceItem[r].used
				}
				temp.x[i/nClusters][i%nClusters] = 1
			} else if flag == 0 {
				temp.x[i/nClusters][i%nClusters] = 0
			}

			insert((flag != 0), i, CLoad, ULoad, temp)
		}
	}
}

func findMaxCluster() {
	var umax float32 = -1
	var uResourceResult float32 = -1
	for ml := 0; ml < nClusters; ml++ {
		for rl := 0; rl < cnResource; rl++ {
			uResourceResult = minn.uLoad[ml][rl] / gCluster[ml].resourcePhy[rl].capacity
			if uResourceResult > umax {
				umax = uResourceResult
				maxCluster = ml + 1
			}
		}
	}
}

//Output is display the scaling result
// @Title Output
// @Description
func Output(x [nMax][nMax]int) {
	fmt.Println("Application：")
	for i := 0; i < nApps; i++ {
		fmt.Println("App", gApplication[i].name, ":")
		fmt.Print("Demand:")
		for j := 0; j < cnResource; j++ {
			fmt.Print(gApplication[i].resourceItem[j].demand, "\t")
		}
		fmt.Println()
		fmt.Print("Used:")
		for j := 0; j < cnResource; j++ {
			fmt.Print(gApplication[i].resourceItem[j].used, "\t")
		}
		fmt.Println()
	}
	fmt.Println("Capacity of clusters：")
	for i := 0; i < nClusters; i++ {
		fmt.Print("Cluster", gCluster[i].name, ":")
		for j := 0; j < cnResource; j++ {
			fmt.Print(gCluster[i].resourcePhy[j].capacity, "\t")
		}
		fmt.Println()
	}

	fmt.Println("The Result of scheduler:")
	if minn.maxLoadRatio == -1 {
		fmt.Println("No Recommendation")
		return
	}

	for i := 0; i < nApps; i++ {
		for j := 0; j < nClusters; j++ {
			if x[i][j] == 1 {
				fmt.Println("App", gApplication[i].name, "------->Cluster", gCluster[j].name)
			}
		}
	}
}
