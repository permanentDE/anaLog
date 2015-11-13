package server

import (
	"net"
	"net/http"
	"net/http/fcgi"

	"github.com/badgerodon/socketmaster/client"
	"github.com/badgerodon/socketmaster/protocol"
	"github.com/codegangsta/negroni"
	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDnegroniLog"

	"go.permanent.de/anaLog/config"
	"go.permanent.de/anaLog/router"
)

func preflight(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

func Listen() {
	logger := iDnegroniLog.NewMiddleware(idl.StandardLogger())

	logger.Stack2Http = config.AnaLog.DevelopmentEnv

	n := negroni.New(logger)

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
