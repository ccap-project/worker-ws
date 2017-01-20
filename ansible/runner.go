package ansible

import (
  "errors"

  "../config"
  "../utils"
)

func RolesInstall(system *config.SystemConfig) error {

  cmd,_,stderr := utils.RunCmd(system.Commands.AnsibleGalaxy, "install", "-r", system.Files.AnsibleRequirements)

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
  }

  return nil
}

func Run(system *config.SystemConfig) error {

  cmd,_,stderr := utils.RunCmd(system.Commands.Ansible, "-i", "hosts", system.Files.AnsiblePlaybook)

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
  }

  return nil
}

func SyntaxCheck(system *config.SystemConfig) error {

  cmd,_,stderr := utils.RunCmd(system.Commands.Ansible, "-i", "hosts", "--syntax-check", system.Files.AnsiblePlaybook)

  if err_wait := cmd.Wait(); err_wait != nil {
    return errors.New(stderr.String())
  }

  return nil
}
