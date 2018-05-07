package algorithm

import (
	"container/list"
	"fmt"
	"sort"
	"time"
)

var removeFlag int
var lowNum int

type dynamicScalingAlgorithom struct{}

func (p *dynamicScalingAlgorithom) Calculate(param interface{}) interface{} {
	fmt.Println("Calculate...")
	args := param.(DynamicSchedulingParam)
	ret := p.powerOffEva(args.Clusters, args.HistoryData)
	return ret
}

//ClustersSort is the Cluster list Node of sorted
type ClustersSort struct {
	usedCPU             float32
	totalCPU            float32
	usedMEM             float32
	totalMEM            float32
	ClusterName         string
	ClusterUUID         string
	ClustersUtilization float32
}

//ClustersSortList is the list of the sorted Clusters
type ClustersSortList []ClustersSort

func (cl ClustersSortList) Len() int {
	return len(cl)
}

func (cl ClustersSortList) Less(i, j int) bool {
	return cl[i].ClustersUtilization < cl[j].ClustersUtilization
}

func (cl ClustersSortList) Swap(i, j int) {
	cl[i].usedCPU, cl[j].usedCPU = cl[j].usedCPU, cl[i].usedCPU
	cl[i].totalCPU, cl[j].totalCPU = cl[j].totalCPU, cl[i].totalCPU
	cl[i].usedMEM, cl[j].usedMEM = cl[j].usedMEM, cl[i].usedMEM
	cl[i].totalMEM, cl[j].totalMEM = cl[j].totalMEM, cl[i].totalMEM
	cl[i].ClusterName, cl[j].ClusterName = cl[j].ClusterName, cl[i].ClusterName
	cl[i].ClusterUUID, cl[j].ClusterUUID = cl[j].ClusterUUID, cl[i].ClusterUUID
	cl[i].ClustersUtilization, cl[j].ClustersUtilization = cl[j].ClustersUtilization, cl[i].ClustersUtilization
}

