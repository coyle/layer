# Layer API for Golang

A Golang library that wraps the [Layer](https://layer.com) Platform API.


## Documentation

You can find full documentation on Platform API at [developer.layer.com/docs/platform](https://developer.layer.com/docs/platform).

## Installation

    go get github.com/coyle/layer

## Simple example

```Go
package main

import "github.com/coyle/layer"

var (
  token    = os.Getenv("LAYER_TOKEN")
	appID    = os.Getenv("LAYER_APPID")
	version  = "1.0"
	timeout  = 30 * time.Second
	l        = NewLayer(token, appID, version, timeout)
)

// Metadata struct can contain any metadata you would like to pass into a conversation
type Metadata struct {

}


func main() {
  // Create a conversation
  distinct := true
  conversation, err := l.CreateConversation([]string{"user1", "user2"}, distinct, Metadata{})

  // Send a message
  cid := conversation.GetID()

  p := Parts{
		Body:     "Hello World",
		MimeType: "text/plain",
	}
	n := Notification{
		Text:  "Hello World Notification",
		Sound: "chime.aiff",
	}

	res, err := l.SendMessage(cid, "user1", []Parts{p}, n)
  if err != nil {
    fmt.Errorf("ERROR=%v", err)
  }
  fmt.Printf("RESULT=%#v", res)
}
```

### LayerAPI (config)

Layer API constructor is initialized with the following configuration values:

 - `token` - Layer Platform API token which can be obtained from [Developer Dashboard](https://developer.layer.com/projects/keys)
 - `appId` - Layer application ID
 - `version` - API version to use
 - `timeout` - Request timeout

## Conversations

Conversations coordinate messaging within Layer and can contain up to 25 participants. All Messages sent are sent within the context of a conversation.

## Messages

Messages can be made up of one or many individual pieces of content.

 - Message `sender` can be specified by `userID`
 - Message `parts` are the atomic object in the Layer universe. They represent the individual pieces of content embedded within a message.
 - Message `notification` object represents [push notification](https://developer.layer.com/docs/platform#push-notifications) payload.

## Announcements

Announcements are messages sent to all users of the application or to a list of users.

Payload property `recipients` can contain one or more user IDs or the literal string "everyone" in order to message the entire userbase.

Send an [Announcement](https://developer.layer.com/docs/platform#send-an-announcement) by providing an AnnouncementRequest.


## Block list

Layer Platform API allows you to manage a [block list](https://developer.layer.com/docs/platform#managing-user-block-lists) in order to align with your own application level blocking. A block list is maintained for each user, enabling users to manage a list of members they don't want to communicate with.


## Testing
  To run tests you must first get a Layer token ([Developer Dashboard](https://developer.layer.com/projects/keys)) and appID.
  The appID must be set as the environment variable `LAYER_TEST_APPID` and token must be set on the environment as `LAYER_TEST_TOKEN`

  go test .

## Contributing

Feedback and contributions are always welcome. Feel free to open up a Pull Request or Issue on Github.

## TODO
1. Extend SendMessage to accept a name in addition to the layerID
2. Implement UploadRichContent
3. Implement SendMessageWithRichContent

## Author

[Coyle](https://github.com/coyle)
