package vpc_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type AddressPair struct {
	// IP address.
	IpAddress string `json:"ip_address"`
	// MAC address.
	MacAddress string `json:"mac_address"`
}

type OptionPair struct {
	// The option value.
	OptValue string `json:"opt_value"`
	// The option name.
	OptName string `json:"opt_name"`
}

type ExtraDhcpOpt struct {
	// opetion name
	OptName string `json:"opt_name"`
	// option value
	OptValue string `json:"opt_value"`
}

// If you specify only a subnet UUID, system allocates an available IP from that subnet to the port.
// If you specify both a subnet UUID and an IP address, system tries to allocate the address to the port.
type FixedIp struct {
	// The UUID of a subnet.
	SubnetId string `json:"subnet_id"`
	// The IP address.
	IpAddress string `json:"ip_address,omitempty"`
}

// FusionCompute only support ID,Name fields.
type Port struct {
	// The UUID of the port.
	ID string `json:"id"`
	// The port name.
	Name string `json:"name"`
	// The port status.
	Status string `json:"status"`
	// The UUID of the network.
	NetworkId string `json:"network_id"`
	// The ID of the tenant who owns the resource.
	TenantId string `json:"tenant_id"`
	// The MAC Address of the port.
	MacAddress string `json:"mac_address"`
	// A set of zero or more allowed address pairs.
	AllowedAddressPairs []AddressPair `json:"allowed_address_pairs"`
	// A set of zero or more extra DHCP option pairs.
	ExtraDhcpOpts []OptionPair `json:"extra_dhcp_opts"`
	// The administrative state of the resource, which is up (true) or down (false).
	AdminStateUp bool `json:"admin_state_up"`
	// The UUID of the entity that uses this port.
	DeviceOwner string `json:"device_owner"`
	// The port security status. A valid value is enabled (true) or disabled (false).
	PortSecurityEnabled bool `json:"port_security_enabled"`
	// The fixed ip addresses of this port.
	FixedIps []FixedIp `json:"fixed_ips"`
	// One or more security group UUIDs.
	SecurityGroups []string `json:"security_groups"`
	// The UUID of the device that uses this port.
	DeviceId string `json:"device_id"`

	// The date and time when the resource was created. The date and time stamp format is ISO 8601.
	ChangedAt string `json:"changed_at"`
	// The date and time when the resource was updated. The date and time stamp format is ISO 8601.
	UpdatedAt string `json:"updated_at"`

	// The UUID of the host which bind the port.
	HostId string `json:"binding:host_id"`
	// The user profile.
	Profile interface{} `json:"binding:profile"`
	// The port interface detail information.
	VifDetails struct {
		PortFilter    bool `json:"port_filter"`
		OvsHybridPlug bool `json:"ovs_hybrid_plug"`
	} `json:"binding:vif_details"`
	// The type of the port.
	VnicType string `json:"binding:vnic_type"`
	// The type of the port interface.
	VifType string `json:"binding:vif_type"`
}

type Ports struct {
	Ports []Port `json:"ports"`
}

type PortDetail struct {
	Port Port `json:"port"`
}

func (vpc *Vpcclient) CreatPort(project_id string, subnet_id string, securitygroups string, token string) (*PortDetail, error) {
	Securitygroups := []string{securitygroups}
	router := fmt.Sprintf("%s/v2.0/ports", vpc.endpoint, project_id)
	requstbody := &PortDetail{
		Port: Port{
			NetworkId:      subnet_id,
			TenantId:       project_id,
			SecurityGroups: Securitygroups,
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
	var result PortDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (vpc *Vpcclient) DeletePort(portID string, token string) error {
	router := fmt.Sprintf("%s/v2.0/ports/%s", vpc.endpoint, portID)
	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 204 && status_code != 404 {
		return errors.New("DeletePort: return code is not 204 or 404:" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}