package controler

import (
	"encoding/json"
	"fmt"
	alg "scheduler/algorithm"
	"testing"
)

func InitPath(JSON string) interface{} {
	JSONString := JSON
	if JSONString == "" {
		JSONString = `{"App":{"apiVersion": "v1beta3", "kind": "Pod", "metadata": { "name": "explorer", "namespace": "" }, "spec": { "containers": [ { "args": [ "-port=8080" ], "image": "gcr.io/google_containers/explorer:1.0", "name": "explorer", "ports": [ { "containerPort": 8080, "protocol": "TCP" } ], "resources": { "limits": { "cpu": "0.5", "disk": "1000", "memory": "100" } }, "volumeMounts": [ { "mountPath": "/mount/test-volume", "name": "test-volume" } ] } ], "volumes": [ { "emptyDir": {}, "name": "test-volume" } ] }},"Path":{"AppAffinity":["app1","app3"],"AppUnaffinity":["app4","app5","app6"],"ClusterAffinity":["app7","clst1"],"ClusterUnaffinity":["app8","clst2"],"Label":"highio","Storage":"300"}}`
	}
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(JSONString), &data); err != nil {
		fmt.Println(err)
	}
	return data["Path"]
}
func TestParseAppRules(t *testing.T) {
	path := InitPath("")
	if path == nil {
		t.Fatal("InitPath")
	}
	rules := ParseAppRules("app8", path)
	if len(rules) != 4 {
		t.Fatal("The Number of rules is Not Correct")
	}
	if rules[0].RuleType == 0 {
		if rules[0].Name != "AppAffinity" {
			t.Fatal(rules[0].Name)
		}
		if rules[0].Apps == "App1,App3" {
			t.Fatal(rules[0].Name)
		}
	}
	if rules[1].RuleType == 1 {
		if rules[1].Name != "AppUnaffinity" {
			t.Fatal(rules[1].Name)
		}
		if rules[1].Apps == "app4,app5,app6" {
			t.Fatal(rules[1].Name)
		}
	}
	if rules[2].RuleType == 2 {
		if rules[2].Name != "ClusterAffinity" {
			t.Fatal(rules[2].Name)
		}
		if rules[2].Apps == "app7,clst1" {
			t.Fatal(rules[2].Name)
		}
	}
	if rules[3].RuleType == 3 {
		if rules[3].Name != "ClusterUnaffinity" {
			t.Fatal(rules[3].Name)
		}
		if rules[3].Apps == "app8,clst2" {
			t.Fatal(rules[3].Name)
		}
	}

	path = nil
	rules = ParseAppRules("app8", path)
	if len(rules) != 0 {
		t.Fatal("The Number of rules should be 0")
	}
}
func TestTformRuleData(t *testing.T) {
	key := "AppAffinity"
	value := "value_test"

	rule := formRuleData(key, value)
	if rule.RuleType != alg.AAffinity && rule.Apps != value {
		t.Fatal(key)
	}

	key = "AppUnaffinity"
	rule = formRuleData(key, value)
	if rule.RuleType != alg.ANonAffinity && rule.Apps != value {
		t.Fatal(key)
	}

	key = "ClusterAffinity"
	rule = formRuleData(key, value)
	if rule.RuleType != alg.CAffinity && rule.Clusters != value {
		t.Fatal(key)
	}

	key = "ClusterUnaffinity"
	rule = formRuleData(key, value)
	if rule.RuleType != alg.CNonAffinity && rule.Clusters != value {
		t.Fatal(key)
	}
}

