package hws_client

import (
	"encoding/json"
	"hws_client/ecs_client"
	"hws_client/elb_client"
	"hws_client/evs_client"
	"hws_client/iam_client"
	"hws_client/ims_client"
	"hws_client/vpc_client"
	"io/ioutil"
)

type hws_config struct {
	ECS_Endpoint string
	EVS_Endpoint string
	IMS_Endpoint string
	VPC_Endpoint string
	IAM_Endpoint string
	ELB_Endpoint string
}
type Hws_Client struct {
	Ecsclient *ecs_client.Ecsclient
	Evsclient *evs_client.Evsclient
	Iamclient *iam_client.Iamclient
	Imsclient *ims_client.Imsclient
	Vpcclient *vpc_client.Vpcclient
	Elbclient *elb_client.Elbclient
}

func GetHwsClient(config string) (*Hws_Client, error) {
	content, err := ioutil.ReadFile(config)
	var conf hws_config
	if err == nil {
		err = json.Unmarshal(content, &conf)
		if err != nil {
			return nil, err
		}

	}
	hws_client := &Hws_Client{
		Ecsclient: nil,
		Evsclient: nil,
		Iamclient: nil,
		Imsclient: nil,
		Vpcclient: nil,
		Elbclient: nil,
	}
	if conf.ECS_Endpoint != "" {
		hws_client.Ecsclient = ecs_client.NewEcsClient(conf.ECS_Endpoint)
	}
	if conf.EVS_Endpoint != "" {
		hws_client.Evsclient = evs_client.NewEvsClient(conf.EVS_Endpoint)
	}
	if conf.IMS_Endpoint != "" {
		hws_client.Imsclient = ims_client.NewImsClient(conf.IMS_Endpoint)
	}
	if conf.VPC_Endpoint != "" {
		hws_client.Vpcclient = vpc_client.NewVpcClient(conf.VPC_Endpoint)
	}
	if conf.IAM_Endpoint != "" {
		hws_client.Iamclient = iam_client.NewIamClient(conf.IAM_Endpoint)
	}
	if conf.ELB_Endpoint != "" {
		hws_client.Elbclient = elb_client.NewElbClient(conf.ELB_Endpoint)
	}
	return hws_client, nil
}
