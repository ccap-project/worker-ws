package ansible

import (
  "bytes"
  "../config/"
  "../utils"
)


const play_tmpl = `{{range .}}- hosts: {{.Name}}
  roles:
    {{range .Roles}}- { role: '{{.Name}}', tags: [ '{{.Name}}' ]}{{end}}

{{end}}
`


func playbook(config *config.Config) (*bytes.Buffer, error) {

  var plays bytes.Buffer

  plays.Write([]byte("---\n"))

  p, err := utils.Template(play_tmpl, config.Hostgroups)
  if err != nil {
    return nil, err
  }

  plays.Write(p.Bytes())

  return &plays, nil
}
