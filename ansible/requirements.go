package ansible

import "bytes"
import "../config/"
import "../utils"

const require_tmpl = `- src: {{.Source}}
  version: {{.Version}}
  name: {{.Name}}

`

func requirements (config *config.Config) (*bytes.Buffer) {

  var requirements bytes.Buffer

  requirements.Write([]byte("---\n"))

  for _,hostgroup := range config.Hostgroups {

    if hostgroup.Roles != nil {
      for _,role := range hostgroup.Roles {
        req := utils.Template(require_tmpl, role)
        requirements.Write(req.Bytes())
      }
    }
  }

  return(&requirements)
}
