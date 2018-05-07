package vpc_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

//VpcQuotas
type VpcQuotas struct {
	Quotas struct {
		Resources []struct {
			Type  string `json:"type"`
			Used  int    `json:"used"`
			Quota int    `json:"quota"`
			Min   int    `json:"min"`
		} `json:"resources"`
	} `json:"quotas"`
}

func (vpc *Vpcclient) GetVpcQuotas(types string, project_id string, token string) (*VpcQuotas, error) {
	var router string
	if types == "" {
		router = fmt.Sprintf("%s/v1/%s/quotas", vpc.endpoint, project_id)
	} else {
		router = fmt.Sprintf("%s/v1/%s/quotas?type=%s", vpc.endpoint, project_id, types)
	}

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result VpcQuotas
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
