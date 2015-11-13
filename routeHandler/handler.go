package routeHandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/webapp"

	"go.permanent.de/anaLog/config"
	"go.permanent.de/anaLog/frontend"
	"go.permanent.de/anaLog/hostnamesec"
	"go.permanent.de/anaLog/nagios"
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

func NagiosStatus(w http.ResponseWriter, req *http.Request) *webapp.Error {
	fmt.Fprintln(w, nagios.Status())
	return nil
}

func NagiosReset(w http.ResponseWriter, req *http.Request) *webapp.Error {
	if nagios.Status() == nagios.OkStatus {
		fmt.Fprintln(w, "OK")
		return nil
	}

	as := req.FormValue("admin-secret")
	if as != "" {
		if as != config.AnaLog.AdminSecret {
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
					<input type="text" name="admin-secret">
					<input type="submit">
				</form>
			</body>
		</html>
	`)

	return nil
}

func FrontendFile(w http.ResponseWriter, req *http.Request) *webapp.Error {

	str, ok := frontend.File(mux.Vars(req)["file"])
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return nil
	}
	fmt.Fprint(w, str)
	return nil
}
