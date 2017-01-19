package ansible

import "bytes"
import "../config/"
import "../utils"

const require_tmpl = `- src: {{.Source}}
  version: {{.Version}}
  name: {{.Name}}

`

func requirements (config *config.Config) (*bytes.Buffer, error) {

  var requirements bytes.Buffer

  requirements.Write([]byte("---\n"))

  // XXX: Move this loop to template
  for _,hostgroup := range config.Hostgroups {

    if hostgroup.Roles != nil {
      for _,role := range hostgroup.Roles {
        req, err := utils.Template(require_tmpl, role)
        if err != nil {
          return nil,err
        }
        requirements.Write(req.Bytes())
      }
    }
  }

  return &requirements,nil
}
