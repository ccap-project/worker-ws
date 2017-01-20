package ansible

import (
  "errors"
  "../config/"
  "../utils/"
)

func RolesInstall(SystemConfig *config.Config) error {

  cmd,_,stderr := utils.RunCmd(SystemConfig.Commands.AnsibleGalaxy, "install", "-r", "requirements.yml")

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
  }

  return nil
}

func Run(SystemConfig *config.Config) error {

  cmd,_,stderr := utils.RunCmd(SystemConfig.Commands.Ansible, "-i", "hosts", "site.yml")

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
  }

  return nil
}

func SyntaxCheck(SystemConfig *config.Config) error {

  cmd,_,stderr := utils.RunCmd(SystemConfig.Commands.Ansible, "-i", "hosts", "--syntax-check", "site.yml")

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
  }

  return nil
}
