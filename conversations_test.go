package layer

import (
	"encoding/json"
	"fmt"
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

type testMetaAdmin struct {
	Admin testMeta `json:"admin"`
}

type testMeta struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func TestCreateConversation(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	defer cleanUpConversation(getID(res.ID))
	require.Contains(t, res.Participants, user1)
	require.Contains(t, res.Participants, user2)

}

func TestGetAllConversationsForUser(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	defer cleanUpConversation(getID(res.ID))

	resp, err := l.GetAllConversationsForUser(user1, nil)
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Contains(t, resp[0].Participants, user1)
	require.Contains(t, resp[0].Participants, user2)
}

func TestGetConversationForUser(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convoID := getID(res.ID)
	defer cleanUpConversation(convoID)

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
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convoID := getID(res.ID)
	defer cleanUpConversation(convoID)

	res2, err := l.GetConversation(convoID)
	require.NoError(t, err)
	convoID2 := getID(res2.ID)
	require.Equal(t, convoID, convoID2)
}

func TestAddParticipants(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	convoID := getID(res.ID)
	require.NoError(t, err)
	defer cleanUpConversation(convoID)

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
	res, err := l.CreateConversation([]string{user1, user2, user3}, true, Conversation{})
	require.NoError(t, err)
	convoID := getID(res.ID)
	defer cleanUpConversation(convoID)

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
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convoID := getID(res.ID)
	defer cleanUpConversation(convoID)

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
	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
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

func TestDeleteMetadata(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	md := testMeta{Name: "fred", UserID: user1}

	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convID := getID(res.ID)
	defer cleanUpConversation(convID)

	ok, err := l.SetMetadata(convID, "metadata.admin", md)
	require.NoError(t, err)
	require.True(t, ok)

	res2, err := l.GetConversation(convID)
	require.NoError(t, err)

	md2 := testMetaAdmin{}
	d, err := json.Marshal(res2.MetaData)
	require.NoError(t, err)
	json.Unmarshal(d, &md2)
	require.Equal(t, md.Name, md2.Admin.Name)
	require.Equal(t, md.UserID, md2.Admin.UserID)

	ok, err = l.DeleteMetadata(convID, "metadata.admin.name")
	require.NoError(t, err)
	require.True(t, ok)

	res3, err := l.GetConversation(convID)
	require.NoError(t, err)
	d3, err := json.Marshal(res3.MetaData)
	require.NoError(t, err)
	md3 := testMetaAdmin{}
	json.Unmarshal(d3, &md3)
	require.Equal(t, "", md3.Admin.Name)
	require.Equal(t, md.UserID, md3.Admin.UserID)
}

func TestSetMetadata(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	md := testMeta{Name: "fred", UserID: user1}

	res, err := l.CreateConversation([]string{user1, user2}, true, Conversation{})
	require.NoError(t, err)
	convID := getID(res.ID)
	defer cleanUpConversation(convID)

	ok, err := l.SetMetadata(convID, "metadata.admin", md)
	require.NoError(t, err)
	require.True(t, ok)

	res2, err := l.GetConversation(convID)
	require.NoError(t, err)

	md2 := testMetaAdmin{}
	d, err := json.Marshal(res2.MetaData)
	require.NoError(t, err)
	json.Unmarshal(d, &md2)
	require.Equal(t, md.Name, md2.Admin.Name)
	require.Equal(t, md.UserID, md2.Admin.UserID)
}

func getID(s string) string {
	return strings.Replace(s, convHead, "", -1)
}

func cleanUpConversation(convID string) {
	ok, err := l.DeleteConversation(convID)
	if err != nil || !ok {
		panic(fmt.Errorf("Failed to cleanup conversation: %s", convID))
	}
}
