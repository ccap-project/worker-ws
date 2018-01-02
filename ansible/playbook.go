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
	"sort"
	"strings"

	"worker-ws/config"
	"worker-ws/utils"
)

type BootstrapGroup struct {
	BootstrapCommand string
	Hosts            []string
}

func (bg *BootstrapGroup) GetHosts() string {
	return strings.Join(bg.Hosts, ", ")
}

func groupBootstrapCommand(H []*config.Hostgroup) []*BootstrapGroup {

	idx := -1
	var lastBootstrapCommand string
	var bootstrapCommandList []*BootstrapGroup

	sort.Sort(config.HostgroupByName{H})

	for _, hg := range H {

		var current *BootstrapGroup

		if len(hg.BootstrapCommand) <= 0 {
			continue
		}

		if strings.Compare(lastBootstrapCommand, hg.BootstrapCommand) != 0 {
			lastBootstrapCommand = hg.BootstrapCommand
			idx++
			current = new(BootstrapGroup)
			bootstrapCommandList = append(bootstrapCommandList, current)
			bootstrapCommandList[idx].BootstrapCommand = lastBootstrapCommand
		}

		bootstrapCommandList[idx].Hosts = append(bootstrapCommandList[idx].Hosts, hg.Name)
	}

	return bootstrapCommandList
}

const bootstrap_tmpl = `
{{range .}}- hosts: "{{.GetHosts}}"
  gather_facts: False
  pre_tasks:
    - name: Bootstrap Ansible
      raw: {{.BootstrapCommand}}
      register: output
      changed_when: output.stdout != ""
    - name: Gathering Facts
      setup:
{{end}}
`

const play_tmpl = `{{range .}}{{if .Roles}}- hosts: {{.Name}}
{{$Component := .Component}}  roles:{{range .Roles}}
    - { role: '{{.Name}}', tags: [ '{{.Name}}'{{if $Component}}, '{{$Component}}'{{end}} ]}

{{end}}{{end}}{{end}}`

func playbook(config *config.Cell) (*bytes.Buffer, error) {

	var plays bytes.Buffer

	plays.Write([]byte("---\n"))

	bootstrapCommands := groupBootstrapCommand(config.Hostgroups)

	if bootstrapCommands != nil {
		p, err := utils.Template(bootstrap_tmpl, bootstrapCommands)
		if err != nil {
			return nil, err
		}
		plays.Write(p.Bytes())
	}

	p, err := utils.Template(play_tmpl, config.Hostgroups)
	if err != nil {
		return nil, err
	}
	plays.Write(p.Bytes())

	return &plays, nil
}
