package aws

import "fmt"
import "../../config/"

type hostgroup config.Hostgroup

func (h *hostgroup) Marshall() {
  fmt.Println("Here !")
}
