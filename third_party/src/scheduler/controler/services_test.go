package controler

import (
	"fmt"
	. "scheduler/algorithm"
	"strings"
	"testing"
)

func InitRule() []RuleData {
	rules := []RuleData{
		{RuleType: ANonAffinity, Name: "rule1", Apps: "app1,app2", Clusters: ""},
		{RuleType: AAffinity, Name: "rule2", Apps: "app3,app2", Clusters: ""},
		{RuleType: CNonAffinity, Name: "rule3", Apps: "", Clusters: "cluster3,cluster2"},
		{RuleType: CAffinity, Name: "rule4", Apps: "", Clusters: "cluster1,cluster2"},
	}
	return rules
}
func TestConvert2algRule(t *testing.T) {
	rules := []RuleData{
		{RuleType: ANonAffinity, Name: "rule1", Apps: "app1,app2", Clusters: ""},
		{RuleType: AAffinity, Name: "rule2", Apps: "app3,app2", Clusters: ""},
		{RuleType: CNonAffinity, Name: "rule3", Apps: "", Clusters: "cluster3,cluster2"},
		{RuleType: CAffinity, Name: "rule4", Apps: "", Clusters: "cluster1,cluster2"},
	}
	i := 0
	for _, rule := range rules {
		orule := Convert2algRule(rule)
		if orule.RuleType == ANonAffinity && orule.Name == "rule1" {
			if orule.Apps != "app1,app2" {
				t.Fatal(orule.Name)
			}
		}
		if orule.RuleType == AAffinity && orule.Name == "rule2" {
			if orule.Apps != "app3,app2" {
				t.Fatal(orule.Name)
			}
		}
		if orule.RuleType == CNonAffinity && orule.Name == "rule3" {
			if orule.Clusters != "cluster3,cluster2" {
				t.Fatal(orule.Name)
			}
		}
		if orule.RuleType == CAffinity && orule.Name == "rule4" {
			if orule.Clusters != "cluster1,cluster2" {
				t.Fatal(orule.Name)
			}
		}
		i++
	}
	if i != len(rules) {
		t.Fatal(i)
	}
}

func TestComposeRuleList_Empty(t *testing.T) {
	rules := []RuleData{}
	var ctrl Controler
	ctrl.ComposeRuleList(rules)

	i := 0
	for rule := ctrl.RuleList.Front(); rule != nil; rule = rule.Next() {
		i++
	}
	if i != 0 {
		t.Fatal(i)
	}

}
func TestComposeRuleList_Create(t *testing.T) {
	rules := InitRule()
	var ctrl Controler
	ctrl.ComposeRuleList(rules)

	i := 0
	for rule := ctrl.RuleList.Front(); rule != nil; rule = rule.Next() {
		rl := (rule.Value.(*Rule))
		if rl.RuleType == ANonAffinity {
			if rl.Name != "rule1" {
				t.Fatal(rl.Name)
			}
		}
		if rl.RuleType == AAffinity {
			if rl.Name != "rule2" {
				t.Fatal(rl.Name)
			}
			if rl.Name == "rule1" {
				t.Fatal(rl.Name)
			}
		}
		if rl.RuleType == CNonAffinity {
			if rl.Name != "rule3" {
				t.Fatal(rl.Name)
			}
		}
		if rl.RuleType == CAffinity {
			if rl.Name != "rule4" {
				t.Fatal(rl.Name)
			}
		}
		i++
	}
	if i != len(rules) {
		t.Fatal(i)
	}
}

//DisplayRuleList is output the rules in linklist
// @Title DisplayRuleList
// @Description  Display rules in RuleList
// @Accept  NULL
// @Param   NULL
func (ctrl *Controler) DisplayRuleList(str string) {
	fmt.Println(str)
	i := 0
	for rule := ctrl.RuleList.Front(); rule != nil; rule = rule.Next() {
		rl := (rule.Value.(Rule))
		fmt.Printf("[%d]rule:  %v \t%v \t%v \t%v\n", i, rl.RuleType, rl.Name, rl.Apps, rl.Clusters)
		i++
	}
}

func TestCheckClusterLabels(t *testing.T) {
	als_null := ""
	cl_null := ""

	als := "highio,highcpu"
	cl := "highio,highcpu,default"

	rt := CheckClusterLabels(als_null, cl, true)
	if rt != true {
		t.Fatal(als_null)
	}

	rt = CheckClusterLabels(als, cl_null, true)
	if rt != false {
		t.Fatal(cl_null)
	}

	rt = CheckClusterLabels(als, cl_null, false)
	if rt != true {
		t.Fatal(cl_null)
	}

	rt = CheckClusterLabels(als, cl, false)
	if rt != true {
		t.Fatal(cl)
	}

	als = "lowio"
	rt = CheckClusterLabels(als, cl, true)
	if rt != false {
		t.Fatal(cl)
	}

	als = "highio"
	rt = CheckClusterLabels(als, cl, true)
	if rt != true {
		t.Fatal(cl)
	}

	cl = "highio"
	als = "highio,highcpu"
	rt = CheckClusterLabels(als, cl, false)
	if rt != true {
		t.Fatal(cl)
	}

	cl = "default"
	rt = CheckClusterLabels(als, cl, false)
	if rt != false {
		t.Fatal(cl)
	}
}

