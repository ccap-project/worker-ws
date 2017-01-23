package config

import (
  "encoding/json"
  log "github.com/Sirupsen/logrus"
  "os"
)

type Commands struct {
  Terraform       string
  Ansible         string
  AnsibleGalaxy   string
}

type WebService struct {
  ApiPrefix     string
  BodyLimit     int64
  Address       string
  Port          int
}

type Files struct {
  AnsibleHosts        string
  AnsiblePlaybook     string
  AnsibleRequirements string
  TerraformSite       string
  TerraformState      string
}

type SystemConfig struct {
  Commands      Commands
  Files         Files
  Log           *log.Logger
  WebService    WebService
}

func ReadFile(configFilePath string) *SystemConfig {

  file, err := os.Open(configFilePath)
  if err != nil {
    log.Fatalf("Can't open config file(%s), %s", configFilePath, err)
  }
  defer file.Close()

  decoder := json.NewDecoder(file)
  config := SystemConfig{Files: Files{AnsibleHosts:         "hosts",
                                      AnsiblePlaybook:      "site.yml",
                                      AnsibleRequirements:  "requirements.yml",
                                      TerraformSite:        "site.tf",
                                      TerraformState:        "terraform.tfstate"},
                        WebService: WebService{ApiPrefix:   "/v1",
                                                BodyLimit:   1048576,
                                                Address:    "0.0.0.0",
                                                Port:       8080}}

  err = decoder.Decode(&config)

  if err != nil {
    log.Fatalf("Can't decode config file(%s), %s", configFilePath, err)
  }

  return &config
}
