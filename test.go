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
func MakeInstance(svc ec2iface.EC2API, name, value *string) (*ec2.Reservation, error) {
	key := "hk_region"
	ud := `#!/bin/bash
	export JENKINS_MASTER=18.162.47.230:80
	nohup /root/command.sh &
	`

	userData := base64.StdEncoding.EncodeToString([]byte(ud))
	sg := "chris"
	sgs := []*string{&sg}
	// subnet := "subnet-c4686fbc"
	// snippet-start:[ec2.go.create_instance_with_tag.call]
	result, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String("ami-01f41038383dc862b"),
		InstanceType:     aws.String("t3.micro"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          &key,
		UserData:         &userData,
		SecurityGroupIds: sgs,
		//SubnetId:         &subnet,
	})
	// snippet-end:[ec2.go.create_instance_with_tag.call]
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return result, nil
}

func main() {

	//------------------------------------
	// snippet-start:[ec2.go.create_instance_with_tag.args]
	name := flag.String("n", "Name", "The name of the tag to attach to the instance")
	value := flag.String("v", "test", "The value of the tag to attach to the instance")
	flag.Parse()

	if *name == "" || *value == "" {
		fmt.Println("You must supply a name and value for the tag (-n NAME -v VALUE)")
		return
	}
	// snippet-end:[ec2.go.create_instance_with_tag.args]

	// snippet-start:[ec2.go.create_instance_with_tag.session]
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)
	// snippet-end:[ec2.go.create_instance_with_tag.session]

	result, err := MakeInstance(svc, name, value)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Got an error creating an instance with tag " + *name)
		return
	}

	fmt.Println("Created tagged instance with ID " + *result.Instances[0].InstanceId)
}

// snippet-end:[ec2.go.create_instance_with_tag]
