package config

import (
  "encoding/json"
  "log"
  "os"
)

type Commands struct {
  Terraform       string
  Ansible         string
  AnsibleGalaxy   string
}

type WebService struct {
  ApiPrefix     string
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
                                      TerraformState:        "terraform.tfstate"}}

  err = decoder.Decode(&config)

  if err != nil {
    log.Fatalf("Can't decode config file(%s), %s", configFilePath, err)
  }

  config.Log = log.New(os.Stderr, "roles-ws: ", log.Llongfile)

  return &config
}
