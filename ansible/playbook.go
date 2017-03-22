package ansible

import (
	"bytes"

	"../config"
	"../utils"
)

const play_tmpl = `{{range .}}{{if .Roles}}- hosts: {{.Name}}
{{if .BootstrapCommand}}  pre_tasks:
    - name: Bootstrap Ansible
      raw: {{.BootstrapCommand}}
      register: output
      changed_when: output.stdout != ""
{{- end}}

  roles:{{range .Roles}}
    - { role: '{{.Name}}', tags: [ '{{.Name}}' ]}
{{- end}}
{{end}}{{end}}
`

func playbook(config *config.Cell) (*bytes.Buffer, error) {

	var plays bytes.Buffer

	plays.Write([]byte("---\n"))

	p, err := utils.Template(play_tmpl, config.Hostgroups)
	if err != nil {
		return nil, err
	}

	plays.Write(p.Bytes())

	return &plays, nil
}
