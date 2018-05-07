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
	OUTBOUND = "egress"
	INBOUND  = "ingress"
)

type SecurityGroupRule struct {
	// The UUID of the resource.
	ID string `json:"id,omitempty"`
	// The ID of the tenant who owns the resource.
	TenantId string `json:"tenant_id,omitempty"`
	// The remote group UUID to associate with this security group rule.
	SecurityGroupId string `json:"security_group_id,omitempty"`
	// The human-readable description for the resource.
	Description string `json:"description,omitempty"`
	// Ingress or egress, which is the direction in which the metering rule is applied.
	Direction string `json:"direction,omitempty"`
	// Must be IPv4 or IPv6, and addresses represented in CIDR must match the ingress or egress rules.
	EtherType string `json:"ethertype,omitempty"`
	// The IP protocol. Valid value is icmp, tcp, udp, or null. No default.
	Protocol string `json:"protocol,omitempty"`
	// The maximum port number in the range that is matched by the security group rule.
	PortRangeMax int `json:"port_range_max"`
	// The minimum port number in the range that is matched by the security group rule.
	PortRangeMin int `json:"port_range_min"`
	// The remote group ID to be associated with this security group rule.
	// You can specify either RemoteGroupID or RemoteIPPrefix.
	RemoteGroupId string `json:"remote_group_id,omitempty"`
	// The remote IP prefix to be associated with this security group rule. You
	// can specify either RemoteGroupID or RemoteIPPrefix . This attribute
	// matches the specified IP prefix as the source IP address of the IP packet.
	RemoteIpPrefix string `json:"remote_ip_prefix,omitempty"`
}

type SecurityGroupRuleDetails struct {
	SecurityGroupRule SecurityGroupRule `json:"security_group_rule"`
}

type SecurityGroupRules struct {
	SecurityGroupRules []SecurityGroupRule `json:"security_group_rules"`
}

func (vpc *Vpcclient) CreatSecurityGroupRule(sec_id string, direction string, protocol string, ip_prefix string, portmin int, portmax int, token string) (*SecurityGroupRuleDetails, error) {
	router := fmt.Sprintf("%s/v2.0/security-group-rules", vpc.endpoint)
	requstbody := SecurityGroupRuleDetails{
		SecurityGroupRule: SecurityGroupRule{
			SecurityGroupId: sec_id,
			Protocol:        protocol,
			Direction:       direction,
			RemoteIpPrefix:  ip_prefix,
			PortRangeMin:    portmin,
			PortRangeMax:    portmax,
		},
	}

	str, _ := json.Marshal(requstbody)
	body := bytes.NewBuffer([]byte(str))
	fmt.Println(string(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 201 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result SecurityGroupRuleDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListImages lists the images according with ImageFilterOptions.
func (vpc *Vpcclient) ListSecurityGroupRules(sec_id string, token string) (*SecurityGroupRules, error) {
	router := fmt.Sprintf("%s/v2.0/security-group-rules?security_group_id=%s", vpc.endpoint, sec_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result SecurityGroupRules
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) GetSecurityGroupRule(sec_rules_id string, token string) (*SecurityGroupRuleDetails, error) {
	router := fmt.Sprintf("%s/v2.0/security-group-rules/%s", vpc.endpoint, sec_rules_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result SecurityGroupRuleDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) DeleteSecurityGroupRule(sec_rules_id string, token string) error {
	router := fmt.Sprintf("%s/v2.0/security-group-rules/%s", vpc.endpoint, sec_rules_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 204 {
		return errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
