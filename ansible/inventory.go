package ansible

import (
	"bytes"
	"fmt"

	"../config"
	"../utils"
)

const files_tmpl = `{{range .}}{{if .Files}}{{.Name}}_files={ {{range .Files}}'{{.Key}}': { 'filename': '{{.Filename}}'{{if .DontCopy}}, 'dont_copy': 'true'{{end}}  }, {{end}} }
{{end}}{{end}}
`

const hostgroupTmpl = `{{range .}}[{{.Name}}]
{{.Name}}[1:{{.Count}}]
{{end}}
`

const hostgroupRolesTmpl = `{{range .}}[{{.Name}}:vars]
{{if .Username}}ansible_ssh_user={{.Username}}
{{- end -}}
{{range .Roles}}
# {{.Name -}}{{range .Params}}
{{.Name}}={{.Value}}
{{- end}}
{{- end}}
{{end}}
`

func hosts(config *config.Cell) *bytes.Buffer {

	var hosts bytes.Buffer

	// XXX: Move this loop to template
	for _, host := range config.Hosts {

		fmt.Fprintf(&hosts, "%s", host.Name)

		for _, opt := range host.Options {
			for k, v := range opt {
				fmt.Fprintf(&hosts, " %s=%s", k, v)
			}
		}
		fmt.Fprintf(&hosts, "\n")
	}
	fmt.Fprintf(&hosts, "\n")

	return (&hosts)
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
