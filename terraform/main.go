package terraform

import (
  "../config/"
  "./openstack"
)

type Terraform interface {
  Apply(*config.Config)     error
  Serialize(*config.Config) error
  Validate(*config.Config)  error
}

func Init(provider string) Terraform {

  if provider == "Openstack" {
    return &openstack.Openstack{}
  }

  return nil
}
