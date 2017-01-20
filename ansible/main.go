package ansible

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "../config/"
)



func Serializer(config *config.SystemConfig, cell *config.Cell) (error) {

  var inventory bytes.Buffer

  hosts         := hosts(cell)
  hostgroups    := hostgroups(cell)
  group_vars    := group_vars(cell)

  inventory.Write([]byte("---\n"))
  inventory.Write(hosts.Bytes())
  inventory.Write(hostgroups.Bytes())
  inventory.Write(group_vars.Bytes())

  ioutil.WriteFile(config.Files.AnsibleHosts, inventory.Bytes(), 0644)

  playbook, err := playbook(cell)
  if err != nil {
    return(err)
  }
  ioutil.WriteFile(config.Files.AnsiblePlaybook, playbook.Bytes(), 0644)

  requirements, err  := requirements(cell)
  if err != nil {
    return(err)
  }
  ioutil.WriteFile(config.Files.AnsibleRequirements, requirements.Bytes(), 0644)

  return(nil)
}

func Deploy(system *config.SystemConfig, cell *config.Cell) (error) {

  if err := Serializer(system, cell); err != nil {
    return fmt.Errorf("Failure serializing Ansible Openstack file, %v", err)
  }

  if err := RolesInstall(system); err != nil {
    return fmt.Errorf("Failure downloading Ansible galaxy roles, %v", err)
  }

  if err := SyntaxCheck(system); err != nil {
    return fmt.Errorf("Failure checking Ansible file, %v", err)
  }

  if err := Run(system); err != nil {
    return fmt.Errorf("Failure running Ansible, %v", err)
  }

  return nil
}
