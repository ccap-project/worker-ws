package terraform

import (
  "errors"
  "../config/"
  "../utils/"
  "./openstack"
)


func Init(provider string) func(*config.Config) error {

  if provider == "Openstack" {
    return openstack.Serialize
  }

  return nil
}

func Validate(SystemConfig *config.Config) error {

  cmd,_,stderr := utils.RunCmd(SystemConfig.Commands.Terraform, "validate")

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
  }

  return nil
}
