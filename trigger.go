package slackrecv

import (
	"context"
	"log"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/nlopes/slack"
)

var flogolog = logger.GetLogger("trigger-flogo-slackrecv")

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
	t.handlers = ctx.GetHandlers()
	return nil
}

// Start implements ext.Trigger.Start
func (t *SlackRecvTrigger) Start() error {

	flogolog.Debugf("Starting slack RTM..")
	handlers := t.handlers

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
			log.Debugf("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello
				//fmt.Println("Hello")

			case *slack.ConnectedEvent:
				flogolog.Debugf("Infos:", ev.Info)
				flogolog.Debugf("Connection counter:", ev.ConnectionCount)
				// Replace C2147483705 with your Channel ID
				//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "CA6BXNMPC"))

			case *slack.MessageEvent:
				flogolog.Debugf("Message: %v\n", msg.Data.(string))
				t.RunHandler(handler, msg.Data.(string))

			case *slack.PresenceChangeEvent:
				//fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				//fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				flogolog.Debugf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				flogolog.Debugf("Invalid credentials")
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

	flogolog.Debugf("Stopping RTM")

	return nil
}

// RunHandler action on new Slack RTM message
func (t *SlackRecvTrigger) RunHandler(handler *trigger.Handler, payload string) {

	trgData := make(map[string]interface{})
	trgData["message"] = payload

	_, err := handler.Handle(context.Background(), trgData)

	if err != nil {
		log.Error("Error starting action: ", err.Error())
	}

	flogolog.Debugf("Ran Handler: [%s]", handler)
	
}