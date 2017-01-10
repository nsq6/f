package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

func createConfig(filename *string) (*oauth2.Config, error) {
	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		pwd, _ := os.Getwd()
		fullPath := filepath.Join(pwd, *filename)
		return nil, fmt.Errorf("Cannot read config file from %v", fullPath)
	}
	config := new(Config)
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("Cannot read config file from %v", filename)
	}

	return &oauth2.Config{
		Scopes:       config.Web.Scopes,
		ClientID:     config.Web.ClientID,
		ClientSecret: config.Web.ClientSecret,
		RedirectURL:  config.Web.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Web.AuthURI,
			TokenURL: config.Web.TokenURI,
		},
	}, nil
}

// Config struct
type Config struct {
	Web    ConfigSection `json:"web"`
	Stored ConfigSection `json:"stored"`
}

// ConfigSection struct
type ConfigSection struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectURL             string   `json:"redirect_url"`
	Scopes                  []string `json:"scopes"`
}
