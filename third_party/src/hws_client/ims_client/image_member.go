package ims_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

const (
	ACCEPT = "accepted"
	REJECT = "rejected"
)

// Image member detail information
type Member struct {
	// The image ID
	ImageID string `json:"image_id"`
	// The image member ID
	MemberID string `json:"member_id"`
	// The schema URL
	Schema string `json:"schema"`
	// The date and time when the resource was created.
	Created string `json:"created_at"`
	// The date and time when the resource was updated.
	Updated string `json:"updated_at"`
	// The status of create image member.
	Status string `json:"status"`
}

type ImageMemberCreate struct {
	Member string `json:"member"`
}
type ImageMemberUpdate struct {
	Status string `json:"status"`
}

func (ims *Imsclient) CreateImageMember(image_id string, member_id string, token string) (*Member, error) {

	router := fmt.Sprintf("%s/v2/images/%s/members", ims.endpoint, image_id)

	requstbody := ImageMemberCreate{
		Member: member_id,
	}

	str, _ := json.Marshal(requstbody)
	body := bytes.NewBuffer([]byte(str))
	fmt.Println(string(str))
	rsp, status_code, _, err := common.DoHttpRequest("POST", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 201 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Member
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (ims *Imsclient) UpdateImageMember(image_id string, member_id string, token string) (*Member, error) {

	router := fmt.Sprintf("%s/v2/images/%s/members/%s", ims.endpoint, image_id, member_id)

	requstbody := ImageMemberUpdate{
		Status: ACCEPT,
	}

	str, _ := json.Marshal(requstbody)
	body := bytes.NewBuffer([]byte(str))
	fmt.Println(string(str))
	rsp, status_code, _, err := common.DoHttpRequest("PUT", router, "application/json", body, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 201 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Member
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (ims *Imsclient) DeleteImageMember(image_id string, member_id string, token string) (*Member, error) {

	router := fmt.Sprintf("%s/v2/images/%s/members/%s", ims.endpoint, image_id, member_id)

	rsp, status_code, _, err := common.DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if err != nil {
		return nil, err
	}
	if status_code != 204 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result Member
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}
