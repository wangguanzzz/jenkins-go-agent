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
func MakeInstance(jenkins, subnetId, sg, key, ami, vmtype string) (string, error) {
	name := flag.String("n", "Name", "The name of the tag to attach to the instance")
	value := flag.String("v", "JenkinsAgent", "The value of the tag to attach to the instance")
	flag.Parse()

	svc := getSvc()

	ud := `#!/bin/bash
	export JENKINS_MASTER=` + jenkins + `nohup /root/command.sh &`

	userData := base64.StdEncoding.EncodeToString([]byte(ud))

	sgs := []*string{&sg}
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
		//SubnetId:         &subnet,
	})
	// snippet-end:[ec2.go.create_instance_with_tag.call]
	if err != nil {
		return "", err
	}

	// snippet-start:[ec2.go.create_instance_with_tag.tag]
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{result.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   name,
				Value: value,
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
	result, err := MakeInstance("18.163.184.77:80", "subnet-c4686fbc", "chris", "hk_region", "ami-02986db8fa9f47e57", "t3.micro")
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(20 * time.Second)
	TerminateInstance(result)
}

// snippet-end:[ec2.go.create_instance_with_tag]
