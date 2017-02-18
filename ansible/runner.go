package ansible

import (
	"errors"

	"../config"
	"../utils"
)

func RolesInstall(system *config.SystemConfig, cell *config.Cell) error {

	cmd, _, stderr := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.AnsibleGalaxy, "install", "-r", system.Files.AnsibleRequirements)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func Run(system *config.SystemConfig, cell *config.Cell) error {

	cmd, _, stderr := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.Ansible, system.Files.AnsiblePlaybook)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func SyntaxCheck(system *config.SystemConfig, cell *config.Cell) error {

	cmd, _, stderr := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.Ansible, "--syntax-check", system.Files.AnsiblePlaybook)

	if err_wait := cmd.Wait(); err_wait != nil {
		return errors.New(stderr.String())
	}

	return nil
}
