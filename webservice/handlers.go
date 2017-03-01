package webservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

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

		if err != nil {
			ctx.Log.Error("makeHandler failed, ", err)
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
			ctx.Log.Error("makeHandler failed, ", err)
			panic(err)
		}

		if mustReadState(ctx) {

			err = terraform.ReadState(ctx)

			if err != nil {
				ctx.Log.Error("makeHandler failed, ", err)
				panic(err)
			}
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

func deploy(w http.ResponseWriter, r *http.Request, ctx *config.RequestContext) {

	deployInfrastructure(w, r, ctx)
	deployApplication(w, r, ctx)
}

func checkInfrastructure(w http.ResponseWriter, r *http.Request, ctx *config.RequestContext) {

	ctx.Log.Debugf("running checkInfrastructure")

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

	ctx.Log.Debugf("running deployInfrastructure")

	if err := terraform.Deploy(ctx); err != nil {
		ctx.Log.Errorf("deployInfrastructure failed, %v", err)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			ctx.Log.Errorf("deployInfrastructure failed, %v", err)
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(fmt.Sprintf("Ok")); err != nil {
		panic(err)
	}
}

func checkApplication(w http.ResponseWriter, r *http.Request, ctx *config.RequestContext) {

	ctx.Log.Debugf("running checkApplication")

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

	vars := mux.Vars(r)

	if vars == nil {
		if err := ansible.Check(ctx); err != nil {
			ctx.Log.Error("deployApplication failed, ", err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
	}

	ctx.Log.Debugf("running deployApplication")

	if err := ansible.Deploy(ctx); err != nil {
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

func mustReadState(ctx *config.RequestContext) bool {

	var inventory, tf os.FileInfo
	var err error

	tfStateFile := terraform.GetStateFilename(ctx)
	inventoryFile := ansible.GetInventoryFilename(ctx.SystemConfig, ctx.Cell)

	tf, err = os.Lstat(tfStateFile)
	if err != nil {
		ctx.Log.Debugf("tfStateFile(%s) does not exists", tfStateFile)
		return false
	}

	inventory, err = os.Lstat(inventoryFile)
	if err != nil {
		ctx.Log.Debugf("inventoryFile(%s) does not exists", inventoryFile)
		return true
	}
	inventoryMtime := inventory.ModTime()
	tfMtime := tf.ModTime()

	if tfMtime.After(inventoryMtime) {
		ctx.Log.Debugf("inventoryFile(%s) is oldter than tfStateFile(%s)", inventoryFile, tfStateFile)
		return true
	}

	return false
}
