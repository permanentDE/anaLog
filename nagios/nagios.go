package nagios

const OkStatus string = "NAGIOS_OK"

var currentStatus = OkStatus

func Status() string {
	return currentStatus
}

func SetOK() {
	currentStatus = OkStatus
}

func SetFailed(msg string) {
	currentStatus = msg
}
