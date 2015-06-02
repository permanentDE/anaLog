package router

import (
	"github.com/gorilla/mux"
	//idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/webapp"
	handler "go.permanent.de/anaLog/v1/routeHandler"
	"net/http"
)

func New() *mux.Router {
	return provision(mux.NewRouter().StrictSlash(true))
}

func provision(r *mux.Router) *mux.Router {
	r.NotFoundHandler = webapp.Handler(handler.NotFound)

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/v1/", http.StatusFound)
	})

	v1 := r.PathPrefix("/v1").Subrouter()

	v1Push := v1.PathPrefix("/push").Subrouter()
	v1Push.Handle("/recurring/{task}", webapp.Handler(handler.PushRecurringBegin)).Methods("POST")
	v1Push.Handle("/recurring/{task}/{identifier}/{state}", webapp.Handler(handler.PushRecurringEnd)).Methods("PUT")

	v1Push.Handle("/recurring/analyze", webapp.Handler(handler.AnalyzeRecurring)).Methods("GET")

	return r
}
