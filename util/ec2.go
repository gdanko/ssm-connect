package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetInstances(profile, region, key, secret, token string) (instances []map[string]string, err error) {
	var (
		name string = ""
		tag  types.Tag
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return instances, err
	}

	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(key, secret, token)
	})

	params := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
				},
			},
		},
	}
	result, err := client.DescribeInstances(context.TODO(), params)
	if err != nil {
		if strings.Contains(err.Error(), "AWS was not able to validate the provided access credentials") {
			return instances, fmt.Errorf("invalid credentials - please check the credentials for this account and try again")
		}
		return instances, err
	}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			id := *instance.InstanceId
			ip := *instance.PrivateIpAddress
			for _, tag = range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					instances = append(instances, map[string]string{
						"name": name,
						"id":   id,
						"ip":   ip,
					})
				}
			}
		}
	}
	return instances, nil
}
