package main

import (
	"fmt"
	"hws_client/ecs_client"
	"hws_client/elb_client"
	"hws_client/evs_client"
	"hws_client/iam_client"
	"hws_client/vpc_client"
	"time"
)

func main() {
	iam := iam_client.NewIamClient("https://10.177.25.248:443")
	ecs := ecs_client.NewEcsClient("https://10.177.25.248:443")
	vpc := vpc_client.NewVpcClient("https://10.177.25.248:443")
	evs := evs_client.NewEvsClient("https://10.177.25.248:443")
	_ = elb_client.NewElbClient("https://10.177.25.248:443")
	token, err := iam.CreatUserToken("Huawei@123", "newbee", "newbee", "eu-de")
	fmt.Println("==============create token===============")
	fmt.Println("%v", err)
	aa, _ := iam.ConveteToken(token)
	projectid := aa.Token.Project.ID
	fmt.Println("==============convet token===============")
	fmt.Println(projectid)
	fmt.Println("%v", aa)
	evsquotas, err := ecs.GetInterfaceAttachments(projectid, "0be572c3-ca01-4ac5-b793-07f253b51809",token)
	fmt.Println("%v", err)
	fmt.Println("%v", evsquotas.InterfaceAttachments)
	return
	quotas, err := vpc.GetVpcQuotas("", projectid, token)
	fmt.Println("==============get vpc quota===============")
	fmt.Println("%v", err)
	fmt.Println("%v", quotas)

	vpcs, err := vpc.ListVpcs(projectid, token)
	fmt.Println("==============list vpcs===============")
	fmt.Println("%v", err)
	fmt.Println("%v", vpcs)
	for _, tmpvpc := range vpcs.Vpcs {
		if tmpvpc.Name == "ccetestvpc" {
			fmt.Println("==============list subnets===============")
			subnets, err := vpc.ListSubnets(projectid, tmpvpc.ID, token)
			fmt.Println("%v", err)
			fmt.Println("%v", subnets)
			for _, tmpsubnet := range subnets.Subnets {
				if tmpsubnet.Name == "ccetestsubnet" {
					fmt.Println("==============delete subnet===============")
					err := vpc.DeleteSubnet(projectid, tmpvpc.ID, tmpsubnet.ID, token)
					fmt.Println("%v", err)
				}
			}
			for index := 0; index < 10; index++ {
				subnets, err := vpc.ListSubnets(projectid, tmpvpc.ID, token)
				fmt.Println("%v", err)
				fmt.Println("%v", subnets)
				if len(subnets.Subnets) == 0 {
					break
				}
				time.Sleep(1)
			}
			fmt.Println("==============delete vpcs===============")
			err = vpc.DeleteVpc(projectid, tmpvpc.ID, token)
			fmt.Println("%v", err)
		}
	}
	fmt.Println("==============create vpcs===============")
	vpcde, err := vpc.CreatVpc(projectid, "ccetestvpc", token)
	fmt.Println("%v", err)
	fmt.Println("%v", vpcde)
	fmt.Println("==============get vpcs===============")
	for index := 0; index < 10; index++ {
		tmp, err := vpc.GetVpc(projectid, vpcde.Vpc.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", tmp)
		if tmp.Vpc.Status == vpc_client.VPC_OK {
			break
		}
		time.Sleep(1)
	}

	fmt.Println("==============list subnets===============")
	subnets, err := vpc.ListSubnets(projectid, vpcde.Vpc.ID, token)
	fmt.Println("%v", err)
	fmt.Println("%v", subnets)
	for _, tmpsubnet := range subnets.Subnets {
		if tmpsubnet.Name == "ccetestsubnet" {
			fmt.Println("==============delete subnet===============")
			err := vpc.DeleteSubnet(projectid, vpcde.Vpc.ID, tmpsubnet.ID, token)
			fmt.Println("%v", err)
		}
	}
	for index := 0; index < 10; index++ {
		subnets, err := vpc.ListSubnets(projectid, vpcde.Vpc.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", subnets)
		if len(subnets.Subnets) == 0 {
			break
		}
		time.Sleep(1)
	}
	fmt.Println("==============create subnet===============")
	subnet, err := vpc.CreatSubnet(projectid, "ccetestsubnet", "192.168.20.0/24", "192.168.20.1", vpcde.Vpc.ID, token)
	fmt.Println("%v", err)
	fmt.Println("%v", subnet)
	fmt.Println("==============Get subnet===============")
	for index := 0; index < 10; index++ {
		tmp, err := vpc.GetSubnet(projectid, subnet.Subnet.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", tmp)
		if tmp.Subnet.Status == vpc_client.SUBNET_ACTIVE {
			break
		}
		time.Sleep(1)
	}

	fmt.Println("==============list secgroup===============")
	secgroups, err := vpc.ListSecurityGroups(token)
	fmt.Println("%v", err)
	fmt.Println("%v", secgroups)
	for _, tmpsecgroup := range secgroups.SecurityGroups {
		if tmpsecgroup.Name == "ccetestsecgroup" {
			fmt.Println("==============delete secgroup===============")
			err = vpc.DeleteSecurityGroup(tmpsecgroup.ID, token)
			fmt.Println("%v", err)
		}
	}
	fmt.Println("==============create secgroup===============")
	secgroup, err := vpc.CreatSecurityGroup("ccetestsecgroup", token)
	fmt.Println("%v", err)
	fmt.Println("%v", secgroup)

	fmt.Println("==============get secgroup===============")
	tmpgroup, err := vpc.GetSecurityGroup(secgroup.SecurityGroup.ID, token)
	fmt.Println("%v", err)
	fmt.Println("%v", tmpgroup)

	fmt.Println("==============list secgrouprules===============")
	secgrouprules, err := vpc.ListSecurityGroupRules(secgroup.SecurityGroup.ID, token)
	fmt.Println("%v", err)
	fmt.Println("%v", secgroups)
	for _, tmpsecgroupsule := range secgrouprules.SecurityGroupRules {
		fmt.Println("==============delete secgrouprules===============")
		err := vpc.DeleteSecurityGroupRule(tmpsecgroupsule.ID, token)
		fmt.Println("%v", err)
	}
	fmt.Println("==============create secgrouprule===============")

	secgrouprule, err := vpc.CreatSecurityGroupRule(secgroup.SecurityGroup.ID, vpc_client.INBOUND, "tcp", "0.0.0.0/0", 22, 22, token)
	fmt.Println("%v", err)
	fmt.Println("%v", secgrouprule)
	fmt.Println("==============get secgrouprule===============")
	tmpgrouprule, err := vpc.GetSecurityGroupRule(secgrouprule.SecurityGroupRule.ID, token)
	fmt.Println("%v", err)
	fmt.Println("%v", tmpgrouprule)

	fmt.Println("==============get ecsquota===============")
	ecsquota, err := ecs.GetECSQuotas(projectid, token)
	fmt.Println("%v", err)
	fmt.Println("%v", ecsquota)
	fmt.Println("==============get falavors===============")
	flavors, err := ecs.ListFlavors(projectid, token)
	fmt.Println("%v", err)
	fmt.Println("%v", flavors)
	for _, tmpf := range flavors.Flavors {
		fmt.Println("==============get falavor===============")
		flavor, err := ecs.GetFlavor(projectid, tmpf.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", flavor)
	}
	fmt.Println("==============list key===============")
	keys, err := ecs.ListKeypairs(projectid, token)
	fmt.Println("%v", err)
	fmt.Println("%v", keys)
	for _, tmpkey := range keys.Keypairs {
		if tmpkey.Keypair.Name == "ccetestkey" {
			fmt.Println("==============delete key===============")
			err = ecs.DeleteKeypair(projectid, tmpkey.Keypair.Name, token)
			fmt.Println("%v", err)
		}
	}
	fmt.Println("==============create key===============")
	key, err := ecs.CreatKeypair(projectid, "ccetestkey", "", token)
	fmt.Println("%v", err)
	fmt.Println("%v", key)
	fmt.Println("==============get key===============")
	tmp, err := ecs.GetKeypair(projectid, "ccetestkey", token)
	fmt.Println("%v", err)
	fmt.Println("%v", tmp)
	fmt.Println("==============create server===============")
	requstbody := &ecs_client.CreateServerOpts{

		Name:             "cce_test",
		ImageRef:         "c66d2d1b-785c-4dbf-ae40-e2bde716a61b",
		FlavorRef:        "highmem1",
		KeyName:          key.Keypair.Name,
		VpcId:            vpcde.Vpc.ID,
		RootVolume:       &ecs_client.Volume{VolumeType: "SATA", Size: 40},
		PublicIp:         &ecs_client.PublicIp{EIP: ecs_client.Eip{Iptype: "5_bgp", Bandwidth: ecs_client.Bandwidth{Size: 5, Sharetype: "PER"}}},
		AvailabilityZone: "eu-de-02",
		Count:            1,
	}
	requstbody.SecurityGroups = append(requstbody.SecurityGroups, ecs_client.SecurityGroup{ID: secgroup.SecurityGroup.ID})
	requstbody.Nics = append(requstbody.Nics, ecs_client.Nic{SubnetId: subnet.Subnet.ID})
	requstbody.DataVolumes = append(requstbody.DataVolumes, ecs_client.Volume{VolumeType: "SATA", Size: 100})
	job, err := ecs.CreatServers(projectid, requstbody, token)
	fmt.Println("%v", err)
	fmt.Println("%v", job)
	server_id := ""
	for index := 0; index < 40; index++ {
		jobdeta, err := ecs.GetJob(projectid, job.JobID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", jobdeta)
		if jobdeta.Entities.SubJobsTotal != 0 && jobdeta.Entities.SubJobsTotal == len(jobdeta.Entities.SubJobs) {
			server_id = jobdeta.Entities.SubJobs[0].Entities["server_id"]
		}
		if jobdeta.Status == ecs_client.JobSuccess {
			break
		}
		time.Sleep(10 * time.Second)
	}
	fmt.Println("==============list volumes===============")
	volumes, err := evs.ListVolumes(projectid, token)
	fmt.Println("%v", err)
	fmt.Println("%v", volumes)
	for _, tmpvolume := range volumes.Volumes {
		if tmpvolume.Name == "ccetestvolume" {
			err := evs.DeleteVolume(projectid, tmpvolume.ID, token)
			fmt.Println("%v", err)
		}
	}
	fmt.Println("==============create volumes===============")
	volume, err := evs.CreatVolume(projectid, "ccetestvolume", 10, "eu-de-02", "SATA", true, token)
	fmt.Println("%v", err)
	fmt.Println("%v", volume)
	for index := 0; index < 10; index++ {
		tmp, err := evs.GetVolume(projectid, volume.Volume.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", tmp)
		if tmp.Volume.Status == "available" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	fmt.Println("==============add tag volumes===============")
	tag := make(map[string]string)
	tag["cce"] = "cce"
	tagde, err := evs.CreatVolumeTag(projectid, tag, volume.Volume.ID, token)
	fmt.Println("%v", err)
	fmt.Println("%v", tagde)

	fmt.Println("==============attach volumes===============")
	attach, err := ecs.AttachVolume(projectid, server_id, volume.Volume.ID, "/dev/sdc", token)
	fmt.Println("%v", err)
	fmt.Println("%v", attach)
	fmt.Println("==============get attach===============")
	tmpattach, err := ecs.GetVolumeAttachment(projectid, server_id, volume.Volume.ID, token)
	fmt.Println("%v", err)
	fmt.Println("%v", tmpattach)
	for index := 0; index < 10; index++ {
		tmp, err := evs.GetVolume(projectid, volume.Volume.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", tmp)
		if tmp.Volume.Status == "in-use" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	fmt.Println("==============dettach attach===============")
	err = ecs.DetachDisk(projectid, server_id, volume.Volume.ID, token)
	fmt.Println("%v", err)
	for index := 0; index < 10; index++ {
		tmp, err := evs.GetVolume(projectid, volume.Volume.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", tmp)
		if tmp.Volume.Status == "available" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	fmt.Println("==============delete volume===============")
	err = evs.DeleteVolume(projectid, volume.Volume.ID, token)
	fmt.Println("%v", err)
	if server_id != "" {
		deletebody := &ecs_client.DeleteServerOpts{

			DeletePublicIP: true,
			DeleteVolume:   true,
		}
		fmt.Println("==============delete servers===============")
		deletebody.Servers = append(deletebody.Servers, ecs_client.ServerID{ID: server_id})
		deletejob, err := ecs.DeleteServers(projectid, deletebody, token)
		fmt.Println("%v", err)
		fmt.Println("%v", deletejob)

		for index := 0; index < 20; index++ {
			jobdeta, err := ecs.GetJob(projectid, deletejob.JobID, token)
			fmt.Println("%v", err)
			fmt.Println("%v", jobdeta)
			if jobdeta.Status == ecs_client.JobSuccess {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	fmt.Println("==============delete key===============")
	err = ecs.DeleteKeypair(projectid, key.Keypair.Name, token)
	fmt.Println("%v", err)
	fmt.Println("==============delete secgrouprules===============")
	err = vpc.DeleteSecurityGroupRule(secgrouprule.SecurityGroupRule.ID, token)
	fmt.Println("%v", err)
	fmt.Println("==============delete secgroup===============")
	err = vpc.DeleteSecurityGroup(secgroup.SecurityGroup.ID, token)
	fmt.Println("%v", err)
	fmt.Println("==============delete subnet===============")
	err = vpc.DeleteSubnet(projectid, vpcde.Vpc.ID, subnet.Subnet.ID, token)
	fmt.Println("%v", err)
	for index := 0; index < 10; index++ {
		subnets, err := vpc.ListSubnets(projectid, vpcde.Vpc.ID, token)
		fmt.Println("%v", err)
		fmt.Println("%v", subnets)
		if len(subnets.Subnets) == 0 {
			break
		}
		time.Sleep(1)
	}
	fmt.Println("==============delete vpcs===============")
	err = vpc.DeleteVpc(projectid, vpcde.Vpc.ID, token)
	fmt.Println("%v", err)

}
