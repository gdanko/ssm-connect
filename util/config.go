package util

import (
	"fmt"

	"github.com/bigkevmcd/go-configparser"
)

func ParseCredentials(credsFile string) (profiles *configparser.ConfigParser, err error) {
	profiles, err = configparser.NewConfigParserFromFile(credsFile)
	if err != nil {
		return profiles, fmt.Errorf("Failed to parse the credentials file: %s", err)
	}
	return profiles, nil
}

func GetProfileCredentials(profiles *configparser.ConfigParser, profile string) (key string, secret string, token string, err error) {
	key, err = profiles.Get(profile, "aws_access_key_id")
	if err != nil {
		return "", "", "", err
	}
	secret, err = profiles.Get(profile, "aws_secret_access_key")
	if err != nil {
		return "", "", "", err
	}
	token, err = profiles.Get(profile, "aws_session_token")
	if err != nil {
		token = ""
	}

	return key, secret, token, nil
}
