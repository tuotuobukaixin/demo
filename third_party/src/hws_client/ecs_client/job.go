package ecs_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

const (
	JobInit    = "INIT"
	JobSuccess = "SUCCESS"
	JobFailed  = "FAILED"
	JobRunning = "RUNNING"
)

type SubJob struct {
	BeginTime  string            `json:"begin_time,omitempty"`
	EndTime    string            `json:"end_time,omitempty"`
	ErrorCode  string            `json:"error_code,omitempty"`
	FailReason string            `json:"fail_reason,omitempty"`
	JobID      string            `json:"job_id,omitempty"`
	JobType    string            `json:"job_type,omitempty"`
	Status     string            `json:"status,omitempty"`
	Entities   map[string]string `json:"entities,omitempty"`
}

type Entities struct {
	SubJobsTotal int      `json:"sub_jobs_total,omitempty"`
	SubJobs      []SubJob `json:"sub_jobs,omitempty"`
}

type JobDetail struct {
	Status     string   `json:"status,omitempty"`
	JobID      string   `json:"job_id,omitempty"`
	JobType    string   `json:"job_type,omitempty"`
	BeginTime  string   `json:"begin_time,omitempty"`
	EndTime    string   `json:"end_time,omitempty"`
	ErrorCode  string   `json:"error_code,omitempty"`
	FailReason string   `json:"fail_reason,omitempty"`
	Message    string   `json:"message,omitempty"`
	Code       string   `json:"code,omitempty"`
	Entities   Entities `json:"entities,omitempty"`
}

func (ecs *Ecsclient) GetJob(project_id string, job_id string, token string) (*JobDetail, error) {
	router := fmt.Sprintf("%s/v1/%s/jobs/%s", ecs.endpoint, project_id, job_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result JobDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
