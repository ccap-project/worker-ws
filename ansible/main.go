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
}
