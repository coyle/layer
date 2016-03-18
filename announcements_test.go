package layer

import (
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

func TestSendAnnouncement(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()
	ar := AnnouncementRequest{
		Recipients:   []string{user1, user2},
		Sender:       Sender{Name: user1},
		Parts:        []Parts{Parts{Body: "Hello World!", MimeType: "text/plain"}},
		Notification: Notification{Text: "This is the alert text to include with the Push Notification.", Sound: "chime.aiff"},
	}

	res, err := l.SendAnnouncement(ar)
	require.NoError(t, err)
	// test Recipients
	require.Contains(t, res.Recipients, user1)
	require.Contains(t, res.Recipients, user2)
	//test sender
	require.Equal(t, res.Sender.Name, user1)
	// test parts
	require.Len(t, res.Parts, 1)
	require.Equal(t, res.Parts[0].MimeType, "text/plain")
	require.Equal(t, res.Parts[0].Body, "Hello World!")

}
