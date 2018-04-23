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

	fmt.Printf("Starting slack RTM..")
	handlers := t.handlers
	
	flogolog.Debug("Processing handlers")
	for _, handler := range handlers {

		channel := handler.GetStringSetting("Channel")
		accessToken := handler.GetStringSetting("AccessToken")
		//accessToken := t.config.GetSetting("AccessToken")
		//flogolog.Debug("AccessToken: ", accessToken)
		api := slack.New(accessToken)
		channelid := t.GetChannelID(accessToken, channel)
		logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
		slack.SetLogger(logger)
		api.SetDebug(true)

		rtm := api.NewRTM()
		go rtm.ManageConnection()

		for msg := range rtm.IncomingEvents {
			flogolog.Debugf("Event Received: ")
			
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				//t.RunHandler(handler, "Hello!")
				// Ignore hello
				//fmt.Println("Hello")

			case *slack.ConnectedEvent:
				fmt.Printf("Infos:", ev.Info)
				fmt.Printf("Connection counter:", ev.ConnectionCount)
				// Replace C2147483705 with your Channel ID
				//rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "CA6BXNMPC"))

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				if channelid == ev.Channel {
					t.RunHandler(handler, ev.Text)
				}

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
	
	fmt.Printf("Stopping RTM...")

	return nil
}

//GetChannelID returns channel ID for channel name
func(t * SlackRecvTrigger) GetChannelID(accessToken string , channelName string) (string) {
	api_var := slack.New(accessToken)
	channels, err := api_var.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}
	for _, channel := range channels {
		fmt.Println("Channel :  %v", channel)
		if channel.Name == channelName {
			fmt.Println("Found Channel:", channel.Name)
			return channel.ID
		}
	}
	return ""
}

// RunHandler action on new Slack RTM message
func (t *SlackRecvTrigger) RunHandler(handler *trigger.Handler, payload string) {

	trgData := make(map[string]interface{})
	trgData["message"] = payload

	_, err := handler.Handle(context.Background(), trgData)

	if err != nil {
		fmt.Printf("Error starting action: ", err.Error())
	}

	fmt.Printf("Ran Handler: [%s]", handler)
	
}
