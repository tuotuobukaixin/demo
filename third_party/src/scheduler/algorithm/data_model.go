package algorithm

//Rules's types
const (
	AAffinity    = 0 //Affinity Among Apps
	ANonAffinity = 1 //Non-affinity Among Apps
	CAffinity    = 2 //Affinity between App and Clusters
	CNonAffinity = 3 //Non-affinity between App and Clusters
)

//Placement algorithm use. Now we finished the BLANCE
const (
	BLANCE     = 0 //LoadBalance Policy
	FASTEN     = 1 //Aggregate Policy
	DISPERSION = 2 //Dispersal Policy
)

//The Threshold of Load Balance Policy
const (
	CONSERVATIVE     = 0.4  // Small Threshold 0
	LESSCONSERVATIVE = 0.3  // Low Threshold 1
	MODERATE         = 0.2  // Middle Threshold 2
	LESSAGGRESSIVE   = 0.1  // High Threshold 3
	AGGRESSIVE       = 0.05 // Large Threshold 4
)

//Resoure Type
const (
	CPU  = 0 //cpu
	MEM  = 1 //内存
	DISK = 2 //磁盘
)

//Resource is The Veidoo of Resource
type Resource struct {
	ResourceType  int     //Resource Type: CPU，MEM，DISK
	Used          float32 //Cluster note allocation: CPU->Core/mCore MEM->MB DISK->GB. App is not Used
	Total         float32 //All Resource: cpu->Core/mCore, MEM->MB, DISK->GB
	Utilization   float32 //The Usage of Resource. The dynamic algorithm used, and placement algorithm not use.
	ReservedRatio float32 //Threshold level, will be used next
}

//Cluster  Information
type Cluster struct {
	UUID         string            //Cluster's UUID
	Name         string            //Cluster's Name
	Label        string            //The Label of Cluster, such as highio/highcpu/default
	EndPoint     string            //The endpoint of Cluster, such as http://ip:port
	State        int               //The status of Cluster
	AppList      map[string]string //The Apps deployed in the Cluster, map[Appname]Appuuid
	ResourceList map[int]Resource  //The Resource of Cluster, map[ResourceType]Resource
}

//Application  Information
type Application struct {
	UUID         string           //The App's UUID, will be set by Adapter
	Name         string           //The App's Name, which User's set
	Label        string           //The App's label, such as highio/highcpu
	State        int              //The status of App
	ResourceList map[int]Resource //The Resource of App, map[ResourceType]Resource
}

//Rule Infomation
type Rule struct {
	RuleType int    //Type of Rule
	Level    int    //Reserve
	Name     string //Name of Rule
	Apps     string //Among Apps, such as: "app1,app2,app3"
	Clusters string //Between App and Cluster, such as: "cluster1,cluster2,cluster2"
}

//PerfSnapshot is a time collection point
type PerfSnapshot struct {
	TimeStamp       int64                  //The timestamp of snapshot
	AppPerfData     map[string]Application //The history of App's
	ClusterPerfData map[string]Cluster     //The history of Cluster's, time should be same as app's
}
