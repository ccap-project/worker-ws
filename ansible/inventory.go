package ansible

import "bytes"
import "fmt"
//import "strings"
//import "text/template"
//import "os"

import "../config/"
//import "../utils"

func hosts (config *config.Config) (*bytes.Buffer) {

  var hosts bytes.Buffer

  for _,host := range config.Hosts {

    fmt.Fprintf(&hosts, "%s", host.Name)

    for _,opt := range host.Options {
      for k,v := range opt {
        fmt.Fprintf(&hosts, " %s=%s", k,v)
        //host_line = fmt.Sprintf("%s %s=%s", host_line, k, v)
      }
    }
    fmt.Fprintf(&hosts, "\n")
  }
  fmt.Fprintf(&hosts, "\n")

  return(&hosts)
}

func hostgroups (config *config.Config) (*bytes.Buffer) {

  var hostgroups bytes.Buffer

  for _,hostgroup := range config.Hostgroups {
    fmt.Fprintf(&hostgroups, "[%s]\n", hostgroup.Name)
    fmt.Fprintf(&hostgroups, "%s[1:%s]\n\n", hostgroup.Name, hostgroup.Count)
  }

  return(&hostgroups)
}

func group_vars (config *config.Config) (*bytes.Buffer) {

  var group_vars bytes.Buffer

  for _,hostgroup := range config.Hostgroups {
    fmt.Fprintf(&group_vars, "[%s:vars]\n", hostgroup.Name)

    for _,vars := range hostgroup.Vars {
      for k,v := range vars {
        fmt.Fprintf(&group_vars, "%s=%s\n", k, v)
      }
    }
    fmt.Fprintf(&group_vars, "\n")
  }

  return(&group_vars)
}
