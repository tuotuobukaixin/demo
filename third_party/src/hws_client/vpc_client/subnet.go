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
	SUBNET_ACTIVE = "ACTIVE"
	SUBNET_OK     = "DOWN"
	SUBNET_DOWN   = "UNKNOWN"
	SUBNET_ERROR  = "ERROR"
)

type Subnet struct {
	// UUID for the subnet
	ID string `mapstructure:"id" json:"id"`
	// Human-readable name for the subnet. Might not be unique.
	Name string `mapstructure:"name" json:"name"`
	// The CIDR representing IP range for this subnet, based on IP version.
	CIDR string `mapstructure:"cidr" json:"cidr"`
	// Default gateway IP used by devices in this subnet.
	GatewayIp string `mapstructure:"gateway_ip" json:"gateway_ip"`
	// Indicates whether DHCP is enabled for this subnet or not.
	DHCPEnable       bool   `mapstructure:"dhcp_enable" json:"dhcp_enable"`
	PrimaryDNS       string `mapstructure:"primary_dns" json:"primary_dns"`
	SecondaryDNS     string `mapstructure:"secondary_dns" json:"secondary_dns"`
	AvailabilityZone string `mapstructure:"availability_zone" json:"availability_zone"`
	VpcID            string `mapstructure:"vpc_id" json:"vpc_id"`
	// Indicates whether subnet is currently operational. Possible values include
	// `ACTIVE', `DOWN', `UNKNOWN', or `ERROR'.
	Status string `mapstructure:"status" json:"status"`
}

type SubnetDetails struct {
	Subnet Subnet `mapstructure:"subnet" json:"subnet"`
}

type Subnets struct {
	Subnets []Subnet `mapstructure:"subnets" json:"subnets"`
}

// ListImages lists the images according with ImageFilterOptions.
func (vpc *Vpcclient) CreatSubnet(project_id string, name string, cidr string, gataway_ip string, vpc_id string, token string) (*SubnetDetails, error) {
	router := fmt.Sprintf("%s/v1/%s/subnets", vpc.endpoint, project_id)
	requstbody := &SubnetDetails{
		Subnet: Subnet{
			Name:       name,
			CIDR:       cidr,
			GatewayIp:  gataway_ip,
			VpcID:      vpc_id,
			DHCPEnable: true,
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
	var result SubnetDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListImages lists the images according with ImageFilterOptions.
func (vpc *Vpcclient) ListSubnets(project_id string, vpc_id string, token string) (*Subnets, error) {
	router := fmt.Sprintf("%s/v1/%s/subnets?vpc_id=%s", vpc.endpoint, project_id, vpc_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Subnets
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) GetSubnet(project_id string, subnet_id string, token string) (*SubnetDetails, error) {
	router := fmt.Sprintf("%s/v1/%s/subnets/%s", vpc.endpoint, project_id, subnet_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result SubnetDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) DeleteSubnet(project_id string, vpc_id string, subnet_id string, token string) error {
	router := fmt.Sprintf("%s/v1/%s/vpcs/%s/subnets/%s", vpc.endpoint, project_id, vpc_id, subnet_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 204 {
		return errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