//power of the action
func (p *dynamicScalingAlgorithom) powerOffAction(clusters *list.List, rec *Recommendation) int32 {
	// clustersSort := make(map[string]float32)
	var cs ClustersSort
	var dstTotalCPU float32
	var dstUsedCPU float32
	var dstTotalMEN float32
	var dstUsedMEN float32

	var dstCluster interface{}

	var beforeUtilization float32
	var afterUtilization float32

	removeFlag = 0
	clusterUtilizationList := list.New()
	lowList := list.New()
	if clusters.Len() == 1 {
		return -1
	}
	csList := []ClustersSort{}
	//note the clusters' Usage
	for iterCluster := clusters.Front(); iterCluster != nil; iterCluster = iterCluster.Next() {
		cs.totalCPU = (iterCluster.Value.(Cluster)).ResourceList[0].Total
		cs.usedCPU = (iterCluster.Value.(Cluster)).ResourceList[0].Used
		cs.totalMEM = (iterCluster.Value.(Cluster)).ResourceList[1].Total
		cs.usedMEM = (iterCluster.Value.(Cluster)).ResourceList[1].Used
		cs.ClusterName = (iterCluster.Value.(Cluster)).Name
		cs.ClusterUUID = (iterCluster.Value.(Cluster)).UUID
		cs.ClustersUtilization = iterCluster.Value.(Cluster).getWeightUtilization()
		csList = append(csList, cs)
	}

	sort.Sort(ClustersSortList(csList))

	for i := 0; i < len(csList); i++ {
		clusterUtilizationList.PushBack(csList[i])
	}
	fmt.Println("list", clusterUtilizationList.Front().Value.(ClustersSort).ClusterName)

	//Seperated the Clusters by Usage into low/middle/high zone, and note the low/middle zone
	for iter := clusterUtilizationList.Front(); iter != nil; iter = iter.Next() {
		if (iter.Value.(ClustersSort)).ClustersUtilization < LoadHighThreshold {
			lowList.PushBack(iter.Value)
		}
	}

	if lowList.Len() == 1 || lowList.Len() == 0 {
		return -1
	}

	//Have a try to scale, move the apps in low zone to middle zone
	fmt.Println("-----------------------------------------------------")
	fmt.Println(lowList.Front().Value.(ClustersSort).ClustersUtilization)
	fmt.Println((lowList.Back().Value.(ClustersSort)).ClustersUtilization)
	fmt.Println(LoadHighThreshold)
	fmt.Println(lowList.Front().Value.(ClustersSort).usedCPU)
	fmt.Println(lowList.Back().Value.(ClustersSort).usedCPU)
	fmt.Println(lowList.Back().Value.(ClustersSort).totalCPU)
	fmt.Println(lowList.Front().Value.(ClustersSort).usedMEM)
	fmt.Println(lowList.Back().Value.(ClustersSort).usedMEM)
	fmt.Println(lowList.Back().Value.(ClustersSort).totalMEM)
	fmt.Println("-----------------------------------------------------")
	if (lowList.Front().Value.(ClustersSort).ClustersUtilization+lowList.Back().Value.(ClustersSort).ClustersUtilization < LoadHighThreshold) && (lowList.Front().Value.(ClustersSort).usedCPU+lowList.Back().Value.(ClustersSort).usedCPU < lowList.Back().Value.(ClustersSort).totalCPU) && (lowList.Front().Value.(ClustersSort).usedMEM+lowList.Back().Value.(ClustersSort).usedMEM < lowList.Back().Value.(ClustersSort).totalMEM) {
		fmt.Println("success in")
		for iterDst := clusters.Front(); iterDst != nil; iterDst = iterDst.Next() {
			if (lowList.Back().Value.(ClustersSort)).ClusterName == (iterDst.Value.(Cluster)).Name {
				dstCluster = iterDst.Value

				fmt.Println("dstCluster", dstCluster.(Cluster).Name)
				beforeUtilization = dstCluster.(Cluster).getWeightUtilization()
				break
			}
		}
		for iterRemove := clusters.Front(); iterRemove != nil; iterRemove = iterRemove.Next() {
			if (lowList.Front().Value.(ClustersSort)).ClusterName == (iterRemove.Value.(Cluster)).Name {
				removeFlag = 1
				for k, v := range (iterRemove.Value.(Cluster)).AppList {
					rec.AppToCluster[v] = lowList.Back().Value.(ClustersSort).ClusterUUID
					dstCluster.(Cluster).AppList[k] = v
				}
				dstTotalCPU = dstCluster.(Cluster).ResourceList[0].Total
				dstTotalMEN = dstCluster.(Cluster).ResourceList[1].Total
				dstUsedCPU = iterRemove.Value.(Cluster).ResourceList[0].Used + dstCluster.(Cluster).ResourceList[0].Used
				dstUsedMEN = iterRemove.Value.(Cluster).ResourceList[1].Used + dstCluster.(Cluster).ResourceList[1].Used
				afterUtilization = (iterRemove.Value.(Cluster).ResourceList[0].Utilization*iterRemove.Value.(Cluster).ResourceList[0].Total+dstCluster.(Cluster).ResourceList[0].Utilization*dstCluster.(Cluster).ResourceList[0].Total)/dstCluster.(Cluster).ResourceList[0].Total*CPUWeight + (iterRemove.Value.(Cluster).ResourceList[1].Utilization*iterRemove.Value.(Cluster).ResourceList[1].Total+dstCluster.(Cluster).ResourceList[1].Utilization*dstCluster.(Cluster).ResourceList[1].Total)/dstCluster.(Cluster).ResourceList[1].Total*MEMWeight
				fmt.Println("after", afterUtilization)
				dstCluster.(Cluster).ResourceList[0] = Resource{
					ResourceType:  0,
					Used:          dstUsedCPU,
					Total:         dstTotalCPU,
					Utilization:   (iterRemove.Value.(Cluster).ResourceList[0].Utilization*iterRemove.Value.(Cluster).ResourceList[0].Total + dstCluster.(Cluster).ResourceList[0].Utilization*dstCluster.(Cluster).ResourceList[0].Total) / dstCluster.(Cluster).ResourceList[0].Total,
					ReservedRatio: 0,
				}
				dstCluster.(Cluster).ResourceList[1] = Resource{
					ResourceType:  1,
					Used:          dstUsedMEN,
					Total:         dstTotalMEN,
					Utilization:   (iterRemove.Value.(Cluster).ResourceList[1].Utilization*iterRemove.Value.(Cluster).ResourceList[1].Total + dstCluster.(Cluster).ResourceList[1].Utilization*dstCluster.(Cluster).ResourceList[1].Total) / dstCluster.(Cluster).ResourceList[1].Total,
					ReservedRatio: 0,
				}

				if beforeUtilization < LoadLowThreshold && afterUtilization >= LoadLowThreshold {
					lowNum--
				}
				rec.ClusterOffID.PushBack(iterRemove.Value.(Cluster).UUID)
				clusters.Remove(iterRemove)
				lowNum--
				break
			}
		}
	} else {
		for iter := lowList.Back().Prev(); iter != lowList.Front(); iter = iter.Prev() {
			fmt.Println(lowList.Front().Value.(ClustersSort).ClustersUtilization + iter.Value.(ClustersSort).ClustersUtilization)
			if (lowList.Front().Value.(ClustersSort).ClustersUtilization+iter.Value.(ClustersSort).ClustersUtilization < LoadHighThreshold) && (lowList.Front().Value.(ClustersSort).usedCPU+iter.Value.(ClustersSort).usedCPU < iter.Value.(ClustersSort).totalCPU) && (lowList.Front().Value.(ClustersSort).usedMEM+iter.Value.(ClustersSort).usedMEM < iter.Value.(ClustersSort).totalMEM) {
				for iterDst := clusters.Front(); iterDst != nil; iterDst = iterDst.Next() {
					if (iter.Value.(ClustersSort)).ClusterName == (iterDst.Value.(Cluster)).Name {
						dstCluster = iterDst.Value
					}
				}
				for iterRemove := clusters.Front(); iterRemove != nil; iterRemove = iterRemove.Next() {
					if (iter.Value.(ClustersSort)).ClusterName == (iterRemove.Value.(Cluster)).Name {
						removeFlag = 1
						for k, v := range (iterRemove.Value.(Cluster)).AppList {
							rec.AppToCluster[v] = iter.Value.(Cluster).UUID
							dstCluster.(Cluster).AppList[k] = v
						}
						dstTotalCPU = iterRemove.Value.(Cluster).ResourceList[0].Total
						dstTotalMEN = iterRemove.Value.(Cluster).ResourceList[1].Total
						dstUsedCPU = iterRemove.Value.(Cluster).ResourceList[0].Used + dstCluster.(Cluster).ResourceList[0].Used
						dstUsedMEN = iterRemove.Value.(Cluster).ResourceList[1].Used + dstCluster.(Cluster).ResourceList[1].Used
						afterUtilization = (iterRemove.Value.(Cluster).ResourceList[0].Utilization*iterRemove.Value.(Cluster).ResourceList[0].Total+dstCluster.(Cluster).ResourceList[0].Utilization*dstCluster.(Cluster).ResourceList[0].Total)/dstCluster.(Cluster).ResourceList[0].Total*CPUWeight + (iterRemove.Value.(Cluster).ResourceList[1].Utilization*iterRemove.Value.(Cluster).ResourceList[1].Total+dstCluster.(Cluster).ResourceList[1].Utilization*dstCluster.(Cluster).ResourceList[1].Total)/dstCluster.(Cluster).ResourceList[1].Total*MEMWeight

						dstCluster.(Cluster).ResourceList[0] = Resource{
							ResourceType:  0,
							Used:          dstUsedCPU,
							Total:         dstTotalCPU,
							Utilization:   (iterRemove.Value.(Cluster).ResourceList[0].Utilization*iterRemove.Value.(Cluster).ResourceList[0].Total + dstCluster.(Cluster).ResourceList[0].Utilization*dstCluster.(Cluster).ResourceList[0].Total) / dstCluster.(Cluster).ResourceList[0].Total,
							ReservedRatio: 0,
						}
						dstCluster.(Cluster).ResourceList[1] = Resource{
							ResourceType:  1,
							Used:          dstUsedMEN,
							Total:         dstTotalMEN,
							Utilization:   (iterRemove.Value.(Cluster).ResourceList[1].Utilization*iterRemove.Value.(Cluster).ResourceList[1].Total + dstCluster.(Cluster).ResourceList[1].Utilization*dstCluster.(Cluster).ResourceList[1].Total) / dstCluster.(Cluster).ResourceList[1].Total,
							ReservedRatio: 0,
						}

						if beforeUtilization < LoadLowThreshold && afterUtilization >= LoadLowThreshold {
							lowNum--
						}
						rec.ClusterOffID.PushBack(iterRemove.Value.(Cluster).UUID)
						clusters.Remove(iterRemove)
						lowNum--
						break
					}
				}
			}

			if removeFlag == 1 {
				break
			}

		}

	}
	fmt.Println("lowNum", lowNum)
	fmt.Println("rec", rec.AppToCluster)

	return 0
}

