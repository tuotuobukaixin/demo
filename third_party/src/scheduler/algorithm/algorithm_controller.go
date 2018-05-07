/*
Package Algorithm
Copyright(c) Huawei Technologies Co.,Ltd
Created: Fri Sep 18 17:35:39 UTC 2015
Filename: algorithom_controller.go
Description:

*/

package algorithm

import (
	"container/list"
	"fmt"
)

//algorith function's type
const (
	PlacementAlgorithm         = 0 //Placement
	DynamicLoadblanceAlgorithm = 1 //Dynamic Loadblance
	DynamicScalingAlgorithm    = 2 //Dynamic Sacling
)

//
const (
	LoadLowThreshold  = 0.20
	LoadHighThreshold = 0.81
)

//
const FiveMinutesDataCount = 5

//
const (
	CPUWeight = 1
	MEMWeight = 0
)

//PlacementParam is Controler provide
// @Title PlacementParam
// @Description
// New App,
type PlacementParam struct {
	App      *Application //New app's Resource
	Clusters *list.List   //Clusters's List <Cluster>
	Rules    *list.List   //Rules's list<*Rule>
	Policy   int          //The placement policy: balance, aggreage and dispersal
}

//DynamicSchedulingParam is dynamicMonitor provide
// @Title DynamicSchedulingParam
// @Description
type DynamicSchedulingParam struct {
	MigrationThreshold float32    //Dynamic Scheduler Policy
	Clusters           *list.List //Cluster stat: list<Cluster>
	Rules              *list.List //Rules's list:  list<Rule>
	HistoryData        *list.List //History's list:  list<PerfSnapshot>
}

//Recommendation is the algorithm return to controler
// @Title Recommendation
// @Description
type Recommendation struct {
	ID           string            //TimeStamp
	AppToCluster map[string]string //Suggestion: App --> Cluster, map<AppUUID>ClusterUUID
	ClusterOffID *list.List        //The Cluster Should be Power Off: list<ClusterUUID>
}

//BaseAlgorithm is the base interface for algorithm func
// @Title BaseAlgorithm
// @Description
type BaseAlgorithm interface {
	Calculate(interface{}) interface{}
}

//CreateAlgorithm is the interface for user will use
// @Title CreateAlgorithm
// @Description
func CreateAlgorithm(algType int) BaseAlgorithm {

	switch algType {
	case PlacementAlgorithm:
		return &placementAlgorithom{}

	case DynamicLoadblanceAlgorithm:
		return &dynamicLoadBlanceBFDAlgorithom{}

	case DynamicScalingAlgorithm:
		return &dynamicScalingAlgorithom{}
	default:
		return nil
	}
}

//Scheduling is the only interface which algorthm model expose
// @Title Scheduling
// @Description
func Scheduling(param interface{}) interface{} {

	switch param.(type) {

	case *PlacementParam:
		alg := CreateAlgorithm(PlacementAlgorithm)
		return alg.Calculate(*param.(*PlacementParam))
	case PlacementParam:
		alg := CreateAlgorithm(PlacementAlgorithm)
		return alg.Calculate(param)
	case DynamicSchedulingParam:
		return dynamicSchedulerNew(param.(DynamicSchedulingParam))
	case *DynamicSchedulingParam:
		return dynamicSchedulerNew(*param.(*DynamicSchedulingParam))
	default:
		fmt.Println("Error: !!!  Scheduling param type error  !!!")
		return nil
	}
}

