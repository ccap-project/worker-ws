package "webservice"

import (
  "encoding/json"
  "net/http"
  "time"
  "github.com/gorilla/mux"

  "../config/"
)

func Start(SystemConfig *config.SystemConfig) {
  main_router := mux.NewRouter()
    api_router := main_router.PathPrefix(SystemConfig.WebService.ApiPrefix).Subrouter()

    // Deploy endpoints
    api_router.Methods("POST").Path("/deploy/infrastructure").HandlerFunc(terraform)

    // Role endpoints
    api_router.Methods("GET").Path("/roles").HandlerFunc(ListRoles)
    api_router.Methods("GET").Path("/roles/{role_id}/params/{version}").HandlerFunc(GetRoleParams)
    api_router.Methods("GET").Path("/roles/{role_id}/meta/{version}").HandlerFunc(GetRoleMeta)

    server := &http.Server{
            Handler:      main_router,
            Addr:         "0.0.0.0:8000",
            // Good practice: enforce timeouts for servers you create!
            WriteTimeout: 15 * time.Second,
            ReadTimeout:  15 * time.Second,
        }

    //role.Init(SystemConfig)

    SystemConfig.Log.Fatal(server.ListenAndServe())
}
