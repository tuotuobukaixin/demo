package algorithm

import (
	"container/list"
	"fmt"
	"math"
	"sort"
	"time"
)

//The max of iteration
const MaxIterationTimes = 100

//The Stabilization Threshold
const StabilizationThreshold = 0.15

type dynamicLoadBlanceBFDAlgorithom struct{}

type clusterList []*Cluster

func (u clusterList) Len() int {
	return len(u)
}

func (u clusterList) Less(i, j int) bool {
	return u[i].getWeightUtilization() > u[j].getWeightUtilization() //sort by utilization
}

func (u clusterList) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (p dynamicLoadBlanceBFDAlgorithom) Calculate(param interface{}) interface{} {
	fmt.Println("dynamicLoadBlanceAlgorithom Calculate...")
	rec, _ := p.CalculateNew(param.(DynamicSchedulingParam))
	return rec
}

//CalculateNew is a interface
func (p dynamicLoadBlanceBFDAlgorithom) CalculateNew(param DynamicSchedulingParam) (interface{}, []*Cluster) {
	migrationThreshold := param.MigrationThreshold
	//get the average data in 5 min
	appLoad, clusterLoad := get5MinAverageLoad(param)
	clusters := []*Cluster{}
	clusterlist := []*Cluster{}
	for _, clt := range clusterLoad {
		t1, t2 := clt, clt
		clusterlist = append(clusterlist, &t1)
		clusters = append(clusters, &t2)
	}
	cltSize := len(clusterlist)
	appToCluster := map[string]string{}
	iteration := 0
	for {
		//1. Check the iteration's up limites
		if iteration > MaxIterationTimes {
			return nil, clusters //No way to Loadbalance, then return the current
		}

		//2. Check the system's balance by standard deviation
		standardVariance := p.getStdVariance(clusterlist)
		fmt.Println("CalculateNew(", migrationThreshold, "): xxxxxx standardVariance=", standardVariance)
		if standardVariance < migrationThreshold { //Will do SQRT(x*2/n)  with the node's change later
			break
		}

		//3. Sorted Clusters by usage, which the resource weight usage
		sort.Sort(clusterList(clusterlist))

		//4. Single-step migration: Get the lowest one from the high zone/clusters, and move it to the low zone/cluster
		minLoadAppName := ""
		var minLoadUitil float32 = 1
		for appName := range clusterlist[0].AppList {
			loadUitil := appLoad[appName].getWeightUtilization()
			if loadUitil < minLoadUitil && p.checkStabilization(appName, &param) { //Cluster should be stabled
				minLoadUitil = loadUitil
				minLoadAppName = appName
			}
		}

		//5. Move it to lowest zone/cluster which the capacity is larger than the required
		srcClt := clusterlist[0]
		for i := cltSize - 1; i > 0; i-- {
			dstClt := clusterlist[i]
			dstResList := map[int]Resource{}
			srcResList := map[int]Resource{}
			isCanPlaced := true
			for resType, appRes := range appLoad[minLoadAppName].ResourceList {
				dstCltRes := dstClt.ResourceList[resType]
				if (dstCltRes.Used + appRes.Total) > dstCltRes.Total {
					isCanPlaced = false
					break //If the Capacity can't meet, break to check next
				}

				dstCltRes.Utilization = (dstCltRes.Utilization*dstCltRes.Total + appRes.Utilization*appRes.Total) / dstCltRes.Total
				dstCltRes.Used = dstCltRes.Used + appRes.Total //Distribute volum
				dstResList[resType] = dstCltRes

				srcCltRes := srcClt.ResourceList[resType]
				srcCltRes.Utilization = (srcCltRes.Utilization*srcCltRes.Total - appRes.Utilization*appRes.Total) / srcCltRes.Total
				srcCltRes.Used = srcCltRes.Used - appRes.Total //Distribute volum
				srcResList[resType] = srcCltRes
			}
			if isCanPlaced {
				dstClt.ResourceList = dstResList
				srcClt.ResourceList = srcResList
				dstClt.AppList[minLoadAppName] = appLoad[minLoadAppName].UUID
				delete(srcClt.AppList, minLoadAppName)
				appToCluster[appLoad[minLoadAppName].UUID] = dstClt.UUID
				break
			}
		}
		iteration++
	}

	recommendation := Recommendation{
		ID:           time.Now().String(),
		AppToCluster: appToCluster,
		ClusterOffID: list.New(),
	}
	return recommendation, clusterlist
}

