// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
// snippet-start:[ec2.go.create_instance_with_tag]
package main

// snippet-start:[ec2.go.create_instance_with_tag.imports]
import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"time"
)

// snippet-end:[ec2.go.create_instance_with_tag.imports]

// MakeInstance creates an Amazon Elastic Compute Cloud (Amazon EC2) instance
// Inputs:
//     svc is an Amazon EC2 service client
//     key is the name of the tag to attach to the instance
//     value is the value of the tag to attach to the instance
// Output:
//     If success, nil
//     Otherwise, an error from the call to RunInstances or CreateTags

func getSvc() ec2iface.EC2API {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)
	return svc
}

func TerminateInstance(instanceId string) {
	svc := getSvc()
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	result, err := svc.TerminateInstances(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}
func MakeInstance(jenkins, subnetId, sg1, sg2, key, ami, vmtype string) (string, error) {

	svc := getSvc()
	ud := `#!/bin/bash
	export JENKINS_MASTER=` + jenkins + `
	nohup /root/command.sh &`

	userData := base64.StdEncoding.EncodeToString([]byte(ud))

	sgs := []*string{&sg1, &sg2}
	// subnet := "subnet-c4686fbc"
	// snippet-start:[ec2.go.create_instance_with_tag.call]
	result, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(ami),
		InstanceType:     aws.String(vmtype),
		SubnetId:         &subnetId,
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          &key,
		UserData:         &userData,
		SecurityGroupIds: sgs,
	})
	// snippet-end:[ec2.go.create_instance_with_tag.call]
	if err != nil {
		return "", err
	}

	// snippet-start:[ec2.go.create_instance_with_tag.tag]
	name := "Name"
	tag := "JenkinsAgent"
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{result.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   &name,
				Value: &tag,
			},
		},
	})
	// snippet-end:[ec2.go.create_instance_with_tag.tag]
	if err != nil {
		return "", err
	}

	return *result.Instances[0].InstanceId, nil
}

func main() {
	// application setting

	fjenkins := flag.String("jenkins", "127.0.0.1:8080", "Jenkins Master URL like 127.0.0.1:8080")
	fcooldown := flag.Int("cooldown", 60, "cool down period")
	fscan := flag.Int("frequency", 10, "scan frequency")
	fcap := flag.Int("vmcap", 2, "agent number cap")
	fidle := flag.Int("idle", 60, "idle period")

	// aws setting
	//"subnet-c4686fbc"
	subnet := flag.String("subnet", "", "AWS subnet ID")
	//"sg-0891ebcae20b9ea2e", "sg-95d428fc"
	sg1 := flag.String("sg1", "", "AWS security group 1")
	sg2 := flag.String("sg2", "", "AWS security group 2")
	//"hk_region"
	keyName := flag.String("key", "", "AWS key name")
	// "ami-02986db8fa9f47e57"
	ami := flag.String("ami", "", "AWS key name")
	// "t3.micro"
	vmtype := flag.String("vmtype", "t3.micro", "AWS vm type")

	flag.Parse()

	jenkins_master := *fjenkins //"18.163.184.77:80"
	cool_down := *fcooldown
	scan_frequency := *fscan
	vm_cap := *fcap
	idle_cap := *fidle
	// [instance_id] = idle second
	agentMap := make(map[string]int, vm_cap)

	for true {
		time.Sleep(time.Duration(scan_frequency) * time.Second)
		// create vm process
		isQueueStuck := queryQueue(jenkins_master, agentMap)
		fmt.Println("queue status is block? ", isQueueStuck)
		if len(agentMap) < vm_cap && isQueueStuck {
			fmt.Println(" start strigger vm creating")
			instanceID, err := MakeInstance(jenkins_master, *subnet, *sg1, *sg2, *keyName, *ami, *vmtype)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("vm " + instanceID + " is created")
			agentMap[instanceID] = 0
			fmt.Println("cooling down started ...")
			time.Sleep(time.Duration(cool_down) * time.Second)
		}
		// update idle status
		fmt.Println("checking idle status ... agent length is ", len(agentMap))
		for k, v := range agentMap {
			fmt.Println("checking idle status of  " + k)
			idle := queryAgent(jenkins_master, k)
			if idle {
				agentMap[k] += scan_frequency
				// deregister and shutdown idle vms
				if agentMap[k] > idle_cap {
					fmt.Println("vm " + k + " is idle too long " + string(v))
					deregisterAgent(jenkins_master, "jenkins-cli.jar", k)
					TerminateInstance(k)
					delete(agentMap, k)
				}
			} else {
				agentMap[k] = 0
			}
		}

	}

}

// snippet-end:[ec2.go.create_instance_with_tag]