//dynamicSchedulerNew is the interface which the dynamic algorithm entry
// @Title dynamicSchedulerNew
// @Description
func dynamicSchedulerNew(param DynamicSchedulingParam) interface{} {

	//1.Check the param
	if param.HistoryData.Len() == 0 || param.Clusters.Len() == 0 {
		fmt.Println("Error: !!! data is empty !!!")
		return nil
	}

	//2.Load Balance First, and get the recommends
	recommendationList := list.New()
	recommendation := Recommendation{}
	alg := CreateAlgorithm(DynamicLoadblanceAlgorithm)
	retlb, cltList := alg.(*dynamicLoadBlanceBFDAlgorithom).CalculateNew(param)

	//3.Separate the clusters: high/middle/low zone. such as: 0%~20%~81%~100%
	highLoads, midLoads, lowLoads := map[string]float32{}, map[string]float32{}, map[string]float32{}
	for _, cluster := range cltList {
		//Get the usage of Clusters with Weight
		utilization := cluster.getWeightUtilization()
		if utilization < LoadLowThreshold {
			lowLoads[cluster.Name] = utilization
		} else if utilization < LoadHighThreshold {
			midLoads[cluster.Name] = utilization
		} else {
			highLoads[cluster.Name] = utilization
		}
	}
	fmt.Println("highLoad=", highLoads, "middleLoad=", midLoads, "lowLoad=", lowLoads)

	//4. If it has high zones, then tell to startup cluster(s)
	if len(highLoads) > 0 {
		fmt.Println("Will execute new cluster algorithm...")
	}

	//5. If it has low zones, then scale the cluster(s), will tell someone power off
	if len(lowLoads) > 0 {
		fmt.Println("Will execute reduce cluster algorithm...")
		alg = CreateAlgorithm(DynamicScalingAlgorithm)
		ret := alg.Calculate(param)
		if ret != nil && len(ret.(Recommendation).AppToCluster) > 0 {
			recommendation = ret.(Recommendation)
			fmt.Print("****************** result excute scaling:", "AppToCluster:", recommendation.AppToCluster, "closeUUID=")
			for e := recommendation.ClusterOffID.Front(); e != nil; e = e.Next() {
				fmt.Print(e.Value, " ")
			}
			fmt.Println()
			recommendationList.PushBack(ret)
			return *recommendationList
		}
	}

	//6. No recommendation suggested,
	//cause of  that there are 2 app, and they deployed 2 clusters
	if retlb == nil {
		fmt.Println("Info: !!! can't make load blance,scaling,expanding !!!")
		return nil
	}

	//7.If No recommendation to scaling, then start loadbalance
	recommendation = retlb.(Recommendation)
	if len(recommendation.AppToCluster) > 0 {
		fmt.Println("****************** result need blance AppToCluster:", recommendation.AppToCluster)
		recommendationList.PushBack(recommendation)
	}
	return *recommendationList
}

//dynamicScheduler is the old interface which dynamic algorithm entry
// @Title dynamicScheduler
// @Description
func dynamicScheduler(param DynamicSchedulingParam) list.List {
	//1.Load Balance First, and get the recommends
	recommendationList := list.New()
	alg := CreateAlgorithm(DynamicLoadblanceAlgorithm)
	ret := alg.Calculate(param)

	recommendation := Recommendation{}
	if ret != nil {
		recommendation = ret.(Recommendation)
		if len(recommendation.AppToCluster) > 0 {
			fmt.Println("****************** result excute blance:", "AppToCluster:", recommendation.AppToCluster)
			recommendationList.PushBack(recommendation)
			return *recommendationList
		}
	}

	//2.Separate the clusters: high/middle/low zone. such as: 0%~20%~81%~100%
	appLoad, clusterLoad := get5MinAverageLoad(param)
	high, _, low := simulateMigration(appLoad, clusterLoad, recommendation)

	//3.If it has high zones, then tell to startup cluster(s)
	if len(high) > 0 {
		fmt.Println("Will execute new cluster algorithm...")
	}

	//4.If it has low zones, then scale the cluster(s), will tell someone power off
	if len(low) > 0 {
		fmt.Println("Will execute reduce cluster algorithm...")
		alg = CreateAlgorithm(DynamicScalingAlgorithm)
		ret = alg.Calculate(param)
		if ret != nil {
			recommendation = ret.(Recommendation)
			fmt.Print("****************** result excute scaling:", "AppToCluster:", recommendation.AppToCluster, "closeUUID=")
			for e := recommendation.ClusterOffID.Front(); e != nil; e = e.Next() {
				fmt.Println(" ", e.Value.(string))
			}
			fmt.Println()
			recommendationList.PushBack(ret)
		}
	}

	//5. There is no recommendation suggest, the list is nil
	return *recommendationList
}

//get5MinAverageLoad is the average load interface, which use 5min data before
// @Title get5MinAverageLoad
// @Description
func get5MinAverageLoad(param DynamicSchedulingParam) (map[string]Application, map[string]Cluster) {
	el := param.HistoryData.Back()
	currentAppLoad := el.Value.(PerfSnapshot).AppPerfData
	currentClusterLoad := el.Value.(PerfSnapshot).ClusterPerfData

	//Compose a temp complex map to note the average usage of all resource
	appLoadTmp := map[string]map[int]*Resource{} //Should Store the pointer
	for appName, app := range currentAppLoad {
		resMap := map[int]*Resource{}
		for resType, res := range app.ResourceList {
			tmpRes := res
			resMap[resType] = &tmpRes
		}
		appLoadTmp[appName] = resMap
	}

	cltLoadTmp := map[string]map[int]*Resource{} //Should Store the pointer
	for cltName, clt := range currentClusterLoad {
		resMap := map[int]*Resource{}
		for resType, res := range clt.ResourceList {
			tmpRes := res
			resMap[resType] = &tmpRes
		}
		cltLoadTmp[cltName] = resMap
	}

	//Compute the datas of average load in 5 mintues
	i := 0
	el = el.Prev()
	for ; el != nil; el = el.Prev() {
		i++
		if i >= FiveMinutesDataCount {
			break
		}
		for appName, app := range el.Value.(PerfSnapshot).AppPerfData { //Every Application
			for resType, res := range app.ResourceList {
				appLoadTmp[appName][resType].Utilization += res.Utilization
			}
		}

		for cltName, clt := range el.Value.(PerfSnapshot).ClusterPerfData { //Every Cluster
			for resType, res := range clt.ResourceList {
				cltLoadTmp[cltName][resType].Utilization += res.Utilization
			}
		}
	}

	for appName, resMap := range appLoadTmp {
		for resType, res := range resMap {
			res.Utilization = res.Utilization / FiveMinutesDataCount
			currentAppLoad[appName].ResourceList[resType] = *res
		}
	}

	for cltName, resMap := range cltLoadTmp {
		for resType, res := range resMap {
			res.Utilization = res.Utilization / FiveMinutesDataCount
			currentClusterLoad[cltName].ResourceList[resType] = *res
		}
	}

	return currentAppLoad, currentClusterLoad
}

