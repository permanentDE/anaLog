package routeHandler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.iondynamics.net/webapp"

	"go.permanent.de/anaLog/api"
)

func ReadFind(w http.ResponseWriter, req *http.Request) *webapp.Error {
	nStr, ok := mux.Vars(req)["number"]
	if !ok {
		nStr = "1"
	}
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, "invalid number: n", http.StatusBadRequest)
		return nil
	}

	task := req.FormValue("task")
	host := req.FormValue("host")
	state := req.FormValue("state")
	rawRegex := req.FormValue("rawRegex")

	trGTEstr := req.FormValue("timeRangeGTE")
	if trGTEstr == "" {
		trGTEstr = "0"
	}
	trGTEn, err := strconv.Atoi(trGTEstr)
	if err != nil {

		http.Error(w, "invalid number: timeRangeGTE", http.StatusBadRequest)
		return nil
	}

	trLTEstr := req.FormValue("timeRangeLTE")
	if trLTEstr == "" {
		trLTEstr = "0"
	}
	trLTEn, err := strconv.Atoi(trLTEstr)
	if err != nil {
		http.Error(w, "invalid number: timeRangeLTE", http.StatusBadRequest)
		return nil
	}

	trGTE := time.Time{}
	trLTE := time.Time{}

	if trGTEn > 0 {
		trGTE = time.Unix(int64(trGTEn), 0)
	}
	if trLTEn > 0 {
		trLTE = time.Unix(int64(trLTEn), 0)
	}

	lps, err := api.Find(task, host, state, rawRegex, trGTE, trLTE, uint(n))
	if err != nil {
		return webapp.Write(err, err.Error(), http.StatusInternalServerError)
	}

	return writeJson(w, lps)
}

func writeJson(w http.ResponseWriter, i interface{}) *webapp.Error {
	jw := json.NewEncoder(w)
	err := jw.Encode(i)
	if err != nil {
		return webapp.Write(err, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
