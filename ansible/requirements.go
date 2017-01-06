package ansible

import "../config/"
import "../utils"

const require_tmpl = `
- src: {{.Source}}
  version: {{.Version}}
  name: {{.Name}}

`

func requirements (config *config.Config) ([]string) {

  var requirements []string

  for _,hostgroup := range config.Hostgroups {

    if hostgroup.Roles != nil {
      for _,role := range hostgroup.Roles {
        req := utils.Template(require_tmpl, role)
        requirements = append(requirements, req)
      }
    }
  }

  return(requirements)
}
