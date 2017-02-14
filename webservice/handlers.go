package webservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"../config"
	"../repo"
	"../terraform"
)

func makeHandler(SystemConfig *config.SystemConfig, fn func(http.ResponseWriter, *http.Request, *config.SystemConfig, *config.Cell)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		cell, err := config.DecodeJson(io.LimitReader(r.Body, SystemConfig.WebService.BodyLimit))

		SystemConfig.Log.Debugf("Cell(%v)", cell)

		if err != nil {
			SystemConfig.Log.Error("checkInfrastructure failed, ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		err = repo.Build(SystemConfig, cell)

		if err != nil {
			SystemConfig.Log.Error("checkInfrastructure failed, ", err)
			panic(err)
		}

		/*
			if m == nil {
				http.NotFound(w, r)
				return
			}
		*/

		fn(w, r, SystemConfig, cell)
	}
}

/*
func deploy(SystemConfig *config.SystemConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		cell, err := config.DecodeJson(io.LimitReader(r.Body, SystemConfig.WebService.BodyLimit))
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
*/

func checkInfrastructure(w http.ResponseWriter, r *http.Request, SystemConfig *config.SystemConfig, cell *config.Cell) {

	SystemConfig.Log.Debugf("running checkInfrastructure CustomerName(%s) cell(%s)", cell.CustomerName, cell.Name)

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

func deployInfrastructure(w http.ResponseWriter, r *http.Request, SystemConfig *config.SystemConfig, cell *config.Cell) {

	SystemConfig.Log.Debugf("running deployInfrastructure CustomerName(%s) cell(%s)", cell.CustomerName, cell.Name)

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

func checkConfiguration(w http.ResponseWriter, r *http.Request, SystemConfig *config.SystemConfig, cell *config.Cell) {

	SystemConfig.Log.Debugf("running checkConfiguration CustomerName(%s) cell(%s)", cell.CustomerName, cell.Name)

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

func deployConfiguration(w http.ResponseWriter, r *http.Request, SystemConfig *config.SystemConfig, cell *config.Cell) {

	SystemConfig.Log.Debugf("running checkInfrastructure CustomerName(%s) cell(%s)", cell.CustomerName, cell.Name)

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
