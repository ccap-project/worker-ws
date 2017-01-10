package main

import (

  //"encoding/json"
  "fmt"
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

  var err error

  SystemConfig := config.ReadJson("example.json")

  switch SystemConfig.Provider.Name {
  case "aws":
    aws.Serializer(SystemConfig)
  case "openstack":
    err = ansible.Serializer(SystemConfig)
    if err != nil {
      fmt.Println("Failure serializing Ansible Openstack file, ", err)
    }

    err = openstack.Serializer(SystemConfig)
    if err != nil {
      fmt.Println("Failure serializing Terraform Openstack file, ", err)
    }
  }
}
