/*
 *
 * Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package webservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"worker-ws/ansible"
	"worker-ws/config"
	"worker-ws/repo"
	"worker-ws/terraform"
	"worker-ws/utils"
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

func buildContext(SystemConfig *config.SystemConfig) *config.RequestContext {

	ctx := new(config.RequestContext)

	ctx.SystemConfig = SystemConfig

	// XXX: check id gen failure
	runID, _ := utils.GetULID()

	ctx.RunID = runID.String()

	ctx.Log = SystemConfig.Log.WithFields(log.Fields{"rid": ctx.RunID})
	// XXX: Log RemoteIP

	return (ctx)
}

func uploadApplicationFile(SystemConfig *config.SystemConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//var err error

		ctx := buildContext(SystemConfig)

		vars := mux.Vars(r)

		ctx.Cell = new(config.Cell)

		ctx.Cell.Name = vars["cell"]
		ctx.Cell.CustomerName = vars["customer"]

		err := repo.Build(ctx)

		if err != nil {
			ctx.Log.Error("repo build failed, ", err)
			panic(err)
		}

		err = r.ParseMultipartForm(1 * 1024 * 1024)

		if err != nil {
			panic(nil)
		}

		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			panic(nil)
		}
		defer file.Close()

		ctx.Log.Debugf("uploadFile(%s)", handler.Filename)

		destPath := fmt.Sprintf("%s/files/%s/", ctx.Cell.Environment.Ansible.Dir, vars["role"])

		ctx.Log.Debugf("destPath(%s)", destPath)

		if err := os.MkdirAll(destPath, 0755); err != nil {
			panic(fmt.Errorf("Creating %s, %v", destPath, err))
		}

		f, err := os.OpenFile(destPath+vars["key"], os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(nil)
		}
		defer f.Close()

		_, err = io.Copy(f, file)

		if err != nil {
			panic(nil)
		}

		/*
		 * Format output
		 */
		var status run_status

		status.RunID = ctx.RunID

		status.StatusCode = 1

		err = repo.Persist(ctx, ctx.Cell.Environment.Ansible, mustTag(ctx))
		if err != nil {
			ctx.Log.Errorf("Commit error, %v", err)

			status.Stage.App.StatusCode = 1
			status.Stage.App.Message = fmt.Sprint(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}

	}
}

func makeHandler(SystemConfig *config.SystemConfig, fn func(*http.Request, *config.RequestContext, *stages)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var err error

		ctx := buildContext(SystemConfig)

		vars := mux.Vars(r)

		if len(vars) > 0 {
			ctx.TagID = vars["id"]
		}

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
			ctx.Log.Error("repo build failed, ", err)
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

	err := repo.Persist(ctx, ctx.Cell.Environment.Terraform, mustTag(ctx))
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

	err = repo.Persist(ctx, ctx.Cell.Environment.Terraform, true)
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

	err := repo.Persist(ctx, ctx.Cell.Environment.Ansible, mustTag(ctx))
	if err != nil {
		ctx.Log.Errorf("Commit error, %v", err)

		stages.App.StatusCode = 1
		stages.App.Message = fmt.Sprint(err)
		return
	}
}

func deployApplication(r *http.Request, ctx *config.RequestContext, stages *stages) {

	stages.App.StatusCode = 0

	/*
	 * Only run check when a new cfg need to be generated
	 */
	if !mustTag(ctx) {
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

		stages.App.Message = fmt.Sprint(err)

		err := repo.Persist(ctx, ctx.Cell.Environment.Ansible, mustTag(ctx))
		if err != nil {
			ctx.Log.Errorf("Persist Ansible repo, %v", err)
			stages.App.Message = fmt.Sprint(err)
		}

		stages.App.StatusCode = 1
		return
	}

	err := repo.Persist(ctx, ctx.Cell.Environment.Ansible, true)
	if err != nil {
		ctx.Log.Errorf("Persist Ansible repo, %v", err)

		stages.Infra.StatusCode = 1
		stages.Infra.Message = fmt.Sprint(err)
		return
	}
}

func mustTag(ctx *config.RequestContext) bool {

	if len(ctx.TagID) > 0 {
		return true
	}

	return false
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
