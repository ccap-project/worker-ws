package openstack

import "bytes"
import "io/ioutil"

import "../../config/"
import "../../utils"

const provider_resource_tmpl = `provider "openstack" {
  user_name  = "{{.Username}}"
  tenant_name = "{{.Tenantname}}"
  password  = "{{.Password}}"
  auth_url  = "{{.AuthUrl}}"
}
`

func Serialize (config *config.Config) (error) {

  var tf bytes.Buffer

  provider, err := provider(config)
  if err != nil {
    return(err)
  }

  router, err := router(config)
  if err != nil {
    return(err)
  }

  router_interface, err := router_interface(config)
  if err != nil {
    return(err)
  }

  network, err := network(config)
  if err != nil {
    return(err)
  }

  subnet, err  := subnet(config)
  if err != nil {
    return(err)
  }

  instance, err  := instance(config)
  if err != nil {
    return(err)
  }

  tf.Write(provider.Bytes())
  tf.Write(router.Bytes())
  tf.Write(router_interface.Bytes())
  tf.Write(network.Bytes())
  tf.Write(subnet.Bytes())
  tf.Write(instance.Bytes())

  ioutil.WriteFile("site.tf", tf.Bytes(), 0644)

  return(nil)
}

func provider (config *config.Config) (*bytes.Buffer, error) {

  p, err := utils.Template(provider_resource_tmpl, config.Provider)
  if err != nil {
    return p,err
  }

  return p,nil
}
