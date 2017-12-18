package ansible

import (
	"fmt"
	"os"
	"strconv"

	"../config"
	"../utils"
)

func RolesInstall(system *config.SystemConfig, cell *config.Cell) error {

	if needRolesCleanup(system, cell) {
		clearRoles(system, cell)
	}

	out, err := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.AnsibleGalaxy, "install", "-r", system.Files.AnsibleRequirements)

	if err != nil {
		return fmt.Errorf("%v, %s", err, out)
	}

	return nil
}

func Run(system *config.SystemConfig, cell *config.Cell) error {

	os.Remove(cell.Environment.Ansible.Dir + "/log")

	out, err := utils.RunCmd(cell.Environment.Ansible.Dir, cell.Environment.Ansible.Env, system.Commands.Ansible, system.Files.AnsiblePlaybook)

	if err != nil {

		c, _ := strconv.Unquote(fmt.Sprintf("%q", *out))

		fmt.Print(c)
		return fmt.Errorf("%v, %s", err, c)
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

func needRolesCleanup(system *config.SystemConfig, cell *config.Cell) bool {

	requirements, err := os.Stat(cell.Environment.Ansible.Dir + system.Files.AnsibleRequirements)

	roles, err := os.Stat(cell.Environment.Ansible.Dir + "/roles")

	if err != nil {
		return true
	}

	requirements_time := requirements.ModTime()

	if requirements_time.After(roles.ModTime()) {
		return true
	}

	return false
}

func clearRoles(system *config.SystemConfig, cell *config.Cell) error {

	err := os.RemoveAll(cell.Environment.Ansible.Dir + "/roles")

	return err
}
