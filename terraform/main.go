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

package terraform

import (
	"fmt"

	"worker-ws/config"
	terraformcommon "worker-ws/terraform/common"

	"worker-ws/terraform/aws"
	"worker-ws/terraform/gcp"
	"worker-ws/terraform/openstack"
)

type Terraform interface {
	Apply(*config.SystemConfig, *config.Cell) (*[]byte, error)
	ReadState(*config.Cell, string) error
	Serialize(*config.SystemConfig, *config.Cell) error
	Validate(*config.SystemConfig, *config.Cell) error
}

func Init(provider string) Terraform {

	switch provider {
	case "Openstack":
		return &openstack.Openstack{}

	case "AWS":
		return &aws.AWS{}

	case "GCP":
		return &gcp.GCP{}
	}

	return nil
}

func Check(ctx *config.RequestContext) error {

	Env := Init(ctx.Cell.Provider.Type)

	if Env == nil {
		return fmt.Errorf("Terraform support for provider(%s) is not implemented ! \n", ctx.Cell.Provider.Type)
	}

	ctx.Log.Debug("Serializing")
	if err := Env.Serialize(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure serializing Terraform Openstack file, %v", err)
	}

	if _, err := terraformcommon.InstallProviders(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure installing Terraform Providers, %v", err)
	}

	ctx.Log.Debug("Validating")
	if err := Env.Validate(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure validating Terraform file, %v", err)
	}

	return nil
}

func Deploy(ctx *config.RequestContext) error {

	Env := Init(ctx.Cell.Provider.Type)

	if Env == nil {
		return fmt.Errorf("Terraform support for provider(%s) is not implemented ! \n", ctx.Cell.Provider.Type)
	}

	ctx.Log.Debug("Serializing")
	if err := Env.Serialize(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure serializing Terraform Openstack file, %v", err)
	}

	if _, err := terraformcommon.InstallProviders(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure installing Terraform Providers, %v", err)
	}

	ctx.Log.Debug("Validating")
	if err := Env.Validate(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure validating Terraform file, %v", err)
	}

	ctx.Log.Debug("Applying")
	_, err := Env.Apply(ctx.SystemConfig, ctx.Cell)

	if err != nil {
		return fmt.Errorf("Failure applying Terraform, %v", err)
	}

	ctx.Log.Debug("Reading state")

	if err := Env.ReadState(ctx.Cell, ctx.Cell.Environment.Terraform.Dir+ctx.SystemConfig.Files.TerraformState); err != nil {
		return fmt.Errorf("reading terraform state, %v", err)
	}

	return nil
}

func ReadState(ctx *config.RequestContext) error {

	Env := Init(ctx.Cell.Provider.Type)
	if Env == nil {
		return fmt.Errorf("Terraform support for provider(%s) is not implemented ! \n", ctx.Cell.Provider.Type)
	}

	ctx.Log.Debug("Reading state")

	if err := Env.ReadState(ctx.Cell, GetStateFilename(ctx)); err != nil {
		return fmt.Errorf("reading terraform state, %v", err)
	}

	return nil
}

func GetStateFilename(ctx *config.RequestContext) string {
	return fmt.Sprintf("%s%s", ctx.Cell.Environment.Terraform.Dir, ctx.SystemConfig.Files.TerraformState)
}
