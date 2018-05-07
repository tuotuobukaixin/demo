package evs_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

// Vpc represents, well, a vpc.
type Tag struct {
	Tag map[string]string `json:"tag,omitempty"`
}

func (evs *Evsclient) CreatVolumeTag(project_id string, tag map[string]string, volume_id string, token string) (*Tag, error) {
	router := fmt.Sprintf("%s/v2/%s/os-vendor-tags/volumes/%s", evs.endpoint, project_id, volume_id)
	requstbody := Tag{
		Tag: tag,
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
	var result Tag
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
