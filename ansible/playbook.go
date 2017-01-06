package ansible

import "fmt"
//import "strings"
//import "text/template"
//import "os"

import "../config/"
//import "../utils"

func playbook (config *config.Config) ([]string) {

  var plays []string

  for _,hostgroup := range config.Hostgroups {

    if hostgroup.Roles != nil {
      plays = append(plays, fmt.Sprintf("- hosts: %s", hostgroup.Name))
      plays = append(plays, "  roles:")

      for _,role := range hostgroup.Roles {
        plays = append(plays, fmt.Sprintf("    - %s", role.Name))
      }

      plays = append(plays, fmt.Sprintf("\n"))
    }
  }

  return(plays)
}