//The New load balance algorithm, Step-migration
//It will find the best migration application, and move it to the best Cluster
//We will do the exchange migration algorithm later, that will execute when the step-migration failed
func (p dynamicLoadBlanceBFDAlgorithom) CalculateNew1(param DynamicSchedulingParam) (interface{}, []*Cluster) {
	migrationThreshold := param.MigrationThreshold
	//Get the average data in 5 mins ago
	appLoad, clusterLoad := get5MinAverageLoad(param)
	clusters := []*Cluster{}
	clusterlist := []*Cluster{}
	for _, clt := range clusterLoad {
		t1, t2 := clt, clt
		clusters = append(clusters, &t2)
		clusterlist = append(clusterlist, &t1)
	}

	//S1. Computing the current load, if every veidoo status is balanced, then exit
	//On the side, do the iteration func to find suggestion.
	standardVariance := p.getStdVariance(clusterlist)
	fmt.Println("CalculateNew1: xxxxxx standardVariance=", standardVariance)
	if standardVariance < migrationThreshold { //Will do SQRT(x*2/n)  with the node's change later
		return nil, clusters
	}

	appToCluster := map[string]string{}
	iteration := 0
	for {
		//S2. Check the iteration's up limit，if upon, give the recommendation by current status, then renturn,
		//or else go on the iterating
		if iteration > MaxIterationTimes {
			return nil, clusters //There is no way to balanced, such as there is 2 apps, deployed 2 different clusters, and has 80% and 20%
		}

		//S3. Into the Iteration and execute setp-migration algorithm,  to find the step which can let the load down.
		// If find, goto S5, or goto S4
		sort.Sort(clusterList(clusterlist))
		variance, appName, cltName := p.singleStepMigrantion(clusterlist, appLoad, &param)
		if variance >= 0 {
			appToCluster[appName] = cltName
		} else {
			//S4. If the Step-migration failed, then execute the exchange sub-algorithm to find a migration suggestion
			// So if success then goto S5, or goto S2
			variance = p.switchMigrantion()
		}
		fmt.Println("xxxxxx variance=", variance)

		//S5. Computing the current Cluster load, if every veidoo have balanced, then exit the computing
		//or else start the iterating
		if variance < migrationThreshold { //Will do SQRT(x*2/n)  with the node's change later
			break
		}
		iteration++
	}

	recommendation := Recommendation{
		ID:           time.Now().String(),
		AppToCluster: appToCluster,
		ClusterOffID: list.New(),
	}
	return recommendation, clusterlist
}

//Get the clusters' balance degree from stardand deviation
func (p dynamicLoadBlanceBFDAlgorithom) getStdVariance(cList []*Cluster) float32 {
	aveLoad := float32(0)
	for _, clt := range cList {
		aveLoad = aveLoad + clt.getWeightUtilization()
	}
	cltSize := len(cList)
	aveLoad = aveLoad / float32(cltSize)

	var standardVariance float32
	for _, clt := range cList {
		x := clt.getWeightUtilization() - aveLoad
		standardVariance = standardVariance + x*x
	}
	standardVariance = standardVariance / float32(cltSize)
	return float32(math.Sqrt(float64(standardVariance)))
}

//The step-migration
func (p dynamicLoadBlanceBFDAlgorithom) singleStepMigrantion(cList []*Cluster, appLoad map[string]Application, param *DynamicSchedulingParam) (float32, string, string) {
	var aveUitil float32
	var utilizations = map[string]float32{}
	cltSize := len(cList)
	for i := 0; i < cltSize; i++ {
		tmp := cList[i].getWeightUtilization()
		utilizations[cList[i].Name] = tmp
		aveUitil += tmp
	}
	aveUitil = aveUitil / float32(cltSize)

	//S1. Choose the highest usage Cluster(ClusterA) by every veidoo weighted usage, and then we will move a App to other clusters.
	//S2. Choose the best App, which it's good for the cluster(ClusterA) when it is migrated.
	// The cluster's(ClusterA) every veidoo weighted usage will be lowest in all clusters, when the App migrated.
	var minDiff float32 = 1
	var utilization float32
	var minLoadAppName = ""
	var srcCluster = cList[0]
	var srcResList *map[int]Resource
	for appName := range srcCluster.AppList {
		tmpResList := map[int]Resource{}
		for resType, appRes := range appLoad[appName].ResourceList {
			tmpRes := srcCluster.ResourceList[resType]
			tmpRes.Utilization = (tmpRes.Utilization*tmpRes.Total - appRes.Utilization*appRes.Total) / tmpRes.Total
			tmpRes.Used = tmpRes.Used - appRes.Total //distributed volum, not usage
			tmpResList[resType] = tmpRes
		}
		tmp := getWeightUtilization(tmpResList)
		diff := tmp - aveUitil
		if diff < 0 {
			diff = diff * -1
		}
		if diff < minDiff && p.checkStabilization(appName, param) {
			minDiff = diff
			utilization = tmp
			minLoadAppName = appName
			srcResList = &tmpResList
		}
	}

	newStdVariance := float32(-1)
	if minLoadAppName == "" {
		return newStdVariance, "", ""
	}
	utilizations[srcCluster.Name] = utilization
	oldStdVariance := p.getLoadBlanceDegree(utilizations)
	maxDiff := float32(0)

	//S3. Choose a Cluster(ClusterB) from the Clusters without ClusterA, and move the App to it.
	// And then, the clusters every veidoo will be lower than current by every weighted usage.
	//S4. If we don't choose a destCluster(ClusterB), remove the App from the migrationList, and continue the circulation.
	// The circulation will stop when the migrationList is empty or we find the step-migration.
	var dstCluster *Cluster
	var dstResList *map[int]Resource
	for i := cltSize - 1; i > 0; i-- {
		isCanPlaced := true
		dstClt := cList[i]
		tmpResList := map[int]Resource{}
		for resType, appRes := range appLoad[minLoadAppName].ResourceList {
			tmpRes := dstClt.ResourceList[resType]
			if (tmpRes.Used + appRes.Total) > tmpRes.Total {
				isCanPlaced = false
				break //If undercapacity, then find next
			}

			tmpRes.Utilization = (tmpRes.Utilization*tmpRes.Total + appRes.Utilization*appRes.Total) / tmpRes.Total
			tmpRes.Used = tmpRes.Used + appRes.Total //distributed volum, not usage
			tmpResList[resType] = tmpRes
		}
		if isCanPlaced {
			bak := utilizations[dstClt.Name]
			utilizations[dstClt.Name] = getWeightUtilization(tmpResList)
			variance := p.getLoadBlanceDegree(utilizations)
			diff := oldStdVariance - variance
			if diff > maxDiff {
				maxDiff = diff
				newStdVariance = variance
				dstCluster = dstClt
				dstResList = &tmpResList
			}
			utilizations[dstClt.Name] = bak
		}
	}

	//If we find the Step-migration, do it
	if newStdVariance > 0 {
		srcCluster.ResourceList = *srcResList
		dstCluster.ResourceList = *dstResList
		delete(srcCluster.AppList, minLoadAppName)
		dstCluster.AppList[minLoadAppName] = appLoad[minLoadAppName].UUID
	}
	return newStdVariance, minLoadAppName, dstCluster.Name
}

