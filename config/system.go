/*
 *
 * Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

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
	config.Log.Formatter = &log.JSONFormatter{}

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
