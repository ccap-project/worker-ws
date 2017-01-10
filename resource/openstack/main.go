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

func Serializer (config *config.Config) {

  var tf bytes.Buffer

  provider          := provider(config)
  router            := router(config)
  router_interface  := router_interface(config)
  network           := network(config)
  subnet            := subnet(config)
  instance          := instance(config)

  tf.Write(provider.Bytes())
  tf.Write(router.Bytes())
  tf.Write(router_interface.Bytes())
  tf.Write(network.Bytes())
  tf.Write(subnet.Bytes())
  tf.Write(instance.Bytes())

  ioutil.WriteFile("site.tf", tf.Bytes(), 0644)
}

func provider (config *config.Config) (*bytes.Buffer) {

  return(utils.Template(provider_resource_tmpl, config.Provider))
}
