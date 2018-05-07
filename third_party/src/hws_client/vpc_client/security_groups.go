package vpc_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type SecurityGroup struct {
	// The UUID of the resource.
	ID string `json:"id,omitempty"`
	// Human-readable name of the resource.
	Name string `json:"name,omitempty"`
	// The ID of the tenant who owns the resource.
	TenantId string `json:"tenant_id,omitempty"`
	// The human-readable description for the resource.
	Description string `json:"description,omitempty"`
	// A list of SecurityGroupRule objects.
	SecurityGroupRules []SecurityGroupRule `json:"security_group_rules"`
}

type SecurityGroupDetails struct {
	SecurityGroup SecurityGroup `json:"security_group"`
}

type SecurityGroups struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

func (vpc *Vpcclient) CreatSecurityGroup(name string, token string) (*SecurityGroupDetails, error) {
	router := fmt.Sprintf("%s/v2.0/security-groups", vpc.endpoint)
	requstbody := SecurityGroupDetails{
		SecurityGroup: SecurityGroup{
			Name: name,
		},
	}
	str, _ := json.Marshal(requstbody)
	body := bytes.NewBuffer([]byte(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 201 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result SecurityGroupDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListImages lists the images according with ImageFilterOptions.
func (vpc *Vpcclient) ListSecurityGroups(token string) (*SecurityGroups, error) {
	router := fmt.Sprintf("%s/v2.0/security-groups", vpc.endpoint)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result SecurityGroups
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) GetSecurityGroup(sec_id string, token string) (*SecurityGroupDetails, error) {
	router := fmt.Sprintf("%s/v2.0/security-groups/%s", vpc.endpoint, sec_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result SecurityGroupDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) DeleteSecurityGroup(sec_id string, token string) error {
	router := fmt.Sprintf("%s/v2.0/security-groups/%s", vpc.endpoint, sec_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 204 {
		return errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
