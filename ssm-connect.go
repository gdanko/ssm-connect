package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"sort"

	"github.com/bigkevmcd/go-configparser"
	"github.com/gdanko/ssm-connect/util"
	"github.com/jessevdk/go-flags"
	"github.com/nexidian/gocliselect"
)

const VERSION = "0.2.2"

type Options struct {
	Region  string `short:"r" long:"region" description:"Specify a region." choice:"us-east-1" choice:"us-east-2" choice:"us-west-1" choice:"us-west-2"`
	Version func() `short:"V" long:"version" description:"Display version information and exit."`
}

func main() {
	var (
		awsDir       string
		credsFile    string
		err          error
		flagsErr     *flags.Error
		homeDir      string
		instance     map[string]string
		instanceId   string
		instanceMenu *gocliselect.Menu
		instances    []map[string]string
		key          string
		ok           bool
		option       string
		opts         Options
		parser       *flags.Parser
		profile      string
		profileMenu  *gocliselect.Menu
		profiles     *configparser.ConfigParser
		secret       string
		section      string
		token        string
		username     *user.User
	)

	username, err = user.Current()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	awsDir = path.Join(homeDir, ".aws")
	credsFile = path.Join(awsDir, "credentials")

	if err = util.Verify(awsDir, credsFile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	opts.Version = func() {
		fmt.Printf("ssm-connect version %s\n", VERSION)
		os.Exit(0)
	}
	if opts.Region == "" {
		opts.Region = "us-west-2"
	}

	parser = flags.NewParser(&opts, flags.Default)
	if _, err = parser.Parse(); err != nil {
		if flagsErr, ok = err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	profiles, err = util.ParseCredentials(credsFile)
	if err != nil {
		fmt.Printf("Failed to parse the credentials file: %s\n", err.Error())
		os.Exit(0)
	}

	if len(profiles.Sections()) <= 0 {
		fmt.Printf("No profiles found. Exiting.\n")
		os.Exit(0)
	}

	profileMenu = gocliselect.NewMenu("Select an AWS profile (esc to cancel)")
	for _, section = range profiles.Sections() {
		profileMenu.AddItem(section, section)
	}
	profile = profileMenu.Display()
	if profile == "" {
		fmt.Printf("\nExiting\n")
		os.Exit(0)
	}

	key, secret, token, err = util.GetProfileCredentials(profiles, profile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	instances, err = util.GetInstances(profile, opts.Region, key, secret, token)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(instances) <= 0 {
		fmt.Printf("No instances found for profile \"%s\". Exiting.\n", profile)
		os.Exit(0)
	}

	sort.Slice(instances, func(i, j int) bool { return instances[i]["name"] < instances[j]["name"] })
	instanceMenu = gocliselect.NewMenu("Select an instance (esc to cancel)")

	for _, instance = range instances {
		option = fmt.Sprintf("%s - %-15s (%s)", instance["name"], instance["ip"], instance["id"])
		instanceMenu.AddItem(option, instance["id"])
	}

	instanceId = instanceMenu.Display()
	if instanceId == "" {
		fmt.Printf("\nExiting\n")
		os.Exit(0)
	}

	if err = util.StartSession(username.Username, key, secret, token, instanceId, profile, opts.Region, VERSION); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
