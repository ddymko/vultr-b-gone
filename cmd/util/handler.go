package util

import (
	"context"
	"fmt"
	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
	"os"
)

type VultrBGone struct {
	Config  *govultr.Client
	Options *OptionsScheme
}

func NewVultrBGone() *VultrBGone {

	var token string
	if token == "" {
		token = os.Getenv("VULTR_API_KEY")
	}

	if token == "" {
		fmt.Println("Please export your VULTR API key as an environment variable or add `api-key` to your config file, eg:")
		fmt.Println("export VULTR_API_KEY='<api_key_from_vultr_account>'")
		os.Exit(1)
	}

	config := &oauth2.Config{}
	ts := config.TokenSource(context.Background(), &oauth2.Token{AccessToken: token})
	client := govultr.NewClient(oauth2.NewClient(context.Background(), ts))

	return &VultrBGone{
		Config: client,
	}
}
