package elb_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type ELBListener struct {
	Name          string `json:"name"`
	ID            string `json:"id"`
	TenantID      string `json:"tenant_id"`
	Status        string `json:"status"`
	CreateTime    string `json:"create_time"`
	HealthcheckID string `json:"healthcheck_id"`
}

func (elb *Elbclient) GetListenersList(project_id string, token string) ([]ELBListener, error) {
	router := fmt.Sprintf("%s/v1.0/%s/elbaas/listeners", elb.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result []ELBListener
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (elb *Elbclient) DeleteListener(project_id string, listenid string, token string) error {
	router := fmt.Sprintf("%s/v1.0/%s/elbaas/listeners/%s", elb.endpoint, project_id, listenid)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 204 {
		return errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
