package main

import (

  //"encoding/json"
  //"fmt"
  //"net/http"
  //"time"
  //"github.com/gorilla/mux"

  //"./role/"
  "./config/"
  "./ansible"
  "./resource/aws"
  "./resource/openstack"

)

//var ( SystemConfig *config.Configuration )

func main() {

  //var ( SystemConfig *config.CellMap )
  //var ( SystemConfig *config.Config )

  SystemConfig := config.ReadJson("example.json")

  switch SystemConfig.Provider.Name {
  case "aws":
    aws.Serializer(SystemConfig)
  case "openstack":
    ansible.Serializer(SystemConfig)
    //openstack.Serializer(SystemConfig)
  }
}
