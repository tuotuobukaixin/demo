package ecs_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

type Keypair struct {
	// The keypair name.
	Name string `json:"name,omitempty"`
	// The keypair public key.
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	// The fingerprint for the keypair.
	Fingerprint string `json:"fingerprint,omitempty"`
	// The user id for a keypair.
	UserID string `json:"user_id,omitempty"`
	// Indicates whether this keypai has been deleted or not.
	Deleted bool `json:"deleted"`
	// The date and time when the resource was created. The date and time stamp format is ISO 8601.
	Created string `json:"created,omitempty"`
	// The date and time when the resource was updated. The date and time stamp format is ISO 8601.
	Updated string `json:"updated,omitempty"`
	// The date and time when the resource was deleted. The date and time stamp format is ISO 8601.
	DeletedAt string `json:"deleted_at,omitempty"`
}

type KeypairDetails struct {
	Keypair Keypair `json:"keypair"`
}

type Keypairs struct {
	Keypairs []KeypairDetails `json:"keypairs"`
}

// ListImages lists the images according with ImageFilterOptions.
func (ecs *Ecsclient) CreatKeypair(project_id string, name string, public_key string, token string) (*KeypairDetails, error) {
	router := fmt.Sprintf("%s/v2/%s/os-keypairs", ecs.endpoint, project_id)
	var requstbody *KeypairDetails
	if public_key == "" {
		requstbody = &KeypairDetails{
			Keypair: Keypair{
				Name: name,
			},
		}
	} else {
		requstbody = &KeypairDetails{
			Keypair: Keypair{
				Name:      name,
				PublicKey: public_key,
			},
		}
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
	var result KeypairDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListImages lists the images according with ImageFilterOptions.
func (ecs *Ecsclient) ListKeypairs(project_id string, token string) (*Keypairs, error) {
	router := fmt.Sprintf("%s/v2/%s/os-keypairs", ecs.endpoint, project_id)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Keypairs
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) GetKeypair(project_id string, name string, token string) (*KeypairDetails, error) {
	router := fmt.Sprintf("%s/v2/%s/os-keypairs/%s", ecs.endpoint, project_id, name)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result KeypairDetails
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ecs *Ecsclient) DeleteKeypair(project_id string, name string, token string) error {
	router := fmt.Sprintf("%s/v2/%s/os-keypairs/%s", ecs.endpoint, project_id, name)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return err
	}
	if status_code != 202 {
		return errors.New("return code is not 202 :" + strconv.Itoa(status_code) + string(rsp))
	}

	return nil
}
