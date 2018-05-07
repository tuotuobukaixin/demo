package evs_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

// Attachment denotes instance onto which the volume is attached.
type Attachment struct {
	ID           string `json:"id,omitempty"`
	ServerId     string `json:"server_id,omitempty"`
	AttachmentId string `json:"attachment_id,omitempty"`
	HostName     string `json:"host_name,omitempty"`
	VolumeId     string `json:"volume_id,omitempty"`
	Device       string `json:"device,omitempty"`
}

// VolumeDetail contains all the information about volume.
type Volume struct {
	// The UUID of the resource.
	ID string `json:"id,omitempty"`
	// The name of the resource.
	Name string `json:"name,omitempty"`
	// Instances onto which the volume is attached.
	Attachments []Attachment `json:"attachments,omitempty"`
	// AvailabilityZone is which availability zone the volume is in.
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// Indicates whether this is a bootable volume.
	Bootable string `json:"bootable,omitempty"`
	// ConsistencyGroupID is the consistency group ID.
	ConsistencyGroupID string `json:"consistencygroup_id,omitempty"`
	// The date and time when the resource was created. The date and time stamp format is ISO 8601
	CreatedAt string `json:"created_at,omitempty"`
	// Human-readable description for the volume.
	Description string `json:"description,omitempty"`
	// Encrypted denotes if the volume is encrypted.
	Encrypted bool `json:"encrypted,omitempty"`
	// The type of volume to create, either SATA or SSD.
	VolumeType string `json:"volume_type,omitempty"`
	// ReplicationDriverData contains data about the replication driver.
	ReplicationDriverData string `json:"os-volume-replication:driver_data,omitempty"`
	// ReplicationExtendedStatus contains extended status about replication.
	ReplicationExtendedStatus string `json:"os-volume-replication:extended_status,omitempty"`
	// ReplicationStatus is the status of replication.
	ReplicationStatus string `json:"replication_status,omitempty"`
	// The ID of the snapshot from which the volume was created.
	SnapshotID string `json:"snapshot_id,omitempty"`
	// The ID of another block storage volume from which the current volume was created.
	SourceVolID string `json:"source_volid,omitempty"`
	// Current status of the volume.
	Status string `json:"status,omitempty"`
	// TenantID is the id of the project that owns the volume.
	TenantID string `json:"os-vol-tenant-attr:tenant_id,omitempty"`
	// Arbitrary key-value pairs defined by the user.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Multiattach denotes if the volume is multi-attach capable.
	Multiattach bool `json:"multiattach,omitempty"`
	// Shareable denotes if the volume is shareble.
	// Using type interface, in order to support unmarshal both string and bool.
	Shareable interface{} `json:"shareable,omitempty"`
	// Size of the volume in GB.
	Size int `json:"size,omitempty"`
	// UserID is the id of the user who created the volume.
	UserID string `json:"user_id,omitempty"`
	// The volume migration status.
	MigrationStatus string `json:"migration_status,omitempty"`
	// The date and time when the resource was created. The date and time stamp format is ISO 8601
	UpdatedAt string `json:"updated_at,omitempty"`
	// Current back-end of the volume.
	Host string `json:"os-vol-host-attr:host,omitempty"`
	// The status of this volume migration (None means that a migration is not currently in progress).
	MigStat string `json:"os-vol-mig-status-attr:migstat,omitempty"`
	// The volume ID that this volume name on the back-end is based on.
	MigStatNameID string `json:"os-vol-mig-status-attr:name_id,omitempty"`
}

type Volumes struct {
	Volumes []Volume `json:"volumes,omitempty"`
}

type VolumeDetail struct {
	Volume Volume `json:"volume,omitempty"`
}

func (evs *Evsclient) CreatVolume(project_id string, name string, size int, az string, volume_type string, shareable bool, token string) (*VolumeDetail, error) {
	router := fmt.Sprintf("%s/v2/%s/volumes", evs.endpoint, project_id)
	requstbody := &VolumeDetail{
		Volume: Volume{
			Name:             name,
			Size:             size,
			AvailabilityZone: az,
			VolumeType:       volume_type,
			Shareable:        shareable,
			Multiattach:      shareable,
		},
	}
	str, _ := json.Marshal(requstbody)
	body := bytes.NewBuffer([]byte(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 202 {
		return nil, errors.New("return code is not 202 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result VolumeDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListImages lists the images according with ImageFilterOptions.
func (evs *Evsclient) ListVolumes(project_id string, token string) (*Volumes, error) {
	router := fmt.Sprintf("%s/v2/%s/volumes", evs.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Volumes
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (evs *Evsclient) GetVolume(project_id string, volume_id string, token string) (*VolumeDetail, error) {
	router := fmt.Sprintf("%s/v2/%s/volumes/%s", evs.endpoint, project_id, volume_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result VolumeDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (evs *Evsclient) DeleteVolume(project_id string, volume_id string, token string) error {
	router := fmt.Sprintf("%s/v2/%s/volumes/%s", evs.endpoint, project_id, volume_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 202 {
		return errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
