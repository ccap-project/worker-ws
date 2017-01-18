package main

import (

  "fmt"
  "os"
  "time"

  "./config/"
  //"./ansible"
  "./terraform"
)

func main() {

  //var err error
  //terraform := terraform.Serializer

  SystemConfig  := config.ReadJson("example.json")
  SystemConfig.Commands.Terraform = "/Users/ale/Downloads/terraform"

  TerraformSerializer     := terraform.Init(SystemConfig.Provider.Name)

  if TerraformSerializer == nil {
    fmt.Printf("Terraform serializer for provider(%s) is not supported ! \n", SystemConfig.Provider.Name)
    os.Exit(-1)
  }

  if err := TerraformSerializer(SystemConfig); err != nil {
    fmt.Println("Failure serializing Terraform Openstack file, ", err)
    os.Exit(-1)
  }

  if err := terraform.Validate(SystemConfig); err != nil {
    fmt.Println("Failure validating Terraform file,", err)
    os.Exit(-1)
  }
  os.Exit(0)

  //err = ansible.Serializer(SystemConfig)
  //if err != nil {
  //  fmt.Println("Failure serializing Ansible Openstack file, ", err)
  //}
}
