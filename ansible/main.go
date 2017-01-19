package ansible

import "bytes"
import "io/ioutil"
import "../config/"

func Serializer (config *config.Config) (error) {

  var inventory bytes.Buffer

  hosts         := hosts(config)
  hostgroups    := hostgroups(config)
  group_vars    := group_vars(config)

  inventory.Write([]byte("---\n"))
  inventory.Write(hosts.Bytes())
  inventory.Write(hostgroups.Bytes())
  inventory.Write(group_vars.Bytes())

  ioutil.WriteFile("hosts", inventory.Bytes(), 0644)

  playbook, err := playbook(config)
  if err != nil {
    return(err)
  }
  ioutil.WriteFile("site.yml", playbook.Bytes(), 0644)

  requirements, err  := requirements(config)
  if err != nil {
    return(err)
  }
  ioutil.WriteFile("requirements.yml", requirements.Bytes(), 0644)

  return(nil)

  //fmt.Printf("==========\n%s\n=============\n", inventory.String())
  //fmt.Printf("%s\n", hosts.String())
  //fmt.Printf("%s\n", hostgroups.String())
  //fmt.Printf("%s\n", group_vars.String())
  //fmt.Printf("%s\n", playbook.String())
  //fmt.Printf("%s\n", requirements.String())


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
