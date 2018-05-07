/*Package scheduler/controler
It's view layer for PaaS Deploy-mgr's scheduler
This layer can Convert controler data to algorithm data,
and Compose the linklist that the algorithm is needed.
User use it like this:
a: Get User's Data from somewhere, like DB,File,...
b: Convert it to controler data, and form the arrays
c: Call the Function PrepareData(...) to compose algorithm data,
and get a Controler handler c
d: Call the function c.GetCluster() and receive the recommendation:
cluster's endpoint and uuid
*/

package controler

import (
	"container/list"
	"fmt"
	alg "scheduler/algorithm"
	"strings"
)

//ParseArray is Parse JSON's array interface{}
func ParseArray(id string, rs interface{}) (value string) {
	trs := rs.([]interface{})
	for i, r := range trs {
		value = value + r.(string)
		if i < len(trs)-1 {
			value = value + ","
		}
	}
	if id != "" {
		value = value + "," + id
	}
	return
}

func formRuleData(key string, value string) RuleData {
	var rule RuleData

	rule.Level = 0
	rule.Name = key
	switch key {
	case "AppAffinity":
		rule.RuleType = alg.AAffinity
		rule.Apps = value
		rule.Clusters = ""
	case "AppUnaffinity":
		rule.RuleType = alg.ANonAffinity
		rule.Apps = value
		rule.Clusters = ""
	case "ClusterAffinity":
		rule.RuleType = alg.CAffinity
		rule.Apps = ""
		rule.Clusters = value
	case "ClusterUnaffinity":
		rule.RuleType = alg.CNonAffinity
		rule.Apps = ""
		rule.Clusters = value
	}
	return rule
}

//ParseAppRules parse App Rules from Path Json
func ParseAppRules(id string, data interface{}) []RuleData {

	rules := []RuleData{}
	if data == nil {
		return rules
	}

	path := data.(map[string]interface{})

	if path["AppAffinity"] != nil {
		value := ParseArray(id, path["AppAffinity"])
		rule := formRuleData("AppAffinity", value)
		rules = append(rules, rule)
	}
	if path["AppUnaffinity"] != nil {
		value := ParseArray(id, path["AppUnaffinity"])
		rule := formRuleData("AppUnaffinity", value)
		rules = append(rules, rule)
	}
	if path["ClusterAffinity"] != nil {
		value := ParseArray("", path["ClusterAffinity"])
		rule := formRuleData("ClusterAffinity", value)
		rules = append(rules, rule)
	}
	if path["ClusterUnaffinity"] != nil {
		value := ParseArray("", path["ClusterUnaffinity"])
		rule := formRuleData("ClusterUnaffinity", value)
		rules = append(rules, rule)
	}

	return rules
}

//SetExtend Set extend strings to map[string]string
func SetExtend(ed map[string]string, key string, value string) map[string]string {
	if ed == nil {
		ed = make(map[string]string)
	}
	ed[strings.ToLower(key)] = value
	return ed
}

//GetExtendByKey Get Value from extend string by Key
func GetExtendByKey(ed map[string]string, key string) (value string) {
	if ed == nil {
		value = ""
		return
	}

	value = ed[strings.ToLower(key)]
	return
}

func (c *Controler) initControler() {
	c.ClstrList = list.New()
	c.RuleList = list.New()
	c.AppList = list.New()
}

//PrepareData Prepare the Date algorithm needed.
// @Title Prepare Datas
// @Description  Prepare Datas from array to Linklist
// @Accept  Array
// @Param   rules, clusters, apps array return Controler
func PrepareData(ia *AppData, paths interface{}, clusters []ClusterData, apps []AppData, extend map[string]string) Controler {
	var ctrl Controler

	ctrl.initControler()
	ctrl.Iapp = Convert2algApp(ia)

	if paths != nil {
		rules := ParseAppRules(ia.Name, paths)
		ctrl.ComposeRuleList(rules)
	}

	et := GetExtendByKey(extend, "enginetype")
	ex := GetExtendByKey(extend, "xendpoint")

	ctrl.ComposeClusterList(clusters, apps, ia.Label, et, ex)
	ctrl.ComposeAppList(apps)
	//ctrl.DisplayClstrList("Prepare_ClusterList")

	return ctrl
}

//GetCluster User use it to get Recommendation: cluster's endpoint & uuid
// @Title GetCluster
// @Description
// @Accept Controler's Iapp, ClusterList, RuleList
// @Param  Cluster's endpoint & UUID
func (c *Controler) GetCluster() (string, string) {
	args := alg.PlacementParam{
		App:      c.Iapp,
		Clusters: c.ClstrList,
		Rules:    c.RuleList,
		Policy:   alg.BLANCE,
	}

	//c.DisplayClstrList("Testing_ClusterList")
	//c.DisplayAppList("Test_AppList")
	//c.DisplayRuleList("Testing_RuleList")

	cluster := alg.Scheduling(args)
	if cluster != nil {
		return cluster.(alg.Cluster).EndPoint, cluster.((alg.Cluster)).UUID
	}
	fmt.Println("Recving nil from Scheduling(...), Use Default Cluster ...")
	return c.Default.Endpoint, c.Default.UUID
}
