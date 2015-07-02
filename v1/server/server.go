package server

import (
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"net"
	"net/http"
	"net/http/fcgi"

	"github.com/badgerodon/socketmaster/client"
	"github.com/badgerodon/socketmaster/protocol"
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

	logger.Stack2Http = config.AnaLog.DevelopmentEnv

	n := negroni.New(logger /*negroni.NewStatic(http.Dir(helper.GetFilePath("./public")))*/)

	cookiestore := cookiestore.New([]byte(config.AnaLog.CookieSecret))
	n.Use(sessions.Sessions("perm_analog_session", cookiestore))
	n.Use(negroni.HandlerFunc(preflight))

	n.UseHandler(router.New())

	if config.AnaLog.UseSocketMaster {
		listener, err := client.Listen(protocol.SocketDefinition{
			Port: config.SocketMaster.Port,
			HTTP: &protocol.SocketHTTPDefinition{
				DomainSuffix: config.SocketMaster.DomainSuffix,
				PathPrefix:   config.SocketMaster.PathPrefix,
			},
			/*TLS: &protocol.SocketTLSDefinition{
				Cert: config.SocketMaster.Cert,
				Key:  config.SocketMaster.Key,
			},*/
		})
		if err != nil {
			idl.Emerg(err)
		}
		idl.Notice("Serving via SocketMaster")
		http.Serve(listener, n)
	} else if config.AnaLog.Fcgi {
		listener, err := net.Listen("tcp", config.AnaLog.Listen)
		if err != nil {
			idl.Emerg(err)
		}
		idl.Notice("Serving via FastCGI")
		fcgi.Serve(listener, n)
	} else {
		idl.Notice("Serving via inbuilt HTTP Server")
		n.Run(config.AnaLog.Listen)
	}
}
