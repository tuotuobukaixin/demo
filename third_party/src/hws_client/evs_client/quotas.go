package evs_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

//VpcQuotas
type EvsQuota struct {
	QuotaSet struct {
		BackupGigabytes struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"backup_gigabytes"`
		Backups struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"backups"`
		Gigabytes struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"gigabytes"`
		GigabytesSAS struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"gigabytes_SAS"`
		GigabytesSATA struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"gigabytes_SATA"`
		GigabytesSSD struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"gigabytes_SSD"`
		ID        string `json:"id"`
		Snapshots struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"snapshots"`
		SnapshotsSAS struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"snapshots_SAS"`
		SnapshotsSATA struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"snapshots_SATA"`
		SnapshotsSSD struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"snapshots_SSD"`
		Volumes struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"volumes"`
		VolumesSAS struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"volumes_SAS"`
		VolumesSATA struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"volumes_SATA"`
		VolumesSSD struct {
			InUse int `json:"in_use"`
			Limit int `json:"limit"`
			Inuse int `json:"inuse"`
		} `json:"volumes_SSD"`
	} `json:"quota_set"`
}

func (evs *Evsclient) GetEVSQuotas(project_id string, token string) (*EvsQuota, error) {
	var router string

	router = fmt.Sprintf("%s/v2/%s/os-quota-sets/%s?usage=True", evs.endpoint, project_id, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result EvsQuota
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
