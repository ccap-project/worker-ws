package ansible

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"../config/"
)

func Serializer(config *config.SystemConfig, cell *config.Cell) error {

	var inventory bytes.Buffer

	if hosts := hosts(cell); hosts != nil {
		inventory.Write(hosts.Bytes())
	}

	if hostgroups := hostgroups(cell); hostgroups != nil {
		inventory.Write(hostgroups.Bytes())
	}

	if group_vars := group_vars(cell); group_vars != nil {
		inventory.Write(group_vars.Bytes())
	}

	ioutil.WriteFile(GetInventoryFilename(config, cell), inventory.Bytes(), 0644)

	playbook, err := playbook(cell)
	if err != nil {
		return (err)
	}
	ioutil.WriteFile(cell.Environment.Ansible.Dir+config.Files.AnsiblePlaybook, playbook.Bytes(), 0644)

	requirements, err := requirements(cell)
	if err != nil {
		return (err)
	}
	ioutil.WriteFile(cell.Environment.Ansible.Dir+config.Files.AnsibleRequirements, requirements.Bytes(), 0644)

	return (nil)
}

func Check(ctx *config.RequestContext) error {

	if err := Serializer(ctx.SystemConfig, ctx.Cell); err != nil {
		return fmt.Errorf("Failure serializing Ansible Openstack file, %v", err)
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
