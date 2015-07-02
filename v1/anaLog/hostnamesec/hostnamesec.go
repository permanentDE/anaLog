package hostnamesec

import (
	"errors"
	"net"
	"os"
	"strings"

	"go.permanent.de/anaLog/v1/config"
)

var ownHost string

func GetValidHost(remoteAddrPort string) (string, error) {
	if config.AnaLog.DevelopmentEnv {
		return "permanent.de", nil
	}

	var err error
	if ownHost == "" {
		ownHost, err = os.Hostname()
		if err != nil {
			return "", err
		}
		ownHost = removeSubdomains(ownHost)
	}

	var remoteAddr string

	remoteAddr, _, err = net.SplitHostPort(remoteAddrPort)
	if err != nil {
		return "", err
	}

	names, err := net.LookupAddr(remoteAddr)
	if err != nil {
		return "", err
	}

	for _, name := range names {
		if removeSubdomains(name) == ownHost {
			return name, nil
		}
	}

	return "", errors.New("invalid request")
}

func removeSubdomains(subdomain string) string {
	for {
		if strings.Count(subdomain, ".") == 1 || strings.Count(subdomain, ".") < 1 {
			return subdomain
		}
		subdomain = strings.SplitAfterN(subdomain, ".", 2)[1]
	}
}
