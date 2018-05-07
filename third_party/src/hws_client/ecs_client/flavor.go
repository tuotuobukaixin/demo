package ecs_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)
//os_extra_specs struct
type Os_extra_specs struct {
	Instance_vnic_type   string `json:"instance_vnic:type"`
	Ecs_performance_type string `json:"ecs:performancetype"`
	Resource_type        string `json:"resource_type"`
}
// FlavorDetail contains all the infomations about flavor.
type Flavor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// The amount of RAM a flavor has, in MiB.
	Ram int `json:"ram"`
	// The size of the root disk that will be created in GiB.
	// If 0 the root disk will be set to exactly the size of the image used to deploy the instance.
	Disk int `json:"disk"`
	// The number of virtual CPUs that will be allocated to the server.
	Vcpus int `json:"vcpus"`
	// The size of the ephemeral disk that will be created, in GiB.
	Ephemeral int `json:"OS-FLV-EXT-DATA:ephemeral"`
	// Whether or not the flavor has been administratively disabled.
	Disabled bool `json:"OS-FLV-DISABLED:disabled"`
	// The size of a dedicated swap disk that will be allocated, in GiB.
	// If 0 (the default), no dedicated swap disk will be created.
	Swap string `json:"swap"`
	// The receive / transimit factor that will be set on ports.
	RxtxFactor float64 `json:"rxtx_factor"`
	// Whether the flavor is public (available to all projects) or scoped to a set of projects.
	//Default is True if not specified.
	IsPublic bool `json:"os-flavor-access:is_public"`
	Os_Extra_Specs Os_extra_specs `json:"os_extra_specs"`
}

type FlavorDetails struct {
	Flavor Flavor `json:"flavor"`
}

type Flavors struct {
	Flavors []Flavor `json:"flavors"`
}

// ListImages lists the images according with ImageFilterOptions.
func (ecs *Ecsclient) ListFlavors(project_id string, token string) (*Flavors, error) {
	router := fmt.Sprintf("%s/v2/%s/flavors/detail", ecs.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Flavors
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) GetFlavor(project_id string, flavor_id string, token string) (*FlavorDetails, error) {
	router := fmt.Sprintf("%s/v2/%s/flavors/%s", ecs.endpoint, project_id, flavor_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result FlavorDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