//Assess the power off function
func (p *dynamicScalingAlgorithom) powerOffEva(clusters *list.List, historyData *list.List) Recommendation {

	removeFlag = 0
	rec := Recommendation{
		ID:           time.Now().String(),
		AppToCluster: map[string]string{},
		ClusterOffID: list.New(),
	}

	//Get the clustters' usage by history sanapshot datas
	currentClusterLoad := historyData.Back().Value.(PerfSnapshot).ClusterPerfData


	highLoadClusters := map[string]float32{}
	middleLoadClusters := map[string]float32{}
	lowLoadClusters := map[string]float32{}

	//Get the cluster's weighted usage
	for cltName, cluster := range currentClusterLoad {
		utilization := cluster.getWeightUtilization()
		if utilization < LoadLowThreshold {
			lowLoadClusters[cltName] = utilization
		} else if utilization < LoadHighThreshold {
			middleLoadClusters[cltName] = utilization
		} else {
			highLoadClusters[cltName] = utilization
		}
	}
	//Get the cluster's num in low zone
	lowNum = len(lowLoadClusters)

	if lowNum == 0 {
		return rec
	}
	for {
		if p.powerOffAction(clusters, &rec) < 0 {
			break
		}

		if lowNum == 0 || removeFlag == 0 {
			break
		}
	}
	return rec

}
