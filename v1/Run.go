package v1

import (
	"fmt"
	"os"
	"time"

	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"
	"go.iondynamics.net/iDslackLog"

	"go.permanent.de/anaLog/v1/anaLog"
	"go.permanent.de/anaLog/v1/anaLog/scheduler"
	_ "go.permanent.de/anaLog/v1/anaLog/scheduler/defaultTask"
	"go.permanent.de/anaLog/v1/config"
	"go.permanent.de/anaLog/v1/server"
)

func Run() {
	defer func() {
		anaLog.Close()
		idl.Panic("shutdown")
	}()

	if config.AppLog.SlackLogUrl != "" {

		prio := priority.Warning
		if config.AnaLog.DevelopmentEnv {
			prio = priority.Debugging
		}

		idl.AddHook(&iDslackLog.SlackLogHook{
			AcceptedPriorities: priority.Threshold(prio),
			HookURL:            config.AppLog.SlackLogUrl,
			IconURL:            "",
			Channel:            "",
			IconEmoji:          "",
			Username:           "anaLog",
		})
	}

	idl.StandardLogger().Async = true
	idl.SetPrefix("anaLog")
	idl.SetErrCallback(func(err error) {
		fmt.Fprintln(os.Stderr, err)
		panic("AppLogger caught an internal error")
	})

	if config.AnaLog.DevelopmentEnv {
		go scheduler.StartIn(1 * time.Second)
	} else {
		go scheduler.StartIn(10 * time.Second)
	}

	server.Listen()
}
