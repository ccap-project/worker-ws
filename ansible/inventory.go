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

	"worker-ws/config"
	"worker-ws/utils"
)

const files_tmpl = `{{range .}}{{if .Files}}{{.Name}}_files={ {{range .Files}}'{{.Key}}': { 'filename': '{{.Filename}}'{{if .DontCopy}}, 'dont_copy': 'true'{{end}}  }, {{end}} }
{{end}}{{end}}
`

const hostTmpl = `{{range .}}{{.Name}}{{if .Options}}{{range .Options}} {{.Name}}={{.Value}}{{end}}
{{end}}{{end}}

`

const hostgroupTmpl = `{{range .}}[{{.Name}}]
{{.Name}}[1:{{.Count}}]
{{end}}
`

const hostgroupRolesTmpl = `{{range .}}[{{.Name}}:vars]
{{if .Username}}ansible_ssh_user={{.Username}}
ansible_become=true
{{- end -}}
{{range .Roles}}
# {{.Name -}}{{range .Params}}
{{.Name}}={{.Value}}
{{- end}}
{{- end}}
{{end}}
`

func hosts(config *config.Cell) (*bytes.Buffer, error) {

	var hosts bytes.Buffer

	p, err := utils.Template(hostTmpl, config.Hosts)
	if err != nil {
		return nil, err
	}

	hosts.Write(p.Bytes())

	return &hosts, nil
}

func hostgroups(config *config.Cell) (*bytes.Buffer, error) {

	var hostgroups bytes.Buffer

	p, err := utils.Template(hostgroupTmpl, config.Hostgroups)
	if err != nil {
		return nil, err
	}

	hostgroups.Write(p.Bytes())

	return &hostgroups, nil
}

func group_vars(config *config.Cell) (*bytes.Buffer, error) {

	var group_vars bytes.Buffer

	p, err := utils.Template(hostgroupRolesTmpl, config.Hostgroups)
	if err != nil {
		fmt.Printf("\n=====================%s==================\n", hostgroupRolesTmpl)
		return nil, err
	}

	group_vars.Write(p.Bytes())

	return &group_vars, nil
}

func GetInventoryFilename(config *config.SystemConfig, cell *config.Cell) string {
	return fmt.Sprintf("%s%s", cell.Environment.Ansible.Dir, config.Files.AnsibleHosts)
}
