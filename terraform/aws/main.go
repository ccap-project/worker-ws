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

package aws

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"worker-ws/config"
	terraformcommon "worker-ws/terraform/common"
	"worker-ws/utils"
)

type AWS struct{}

const provider_resource_tmpl = `
provider "aws" {
  version = "= 1.17.0"
  access_key  = "{{.AccessKey}}"
  secret_key  = "{{.SecretKey}}"
  region  = "{{.Region}}"
}

resource "aws_default_vpc" "default" {}
`

func (o *AWS) Apply(system *config.SystemConfig, cell *config.Cell) (*[]byte, error) {
	return terraformcommon.Apply(system, cell)
}

func (o *AWS) Validate(system *config.SystemConfig, cell *config.Cell) error {
	return terraformcommon.Validate(system, cell)
}

func (o *AWS) Serialize(system *config.SystemConfig, cell *config.Cell) error {

	var tf bytes.Buffer
	terraformSite := fmt.Sprintf("%s/%s", cell.Environment.Terraform.Dir, system.Files.TerraformSite)

	provider, err := provider(cell)
	if err != nil {
		return (err)
	}

	vpc, err := vpc(cell)
	if err != nil {
		return (err)
	}

	loadbalancer, err := loadbalancer(cell)
	if err != nil {
		return (err)
	}

	subnet, err := subnet(cell)
	if err != nil {
		return (err)
	}

	securitygroup, err := securitygroup(cell)
	if err != nil {
		return (err)
	}

	keypair, err := keypair(cell)
	if err != nil {
		return (err)
	}

	instance, err := instance(cell)
	if err != nil {
		return (err)
	}

	tf.Write(provider.Bytes())
	tf.Write(vpc.Bytes())
	tf.Write(loadbalancer.Bytes())
	tf.Write(subnet.Bytes())
	tf.Write(securitygroup.Bytes())
	tf.Write(keypair.Bytes())
	tf.Write(instance.Bytes())

	ioutil.WriteFile(terraformSite, tf.Bytes(), 0644)

	return (nil)
}

func provider(cell *config.Cell) (*bytes.Buffer, error) {

	p, err := utils.Template(provider_resource_tmpl, cell.Provider)
	if err != nil {
		return p, err
	}

	return p, nil
}
