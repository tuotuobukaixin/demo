package vpc_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

const (
	VPC_CREATING       = "CREATING"
	VPC_OK             = "OK"
	VPC_DOWN           = "DOWN"
	VPC_PENDING_UPDATE = "PENDING_UPDATE"
	VPC_PENDING_DELETE = "PENDING_DELETE"
	VPC_ERROR          = "ERROR"
)

// Vpc represents, well, a vpc.
type Vpc struct {
	// UUID for the vpc
	ID string `mapstructure:"id" json:"id"`

	// Human-readable name for the vpc. Might not be unique.
	Name string `mapstructure:"name" json:"name"`

	CIDR string `mapstructure:"cidr" json:"cidr"`

	// Indicates whether vpc is currently operational. Possible values include
	// `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional values.
	Status string `mapstructure:"status" json:"status"`
}

type VpcDetails struct {
	Vpc Vpc `mapstructure:"vpc" json:"vpc"`
}

type Vpcs struct {
	Vpcs []Vpc `mapstructure:"vpcs" json:"vpcs"`
}

// ListImages lists the images according with ImageFilterOptions.
func (vpc *Vpcclient) CreatVpc(project_id string, name string, token string) (*VpcDetails, error) {
	router := fmt.Sprintf("%s/v1/%s/vpcs", vpc.endpoint, project_id)
	requstbody := VpcDetails{
		Vpc: Vpc{
			Name: name,
		},
	}
	str, _ := json.Marshal(requstbody)
	body := bytes.NewBuffer([]byte(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result VpcDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListImages lists the images according with ImageFilterOptions.
func (vpc *Vpcclient) ListVpcs(project_id string, token string) (*Vpcs, error) {
	router := fmt.Sprintf("%s/v1/%s/vpcs", vpc.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Vpcs
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) GetVpc(project_id string, vpc_id string, token string) (*VpcDetails, error) {
	router := fmt.Sprintf("%s/v1/%s/vpcs/%s", vpc.endpoint, project_id, vpc_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result VpcDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) DeleteVpc(project_id string, vpc_id string, token string) error {
	router := fmt.Sprintf("%s/v1/%s/vpcs/%s", vpc.endpoint, project_id, vpc_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 204 {
		return errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
