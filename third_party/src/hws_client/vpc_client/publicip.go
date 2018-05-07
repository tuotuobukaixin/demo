package vpc_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type PublicIPList struct {
	Publicips []PublicIP `json:"publicips"`
}

type PublicIP struct {
	ID                 string `json:"id"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	PortID             string `json:"port_id,omitempty"`
	PublicIPAddress    string `json:"public_ip_address"`
	PrivateIPAddress   string `json:"private_ip_address,omitempty"`
	TenantID           string `json:"tenant_id"`
	CreateTime         string `json:"create_time"`
	BandwidthID        string `json:"bandwidth_id"`
	BandwidthShareType string `json:"bandwidth_share_type"`
	BandwidthSize      int    `json:"bandwidth_size"`
}

type PublicIPDetail struct {
	PublicIP PublicIP `json:"publicip"`
}

func (vpc *Vpcclient) GetPublicIPS(project_id string, token string) (*PublicIPList, error) {

	router := fmt.Sprintf("%s/v1/%s/publicips", vpc.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result PublicIPList
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) GetPublicIP(project_id string, publicip_id string, token string) (*PublicIPDetail, int, error) {

	router := fmt.Sprintf("%s/v1/%s/publicips/%s", vpc.endpoint, project_id, publicip_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, status_code, err
	}
	if status_code != 200 {
		return nil, status_code, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result PublicIPDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, status_code, err
	}

	return &result, status_code, nil
}
