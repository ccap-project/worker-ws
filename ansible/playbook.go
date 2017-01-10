package ansible

import "bytes"
import "fmt"

import "../config/"

func playbook (config *config.Config) (*bytes.Buffer) {

  var plays bytes.Buffer

  plays.Write([]byte("---\n"))

  for _,hostgroup := range config.Hostgroups {

    if hostgroup.Roles != nil {
      fmt.Fprintf(&plays, "- hosts: %s\n", hostgroup.Name)
      fmt.Fprintf(&plays, "  roles:\n")

      for _,role := range hostgroup.Roles {
        fmt.Fprintf(&plays, "    - %s\n", role.Name)
      }

      fmt.Fprintf(&plays, "\n")
    }
  }

  return(&plays)
}
