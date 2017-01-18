package openstack

import (
  "fmt"
  "../common"
)

func (o *Openstack) ReadState(file string) error {

  state, err := terraformcommon.ReadState(file)

  fmt.Println(state)

  return err
}
