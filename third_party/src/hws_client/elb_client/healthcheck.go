package elb_client

import (
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

func (elb *Elbclient) Deletehealthcheck(project_id string, healthcheck_id string, token string) error {
	router := fmt.Sprintf("%s/v1.0/%s/elbaas/healthcheck/%s", elb.endpoint, project_id, healthcheck_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 204 {
		return errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
