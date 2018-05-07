package algorithm

import (
	"container/list"
	"fmt"
	"testing"
)

func TestTfiltByRule(t *testing.T) {
	// app configuration
	TestAppName := "myTestApp"

	// clusters configuration
	myClustersList := list.New()
	clust1 := Cluster{
		UUID:         "3-c1",
		Name:         "clust1",
		Label:        "highCpu",
		EndPoint:     "http://123.23.2.1:8080",
		AppList:      map[string]string{"app1": "a1", "app2": "a2"},
		ResourceList: map[int]Resource{CPU: {CPU, 30, 100, 0, 0}, MEM: {MEM, 50, 100, 0, 0}},
	}
	clust2 := Cluster{
		UUID:         "3-c2",
		Name:         "clust2",
		Label:        "highCpu",
		EndPoint:     "http://123.23.2.2:8080",
		AppList:      map[string]string{"app3": "a3", "app4": "a4", "app5": "a5"},
		ResourceList: map[int]Resource{CPU: {CPU, 40, 100, 0, 0}, MEM: {MEM, 10, 100, 0, 0}},
	}
	clust3 := Cluster{
		UUID:         "3-c3",
		Name:         "clust3",
		Label:        "hihgIO",
		EndPoint:     "http://123.23.2.3:8080",
		AppList:      map[string]string{"app6": "a6"},
		ResourceList: map[int]Resource{CPU: {CPU, 30, 100, 0, 0}, MEM: {MEM, 20, 100, 0, 0}},
	}

	myClustersList.PushBack(clust1)
	myClustersList.PushBack(clust2)
	myClustersList.PushBack(clust3)

	// Rules configuration
	myRulesList := list.New()
	rule1 := Rule{
		RuleType: ANonAffinity, Name: "rule1", Apps: "myTestApp,app2", Clusters: "", //myTestApp,app2反亲和性
	}
	rule2 := Rule{
		RuleType: CNonAffinity, Name: "rule2", Apps: "", Clusters: "clust2", //myTestApp,clust2反亲和性
	}
	rule3 := Rule{
		RuleType: AAffinity, Name: "rule3", Apps: "myTestApp,app6", Clusters: "", //myTestApp,app6亲和性
	}
	myRulesList.PushBack(&rule1)
	myRulesList.PushBack(&rule2)
	myRulesList.PushBack(&rule3)

	var p placementAlgorithom
	filterClusters := p.filtByRule(TestAppName, myClustersList, myRulesList)
	for clt := filterClusters.Front(); clt != nil; clt = clt.Next() {
		fmt.Println("There is selected cluster")
	}
}

func TestTgetLocation(t *testing.T) {

	TestApp := Application{
		Name:         "myTestApp",
		Label:        "highio",
		ResourceList: map[int]Resource{CPU: {CPU, 0, 10, 0, 0}, MEM: {MEM, 0, 10, 0, 0}},
	} //appName

	// clusters configuration
	myClustersList := list.New()
	clust1 := Cluster{
		UUID:         "3-c1",
		Name:         "clust1",
		Label:        "highCpu",
		EndPoint:     "http://123.23.2.1:8080",
		AppList:      map[string]string{"app1": "a1", "app2": "a2"},
		ResourceList: map[int]Resource{CPU: {CPU, 30, 100, 0, 0}, MEM: {MEM, 50, 100, 0, 0}},
	}
	clust2 := Cluster{
		UUID:         "3-c2",
		Name:         "clust2",
		Label:        "highCpu",
		EndPoint:     "http://123.23.2.2:8080",
		AppList:      map[string]string{"app3": "a3", "app4": "a4", "app5": "a5"},
		ResourceList: map[int]Resource{CPU: {CPU, 40, 100, 0, 0}, MEM: {MEM, 10, 100, 0, 0}},
	}
	clust3 := Cluster{
		UUID:         "3-c3",
		Name:         "clust3",
		Label:        "hihgIO",
		EndPoint:     "http://123.23.2.3:8080",
		AppList:      map[string]string{"app6": "a6"},
		ResourceList: map[int]Resource{CPU: {CPU, 30, 100, 0, 0}, MEM: {MEM, 20, 100, 0, 0}},
	}

	myClustersList.PushBack(clust1)
	myClustersList.PushBack(clust2)
	myClustersList.PushBack(clust3)

	// Rules configuration
	myRulesList := list.New()
	rule1 := Rule{
		RuleType: ANonAffinity, Name: "rule1", Apps: "myTestApp,app2", Clusters: "", //app1,app2反亲和性
	}
	rule2 := Rule{
		RuleType: CNonAffinity, Name: "rule2", Apps: "", Clusters: "clust2", //myTestApp,clust2反亲和性
	}
	rule3 := Rule{
		RuleType: AAffinity, Name: "rule3", Apps: "myTestApp,app6", Clusters: "", //myTestApp,app6亲和性
	}
	myRulesList.PushBack(&rule1)
	myRulesList.PushBack(&rule2)
	myRulesList.PushBack(&rule3)

	var p placementAlgorithom
	locateClust := p.getLocation(&TestApp, myClustersList, myRulesList, 0)
	fmt.Println("TargetCluster UUID is:", locateClust.UUID)
}
