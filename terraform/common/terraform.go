package terraformcommon

import (
  "errors"
  "../../config/"
  "../../utils/"
)

func Apply(SystemConfig *config.Config) error {

  cmd,_,stderr := utils.RunCmd(SystemConfig.Commands.Terraform, "apply")

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
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
