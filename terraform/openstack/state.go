package openstack

import (
  "fmt"
  "../common"
)

func (o *Openstack) ReadState(file string) error {

  state, err := terraformcommon.ReadState(file)

  for m := range state.Modules {
    for r,rv := range state.Modules[m].Resources {

      switch rv.Type {
        case "openstack_compute_instance_v2":
            fmt.Printf("(%s) (%s) (%s) (%s)\n", r, rv.Type, rv.Primary.Attributes["name"], rv.Primary.Attributes["access_ip_v4"])

        default:

      }
    }
  }

  return err
}
