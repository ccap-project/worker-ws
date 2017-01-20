package openstack

import (
  "bytes"
  "io/ioutil"
  "../../config/"
  "../common"
  "../../utils"
)

type Openstack struct {}

const provider_resource_tmpl = `provider "openstack" {
  user_name  = "{{.Username}}"
  tenant_name = "{{.Tenantname}}"
  password  = "{{.Password}}"
  auth_url  = "{{.AuthUrl}}"
}
`
func (o *Openstack) Apply(system *config.SystemConfig) (error) {
  return terraformcommon.Apply(system)
}

func (o *Openstack) Validate(system *config.SystemConfig) (error) {
  return terraformcommon.Validate(system)
}

func (o *Openstack) Serialize(system *config.SystemConfig, cell *config.Cell) (error) {

  var tf bytes.Buffer

  provider, err := provider(cell)
  if err != nil {
    return(err)
  }

  router, err := router(cell)
  if err != nil {
    return(err)
  }

  router_interface, err := router_interface(cell)
  if err != nil {
    return(err)
  }

  network, err := network(cell)
  if err != nil {
    return(err)
  }

  subnet, err  := subnet(cell)
  if err != nil {
    return(err)
  }

  instance, err  := instance(cell)
  if err != nil {
    return(err)
  }

  tf.Write(provider.Bytes())
  tf.Write(router.Bytes())
  tf.Write(router_interface.Bytes())
  tf.Write(network.Bytes())
  tf.Write(subnet.Bytes())
  tf.Write(instance.Bytes())

  ioutil.WriteFile(system.Files.TerraformSite, tf.Bytes(), 0644)

  return(nil)
}


func provider (cell *config.Cell) (*bytes.Buffer, error) {

  p, err := utils.Template(provider_resource_tmpl, cell.Provider)
  if err != nil {
    return p,err
  }

  return p,nil
}
