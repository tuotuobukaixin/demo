package ecs_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type StopOpt struct {
	Type    string     `json:"type"`
	Servers []ServerID `json:"servers,omitempty"`
}
type StopEcsOpts struct {
	StopOpt StopOpt `json:"os-stop"`
}

// ListImages lists the images according with ImageFilterOptions.
func (ecs *Ecsclient) StopServer(project_id string, opt *StopEcsOpts, token string) (*JobDetail, error) {
	router := fmt.Sprintf("%s/v1/%s/cloudservers/action", ecs.endpoint, project_id)

	str, _ := json.Marshal(opt)
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
