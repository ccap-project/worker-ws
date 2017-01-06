package ansible

import "fmt"
import "../config/"

func Serializer (config *config.Config) {

  output := make([]string, 1)

  output = append(output, hosts(config)...)
  output = append(output, hostgroups(config)...)
  output = append(output, group_vars(config)...)
  output = append(output, playbook(config)...)
  output = append(output, requirements(config)...)

  /*
  output = append(output, router(config)...)
  output = append(output, router_interface(config)...)
  output = append(output, network(config)...)
  output = append(output, subnet(config)...)
  output = append(output, instance(config)...)
  */

  for _,v := range output {
    fmt.Printf("%s\n", v)
  }
}
