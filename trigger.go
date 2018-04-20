package slackrecv

import (
	"context"
	"log"
	"os"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/nlopes/slack"
)

var flogolog = logger.GetLogger("trigger-flogo-slackrecv")
var exitflag = 0

// SlackRecvTrigger is Slack RTM message trigger
type SlackRecvTrigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	handlers []*trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &SlackRecvFactory{metadata: md}
}

// SlackRecvFactory Timer Trigger factory
type SlackRecvFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *SlackRecvFactory) New(config *trigger.Config) trigger.Trigger {
	return &SlackRecvTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *SlackRecvTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Initialize implements trigger.Init
func (t *SlackRecvTrigger) Initialize(ctx trigger.InitContext) error {
	flogolog.Debugf("Initializing slack recv trigger...")
	t.handlers = ctx.GetHandlers()
	return nil
}

// Start implements ext.Trigger.Start
func (t *SlackRecvTrigger) Start() error {

	flogolog.Debugf("Starting slack RTM..")
	handlers := t.handlers
	exitflag = 0
	flogolog.Debug("Processing handlers")
	for _, handler := range handlers {

		accessToken := handler.GetStringSetting("AccessToken")
		//accessToken := t.config.GetSetting("AccessToken")
		flogolog.Debug("AccessToken: ", accessToken)
		api := slack.New(accessToken)
		logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
		slack.SetLogger(logger)
		api.SetDebug(true)

		rtm := api.NewRTM()
		go rtm.ManageConnection()

		for msg := range rtm.IncomingEvents {
			flogolog.Debugf("Event Received: ")
			
			if exitflag == 1
				return nil
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				t.RunHandler(handler, "Hello!")
				// Ignore hello
				//fmt.Println("Hello")

			case *slack.ConnectedEvent:
				fmt.Printf("Infos:", ev.Info)
				fmt.Printf("Connection counter:", ev.ConnectionCount)
				// Replace C2147483705 with your Channel ID
				//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "CA6BXNMPC"))

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				t.RunHandler(handler, ev.Text)

			case *slack.PresenceChangeEvent:
				//fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				//fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				return nil

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
		//log.Debugf("Processing Handler: %s", handler.ActionId)
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *SlackRecvTrigger) Stop() error {
	exitflag = 1
	flogolog.Debugf("Stopping RTM")

	return nil
}

// RunHandler action on new Slack RTM message
func (t *SlackRecvTrigger) RunHandler(handler *trigger.Handler, payload string) {

	trgData := make(map[string]interface{})
	trgData["message"] = payload

	_, err := handler.Handle(context.Background(), trgData)

	if err != nil {
		flogolog.Error("Error starting action: ", err.Error())
	}

	flogolog.Debugf("Ran Handler: [%s]", handler)
	
}
