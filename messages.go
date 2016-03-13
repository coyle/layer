package layer

import (
	"encoding/json"
	"fmt"
	"time"
)

// MessageResponse is the struct containing fields returned from a successful message response
type MessageResponse struct {
	ID              string            `json:"id,omitempty"`
	URL             string            `json:"url,omitempty"`
	Conversation    Conversation      `json:"conversation,omitempty"`
	Parts           []Parts           `json:"parts,omitempty"`
	SentAt          time.Time         `json:"sent_at,omitempty"`
	Sender          Sender            `json:"sender,omitempty"`
	RecipientStatus map[string]string `json:"recipient_status,omitempty"`
	IsUnread        bool              `json:"is_unread,omitempty"`
	Recieved        time.Time         `json:"received_at,omitempty"`
}

// MessageRequest contains the response from a message request
type MessageRequest struct {
	Sender       Sender       `json:"sender,omitempty"`
	Parts        []Parts      `json:"parts,omitempty"`
	Notification Notification `json:"notification,omitempty"`
}

// SendMessage creates a new message in a conversation
func (l *Layer) SendMessage(convID string, sender string, parts []Parts, n Notification) (MessageResponse, error) {
	b := MessageRequest{Sender: Sender{UserID: sender}, Parts: parts, Notification: n}
	body, err := json.Marshal(&b)
	if err != nil {
		return MessageResponse{}, err
	}
	p := Parameters{Path: fmt.Sprintf("conversations/%s/messages", convID), Body: body}
	resp, err := l.request("POST", &p)
	if err != nil {
		return MessageResponse{}, err
	}
	defer resp.Body.Close()

	m := MessageResponse{}
	json.NewDecoder(resp.Body).Decode(&m)
	return m, err
}

// UploadRichContent allows messages whose body is larger than 2KB to be sent. Must be called prior to sending
// message with rich content
func UploadRichContent() {}

// SendMessageWithRichContent will send a Message that includes the Rich Content part of a messgae
// once the Rich Content upload has completed from a send message request
func SendMessageWithRichContent() {}

// GetMessagesForUser requests all messages in a conversation from a specific user's perspective
func GetMessagesForUser() {}

// GetAllMessages requests all messages in a conversation from the System's perspective
func (l *Layer) GetAllMessages(convID string) ([]MessageResponse, error) {
	p := Parameters{Path: fmt.Sprintf("conversations/%s/messages", convID)}
	resp, err := l.request("GET", &p)
	if err != nil {
		return []MessageResponse{}, err
	}
	defer resp.Body.Close()

	m := []MessageResponse{}
	json.NewDecoder(resp.Body).Decode(&m)
	return m, err
}

// GetMessageForUser requests a single message from a conversation from a specific user's perspective
func GetMessageForUser() {}

// GetMessage request a single message from a conversation from the System's perspective
func GetMessage() {

}

// DeleteMessage causes the message to be destroyed for all recipients.
func DeleteMessage() {}
