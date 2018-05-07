package algorithm

import (
	"container/list"
	"fmt"
	"strings"
)

type placementAlgorithom struct{}

func (p *placementAlgorithom) Calculate(param interface{}) interface{} {
	fmt.Println("Calculate...")
	if param == nil {
		fmt.Println("Error : placement param is nil")
	}

	args := param.(PlacementParam)
	if args.App == nil || args.Clusters == nil {
		fmt.Println("Error : placement param App or Clusters is nil")
	}

	ret := p.getLocation(args.App, args.Clusters, args.Rules, 0)
	return ret
}

//filtByRule is a interface to filter Clusters by rules which related to NewApp
func (p *placementAlgorithom) filtByRule(appName string, clusters *list.List, rules *list.List) *list.List {
	mapClusters := map[string]*Cluster{}
	deployedApps := map[string]string{}
	//Traverse the ClusterList and to Form a map,
	//which the key is AppName, and the value is ClusterName
	for ec := clusters.Front(); ec != nil; ec = ec.Next() {
		cluster := ec.Value.(Cluster)
		mapClusters[cluster.Name] = &cluster
		for key := range cluster.AppList {
			deployedApps[key] = cluster.Name
		}
	}

	//includeMap is Storing the clusters will be recommended
	includeMap := map[string]*Cluster{}
	//excludeMap is Storing the clusters excluded by rules
	excludeMap := map[string]*Cluster{}
	//ruleMap is the Union by RuleType, Such as affinity/non-affinity with apps/clusters
	ruleMap := map[int]map[string]*Cluster{
		AAffinity:    includeMap,
		CAffinity:    includeMap,
		ANonAffinity: excludeMap,
		CNonAffinity: excludeMap,
	}
	if rules == nil {
		fmt.Println("Info : placement rules is nil.")
	} else {
		for e := rules.Front(); e != nil; e = e.Next() {
			rule := e.Value.(*Rule)

			//1. Remove the rules is not related to the new App
			if !strings.Contains(rule.Apps, appName) {
				continue
			}

			//2. Use the rules among apps, and put the affinity clusters to includeMap,
			// and put non-affinity clusters to excludemap
			if rule.RuleType == AAffinity || rule.RuleType == ANonAffinity {
				apps := strings.Split(rule.Apps, ",")
				for index := 0; index < len(apps); index++ {
					if apps[index] != appName {
						if cltName, ok := deployedApps[apps[index]]; ok {
							ruleMap[rule.RuleType][cltName] = mapClusters[cltName]
						}
					}
				}
			}

			//3. Use the rules between app and cluster,
			// and move clusters input maps which like the step 2
			if rule.RuleType == CAffinity || rule.RuleType == CNonAffinity {
				clusterNames := strings.Split(rule.Clusters, ",")
				for index := 0; index < len(clusterNames); index++ {
					cltName := clusterNames[index]
					ruleMap[rule.RuleType][cltName] = mapClusters[cltName]
				}
			}
		}
	}

	//If the includeMap is empty, and then put the all clusters to it;
	if len(includeMap) == 0 {
		includeMap = mapClusters
	}

	//Form the filterList by through includeMap and throw away the excludeMap
	filterList := list.New()
	for key, value := range includeMap {
		if _, ok := excludeMap[key]; !ok {
			filterList.PushBack(value)
		}
	}

	return filterList
}

//getLocation is a interface, which will be use by placement Calculate
func (p *placementAlgorithom) getLocation(app *Application, clusters *list.List, rules *list.List, policy int) (clt Cluster) {
	minUsedRate := float32(1)
	//filter the clusters by rules, which is related to newApp
	filterClusters := p.filtByRule(app.Name, clusters, rules)
	if filterClusters.Len() == 0 {
		fmt.Println("Error : !!!can't get clusters by rules.!!!")
	}

	for e := filterClusters.Front(); e != nil; e = e.Next() {
		maxUsedRate := float32(0)
		res := e.Value.(*Cluster)

		//1. Choose the max UseRate in CPU and MEM, and use it to involved computing
		for key, value := range app.ResourceList {
			var usedRate = (res.ResourceList[key].Used + value.Total) / res.ResourceList[key].Total
			if maxUsedRate < usedRate {
				maxUsedRate = usedRate
			}
		}

		//2. Use the minimum UseRate Clusters to Recommended
		if minUsedRate > maxUsedRate {
			minUsedRate = maxUsedRate
			clt = *res
		}
	}
	//Testing ...
	if T.GTflag {
		p.AlgStatistics(app, filterClusters)
	}

	return
}

//AlgStatistics is a statistics interface
func (p *placementAlgorithom) AlgStatistics(app *Application, cl *list.List) {
	T.GSlistReset(T.GTflag)
	minUsedRate := float32(1)

	for e := cl.Front(); e != nil; e = e.Next() {
		maxUsedRate := float32(0)
		res := e.Value.(*Cluster)
		var mUsage float32
		var cUsage float32
		for key, value := range app.ResourceList {
			var usedRate = (res.ResourceList[key].Used + value.Total) / res.ResourceList[key].Total
			if key == 1 {
				mUsage = usedRate
			} else {
				cUsage = usedRate
			}

			if maxUsedRate < usedRate {
				maxUsedRate = usedRate
			}
		}
		T.MakeStatisticsNode(res.UUID, res.EndPoint, cUsage, mUsage, len(res.AppList))
		if minUsedRate > maxUsedRate {
			minUsedRate = maxUsedRate
		}
	}
	T.DisplayStatistics()
}
