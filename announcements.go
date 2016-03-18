package layer

import (
	"encoding/json"
	"time"
)

// AnnouncementRequest represents the body for an announcement
type AnnouncementRequest struct {
	Recipients   []string     `json:"recipients"`
	Sender       Sender       `json:"sender"`
	Parts        []Parts      `json:"parts"`
	Notification Notification `json:"notification"`
}

// AnnouncementResponse represents the Layer response body from sending an announcement
type AnnouncementResponse struct {
	ID         string    `json:"id"`
	URL        string    `json:"url"`
	Sent       time.Time `json:"sent_at"`
	Recipients []string  `json:"recipients"`
	Sender     Sender    `json:"sender"`
	Parts      []Parts   `json:"parts"`
}

// SendAnnouncement messages are sent to all users of the application or to a list of users.
// These Messages will arrive outside of the context of a conversation
func (l *Layer) SendAnnouncement(req AnnouncementRequest) (AnnouncementResponse, error) {
	body, err := json.Marshal(&req)
	if err != nil {
		return AnnouncementResponse{}, err
	}
	p := Parameters{Path: "announcements", Body: body}
	resp, err := l.request("POST", &p)

	if err != nil {
		return AnnouncementResponse{}, err
	}
	defer resp.Body.Close()

	ar := AnnouncementResponse{}
	json.NewDecoder(resp.Body).Decode(&ar)
	return ar, err
}
