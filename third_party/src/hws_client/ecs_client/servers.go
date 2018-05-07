package ecs_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type Nic struct {
	SubnetId  string `json:"subnet_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
}

type SecurityGroup struct {
	ID string `json:"id,omitempty"`
}

// Personality is an array of files that are injected into the server at launch.
type Personality []*File

// File is used within CreateOpts and RebuildOpts to inject a file into the server at launch.
// File implements the json.Marshaler interface, so when a Create or Rebuild operation is requested,
// json.Marshal will call File's MarshalJSON method.
type File struct {
	// Path of the file
	Path string `json:"path,omitempty"`
	// Contents of the file. Maximum content size is 255 bytes.
	Contents string `json:"contents,omitempty"`
}
type Bandwidth struct {
	Size       int    `json:"size,omitempty"`
	Sharetype  string `json:"sharetype,omitempty"`
	Chargemode string `json:"chargemode,omitempty"`
}
type Eip struct {
	Iptype    string    `json:"iptype,omitempty"`
	Bandwidth Bandwidth `json:"bandwidth,omitempty"`
}
type PublicIp struct {
	ID  string `json:"id,omitempty"`
	EIP Eip    `json:"eip,omitempty"`
}
type Volume struct {
	VolumeType string `json:"volumetype,omitempty"`
	Size       int    `json:"size,omitempty"`
}

type CreateServerOpts struct {
	// [required]
	Name string `json:"name,omitempty"`

	// [required]
	ImageRef string `json:"imageRef,omitempty"`

	// [required]
	FlavorRef string `json:"flavorRef,omitempty"`

	// [optional]
	Personality Personality `json:"personality,omitempty"`

	// [optional]
	UserData string `json:"user_data,omitempty"`

	// [optional]
	AdminPass string `json:"key_name,omitempty"`

	// [optional]
	KeyName string `json:"key_name,omitempty"`

	// [required]
	VpcId string `json:"vpcid,omitempty"`

	// [required]
	Nics []Nic `json:"nics,omitempty"`

	// [optional]
	PublicIp *PublicIp `json:"publicip,omitempty"`

	// [optional]
	Count int `json:"count,omitempty"`

	// [required]
	RootVolume *Volume `json:"root_volume,omitempty"`

	// [optional]
	DataVolumes []Volume `json:"data_volumes,omitempty"`

	// [optional]
	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`

	// [optional]
	AvailabilityZone string `json:"availability_zone,omitempty"`

	// [optional]
	ExtendParam map[string]interface{} `json:"extendparam,omitempty"`

	// [optional]
	Metadata map[string]string `json:"metadata,omitempty"`
}
type ServerID struct {
	ID string `json:"id,omitempty"`
}

type ResourceLink struct {
	// Link includes HTTP references to the itself, useful for passing along to other APIs that might want a resource reference.
	Href string `json:"href"`

	// Types of link relations associated with resources, like "self","bookmark","alternate"
	// A self link contains a versioned link to the resource. Use these links when the link is followed immediately.
	// A bookmark link provides a permanent link to a resource that is appropriate for long term storage.
	// An alternate link can contain an alternate representation of the resource.
	Rel string `json:"rel"`

	// The type attribute provides a hint as to the type of representation to expect when following the link.
	Type string `json:"type,omitempty"`
}
type Image struct {
	// The UUID of the image.
	Id string `json:"Id"`
	// The links of the image.
	Links []ResourceLink `json:"links"`
}
type SecurityGroupInMachine struct {
	// Security group name.
	Name string `json:"name"`
}

