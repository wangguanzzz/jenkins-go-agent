// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0
// snippet-start:[ec2.go.create_instance_with_tag]
package main

// snippet-start:[ec2.go.create_instance_with_tag.imports]
import (
	"encoding/base64"
	// "flag"
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

func TerminateInstance(svc ec2iface.EC2API, instanceId string) {

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
func MakeInstance(svc ec2iface.EC2API, jenkins, subnetId, sg1, sg2, key, ami, vmtype string) (string, error) {
	// name := flag.String("n", "Name", "The name of the tag to attach to the instance")
	// value := flag.String("v", "JenkinsAgent", "The value of the tag to attach to the instance")
	// flag.Parse()

	ud := `#!/bin/bash
	export JENKINS_MASTER=` + jenkins + `
	nohup /root/command.sh &`

	userData := base64.StdEncoding.EncodeToString([]byte(ud))

	sgs := []*string{&sg1, &sg2}
	// subnet := "subnet-c4686fbc"
	// snippet-start:[ec2.go.create_instance_with_tag.call]
	fmt.Println("here ..................................................................")
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
	fmt.Println("here2 ..................................................................")
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
	jenkins_master := "18.163.184.77:80"
	cool_down := 60
	scan_frequency := 10
	vm_cap := 2
	idle_cap := 60
	// [instance_id] = idle second
	agentMap := make(map[string]int, vm_cap)
	svc := getSvc()

	for true {
		time.Sleep(time.Duration(scan_frequency) * time.Second)
		// create vm process
		isQueueStuck := queryQueue(jenkins_master, agentMap)
		fmt.Println("queue status is block? ", isQueueStuck)
		if len(agentMap) < vm_cap && isQueueStuck {
			fmt.Println(" start strigger vm creating")
			instanceID, err := MakeInstance(svc, jenkins_master, "subnet-c4686fbc", "sg-0891ebcae20b9ea2e", "sg-95d428fc", "hk_region", "ami-02986db8fa9f47e57", "t3.micro")
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
					TerminateInstance(svc, k)
					delete(agentMap, k)
				}
			} else {
				agentMap[k] = 0
			}
		}

	}

}

// snippet-end:[ec2.go.create_instance_with_tag]
