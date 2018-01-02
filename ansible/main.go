/*
 *
 * Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

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
