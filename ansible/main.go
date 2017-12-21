package ansible

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"worker-ws/config"
)

func Serializer(config *config.SystemConfig, cell *config.Cell) error {

	var inventory bytes.Buffer

	hosts, err := hosts(cell)

	if err != nil {
		return fmt.Errorf("inventory hosts, %v", err)
	}
	inventory.Write(hosts.Bytes())

	hostgroups, err := hostgroups(cell)
	if err != nil {
		return fmt.Errorf("inventory hostgroups, %v", err)
	}

	inventory.Write(hostgroups.Bytes())

	group_vars, err := group_vars(cell)
	if err != nil {
		return fmt.Errorf("inventory groupvars, %v", err)
	}

	inventory.Write(group_vars.Bytes())

	ioutil.WriteFile(GetInventoryFilename(config, cell), inventory.Bytes(), 0644)

	playbook, err := playbook(cell)
	if err != nil {
		return fmt.Errorf("playbook, %v", err)
	}
	ioutil.WriteFile(cell.Environment.Ansible.Dir+config.Files.AnsiblePlaybook, playbook.Bytes(), 0644)

	requirements, err := requirements(cell)
	if err != nil {
		return fmt.Errorf("requirements, %v", err)
	}
	ioutil.WriteFile(cell.Environment.Ansible.Dir+config.Files.AnsibleRequirements, requirements.Bytes(), 0644)

	configuration, err := configuration(cell)
	if err != nil {
		return fmt.Errorf("configuration, %v", err)
	}
	ioutil.WriteFile(cell.Environment.Ansible.Dir+"ansible.cfg", configuration.Bytes(), 0644)

	return (nil)
}

func Check(ctx *config.RequestContext) error {

	if err := Serializer(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure serializing Ansible %v", err)
	}

	if err := RolesInstall(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure downloading Ansible galaxy roles, %v", err)
	}

	if err := SyntaxCheck(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure checking Ansible file, %v", err)
	}

	return nil
}

func Deploy(ctx *config.RequestContext) error {

	if err := Run(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure running Ansible, %v", err)
	}

	return nil
}
