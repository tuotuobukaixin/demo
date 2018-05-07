package ecs_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

//VpcQuotas
type EcsQuota struct {
	QuotaSet struct {
		MaxImageMeta            int `json:"maxImageMeta"`
		MaxPersonality          int `json:"maxPersonality"`
		MaxPersonalitySize      int `json:"maxPersonalitySize"`
		MaxSecurityGroupRules   int `json:"maxSecurityGroupRules"`
		MaxSecurityGroups       int `json:"maxSecurityGroups"`
		MaxServerGroupMembers   int `json:"maxServerGroupMembers"`
		MaxServerGroups         int `json:"maxServerGroups"`
		MaxServerMeta           int `json:"maxServerMeta"`
		MaxTotalCores           int `json:"maxTotalCores"`
		MaxTotalFloatingIps     int `json:"maxTotalFloatingIps"`
		MaxTotalInstances       int `json:"maxTotalInstances"`
		MaxTotalKeypairs        int `json:"maxTotalKeypairs"`
		MaxTotalRAMSize         int `json:"maxTotalRAMSize"`
		TotalCoresUsed          int `json:"totalCoresUsed"`
		TotalFloatingIpsUsed    int `json:"totalFloatingIpsUsed"`
		TotalInstancesUsed      int `json:"totalInstancesUsed"`
		TotalRAMUsed            int `json:"totalRAMUsed"`
		TotalSecurityGroupsUsed int `json:"totalSecurityGroupsUsed"`
		TotalServerGroupsUsed   int `json:"totalServerGroupsUsed"`
	} `json:"absolute"`
}

func (ecs *Ecsclient) GetECSQuotas(project_id string, token string) (*EcsQuota, error) {
	var router string

	router = fmt.Sprintf("%s/v1/%s/cloudservers/limits", ecs.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result EcsQuota
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
