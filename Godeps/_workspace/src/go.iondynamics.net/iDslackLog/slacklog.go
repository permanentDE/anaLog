package iDslackLog

import (
	"github.com/johntdyer/slack-go"
	"go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"
)

type SlackLogHook struct {
	// Messages with a log Priority not contained in this array
	// will not be dispatched. If nil, all messages will be dispatched.
	AcceptedPriorities []priority.Priority
	HookURL            string
	IconURL            string
	Channel            string
	IconEmoji          string
	Username           string
	c                  *slack.Client
}

// Prioritys sets which Prioritys to sent to slack
func (sh *SlackLogHook) Priorities() []priority.Priority {
	if sh.AcceptedPriorities == nil {
		return priority.Threshold(priority.Debugging)
	}
	return sh.AcceptedPriorities
}

func (sh *SlackLogHook) Fire(e *iDlogger.Event) error {
	if sh.c == nil {
		if err := sh.initClient(); err != nil {
			return err
		}
	}

	color := ""
	switch e.Priority {
	case priority.Debugging:
		color = "#9B30FF"
	case priority.Informational:
		color = "good"
	case priority.Error, priority.Critical, priority.Alert, priority.Emergency:
		color = "danger"
	default:
		color = "warning"
	}

	msg := &slack.Message{
		Username: sh.Username,
		Channel:  sh.Channel,
	}

	msg.IconEmoji = sh.IconEmoji
	msg.IconUrl = sh.IconURL

	attach := msg.NewAttachment()

	// If there are fields we need to render them at attachments
	if len(e.Data) > 0 {

		// Add a header above field data
		attach.Text = "Message fields"

		for k, v := range e.Data {
			slackField := &slack.Field{}

			if str, ok := v.(string); ok {
				slackField.Title = k
				slackField.Value = str
				// If the field is <= 20 then we'll set it to short
				if len(str) <= 20 {
					slackField.Short = true
				}
			}
			attach.AddField(slackField)

		}
		attach.Pretext = e.Message
	} else {
		attach.Text = e.Message
	}
	attach.Fallback = e.Message
	attach.Color = color

	return sh.c.SendMessage(msg)
}

func (sh *SlackLogHook) initClient() error {
	sh.c = &slack.Client{sh.HookURL}

	if sh.Username == "" {
		sh.Username = "SlackLog"
	}

	return nil
}