func InitClusters() []ClusterData {
	clusters := []ClusterData{
		{UUID: "c111111111111111111", Name: "cluster1", Label: "highio,default", Enginetype: "K8S", Endpoint: "http://9.91.18.84:8080", CPU: 100, Memory: 100000, Disk: 100000000000000, Status: "AVAILABLE"},
		{UUID: "c222222222222222222", Name: "cluster2", Label: "highcpu,highio", Enginetype: "K8S", Endpoint: "http://9.91.18.84:8080", CPU: 200, Memory: 200000, Disk: 200000000000000, Status: "AVAILABLE"},
		{UUID: "c333333333333333333", Name: "cluster3", Label: "highio", Enginetype: "K8S", Endpoint: "http://9.91.18.84:8080", CPU: 300, Memory: 300000, Disk: 300000000000000, Status: "AVAILABLE"},
		{UUID: "c444444444444444444", Name: "cluster4", Label: "highio", Enginetype: "K8S", Endpoint: "http://9.91.18.84:8080", CPU: 400, Memory: 400000, Disk: 400000000000000, Status: "AVAILABLE"},
		{UUID: "c555555555555555555", Name: "cluster5", Label: "highio", Enginetype: "K8S", Endpoint: "http://9.91.18.84:8080", CPU: 400, Memory: 400000, Disk: 400000000000000, Status: "UNAVAILABLE"},
		{UUID: "c666666666666666666", Name: "cluster6", Label: "highio", Enginetype: "CF", Endpoint: "http://9.91.18.84:8080", CPU: 400, Memory: 400000, Disk: 400000000000000, Status: "AVAILABLE"},
		{UUID: "c777777777777777777", Name: "cluster7", Label: "highio", Enginetype: "K8S", Endpoint: "http://10.10.1.0:8080", CPU: 100, Memory: 100000, Disk: 100000000000000, Status: "AVAILABLE"},
		{UUID: "c888888888888888888", Name: "cluster8", Label: "lowcpu", Enginetype: "K8S", Endpoint: "http://9.91.18.84:8080", CPU: 200, Memory: 200000, Disk: 200000000000000, Status: "AVAILABLE"},
	}
	return clusters
}
func InitApps() []AppData {
	apps := []AppData{
		{UUID: "a111111", Name: "app1", Label: "highio", ClusterUUID: "c111111111111111111", CPU: 1, Memory: 10000, Disk: 100000, ReservedRatio: 85},
		{UUID: "a222222", Name: "app2", Label: "highcpu", ClusterUUID: "c222222222222222222", CPU: 2, Memory: 20000, Disk: 200000, ReservedRatio: 85},
		{UUID: "a333333", Name: "app3", Label: "highio", ClusterUUID: "c333333333333333333", CPU: 3, Memory: 30000, Disk: 300000, ReservedRatio: 85},
		{UUID: "a444444", Name: "app4", Label: "highcpu", ClusterUUID: "c444444444444444444", CPU: 4, Memory: 40000, Disk: 400000, ReservedRatio: 85},
	}
	return apps
}
func TestPrepareData(t *testing.T) {
	apps := InitApps()
	clusters := InitClusters()
	path := InitPath("")

	var app AppData
	app.Name = "app8"
	app.CPU = 1
	app.Memory = 60
	c := PrepareData(&app, path, clusters, apps, nil)
	if c.AppList == nil {
		t.Fatal("Applist is nil")
	}
	if c.RuleList == nil {
		t.Fatal("RuleList is nil")
	}
	if c.ClstrList == nil {
		t.Fatal("ClstrList is nil")
	}

}
func TestGetCluster(t *testing.T) {
	apps := InitApps()
	clusters := InitClusters()
	path := InitPath("")

	var app AppData
	app.Name = "app8"
	app.CPU = 1
	app.Memory = 60
	c := PrepareData(&app, path, clusters, apps, nil)
	fmt.Println("Default", c.Default)
	//c.DisplayAppList()
	//c.DisplayRuleList()
	//c.DisplayClstrList()
	ep, uuid := c.GetCluster()
	fmt.Println("Recommend Cluster UUId:", uuid, "  Endpoint:", ep)
	if ep == "" {
		t.Fatal("No Cluster Recommend")
	}
	if uuid == "" {
		t.Fatal("No Cluster Recommend")
	}

	JSONString := `{"App":{"apiVersion": "v1beta3", "kind": "Pod", "metadata": { "name": "explorer", "namespace": "" }, "spec": { "containers": [ { "args": [ "-port=8080" ], "image": "gcr.io/google_containers/explorer:1.0", "name": "explorer", "ports": [ { "containerPort": 8080, "protocol": "TCP" } ], "resources": { "limits": { "cpu": "0.5", "disk": "1000", "memory": "100" } }, "volumeMounts": [ {"mountPath": "/mount/test-volume", "name": "test-volume" } ] } ], "volumes": [ { "emptyDir": {}, "name": "test-volume" } ] }},"Path":{"AppAffinity":["app1","app3"],"AppUnaffinity":["app3","app2","app6"],"Label":"highcpu","Storage":"300"}}`
	path = InitPath(JSONString)
	app.Name = "app9"
	app.CPU = 5
	app.Memory = 500
	c = PrepareData(&app, path, clusters, apps, nil)
	ep, uuid = c.GetCluster()
	fmt.Println("Recommend Cluster UUId:", uuid, "  Endpoint:", ep)
	if ep != "http://9.91.18.84:8080" {
		t.Fatal("Error, should have no Cluster to Recommend")
	}
	if uuid != "c111111111111111111" {
		t.Fatal("Error, should have no Cluster to Recommend")
	}

	JSONString = `{"App":{"apiVersion": "v1beta3", "kind": "Pod", "metadata": { "name": "explorer", "namespace": "" }, "spec": { "containers": [ { "args": [ "-port=8080" ], "image": "gcr.io/google_containers/explorer:1.0", "name": "explorer", "ports": [ { "containerPort": 8080, "protocol": "TCP" } ], "resources": { "limits": { "cpu": "0.5", "disk": "1000", "memory": "100" } }, "volumeMounts": [ {"mountPath": "/mount/test-volume", "name": "test-volume" } ] } ], "volumes": [ { "emptyDir": {}, "name": "test-volume" } ] }},"Path":{"AppAffinity":["app1","app3"],"AppUnaffinity":["app3","app1","app6"],"Label":"highio","Storage":"300"}}`
	path = InitPath(JSONString)
	c = PrepareData(&app, path, clusters, apps, nil)
	ep, uuid = c.GetCluster()
	if ep != "" {
		t.Fatal("No Cluster Recommend")
	}
	if uuid != "" {
		t.Fatal("No Cluster Recommend")
	}
}

func TestGetExtendByKey(t *testing.T) {
	tt := SetExtend(nil, "label", "highio")
	tt = SetExtend(tt, "exclude", "http://9.91.18.84/")
	v := GetExtendByKey(tt, "label")
	if v != "highio" {
		t.Fatal(tt)
	}
}