func TestComposeClusterList(t *testing.T) {

	var ctrl Controler
	apps := InitApps()
	clusters := InitClusters()

	ctrl.ComposeClusterList(clusters, apps, "highio,highcpu", "K8S", "http://10.10.1.0:8080")

	i := 0
	for elem := ctrl.ClstrList.Front(); elem != nil; elem = elem.Next() {
		el := (elem.Value.(Cluster))
		if el.UUID == clusters[i].UUID {
			if el.Name != clusters[i].Name {
				t.Fatal(el.Name)
			}
			if el.AppList[apps[i].Name] != apps[i].UUID {
				t.Fatal(el.Name)
			}
		}
		i++
	}
	if i != len(clusters)-4 {
		t.Fatal(i)
	}

}

//DisplayClstrList is output the clusters in linklist
// @Title DisplayClstrList
// @Description  Display clusters in ClusterList
// @Accept  NULL
// @Param   NULL
func (ctrl *Controler) DisplayClstrList(str string) {
	fmt.Println(str)
	i := 0
	fmt.Println("The default Cluster: ", ctrl.Default)
	for clstr := ctrl.ClstrList.Front(); clstr != nil; clstr = clstr.Next() {
		clr := (clstr.Value.(Cluster))
		fmt.Printf("[%d]clstr %s, %s, %s, %s\n", i, clr.UUID, clr.Name, clr.Label, clr.EndPoint)
		fmt.Println("AppList:\n", clr.AppList)
		fmt.Println("Resources:\n", clr.ResourceList)
		i++
		fmt.Println("==========================================")
	}
}

func TestConvert2algApp(t *testing.T) {
	iapp_test := AppData{UUID: "app123456", Name: "myapp", Label: "",
		CPU: 2, Memory: 4, Disk: 100, ReservedRatio: 20, State: 0,
	}
	var oapp *Application
	oapp = Convert2algApp(&iapp_test)
	if !strings.EqualFold(oapp.UUID, iapp_test.UUID) {
		t.Fatal(iapp_test.UUID)
	}
	if !strings.EqualFold(oapp.Name, iapp_test.Name) {
		t.Fatal(iapp_test.Name)
	}
	if !strings.EqualFold(oapp.Label, iapp_test.Label) {
		t.Fatal(iapp_test.Label)
	}
	if oapp.ResourceList[CPU].Total != float32(iapp_test.CPU) {
		t.Fatal("CPU Error")
	}
	if oapp.ResourceList[MEM].Total != float32(iapp_test.Memory) {
		t.Fatal("Memory Error")
	}
	if oapp.ResourceList[CPU].ReservedRatio != float32(iapp_test.ReservedRatio) {
		t.Fatal("ReservedRatio Error")
	}
	if oapp.State != iapp_test.State {
		t.Fatal("State Error")
	}
}

func TestComposeAppList(t *testing.T) {
	var ctrl Controler
	apps := InitApps()
	ctrl.ComposeAppList(apps)
	i := 0
	for app := ctrl.AppList.Front(); app != nil; app = app.Next() {
		el := (app.Value.(*Application))
		if el.UUID == apps[i].UUID {
			if el.Name != apps[i].Name {
				t.Fatal(el.Name)
			}
		}
		i++
	}
	if i != len(apps) {
		t.Fatal(i)
	}
}

// DisplayAppList is output application in linklist
// @Title DisplayAppList
// @Description  Display App in AppList
// @Accept  NULL
// @Param   NULL
func (ctrl *Controler) DisplayAppList(str string) {
	fmt.Println(str)
	i := 0
	for app := ctrl.AppList.Front(); app != nil; app = app.Next() {
		ap := (app.Value.(Application))
		fmt.Println(ap)
		i++
	}
}

/*
func TestExtend(t *testing.T) {
	n := SetExtend(nil, "Enginetype", "k8s")
	if n == nil {
		t.Fatal("SetExtend return nil")
	}
	n = SetExtend(n, "xendpoint", "http://9.91.18.131:8080")

	//fmt.Println(n)

	et := GetExtendByKey(n, "Enginetype")
	if et != "k8s" {
		t.Fatal(et)
	}
	ex := GetExtendByKey(n, "Xendpoint")
	if ex != "http://9.91.18.131:8080" {
		t.Fatal(ex)
	}
}*/
