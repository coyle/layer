package layer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Notification is the payload of a push notification in layer
type Notification struct {
	Title      string                  `json:"title,omitempty"`
	Text       string                  `json:"text,omitempty"`
	Sound      string                  `json:"sound,omitempty"`
	Recipients map[string]Notification `json:"recipients,omitempty"`
}

// SetBadgeRequest holds the number to set the badge count to
type SetBadgeRequest struct {
	Count int `json:"external_unread_count"`
}

// GetBadgeResponse is the Layer response for reading a users badge counts
type GetBadgeResponse struct {
	UnreadExternal     int `json:"external_unread_count"`
	UnreadConversation int `json:"unread_conversation_count"`
	UnreadMessage      int `json:"unread_message_count"`
}

// SetUsersBadge sets an external unread count for a particular user
func (l *Layer) SetUsersBadge(userID string, count int) (ok bool, err error) {
	r := SetBadgeRequest{Count: count}
	body, err := json.Marshal(&r)
	if err != nil {
		return false, err
	}
	p := Parameters{Path: fmt.Sprintf("users/%s/badge", userID), Body: body}
	resp, err := l.request("PUT", &p)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return false, fmt.Errorf("Responded with Error Code %d", resp.StatusCode)
	}

	return true, nil
}

// GetUsersBadge reads the badge for a particular user
func (l *Layer) GetUsersBadge(userID string) (GetBadgeResponse, error) {
	p := Parameters{Path: fmt.Sprintf("users/%s/badge", userID)}
	resp, err := l.request("GET", &p)
	if err != nil {
		return GetBadgeResponse{}, err
	}
	defer resp.Body.Close()

	m := GetBadgeResponse{}
	json.NewDecoder(resp.Body).Decode(&m)
	return m, err
}
