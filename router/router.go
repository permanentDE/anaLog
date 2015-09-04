package router

import (
	"github.com/gorilla/mux"
	//idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/webapp"
	handler "go.permanent.de/anaLog/routeHandler"
	"net/http"
)

func New() *mux.Router {
	return provision(mux.NewRouter().StrictSlash(true))
}

func provision(r *mux.Router) *mux.Router {
	r.NotFoundHandler = webapp.Handler(handler.NotFound)

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/", http.StatusFound)
	})

	v1 := r.PathPrefix("/v1").Subrouter()

	v1Push := v1.PathPrefix("/push").Subrouter()
	v1Push.Handle("/recurring/{task}", webapp.Handler(handler.PushRecurringBegin)).Methods("POST")
	v1Push.Handle("/recurring/{task}/{identifier}/heartbeat", webapp.Handler(handler.PushRecurringHeartbeat)).Methods("GET")
	v1Push.Handle("/recurring/{task}/{identifier}/heartbeat/{subtask}", webapp.Handler(handler.PushRecurringHeartbeat)).Methods("GET")
	v1Push.Handle("/recurring/{task}/{identifier}/{state}", webapp.Handler(handler.PushRecurringEnd)).Methods("PUT")

	v1Nagios := v1.PathPrefix("/nagios").Subrouter()
	v1Nagios.Handle("/status", webapp.Handler(handler.NagiosStatus)).Methods("GET")
	v1Nagios.Handle("/reset", webapp.Handler(handler.NagiosReset)).Methods("GET", "POST")

	return r
}
