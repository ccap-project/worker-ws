package ansible

import "fmt"
//import "strings"
//import "text/template"
//import "os"

import "../config/"
//import "../utils"

func hosts (config *config.Config) ([]string) {

  var hosts []string

  for _,host := range config.Hosts {

    host_line := host.Name

    for _,opt := range host.Options {
      for k,v := range opt {
        host_line = fmt.Sprintf("%s %s=%s", host_line, k, v)
      }
    }

    hosts = append(hosts, host_line)
  }
  hosts = append(hosts, "")

  return(hosts)
}

func hostgroups (config *config.Config) ([]string) {

  var hostgroups []string

  for _,hostgroup := range config.Hostgroups {
    hostgroups = append(hostgroups, fmt.Sprintf("[%s]", hostgroup.Name))
    hostgroups = append(hostgroups, fmt.Sprintf("%s[1:%s]", hostgroup.Name, hostgroup.Count))

    hostgroups = append(hostgroups, fmt.Sprintf("\n"))
  }

  return(hostgroups)
}

func group_vars (config *config.Config) ([]string) {

  var group_vars []string

  for _,hostgroup := range config.Hostgroups {
    group_vars = append(group_vars, fmt.Sprintf("[%s:vars]", hostgroup.Name))

    for _,vars := range hostgroup.Vars {
      for k,v := range vars {
        group_vars = append(group_vars, fmt.Sprintf("%s=%s", k, v))
      }
    }

    group_vars = append(group_vars, fmt.Sprintf("\n"))
  }

  return(group_vars)
}
