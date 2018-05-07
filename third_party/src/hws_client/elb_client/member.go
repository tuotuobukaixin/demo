package elb_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

// Backend host Member
type Member struct {
	// create Member parameters
	ListenerID string `json:"listener_id,omitempty"`
	ServerID   string `json:"server_id"`
	Address    string `json:"address"`
}
type JobResp struct {
	JobID string `json:"job_id"`
	Uri   string `json:"uri"`
}
type MemDetail struct {
	Member
	ServerAddress string              `json:"server_address"`
	ID            string              `json:"id"`
	Status        string              `json:"status"`
	Listeners     []map[string]string `json:"listeners"`
	ServerName    string              `json:"server_name"`
	HealthStatus  string              `json:"health_status"`
}
type MemberRm struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type MembersDel struct {
	RemoveMember []MemberRm `json:"removeMember"`
}

func (elb *Elbclient) GetMemberList(project_id string, listenerid string, token string) ([]MemDetail, error) {
	router := fmt.Sprintf("%s/v1.0/%s/elbaas/listeners/%s/members", elb.endpoint, project_id, listenerid)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result []MemDetail
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (elb *Elbclient) DeleteMember(project_id string, listenerid string, opts *MembersDel, token string) (*JobResp, error) {
	router := fmt.Sprintf("%s/v1.0/%s/elbaas/listeners/%s/members/action", elb.endpoint, project_id, listenerid)
	str, _ := json.Marshal(opts)
	body := bytes.NewBuffer([]byte(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result JobResp
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
