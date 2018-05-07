package iam_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hws_client/common"
	"strconv"
)

//SubjectReq format
type AuthReqProject struct {
	Auth AuthProject `json:"auth"`
}

type AuthProject struct {
	Identity Identity `json:"identity"`
	Scope    ScopeProject    `json:"scope"`
}
type ScopeProject struct {
	Project Project `json:"project,omitempty"`
}

//SubjectReq format
type AuthReqDomain struct {
	Auth AuthDomain `json:"auth"`
}

type AuthDomain struct {
	Identity Identity `json:"identity"`
	Scope    ScopeDomain    `json:"scope"`
}
type ScopeDomain struct {
	Domain  Domain  `json:"domain,omitempty"`
}

type Identity struct {
	Password Password `json:"password"`
	Methods  []string `json:"methods"`
}
type Password struct {
	User User `json:"user"`
}
type User struct {
	Password string `json:"password"`
	Domain   Domain `json:"domain"`
	Name     string `json:"name"`
}
type Domain struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}
type Project struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

//KeyStoneAuth format
type KeyStoneAuth struct {
	Token struct {
		Roles []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
		User struct {
			Domain struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		Project struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Domain struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
		} `json:"project"`
	} `json:"token"`
}

// ListImages lists the images according with ImageFilterOptions.
func (iam *Iamclient) CreatUserToken(passwd string, user_name string, domain_name string, project_name string, project_id string) (string, error) {
	var str []byte
	methods := []string{"password"}
	router := fmt.Sprintf("%s/v3/auth/tokens", iam.endpoint)
	if project_name != "" || project_id != "" {
		requstbody := AuthReqProject{
			Auth: AuthProject{
				Identity: Identity{
					Methods: methods,
					Password: Password{
						User: User{
							Password: passwd,
							Name:     user_name,
							Domain: Domain{
								Name: domain_name,
							},
						},
					},
				},
				Scope: ScopeProject{
					Project: Project{
						Name: project_name,
						ID:   project_id,
					},
				},
			},
		}
		str, _ = json.Marshal(requstbody)
	}else {
		requstbody := AuthReqDomain{
			Auth: AuthDomain{
				Identity: Identity{
					Methods: methods,
					Password: Password{
						User: User{
							Password: passwd,
							Name:     user_name,
							Domain: Domain{
								Name: domain_name,
							},
						},
					},
				},
				Scope: ScopeDomain{
					Domain: Domain{
						Name: domain_name,
					},
				},
			},
		}
		str, _ = json.Marshal(requstbody)
	}
	body := bytes.NewBuffer([]byte(str))
	rsp, status_code, hearder, err := common.DoHttpRequest("POST", router, "application/json", body, "", "")
	if err != nil {
		return "", err
	}
	if status_code != 201 {
		return "", errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	token := hearder.Get("X-Subject-Token")

	return token, nil
}

func (iam *Iamclient) ConveteToken(token string) (*KeyStoneAuth, error) {
	router := fmt.Sprintf("%s/v3/auth/tokens", iam.endpoint)

	rsp, status_code, _, err := common.DoHttpRequest("GET", router, "application/json", nil, token, token)
	if err != nil {
		return nil, err
	}
	if status_code != 200 {
		return nil, errors.New("return code is not 200 :" + strconv.Itoa(status_code) + string(rsp))
	}
	var result KeyStoneAuth
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
