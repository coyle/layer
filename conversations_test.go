package layer

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

var (
	token    = os.Getenv("LAYER_TEST_TOKEN")
	appID    = os.Getenv("LAYER_TEST_APPID")
	version  = "1.0"
	timeout  = 30 * time.Second
	l        = NewLayer(token, appID, version, timeout)
	convHead = "layer:///conversations/"
)

func TestCreateConversation(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, []byte{})
	require.NoError(t, err)
	require.Contains(t, res.Participants, user1)
	require.Contains(t, res.Participants, user2)
}

func TestGetAllConversationsForUser(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	_, err := l.CreateConversation([]string{user1, user2}, true, []byte{})
	require.NoError(t, err)

	resp, err := l.GetAllConversationsForUser(user1, nil)
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Contains(t, resp[0].Participants, user1)
	require.Contains(t, resp[0].Participants, user2)
}

func TestGetConversationForUser(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, []byte{})
	require.NoError(t, err)
	convoID := getID(res.ID)

	res2, err := l.GetConversationForUser(user1, convoID, nil)
	require.NoError(t, err)
	convoID2 := getID(res2.ID)
	require.Equal(t, convoID, convoID2)

	require.Contains(t, res2.Participants, user1)
	require.Contains(t, res2.Participants, user2)
}

func TestGetConversation(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, []byte{})
	require.NoError(t, err)
	convoID := getID(res.ID)

	res2, err := l.GetConversation(convoID)
	require.NoError(t, err)
	convoID2 := getID(res2.ID)
	require.Equal(t, convoID, convoID2)
}

func TestAddParticipants(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, []byte{})
	convoID := getID(res.ID)
	require.NoError(t, err)

	user3 := uuid.New()
	ok, err := l.AddParticipants(convoID, []string{user3})
	require.NoError(t, err)
	require.True(t, ok)

	res3, err := l.GetConversation(convoID)
	require.NoError(t, err)
	require.Len(t, res3.Participants, 3)
	require.Contains(t, res3.Participants, user3)
}

func TestRemoveParticipants(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	user3 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2, user3}, true, []byte{})
	convoID := getID(res.ID)
	require.NoError(t, err)

	ok, err := l.RemoveParticipants(convoID, []string{user3})
	require.NoError(t, err)
	require.True(t, ok)

	res3, err := l.GetConversation(convoID)
	require.NoError(t, err)
	require.Len(t, res3.Participants, 2)
	require.NotContains(t, res3.Participants, user3)
}

func TestSetParticipants(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()

	user3 := uuid.New()
	user4 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, []byte{})
	convoID := getID(res.ID)
	require.NoError(t, err)

	ok, err := l.SetParticipants(convoID, []string{user3, user4})
	require.NoError(t, err)
	require.True(t, ok)

	res3, err := l.GetConversation(convoID)
	require.NoError(t, err)
	require.Len(t, res3.Participants, 2)
	require.Contains(t, res3.Participants, user3)
	require.Contains(t, res3.Participants, user4)
}

func TestDeleteConversation(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, []byte{})
	require.NoError(t, err)
	convoID := getID(res.ID)

	res2, err := l.GetConversation(convoID)
	require.NoError(t, err)
	convoID2 := getID(res2.ID)
	require.Equal(t, convoID, convoID2)

	ok, err := l.DeleteConversation(convoID)
	require.True(t, ok)
	require.NoError(t, err)

	res3, err := l.GetConversation(convoID)
	require.NoError(t, err)
	convoID3 := getID(res3.ID)
	require.Equal(t, "object_deleted", convoID3)
}

func getID(s string) string {
	return strings.Replace(s, convHead, "", -1)
}