type ServerFlavor struct {
	// The UUID of the flavor.
	Id string `json:"Id"`
	// The links of the flavor.
	Links []ResourceLink `json:"links"`
}
type Server struct {
	// The UUID of the resource.
	ID string `json:"id"`

	// The name of the resource.
	Name string `json:"name"`
// The addresses for the server.
	Addresses map[string]interface{} `json:"addresses"`
	// The ID and links for the flavor for your server instance.
	Flavor ServerFlavor `json:"flavor"`
	// The name of associated key pair, if any.
	KeyName string `json:"key_name"`
	// A dictionary of metadata key-and-value pairs, which is maintained for backward compatibility.
	MetaData map[string]string `json:"metadata"`
	// The availability zone name.
	AvailabilityZone string `json:"OS-EXT-AZ:availability_zone"`
	// The security group name.
	SecurityGroups []SecurityGroupInMachine `json:"security_groups"`
	// The server status.
	Status string `json:"status"`
	// A percentage value of the build progress.
	Progress int `json:"progress"`
	// The ID of the host.
	HostId string `json:"hostId"`
	// IPv4 address that should be used to access this server.
	AccessIPv4 string `json:"accessIPv4"`
	// IPv6 address that should be used to access this server.
	AccessIPv6 string `json:"accessIPv6"`
	// Indicates whether a configuration drive enables metadata injection.
	ConfigDrive string `json:"config_drive"`
	// Controls how the API partitions the disk when you create, rebuild, or resize servers.
	DiskConfig string `json:"OS-DCF:diskConfig"`
	// The host name. Appears in the response for administrative users only.
	Host string `json:"OS-EXT-SRV-ATTR:host"`
	// The hypervisor host name. Appears in the response for administrative users only.
	HypervisorHostname string `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	// The instance name. The Compute API generates the instance name from the instance name template. Appears in the response for administrative users only.
	InstanceName string `json:"OS-EXT-SRV-ATTR:instance_name"`
	// The power state of the instance.
	PowerState int `json:"OS-EXT-STS:power_state"`
	// The task state of the instance.
	TaskState string `json:"OS-EXT-STS:task_state"`
	// The VM state.
	VmState string `json:"OS-EXT-STS:vm_state"`
	// The service state.
	ServiceState string `json:"OS-EXT-SERVICE:service_state"`
	// The date and time when the server was launched. The date and time stamp format is ISO 8601.
	LaunchedAt string `json:"OS-SRV-USG:launched_at"`
	// The date and time when the server was deleted. The date and time stamp format is ISO 8601.
	TerminatedAt string `json:"OS-SRV-USG:terminated_at"`
	// The attached volumes, if any.
	VolumesAttached []struct {
		ID string `json:"id"`
	} `json:"os-extended-volumes:volumes_attached"`
	// The hostname set on the instance when it is booted.
	HostName string `json:"OS-EXT-SERV-ATTR:hostname"`
	// The date and time when the resource was created. The date and time stamp format is ISO 8601.
	Created string `json:"created"`
	// The date and time when the resource was updated. The date and time stamp format is ISO 8601.
	Updated string `json:"updated"`
	// The UUID of the tenant in a multi-tenancy cloud.
	TenantId string `json:"tenant_id"`
	// The user ID of the user who owns the server.
	UserId string `json:"user_id"`
// The UUID and links for the image for your server instance.
	Image Image `json:"image"`
	// The host status.
	HostStatus string `json:"host_status"`
	// The reservation id for the server. This is an id that can be useful in tracking groups of servers created with multiple create, that will all have the same reservation_id.
	ReservationId string `json:"OS-EXT-SERV-ATTR:reservation_id"`
	// When servers are launched via multiple create, this is the sequence in which the servers were launched.
	LaunchIndex int `json:"OS-EXT-SERV-ATTR:launch_index"`
	// The UUID of the kernel image when using an AMI. Will be null if not.
	KernelId string `json:"OS-EXT-SERV-ATTR:kernel_id"`
	// The UUID of the ramdisk image when using an AMI. Will be null if not.
	RamdiskId string `json:"OS-EXT-SERV-ATTR:ramdisk_id"`
	// The root device name for the instance
	RootDeviceName string `json:"OS-EXT-SERV-ATTR:root_device_name"`
	// The user_data the instance was created with.
	UserData string `json:"OS-EXT-SERV-ATTR:user_data"`
	// Location means the location of the machine create
	Location string `json:"location"`
	// DataStoreUrns means the related disks urns
	DataStoreUrns []string `json:"dataStoreUrns"`
}
type DeleteServerOpts struct {
	// [required]
	Servers []ServerID `json:"servers,omitempty"`

	// [required]
	DeletePublicIP bool `json:"delete_publicip,omitempty"`

	// [required]
	DeleteVolume bool `json:"delete_volume,omitempty"`
}
type ServerDetails struct {
	Server Server `json:"server"`
}
type Servers struct {
	Servers []Server `json:"servers"`
}
func (ecs *Ecsclient) ToServerCreateMap(opts *CreateServerOpts) (map[string]interface{}, error) {

	if opts.Name == "" {
		return nil, errors.New("Missing field required for cloudserver creation: name")
	}
	if opts.ImageRef == "" {
		return nil, errors.New("Missing field required for cloudserver creation: imageRef")
	}
	if opts.FlavorRef == "" {
		return nil, errors.New("Missing field required for cloudserver creation: flavorRef")
	}
	if opts.VpcId == "" {
		return nil, errors.New("Missing field required for cloudserver creation: vpcid")
	}
	if len(opts.Nics) == 0 {
		return nil, errors.New("Missing field required for cloudserver creation: nics")
	}
	if opts.RootVolume == nil {
		return nil, errors.New("Missing field required for cloudserver creation: root_volume")
	}
	if opts.AvailabilityZone == "" {
		return nil, errors.New("Missing field required for cloudserver creation: availability_zone")
	}

	server := make(map[string]interface{})
	server["name"] = opts.Name
	server["imageRef"] = opts.ImageRef
	server["flavorRef"] = opts.FlavorRef
	server["vpcid"] = opts.VpcId

	var nics []interface{}
	for _, nic := range opts.Nics {
		nics = append(nics, map[string]interface{}{
			"subnet_id":  nic.SubnetId,
			"ip_address": nic.IpAddress,
		})
	}
	server["nics"] = nics

	rootVolume := make(map[string]interface{})
	rootVolume["volumetype"] = opts.RootVolume.VolumeType
	rootVolume["size"] = opts.RootVolume.Size

	server["root_volume"] = rootVolume

	server["availability_zone"] = opts.AvailabilityZone

	if len(opts.Personality) != 0 {
		server["personality"] = opts.Personality
	}

	if opts.AdminPass != "" {
		server["adminPass"] = opts.AdminPass
	}

	if opts.KeyName != "" {
		server["key_name"] = opts.KeyName
	}

	if opts.PublicIp != nil {
		server["publicip"] = map[string]interface{}{
			"id":  opts.PublicIp.ID,
			"eip": opts.PublicIp.EIP,
		}
	}

	if opts.Count != 0 {
		server["count"] = opts.Count
	}

	if len(opts.DataVolumes) != 0 {
		dataVolumes := make([]interface{}, 0)
		for _, dv := range opts.DataVolumes {
			dataVolume := make(map[string]interface{})
			dataVolume["volumetype"] = dv.VolumeType
			dataVolume["size"] = dv.Size
			dataVolumes = append(dataVolumes, dataVolume)
		}
		server["data_volumes"] = dataVolumes
	}

	if len(opts.SecurityGroups) != 0 {
		securitygroups := make([]interface{}, 0)
		for _, sg := range opts.SecurityGroups {
			securitygroup := make(map[string]interface{})
			securitygroup["id"] = sg.ID
			securitygroups = append(securitygroups, securitygroup)
		}
		server["security_groups"] = securitygroups
	}

	if len(opts.Metadata) != 0 {
		server["metadata"] = opts.Metadata
	}

	if len(opts.ExtendParam) != 0 {
		server["extendparam"] = opts.ExtendParam
	}
	serverMap := make(map[string]interface{})
	serverMap["server"] = server
	return serverMap, nil
}

func (ecs *Ecsclient) CreatServers(project_id string, opts *CreateServerOpts, token string) (*JobDetail, error) {
	router := fmt.Sprintf("%s/v1/%s/cloudservers", ecs.endpoint, project_id)
	requstbody, err := ecs.ToServerCreateMap(opts)
	if err != nil {
		return nil, err
	}
	str, _ := json.Marshal(requstbody)
	body := bytes.NewBuffer([]byte(str))
	fmt.Println(string(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 && status_code != 201 && status_code != 202 && status_code != 204 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result JobDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) GetServer(project_id string, server_id string, token string) (*ServerDetails, error) {
	router := fmt.Sprintf("%s/v2/%s/servers/%s", ecs.endpoint, project_id, server_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 && status_code != 201 && status_code != 202 && status_code != 204 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result ServerDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) GetServers(project_id string,  token string) (*Servers, error) {
	router := fmt.Sprintf("%s/v2/%s/servers/detail", ecs.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 && status_code != 201 && status_code != 202 && status_code != 204 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Servers
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) DeleteServers(project_id string, opts *DeleteServerOpts, token string) (*JobDetail, error) {
	router := fmt.Sprintf("%s/v1/%s/cloudservers/delete", ecs.endpoint, project_id)

	str, _ := json.Marshal(opts)
	body := bytes.NewBuffer([]byte(str))
	fmt.Println(string(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 && status_code != 201 && status_code != 202 && status_code != 204 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result JobDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
