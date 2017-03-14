package webservice

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"../config/"
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
	api_router.Methods("GET").Path("/application/deploy/{id:[A-Z0-9]+}").HandlerFunc(makeHandler(SystemConfig, deployApplication))

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
