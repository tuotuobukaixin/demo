package controler

import (
	"container/list"
	alg "scheduler/algorithm"
	"strings"
)

//Convert2algRule is convert controler data to algorithm data
// @Title convert RuleData to Alg.Rule
// @Accept  RuleData
// @Param   Rule
func Convert2algRule(ir RuleData) (or alg.Rule) {
	or.RuleType = ir.RuleType
	or.Level = ir.Level
	or.Name = ir.Name
	or.Apps = ir.Apps
	or.Clusters = ir.Clusters
	return
}

//ComposeRuleList is Compose rules array to alg rule linklist
// @Title Compose Rules Linklist
// @Description  Create a Linklist from a Rule Array
// @Accept  Array
// @Param   rules array
func (ctrl *Controler) ComposeRuleList(rs []RuleData) {
	if ctrl.RuleList == nil {
		ctrl.RuleList = list.New()
	}
	for _, r := range rs {
		rule := Convert2algRule(r)
		ctrl.RuleList.PushBack(&rule)
	}
	return
}

//CheckClusterLabels is check cluster's labels and app's need
//If flag is true, then use the strict match, or the other
func CheckClusterLabels(al string, cl string, flag bool) bool {
	//App's label is null, all clusters can choose
	//Cluster's label is null, can provide all ability
	if al == "" {
		return true
	}
	if cl == "" {
		//The strict match
		if flag {
			return false
		}
		return true
	}
	if flag {
		//If use strict match, that a app's label not in cluster's labels, return false
		for _, a := range strings.Split(al, ",") {
			if strings.Contains(strings.ToLower(cl), strings.ToLower(a)) == false {
				return false
			}
		}
		return true
	}
	//If cluster have a label which app have, then it's ok
	for _, c := range strings.Split(cl, ",") {
		if strings.Contains(strings.ToLower(al), strings.ToLower(c)) {
			return true
		}
	}
	return false
}

//ComposeClusterList is compose clusters array to clusters linklist
// @Title ComposeClusterList
// @Description Compose ClusterList from ClusterData and AppData
// @Accept  ClusterData & AppData
// @Param   clusters & aas ctrl.ClstrList
func (ctrl *Controler) ComposeClusterList(clusters []ClusterData, apps []AppData, label string, et string, ex string) {
	if ctrl.ClstrList == nil {
		ctrl.ClstrList = list.New()
	}
	for _, cluster := range clusters {
		var cpu, mem int
		var clstr alg.Cluster

		if strings.ToLower(cluster.Status) != strings.ToLower(AVAILABLE) {
			continue
		}
		if et != "" {
			//if the cluster is not Enginetype, then filter it
			if strings.ToLower(et) != strings.ToLower(cluster.Enginetype) {
				continue
			}
		}

		if ex != "" {
			//if the endpoint is exclude, then filter it
			if strings.Contains(strings.ToLower(ex), strings.ToLower(cluster.Endpoint)) {
				continue
			}
		}

		//Get Default Clusters
		if strings.Contains(strings.ToLower(cluster.Label), "default") {
			ctrl.Default = cluster
		}

		if label != "" {
			//if the cluster is not The app's Label, then filter it
			rt := CheckClusterLabels(label, cluster.Label, false)
			if !rt {
				continue
			}
		}

		clstr.AppList = make(map[string]string)
		clstr.ResourceList = make(map[int]alg.Resource)

		clstr.UUID = cluster.UUID
		clstr.Name = cluster.Name
		clstr.Label = cluster.Label
		clstr.EndPoint = cluster.Endpoint
		for _, app := range apps {
			if app.ClusterUUID == clstr.UUID {
				clstr.AppList[app.Name] = app.UUID
				cpu += app.CPU
				mem += app.Memory
			}
		}
		clstr.ResourceList[alg.CPU] = alg.Resource{alg.CPU, float32(cpu), float32(cluster.CPU), (float32(cpu) / float32(cluster.CPU)), 90}
		clstr.ResourceList[alg.MEM] = alg.Resource{alg.MEM, float32(mem), float32(cluster.Memory), (float32(mem) / float32(cluster.Memory)), 90}

		ctrl.ClstrList.PushBack(clstr)
	}
	return
}

//Convert2algApp is  convert controler data to algorithm data
// @Title convert AppData to Alg.Application
// @Accept  AppData
// @Param   Application
func Convert2algApp(iapp *AppData) *alg.Application {
	var oapp alg.Application
	oapp.ResourceList = make(map[int]alg.Resource)

	oapp.UUID = iapp.UUID
	oapp.Name = iapp.Name
	oapp.Label = iapp.Label
	oapp.State = iapp.State
	oapp.ResourceList[alg.CPU] = alg.Resource{alg.CPU, 0, float32(iapp.CPU), 0, float32(iapp.ReservedRatio)}
	oapp.ResourceList[alg.MEM] = alg.Resource{alg.MEM, 0, float32(iapp.Memory), 0, float32(iapp.ReservedRatio)}

	return &oapp
}

//ComposeAppList is compose applications array to linklist
// @Title Compose Apps Linklist
// @Description  Create a Linklist from a Application Array
// @Accept  Array
// @Param   application array
func (ctrl *Controler) ComposeAppList(apps []AppData) {
	if ctrl.AppList == nil {
		ctrl.AppList = list.New()
	}
	for _, app := range apps {
		napp := Convert2algApp(&app)
		ctrl.AppList.PushBack(napp)
	}
	return
}
