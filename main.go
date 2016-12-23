package main

import (

  //"encoding/json"
  //"fmt"
  //"net/http"
  //"time"
  //"github.com/gorilla/mux"

  //"./role/"
  "./config/"
  "./resource/aws"

)

//var ( SystemConfig *config.Configuration )

func main() {

  //var ( SystemConfig *config.CellMap )

  SystemConfig := config.ReadJson("example.json")

  /*
  for k, v := range SystemConfig {
      switch vv := v.(type) {
      case string:
          fmt.Println(k, "is string", vv)
      case int:
          fmt.Println(k, "is int", vv)
      case []interface{}:
          fmt.Println(k, "is an array:")
          for i, u := range vv {
              fmt.Println(i, u)
          }
      default:
          fmt.Println(k, "is of a type I don't know how to handle")
      }
  }
  */
  switch SystemConfig.Provider.Name {
  case "aws":
    aws.

  }
  SystemConfig.Hostgroup.Marshall()
  //fmt.Println(SystemConfig)
}
