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
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"worker-ws/config"
)

func Start(SystemConfig *config.SystemConfig) {
	main_router := mux.NewRouter()
	api_router := main_router.PathPrefix(SystemConfig.WebService.ApiPrefix).Subrouter()

	// Deploy endpoints
	api_router.Methods("POST").Path("/deploy").HandlerFunc(makeHandler(SystemConfig, deploy))

	api_router.Methods("POST").Path("/infrastructure/check").HandlerFunc(makeHandler(SystemConfig, checkInfrastructure))
	api_router.Methods("POST").Path("/infrastructure/deploy").HandlerFunc(makeHandler(SystemConfig, deployInfrastructure))
	api_router.Methods("GET").Path("/infrastructure/deploy/{id:[A-Z0-9]+}").HandlerFunc(makeHandler(SystemConfig, deployInfrastructure))

	api_router.Methods("POST").Path("/application/check").HandlerFunc(makeHandler(SystemConfig, checkApplication))
	api_router.Methods("POST").Path("/application/deploy").HandlerFunc(makeHandler(SystemConfig, deployApplication))
	api_router.Methods("GET").Path("/application/deploy/{id:[a-zA-Z0-9]+}").HandlerFunc(makeHandler(SystemConfig, deployApplication))
	api_router.Methods("POST").Path("/application/{customer:[a-zA-Z0-9._]+}/{cell:[a-zA-Z0-9._]+}/{role:[a-zA-Z0-9._]+}/file/{key:[a-zA-Z0-9._]+}").HandlerFunc(uploadApplicationFile(SystemConfig))

	server := &http.Server{
		Handler: handlers.CombinedLoggingHandler(os.Stdout, main_router),
		Addr:    fmt.Sprintf("%s:%d", SystemConfig.WebService.Address, SystemConfig.WebService.Port),
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	SystemConfig.Log.Debugf("Starting socket on %s:%d", SystemConfig.WebService.Address, SystemConfig.WebService.Port)

	SystemConfig.Log.Fatal(server.ListenAndServe())
}
