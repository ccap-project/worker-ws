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
	"fmt"
	"os"
	"strconv"

	"worker-ws/config"
	"worker-ws/utils"
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
