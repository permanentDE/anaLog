package v1

import (
	"time"

	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"
	"go.iondynamics.net/iDslackLog"

	"go.permanent.de/anaLog/v1/config"
	"go.permanent.de/anaLog/v1/server"
)

func Run() {
	if config.Std.AppLog.SlackLogUrl != "" {

		prio := priority.Warning
		if config.Std.AnaLog.DevelopmentEnv {
			prio = priority.Debugging
		}

		idl.AddHook(&iDslackLog.SlackLogHook{
			AcceptedPriorities: priority.Threshold(prio),
			HookURL:            config.Std.AppLog.SlackLogUrl,
			IconURL:            "",
			Channel:            "",
			IconEmoji:          "",
			Username:           "anaLog",
		})
	}
	idl.StandardLogger().Async = true
	idl.SetPrefix("anaLog")
	idl.SetErrCallback(func(err error) {
		idl.StandardLogger().Async = true
		idl.Log(&idl.Event{
			idl.StandardLogger(),
			map[string]interface{}{
				"error": err,
			},
			time.Now(),
			priority.Emergency,
			"AppLogger caught an internal error",
		})
		panic("AppLogger caught an internal error")
	})

	server.Listen()
}
