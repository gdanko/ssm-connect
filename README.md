# ssm-connect
ssm-connect is a small command line utility that attempts to ease the process of connecting to AWS instances via ssm.

## Installation
Clone the repository and from within the repository directory, type `make build`. This will create a directory with the given value of `GOOS` and install the binary there. It will also create a tarball which will eventually be used for Homebrew formulae.

## Features
* Gets the credentials from ~/.aws/credentials -- nothing to configure.
* Notifies you if you need to first install the session-manager-plugin.

## How does it work?
* ssm-connect parses the ~/.aws/credentials file and creates a list of profiles.
* The profile list is displayed on screen and the user is asked to select a profile.
* Once the profile is selected, ssm-connect reads the config to get the token information for the specified profile.
* Once the credentials are parsed, ssm-connect connects to the account and gathers a list of running instances.
* The list of running instances is used to build a second menu and the user is asked to select an instance to connect to.
* Once the instance has been selected, `session-manager-plugin` is invoked with the correct parameters and the shell is connected.

## ssm-connect runtime options
```
$ ssm-connect --help
Usage:
  ssm-connect [OPTIONS]

Application Options:
  -r, --region=  Specify a region.
  -V, --version  Display version information and exit.

Help Options:
  -h, --help     Show this help message
  ```

## CloudTrail Events
When the user uses ssm-connect to connect to an instance, the username and the ssm-connect version are logged to CloudTrail. This is for accountability. Here is a snippet from the event:
```
  "eventTime": "2024-06-26T23:16:06Z",
  "eventSource": "ssm.amazonaws.com",
  "eventName": "StartSession",
  "awsRegion": "us-west-2",
  "sourceIPAddress": "1.2.3.4",
  "userAgent": "aws-sdk-go/1.44.307 (go1.20.5; darwin; arm64)",
  "requestParameters": {
      "target": "i-04827860f68c6dbfa",
      "documentName": "SSM-SessionManagerRunShell",
      "reason": "gdanko connected via ssm-connect version 0.2.1 (https://github.com/gdanko/ssm-connect)"
  },
```
