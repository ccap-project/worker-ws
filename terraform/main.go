package terraform

import (
	"fmt"

	"../config/"
	"./openstack"
)

type Terraform interface {
	Apply(*config.SystemConfig, *config.Cell) (*[]byte, error)
	ReadState(*config.Cell, string) error
	Serialize(*config.SystemConfig, *config.Cell) error
	Validate(*config.SystemConfig, *config.Cell) error
}

func Init(provider string) Terraform {

	if provider == "Openstack" {
		return &openstack.Openstack{}
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