//simulateMigration  is the interface which will do some moving in sandbox
// @Title simulateMigration
// @Description
func simulateMigration(appLoad map[string]Application, clusterLoad map[string]Cluster, recommend Recommendation) (map[string]float32, map[string]float32, map[string]float32) {
	//1.Deal the apps string by Cluster, and compose a map,
	//which key is AppName, and value is ClusterName
	deployedApps := map[string]string{}
	for cltName, cluster := range clusterLoad {
		for appName := range cluster.AppList {
			deployedApps[appName] = cltName
		}
	}

	//2. Try to scale, move someone to the destCluster
	for appName, dstCltName := range recommend.AppToCluster {
		srcCltName := deployedApps[appName] //srcCluster Name
		for resType, appRes := range appLoad[appName].ResourceList {
			//2.1 Add the resource to the destCluster
			dstCltRes := clusterLoad[dstCltName].ResourceList[resType]
			dstCltRes.Utilization = (dstCltRes.Utilization*dstCltRes.Total + appRes.Utilization*appRes.Total) / dstCltRes.Total //practical efficiency
			dstCltRes.Used = dstCltRes.Used + appRes.Total                                                                      //distributed volum
			clusterLoad[dstCltName].ResourceList[resType] = dstCltRes
			clusterLoad[dstCltName].AppList[appName] = appLoad[appName].UUID

			//2.2 Subtract the resource from the srcCluster
			srcCltRes := clusterLoad[srcCltName].ResourceList[resType]
			srcCltRes.Utilization = (srcCltRes.Utilization*srcCltRes.Total - appRes.Utilization*appRes.Total) / srcCltRes.Total
			srcCltRes.Used = srcCltRes.Used - appRes.Total
			clusterLoad[srcCltName].ResourceList[resType] = srcCltRes
			delete(clusterLoad[srcCltName].AppList, appName)
		}
	}

	//3.Separate the clusters: high/middle/low zone. such as: 0%~45%~81%~100%
	highLoadClusters := map[string]float32{}
	middleLoadClusters := map[string]float32{}
	lowLoadClusters := map[string]float32{}
	for cltName, cluster := range clusterLoad {
		utilization := cluster.getWeightUtilization()
		if utilization < LoadLowThreshold {
			lowLoadClusters[cltName] = utilization
		} else if utilization < LoadHighThreshold {
			middleLoadClusters[cltName] = utilization
		} else {
			highLoadClusters[cltName] = utilization
		}
	}
	fmt.Println("highLoadClusters=", highLoadClusters)
	fmt.Println("middleLoadClusters=", middleLoadClusters)
	fmt.Println("lowLoadClusters=", lowLoadClusters)
	return highLoadClusters, middleLoadClusters, lowLoadClusters
}

//getWeightUtilization is the Cluster's weighted Usage
// @Title getWeightUtilization
// @Description
func (p Cluster) getWeightUtilization() (ret float32) {
	weight := map[int]float32{CPU: CPUWeight, MEM: MEMWeight}
	for resType, res := range p.ResourceList {
		ret = ret + res.Utilization*weight[resType]
	}
	return
}

//getWeightUtilization is the Application's weighted Usage
// @Title getWeightUtilization
// @Description
func (p Application) getWeightUtilization() (ret float32) {
	weight := map[int]float32{CPU: CPUWeight, MEM: MEMWeight}
	for resType, res := range p.ResourceList {
		ret = ret + res.Utilization*weight[resType]
	}
	return
}

//getWeightUtilization is the independent weighted Usage
// @Title getWeightUtilization
// @Description
func getWeightUtilization(resources map[int]Resource) (ret float32) {
	weight := map[int]float32{CPU: CPUWeight, MEM: MEMWeight}
	for resType, res := range resources {
		ret = ret + res.Utilization*weight[resType]
	}
	return
}
