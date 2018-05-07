package ecs_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type VolumeAttachmentDetail struct {
	VolumeAttachment VolumeAttachment `json:"volumeAttachment"`
}
type VolumeAttachments struct {
	VolumeAttachments []VolumeAttachment `json:"volumeAttachments"`
}
type VolumeAttachment struct {
	Device   string `json:"device"`
	ID       string `json:"id"`
	ServerId string `json:"serverId"`
	VolumeId string `json:"volumeId"`
}

func (ecs *Ecsclient) GetVolumeAttachments(project_id string, server_id string, token string) (*VolumeAttachments, error) {
	router := fmt.Sprintf("%s/v2/%s/servers/%s/os-volume_attachments", ecs.endpoint, project_id, server_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result VolumeAttachments
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) GetVolumeAttachment(project_id string, server_id string, volume_id string, token string) (*VolumeAttachmentDetail, error) {
	router := fmt.Sprintf("%s/v2/%s/servers/%s/os-volume_attachments/%s", ecs.endpoint, project_id, server_id, volume_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result VolumeAttachmentDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) AttachVolume(project_id string, server_id string, volume_id string, device string, token string) (*VolumeAttachmentDetail, error) {
	router := fmt.Sprintf("%s/v2/%s/servers/%s/os-volume_attachments", ecs.endpoint, project_id, server_id)
	requstbody := &VolumeAttachmentDetail{
		VolumeAttachment: VolumeAttachment{
			VolumeId: volume_id,
			Device:   device,
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
	var result VolumeAttachmentDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) DetachDisk(project_id string, server_id string, volume_id string, token string) error {
	router := fmt.Sprintf("%s/v1/%s/cloudservers/%s/detachvolume/%s", ecs.endpoint, project_id, server_id, volume_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 202 {
		return errors.New("return code is not 202 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
