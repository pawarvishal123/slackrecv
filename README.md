# Slack RTM Message Receiver Trigger
Flogo trigger activity for Slack RTM message receiver


## Installation

```bash
flogo install github.com/pawarvishal123/slackrecv
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "output": [
    {
      "name": "message",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [{
      "name": "AccessToken",
      "type": "string",
	  "required":"true"
    },
    {
      "name": "Channel",
      "type": "string",
	  "required":"true"
    }]
```

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Slack RTM Trigger.

### Start a flow
Provide access token and channel name to receive message from slack channel via RTM. The access token should have scope and permissions configured to allow streaming of messages.

```json
{
  "triggers": [
    {
      "id": "receive_slack_rtm_messages",
      "ref": "https://github.com/pawarvishal123/slackrecv",
      "name": "Receive Slack RTM Messages",
      "description": "Slack RTM Message Trigger",
      "settings": {},
      "handlers": [
        {
          "action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:test_trigger"
            }
          },
          "settings": {
            "AccessToken": "<<YOUR-TOKEN>>",
            "Channel": "<<Your-CHANNEL>>"
          }
        }
      ]
    }
}
```

## Third Party Library
Slack API in Go - [https://github.com/nlopes/slack](https://github.com/nlopes/slack)
