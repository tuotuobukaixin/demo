package ecs_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type InterfaceAttachmentDetail struct {
	InterfaceAttachment InterfaceAttachment `json:"interfaceAttachment"`
}
type InterfaceAttachments struct {
	InterfaceAttachments []InterfaceAttachment `json:"interfaceAttachments"`
}

type FixIP struct {
	Subnet_id string `json:"subnet_id,omitempty"`
	IpAddr    string `json:"ip_address,omitempty"`
}
type InterfaceAttachment struct {
	Port_state string  `json:"port_state,omitempty"`
	FixIP      []FixIP `json:"fixed_ips,omitempty"`
	NetID      string  `json:"net_id,omitempty"`
	PortID     string  `json:"port_id,omitempty"`
	MacAddr    string  `json:"mac_addr,omitempty"`
}

func (ecs *Ecsclient) GetInterfaceAttachments(project_id string, server_id string, token string) (*InterfaceAttachments, error) {
	router := fmt.Sprintf("%s/v2/%s/servers/%s/os-interface", ecs.endpoint, project_id, server_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result InterfaceAttachments
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
