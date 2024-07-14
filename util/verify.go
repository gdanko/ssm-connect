package util

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func Verify(awsDir, credsFile string) (err error) {
	if _, err := os.Stat(awsDir); os.IsNotExist(err) {
		return fmt.Errorf(awsDir + " does not exist.")
	}

	if _, err := os.Stat(credsFile); os.IsNotExist(err) {
		return fmt.Errorf(credsFile + " does not exist.")
	}

	if unix.Access(credsFile, unix.R_OK) != nil {
		return fmt.Errorf(credsFile + " exists but is not readable.")
	}

	_, err = exec.LookPath("session-manager-plugin")
	if err != nil {
		return fmt.Errorf("please install the aws session-manager-plugin and try again")
	}

	return nil
}
