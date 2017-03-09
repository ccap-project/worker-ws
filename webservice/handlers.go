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

type stages struct {
	Infra status `json:"infrastructure"`
	App   status `json:"application"`
}

type status struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"msg"`
}

type run_status struct {
	RunID      string `json:"run_id"`
	StatusCode int    `json:"status_code"`
	Stage      stages `json:"stage"`
}

func makeHandler(SystemConfig *config.SystemConfig, fn func(*http.Request, *config.RequestContext, *stages)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := new(config.RequestContext)

		ctx.SystemConfig = SystemConfig

		// XXX: check id gen failure
		run_id, err := utils.GetULID()

		ctx.RunID = run_id.String()

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

		/*
		 * Format output
		 */
		var status run_status

		status.RunID = ctx.RunID

		fn(r, ctx, &status.Stage)

		status.StatusCode = status.Stage.Infra.StatusCode | status.Stage.App.StatusCode

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
	}
}

func deploy(r *http.Request, ctx *config.RequestContext, stages *stages) {

	deployInfrastructure(r, ctx, stages)
	deployApplication(r, ctx, stages)
}

func checkInfrastructure(r *http.Request, ctx *config.RequestContext, stages *stages) {

	stages.Infra.StatusCode = 0

	ctx.Log.Debugf("running checkInfrastructure")

	if err := terraform.Check(ctx); err != nil {
		ctx.Log.Error("checkInfrastructure failed, ", err)
		stages.Infra.StatusCode = 1
		stages.Infra.Message = fmt.Sprint(err)
		return
	}

	err := repo.Persist(ctx, ctx.Cell.Environment.Terraform)
	if err != nil {
		ctx.Log.Errorf("Persist Terraform repo, %v", err)

		stages.Infra.StatusCode = 1
		stages.Infra.Message = fmt.Sprint(err)
		return
	}
}

func deployInfrastructure(r *http.Request, ctx *config.RequestContext, stages *stages) {

	stages.Infra.StatusCode = 0

	ctx.Log.Debugf("running deployInfrastructure")
	err := terraform.Deploy(ctx)
	if err != nil {
		ctx.Log.Errorf("deployInfrastructure failed, %v", err)

		stages.Infra.StatusCode = 1
		stages.Infra.Message = fmt.Sprint(err)
		return
	}

	err = repo.Persist(ctx, ctx.Cell.Environment.Terraform)
	if err != nil {
		ctx.Log.Errorf("Persist Terraform repo, %v", err)

		stages.Infra.StatusCode = 1
		stages.Infra.Message = fmt.Sprint(err)
		return
	}
}

func checkApplication(r *http.Request, ctx *config.RequestContext, stages *stages) {

	stages.App.StatusCode = 0

	ctx.Log.Debugf("running checkApplication")

	if err := ansible.Check(ctx); err != nil {
		ctx.Log.Error("checkApplication failed, ", err)

		stages.App.StatusCode = 1
		stages.App.Message = fmt.Sprint(err)
		return
	}

	ctx.Log.Infof("Commit Repo(%s)", ctx.Cell.Environment.Ansible.Name)
	if err := repo.Persist(ctx, ctx.Cell.Environment.Ansible); err != nil {
		ctx.Log.Errorf("Commit error, %v", err)

		stages.App.StatusCode = 1
		stages.App.Message = fmt.Sprint(err)
		return
	}
}

func deployApplication(r *http.Request, ctx *config.RequestContext, stages *stages) {

	vars := mux.Vars(r)

	stages.App.StatusCode = 0

	if vars == nil {
		if err := ansible.Check(ctx); err != nil {
			ctx.Log.Error("deployApplication failed, ", err)

			stages.App.StatusCode = 1
			stages.App.Message = fmt.Sprint(err)
			return
		}
	}

	ctx.Log.Debugf("running deployApplication")

	if err := ansible.Deploy(ctx); err != nil {
		ctx.Log.Error("deployApplication failed, ", err)

		stages.App.StatusCode = 1
		stages.App.Message = fmt.Sprint(err)
		return
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
