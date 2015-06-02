package server

import (
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"net"
	"net/http"
	"net/http/fcgi"

	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDnegroniLog"

	"go.permanent.de/anaLog/v1/config"
	"go.permanent.de/anaLog/v1/router"
)

func preflight(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

func Listen() {
	logger := iDnegroniLog.NewMiddleware(idl.StandardLogger())

	logger.Stack2Http = config.Std.AnaLog.DevelopmentEnv

	n := negroni.New(logger /*negroni.NewStatic(http.Dir(helper.GetFilePath("./public")))*/)

	cookiestore := cookiestore.New([]byte(config.Std.AnaLog.CookieSecret))
	n.Use(sessions.Sessions("perm_analog_session", cookiestore))
	n.Use(negroni.HandlerFunc(preflight))

	n.UseHandler(router.New())

	if config.Std.AnaLog.Fcgi {
		listener, err := net.Listen("tcp", config.Std.AnaLog.Listen)
		if err != nil {
			idl.Emerg(err)
		}
		fcgi.Serve(listener, n)
	} else {
		n.Run(config.Std.AnaLog.Listen)
	}
}