//Exchange migration sub-algorithm
func (p dynamicLoadBlanceBFDAlgorithom) switchMigrantion() float32 {
	//S1. Sorted by the Clusters' every veidoo weighted usage
	// then split by the average usage of Clusters into HighLoad Zone and LowLoad Zone.

	//S2. Try to find a App(AppH) from HighLoad zone and a App(AppL) from LowLoad zone, then exchange them.
	// If the Clusters' every veidoo weighted usage lower than the current, we'll give this suggestion. Or else goon to find
	// the AppH and AppL to have a try.

	//S3. If we get it successful, give the suggestion, and then return, or else failed.
	return -1
}

//Get the Cluster's standard deviation
func (p dynamicLoadBlanceBFDAlgorithom) getLoadBlanceDegree(utilizations map[string]float32) float32 {
	aveLoad := float32(0)
	for _, clt := range utilizations {
		aveLoad = aveLoad + clt
	}
	cltSize := len(utilizations)
	aveLoad = aveLoad / float32(cltSize)

	var standardVariance float32
	for _, clt := range utilizations {
		x := clt - aveLoad
		standardVariance = standardVariance + x*x
	}
	standardVariance = standardVariance / float32(cltSize)
	return float32(math.Sqrt(float64(standardVariance)))
}

//Shock proof algorithm:
//S1. Execute the load balance algorithm: call the shock proof algorithm before choose a App to migrate.
//S2. Get the latest 10 points load data from the App's history. We suppose the interval is 20 seconds each point.
//S3. Computing these 10 points data's standard deviation, respectively calculated by CPU/MEM veidoo.
//S4. If we find a veidoo's standard deviation larger than 0.15(empirical value),then the cluster is instablility,
// and it's unfit for migrate, or else maybe migrated
func (p dynamicLoadBlanceBFDAlgorithom) checkStabilization(appName string, param *DynamicSchedulingParam) bool {
	el := param.HistoryData.Back()
	appPerfDatas := []Application{}
	size := 0
	for ; el != nil; el = el.Prev() {
		if size > FiveMinutesDataCount {
			break
		}
		appPerfDatas = append(appPerfDatas, el.Value.(PerfSnapshot).AppPerfData[appName])
		size++
	}

	if size == 0 {
		return true
	}

	//Computing every resource average usage
	average := map[int]float32{}
	for _, data := range appPerfDatas {
		for resType, res := range data.ResourceList {
			average[resType] += res.Utilization
		}
	}
	for resType := range average {
		average[resType] = average[resType] / float32(size)
	}

	//Computing every resource standard deviation
	variances := map[int]float32{}
	for _, app := range appPerfDatas {
		for resType, res := range app.ResourceList {
			x := res.Utilization - average[resType]
			variances[resType] += x * x
		}
	}

	for resType, variance := range variances {
		stdVariance := float32(math.Sqrt(float64(variance / float32(size))))
		//If some veidoo's load standard deviation larger than 0.15(empirical value),
		//then we judge the cluster is instablility, and it's unfit for migrate
		if stdVariance > StabilizationThreshold {
			fmt.Println("Info : app=", appName, "resType=", resType, " not stabilize, stdVariance is ", stdVariance)
			return false
		}
	}
	return true
}
