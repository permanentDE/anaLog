package hostnamesec

import (
	"errors"
	"net"
	"os"
	"strings"

	"go.permanent.de/anaLog/config"
)

var ownDomain string

func GetValidHost(remoteAddrPort string) (string, error) {
	if config.AnaLog.DevelopmentEnv {
		return "permanent.de", nil
	}

	var err error
	if ownDomain == "" {
		if config.AnaLog.Domain == "" {
			ownDomain, err = os.Hostname()
			if err != nil {
				return "", err
			}
			ownDomain = removeSubdomains(ownDomain)
		} else {
			ownDomain = config.AnaLog.Domain
		}
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
		remoteDomain := removeSubdomains(name)
		if remoteDomain == ownDomain {
			return name, nil
		}
	}

	return "", errors.New("invalid request")
}

func removeSubdomains(subdomain string) string {
	subdomain = strings.TrimSuffix(subdomain, ".")
	for {
		if strings.Count(subdomain, ".") == 1 || strings.Count(subdomain, ".") < 1 {
			return subdomain
		}
		subdomain = strings.SplitAfterN(subdomain, ".", 2)[1]
	}
}
