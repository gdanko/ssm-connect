package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Target struct {
	Target string `json:"Target"`
}

func StartSession(username string, key, secret, token, instanceId, profile, region, version string) error {
	t := Target{
		Target: instanceId,
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return err
	}

	// Why is this not in the client?
	endpoint := fmt.Sprintf("https://ssm.%s.amazonaws.com", region)

	client := ssm.NewFromConfig(cfg, func(o *ssm.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(key, secret, token)
		o.APIOptions = append(o.APIOptions, middleware.AddUserAgentKeyValue("ssm-connect", version))
	})
	resp, err := client.StartSession(
		context.TODO(),
		&ssm.StartSessionInput{
			Target:       &instanceId,
			DocumentName: aws.String("SSM-SessionManagerRunShell"),
			Reason:       aws.String(fmt.Sprintf("%s connected via ssm-connect version %s (https://github.com/gdanko/ssm-connect)", username, region)),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start the session: %s", err.Error())
	}

	target, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to create the target JSON: %s", err.Error())
	}

	response, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to parse the response from AWS: %s", err.Error())
	}

	cmd := exec.Command("session-manager-plugin", string(response), region, "StartSession", profile, string(target), endpoint)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// ignore signal(sigint)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-sigs:
			case <-done:
				break
			}
		}
	}()
	defer close(done)

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
