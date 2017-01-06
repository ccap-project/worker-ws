package ansible

import "fmt"
import "../config/"

func Serializer (config *config.Config) {

  var inventory bytes.Buffer

  hosts := hosts(config)...)

  fmt.Printf("%s\n", inventory.String())

/*
  output = append(output, hosts(config)...)
  output = append(output, hostgroups(config)...)
  output = append(output, group_vars(config)...)
  output = append(output, playbook(config)...)
  output = append(output, requirements(config)...)


  for _,v := range output {
    fmt.Printf("%s\n", v)
  }
  */
}
