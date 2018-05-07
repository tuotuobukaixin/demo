package controler

import (
	"container/list"
	alg "scheduler/algorithm"
)

//AppData is a view of user's Application table
// @Title AppData Struct
type AppData struct {
	UUID          string
	Name          string
	Status        string
	URL           string
	Label         string
	ClusterUUID   string
	Labels        string
	AppJSON       string
	Instance      int
	CPU           int
	Memory        int
	Disk          int
	ReservedRatio int
	State         int
}

//ClusterData is a view of runtime engine table
// @Title ClusterData Struct
type ClusterData struct {
	UUID       string
	Name       string
	Label      string
	Enginetype string
	Endpoint   string
	Status     string
	CPU        int
	Memory     int
	Disk       int
}

//RuleData is a view of rules table or user's app rule
// @title Rules from DB or Json
type RuleData struct {
	RuleType int
	Level    int
	Name     string
	Apps     string
	Clusters string
}

//Cluster's Status
const (
	AVAILABLE   = "AVAILABLE"
	POWEROFF    = "POWEROFF"
	TOBEOFF     = "TOBEOFF"
	UNAVAILABLE = "UNAVAILABLE"
)

//App's Status
const (
	UNINITIALIZED = "UNINITIALIZED"
	INITIALIZING  = "INITIALIZING"
	INITIALIZED   = "INITIALIZED"
	ERROR         = "ERROR"
	STARTING      = "STARTING"
	STARTED       = "STARTED"
	STOPPING      = "STOPPING"
	STOPPED       = "STOPPED"
	SCALING       = "SCALING"
	TIMEOUT       = "TIMEOUT"
	DELETING      = "DELETING"
)

//Controler is a handler of alg controler
// @Title Controler Struct
type Controler struct {
	Default   ClusterData
	Iapp      *alg.Application
	ClstrList *list.List
	RuleList  *list.List
	AppList   *list.List
}
