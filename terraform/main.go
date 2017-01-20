package terraform

import (
  "fmt"
  "../config/"
  "./openstack"
)

type Terraform interface {
  Apply(*config.SystemConfig)                   error
  ReadState(*config.Cell, string)          error
  Serialize(*config.SystemConfig, *config.Cell) error
  Validate(*config.SystemConfig)                error
}

func Init(provider string) Terraform {

  if provider == "Openstack" {
    return &openstack.Openstack{}
  }

  return nil
}

func Deploy(config *config.SystemConfig, cell *config.Cell) (error) {

  Env := Init(cell.Provider.Name)

  if Env == nil {
    return fmt.Errorf("Terraform support for provider(%s) is not implemented ! \n", cell.Provider.Name)
  }

  if err := Env.Serialize(config, cell); err != nil {
    return fmt.Errorf("Failure serializing Terraform Openstack file, %v", err)
  }

  if err := Env.Validate(config); err != nil {
    return fmt.Errorf("Failure validating Terraform file, %v", err)
  }

  if err := Env.Apply(config); err != nil {
    return fmt.Errorf("Failure applying Terraform, %v", err)
  }

  if err := Env.ReadState(cell, config.Files.TerraformState); err != nil {
    return fmt.Errorf("Failure reading Terraform state, %v", err)
  }

  return nil
}
