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
