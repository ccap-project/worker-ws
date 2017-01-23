package webservice

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "../terraform"
  "../config"
)


func deploy(SystemConfig *config.SystemConfig) http.HandlerFunc {

  return func(w http.ResponseWriter, r *http.Request) {

    cell,err := config.DecodeJson(io.LimitReader(r.Body, SystemConfig.WebService.BodyLimit))
    if err != nil {
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(422) // unprocessable entity
      if err := json.NewEncoder(w).Encode(err); err != nil {
        panic(err)
      }
    }

    if err := terraform.Deploy(SystemConfig, cell); err != nil {
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(422) // unprocessable entity
      if err := json.NewEncoder(w).Encode(err); err != nil {
        panic(err)
      }
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(fmt.Sprintf("Ok")); err != nil {
      panic(err)
    }
  }
}

func checkInfrastructure(SystemConfig *config.SystemConfig) http.HandlerFunc {

  return func(w http.ResponseWriter, r *http.Request) {

    cell,err := config.DecodeJson(io.LimitReader(r.Body, SystemConfig.WebService.BodyLimit))
    if err != nil {
      SystemConfig.Log.Error("checkInfrastructure failed, ", err)
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(422) // unprocessable entity
      if err := json.NewEncoder(w).Encode(err); err != nil {
        panic(err)
      }
    }

    if err := terraform.Check(SystemConfig, cell); err != nil {
      SystemConfig.Log.Error("checkInfrastructure failed, ", err)
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(422) // unprocessable entity
      if err := json.NewEncoder(w).Encode(err); err != nil {
        panic(err)
      }
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(fmt.Sprintf("Ok")); err != nil {
      panic(err)
    }
  }
}

func deployInfrastructure(SystemConfig *config.SystemConfig) http.HandlerFunc {

  return func(w http.ResponseWriter, r *http.Request) {

    cell,err := config.DecodeJson(io.LimitReader(r.Body, SystemConfig.WebService.BodyLimit))
    if err != nil {
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(422) // unprocessable entity
      if err := json.NewEncoder(w).Encode(err); err != nil {
        panic(err)
      }
    }

    if err := terraform.Deploy(SystemConfig, cell); err != nil {
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(422) // unprocessable entity
      if err := json.NewEncoder(w).Encode(err); err != nil {
        panic(err)
      }
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(fmt.Sprintf("Ok")); err != nil {
      panic(err)
    }
  }
}
