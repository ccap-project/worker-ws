package config

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
)

type GitlabCfg struct {
	Url                   string `json:"Url"`
	Token                 string
	TLSInsecureSkipVerify bool `json:",string"`
}

type Commands struct {
	Terraform     string
	Ansible       string
	AnsibleGalaxy string
}

type WebService struct {
	ApiPrefix string
	BodyLimit int64
	Address   string
	Port      int
}

type Files struct {
	AnsibleHosts        string
	AnsiblePlaybook     string
	AnsibleRequirements string
	TerraformSite       string
	TerraformState      string
	TempDir             string
}

type SystemConfig struct {
	Commands   Commands
	Files      Files
	Gitlab     GitlabCfg `json:"Gitlab"`
	Log        *log.Logger
	WebService WebService
}

func ReadFile(configFilePath string) *SystemConfig {

	config := SystemConfig{Files: Files{AnsibleHosts: "hosts",
		AnsiblePlaybook:     "site.yml",
		AnsibleRequirements: "requirements.yml",
		TerraformSite:       "site.tf",
		TerraformState:      "terraform.tfstate",
		TempDir:             "tmp/"},
		WebService: WebService{ApiPrefix: "/v1",
			BodyLimit: 1048576,
			Address:   "0.0.0.0",
			Port:      8080}}

	config.Log = log.New()

	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Can't open config file(%s), %s", configFilePath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)

	if err != nil {
		log.Fatalf("Can't decode config file(%s), %s", configFilePath, err)
	}

	config.Log.Level = log.DebugLevel

	return &config
}
