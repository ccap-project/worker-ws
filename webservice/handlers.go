package webservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"../ansible"
	"../config"
	"../repo"
	"../terraform"
	"../utils"
)

func makeHandler(SystemConfig *config.SystemConfig, fn func(http.ResponseWriter, *http.Request, *config.RequestContext)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var err error

		ctx := new(config.RequestContext)

		ctx.SystemConfig = SystemConfig

		ctx.RunID, err = utils.GetULID()

		ctx.Log = SystemConfig.Log.WithFields(log.Fields{"rid": ctx.RunID})
		// XXX: Log RemoteIP

		ctx.Cell, err = config.DecodeJson(io.LimitReader(r.Body, SystemConfig.WebService.BodyLimit))

		ctx.Log.Debugf("Cell(%v)", ctx.Cell)

		if err != nil {
			ctx.Log.Error("checkInfrastructure failed, ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err = json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		ctx.Log = SystemConfig.Log.WithFields(log.Fields{"rid": ctx.RunID,
			"cell":     ctx.Cell.Name,
			"customer": ctx.Cell.CustomerName})

		err = repo.Build(ctx)

		if err != nil {
			ctx.Log.Error("checkInfrastructure failed, ", err)
			panic(err)
		}

		/*
			if m == nil {
				http.NotFound(w, r)
				return
			}
		*/

		fn(w, r, ctx)
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

func checkInfrastructure(w http.ResponseWriter, r *http.Request, ctx *config.RequestContext) {

	ctx.Log.Debugf("running checkInfrastructure CustomerName(%s) cell(%s)", ctx.Cell.CustomerName, ctx.Cell.Name)

	if err := terraform.Check(ctx); err != nil {
		ctx.Log.Error("checkInfrastructure failed, ", err)
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

func deployInfrastructure(w http.ResponseWriter, r *http.Request, ctx *config.RequestContext) {

	ctx.Log.Debugf("running deployInfrastructure CustomerName(%s) cell(%s)", ctx.Cell.CustomerName, ctx.Cell.Name)

	if err := terraform.Deploy(ctx); err != nil {
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

func checkApplication(w http.ResponseWriter, r *http.Request, ctx *config.RequestContext) {

	ctx.Log.Debugf("running checkApplication CustomerName(%s) cell(%s)", ctx.Cell.CustomerName, ctx.Cell.Name)

	if err := ansible.Check(ctx); err != nil {
		ctx.Log.Error("checkApplication failed, ", err)
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

func deployApplication(w http.ResponseWriter, r *http.Request, ctx *config.RequestContext) {

	ctx.Log.Debugf("running deployApplication CustomerName(%s) cell(%s)", ctx.Cell.CustomerName, ctx.Cell.Name)

	if err := terraform.Deploy(ctx); err != nil {
		ctx.Log.Error("deployApplication failed, ", err)
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
