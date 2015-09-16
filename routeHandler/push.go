package routeHandler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.iondynamics.net/webapp"

	"go.permanent.de/anaLog/api"
)

func PushRecurringBegin(w http.ResponseWriter, req *http.Request) *webapp.Error {
	host, err := auth(w, req)
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil
	}

	vars := mux.Vars(req)
	task, ok := vars["task"]
	if !ok {
		return webapp.Write(errors.New("Invalid request: no task"), "no task given", http.StatusBadRequest)
	}

	req.Form.Del("override-secret")
	req.Form.Del("override-host")
	data := make(map[string]interface{})
	if len(req.Form) > 0 {
		data["urlValues"] = req.Form
	}

	identifier, err := api.PushRecurringBegin(task, host, data)
	if err != nil {
		return webapp.Write(err, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, identifier)
	return nil
}

func PushRecurringHeartbeat(w http.ResponseWriter, req *http.Request) *webapp.Error {
	host, err := auth(w, req)
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil
	}

	vars := mux.Vars(req)
	task, tOk := vars["task"]
	identifier, iOk := vars["identifier"]
	subtask, sOk := vars["subtask"]
	if !tOk || !iOk {
		return webapp.Write(errors.New("Invalid request"), "invalid request", http.StatusBadRequest)
	}
	if !sOk {
		subtask = "heartbeat"
	}

	err = api.PushRecurringHeartbeat(host, task, identifier, subtask)
	if err != nil {
		return webapp.Write(err, "Internal server error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func PushRecurringEnd(w http.ResponseWriter, req *http.Request) *webapp.Error {
	host, err := auth(w, req)
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil
	}

	vars := mux.Vars(req)
	task, tOk := vars["task"]
	identifier, iOk := vars["identifier"]
	state, sOk := vars["state"]
	if !tOk || !iOk || !sOk {
		return webapp.Write(errors.New("Invalid request"), "invalid request", http.StatusBadRequest)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	defer req.Body.Close()
	body := buf.String()
	if body == "{}" {
		body = ""
	}

	req.Form.Del("override-secret")
	req.Form.Del("override-host")
	data := make(map[string]interface{})
	if len(req.Form) > 0 {
		data["urlValues"] = req.Form
	}

	err = api.PushRecurringEnd(task, host, identifier, state, body, data)
	if err != nil {
		return webapp.Write(err, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}
