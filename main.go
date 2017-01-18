package main

import (

  "fmt"
  "os"

  "./config/"
  //"./ansible"
  "./terraform"
)

func main() {

  //var err error
  //terraform := terraform.Serializer

  SystemConfig  := config.ReadJson("example.json")
  SystemConfig.Commands.Terraform = "/Users/ale/Downloads/terraform"

  Terraform := terraform.Init(SystemConfig.Provider.Name)

  if Terraform == nil {
    fmt.Printf("Terraform support for provider(%s) is not implemented ! \n", SystemConfig.Provider.Name)
    os.Exit(-1)
  }

  if err := Terraform.Serialize(SystemConfig); err != nil {
    fmt.Println("Failure serializing Terraform Openstack file, ", err)
    os.Exit(-1)
  }

  if err := Terraform.Validate(SystemConfig); err != nil {
    fmt.Println("Failure validating Terraform file,", err)
    os.Exit(-1)
  }

  if err := Terraform.Apply(SystemConfig); err != nil {
    fmt.Println("Failure applying Terraform,", err)
    os.Exit(-1)
  }

  //if err := Terraform.ReadState(SystemConfig); err != nil {
  //  fmt.Println("Failure applying Terraform,", err)
  //  os.Exit(-1)
  //}


  os.Exit(0)

  //err = ansible.Serializer(SystemConfig)
  //if err != nil {
  //  fmt.Println("Failure serializing Ansible Openstack file, ", err)
  //}
}
