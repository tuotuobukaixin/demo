package elb_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

// Asynchronous job query response
type Job struct {
	Status   string `json:"status"`
	Entities struct {
		Elb struct {
			LoadbalancerId string `json:"id"`
		} `json:"elb"`
		Members []struct {
			Address string `json:"address"`
			ID      string `json:"id"`
		} `json:"members"`
	} `json:"entities"`
	FailReason string `json:"fail_reason"`
	ErrorCode  string `json:"error_code"`
}

func (elb *Elbclient) GetJob(project_id string, job_id string, token string) (*Job, error) {
	router := fmt.Sprintf("%s/v1.0/%s/jobs/%s", elb.endpoint, project_id, job_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Job
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
