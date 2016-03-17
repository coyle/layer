package layer

import (
	"strings"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

var (
	msgHead = "layer:///messages/"
)

func TestSendMessage(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convID := getConvID(res.ID)
	defer cleanUpConversation(convID)
	require.Contains(t, res.Participants, user1)
	require.Contains(t, res.Participants, user2)

	p := Parts{
		Body:     "Hello World",
		MimeType: "text/plain",
	}
	n := Notification{
		Text:  "Hello World Notification",
		Sound: "chime.aiff",
	}
	res2, err := l.SendMessage(convID, user1, []Parts{p}, n)
	require.NoError(t, err)
	require.Equal(t, res2.Sender.UserID, user1)

	res3, err := l.GetAllMessages(convID)
	require.NoError(t, err)
	require.Len(t, res3, 1)
	require.Equal(t, p.Body, res3[0].Parts[0].Body)
	require.Equal(t, convID, getConvID(res3[0].Conversation.ID))
}

func TestGetMessagesForUser(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convID := getConvID(res.ID)
	defer cleanUpConversation(convID)
	require.Contains(t, res.Participants, user1)
	require.Contains(t, res.Participants, user2)

	p := Parts{
		Body:     "Hello World",
		MimeType: "text/plain",
	}
	n := Notification{
		Text:  "Hello World Notification",
		Sound: "chime.aiff",
	}
	res2, err := l.SendMessage(convID, user1, []Parts{p}, n)
	require.NoError(t, err)

	require.Equal(t, res2.Sender.UserID, user1)

	res3, err := l.GetMessagesForUser(convID, user2)
	require.NoError(t, err)

	require.Len(t, res3, 1)

	require.Equal(t, p.Body, res3[0].Parts[0].Body)
	require.Equal(t, convID, getConvID(res3[0].Conversation.ID))
}

func TestGetMessageForUser(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convID := getConvID(res.ID)
	defer cleanUpConversation(convID)
	require.Contains(t, res.Participants, user1)
	require.Contains(t, res.Participants, user2)

	p := Parts{
		Body:     "Hello World",
		MimeType: "text/plain",
	}
	n := Notification{
		Text:  "Hello World Notification",
		Sound: "chime.aiff",
	}
	res2, err := l.SendMessage(convID, user1, []Parts{p}, n)
	require.NoError(t, err)
	msgID := getMessageID(res2.ID)
	require.Equal(t, res2.Sender.UserID, user1)

	res3, err := l.GetMessageForUser(user2, msgID)
	require.NoError(t, err)
	require.Equal(t, p.Body, res3.Parts[0].Body)
}

func getMessageID(s string) string {
	return strings.Replace(s, msgHead, "", -1)
}
