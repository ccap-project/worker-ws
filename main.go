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

func main() {

  var err error

  SystemConfig := config.ReadJson("example.json")

  err = ansible.Serializer(SystemConfig)
  if err != nil {
    fmt.Println("Failure serializing Ansible Openstack file, ", err)
  }

  switch SystemConfig.Provider.Name {
    case "aws":
      aws.Serializer(SystemConfig)

    case "openstack":
      err = openstack.Serializer(SystemConfig)
      if err != nil {
        fmt.Println("Failure serializing Terraform Openstack file, ", err)
      }
  }
}
