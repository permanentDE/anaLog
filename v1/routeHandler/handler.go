package routeHandler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/webapp"

	"go.permanent.de/anaLog/v1/anaLog"
	"go.permanent.de/anaLog/v1/anaLog/hostnamesec"
	"go.permanent.de/anaLog/v1/anaLog/nagios"
	"go.permanent.de/anaLog/v1/config"
)

func NotFound(w http.ResponseWriter, req *http.Request) *webapp.Error {
	idl.Debug(req)
	return webapp.Write(errors.New("404 - Not found"), "404 - This isn't the page you're looking for", http.StatusNotFound)
}

func auth(w http.ResponseWriter, req *http.Request) (host string, err error) {
	secret := req.FormValue("override-secret")
	if secret == "" {
		host, err = hostnamesec.GetValidHost(req.RemoteAddr)
		return
	} else if secret != config.AnaLog.OverrideSecret {
		err = errors.New("Invalid Secret")
		return
	}

	host = req.FormValue("override-host")

	return

}

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

	identifier, err := anaLog.PushRecurringBegin(task, host)
	if err != nil {
		return webapp.Write(err, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, identifier)
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

	err = anaLog.PushRecurringEnd(task, host, identifier, state, body)
	if err != nil {
		return webapp.Write(err, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

func NagiosStatus(w http.ResponseWriter, req *http.Request) *webapp.Error {
	fmt.Fprintln(w, nagios.Status())
	return nil
}

func NagiosReset(w http.ResponseWriter, req *http.Request) *webapp.Error {
	if nagios.Status() == nagios.OkStatus {
		fmt.Fprintln(w, "OK")
		return nil
	}

	ns := req.FormValue("nagios-secret")
	if ns != "" {
		if ns != config.AnaLog.NagiosSecret {
			http.Error(w, "Invalid secret", http.StatusUnauthorized)
			return nil
		} else {
			nagios.SetOK()
			fmt.Fprintln(w, "OK")
			return nil
		}
	}

	fmt.Fprintln(w, `
		<!DOCTYPE html>
		<html>
			<body>
				<form method="POST">
					<input type="text" name="nagios-secret">
					<input type="submit">
				</form>
			</body>
		</html>
	`)

	return nil
}
