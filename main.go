package main

import (
  "fmt"
  "os"

  "./config/"
  "./ansible"
  "./terraform"
  "./webservice"
)

func main() {

  system := config.ReadFile("etc/system.conf")


  //ReadJson("example.json")

  exit(0)


}
