package layer

import (
	"os"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

var (
	token   = os.Getenv("LAYER_TEST_TOKEN")
	appID   = os.Getenv("LAYER_TEST_APPID")
	version = "1.0"
	timeout = 30 * time.Second
	l       = NewLayer(token, appID, version, timeout)
	uid     = uuid.New()
	uid2    = uuid.New()
)

func TestCreateConversation(t *testing.T) {
	_, err := l.CreateConversation([]string{uid, uid2}, true, []byte{})
	require.NoError(t, err)
}

func TestGetAllConversationsForUser(t *testing.T) {
	resp, err := l.GetAllConversationsForUser(uid, nil)
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Contains(t, resp[0].Participants, uid)
	require.Contains(t, resp[0].Participants, uid2)
}
