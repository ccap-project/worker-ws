package terraform

import (
	"fmt"

	"../config/"
	"./openstack"
)

type Terraform interface {
	Apply(*config.SystemConfig, string) error
	ReadState(*config.Cell, string) error
	Serialize(*config.SystemConfig, *config.Cell, string) error
	Validate(*config.SystemConfig, string) error
}

func Init(provider string) Terraform {

	if provider == "Openstack" {
		return &openstack.Openstack{}
	}

	return nil
}

func Check(SystemConfig *config.SystemConfig, cell *config.Cell, repo string) error {

	Env := Init(cell.Provider.Name)

	if Env == nil {
		return fmt.Errorf("Terraform support for provider(%s) is not implemented ! \n", cell.Provider.Name)
	}

	SystemConfig.Log.Debug("Serializing")
	if err := Env.Serialize(SystemConfig, cell, repo); err != nil {
		return fmt.Errorf("Failure serializing Terraform Openstack file, %v", err)
	}

	SystemConfig.Log.Debug("Validating")
	if err := Env.Validate(SystemConfig, repo); err != nil {
		return fmt.Errorf("Failure validating Terraform file, %v", err)
	}

	return nil
}

func Deploy(SystemConfig *config.SystemConfig, cell *config.Cell, repo string) error {

	Env := Init(cell.Provider.Name)

	if Env == nil {
		return fmt.Errorf("Terraform support for provider(%s) is not implemented ! \n", cell.Provider.Name)
	}

	SystemConfig.Log.Debug("Serializing")
	if err := Env.Serialize(SystemConfig, cell, repo); err != nil {
		return fmt.Errorf("Failure serializing Terraform Openstack file, %v", err)
	}

	SystemConfig.Log.Debug("Validating")
	if err := Env.Validate(SystemConfig, repo); err != nil {
		return fmt.Errorf("Failure validating Terraform file, %v", err)
	}

	SystemConfig.Log.Debug("Applying")
	if err := Env.Apply(SystemConfig, repo); err != nil {
		return fmt.Errorf("Failure applying Terraform, %v", err)
	}

	SystemConfig.Log.Debug("Reading state")
	if err := Env.ReadState(cell, SystemConfig.Files.TerraformState); err != nil {
		return fmt.Errorf("Failure reading Terraform state, %v", err)
	}

	return nil
}
