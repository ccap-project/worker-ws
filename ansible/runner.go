package ansible

import (
	"fmt"

	"../config"
	"../utils"
)

func RolesInstall(system *config.SystemConfig, cell *config.Cell) error {

	out, err := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.AnsibleGalaxy, "install", "-r", system.Files.AnsibleRequirements)

	if err != nil {
		return fmt.Errorf("%v, %s", err, out)
	}

	return nil
}

func Run(system *config.SystemConfig, cell *config.Cell) error {

	out, err := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.Ansible, system.Files.AnsiblePlaybook)

	if err != nil {
		return fmt.Errorf("%v, %s", err, out)
	}

	return nil
}

func SyntaxCheck(system *config.SystemConfig, cell *config.Cell) error {

	out, err := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.Ansible, "--syntax-check", system.Files.AnsiblePlaybook)

	if err != nil {
		return fmt.Errorf("%v, %s", err, out)
	}

	return nil
}
