package layer

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	convHead = "layer:///conversations/"
)

// ConversationResponse contains fields returned in the JSON response of requests made to the conversation endpoint
type ConversationResponse struct {
	ID                 string      `json:"id,omitempty"`
	URL                string      `json:"url,omitempty"`
	MessagesURL        string      `json:"messages_url,omitempty"`
	Created            time.Time   `json:"created_at,omitempty"`
	MetaData           interface{} `json:"metadata,omitempty"`
	Distinct           bool        `json:"distinct,omitempty"`
	LastMessage        LastMessage `json:"last_message,omitempty"`
	UnreadMessageCount int         `json:"recipient_status,omitempty"`
	Participants       []string    `json:"participants,omitempty"`
}

// LastMessage contains information referring to the lastMessage in a ConversationResponse
type LastMessage struct {
	ID              string       `json:"id,omitempty"`
	Position        int          `json:"position,omitempty"`
	Conversation    Conversation `json:"conversation,omitempty"`
	Parts           []Parts      `json:"parts,omitempty"`
	SentAt          time.Time    `json:"sent_at,omitempty"`
	ReceivedAt      time.Time    `json:"received_at,omitempty"`
	Sender          Sender       `json:"sender,omitempty"`
	Unread          bool         `json:"is_unread,omitempty"`
	RecipientStatus []byte       `json:"recipient_status,omitempty"`
}

// Conversation refers to the conversation object in a ConversationResponse
type Conversation struct {
	ID  string `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

// Parts contains information pertaining to the parts of a conversation message
type Parts struct {
	ID       string `json:"id,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Body     string `json:"body,omitempty"`
	Encoding string `json:"encoding,omitempty"`
}

// Sender contains information pertaining to the Sender of the LastMessage within a ConversationResponse
type Sender struct {
	UserID string `json:"user_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

type createConversationBody struct {
	Participants []string    `json:"participants"`
	Distinct     bool        `json:"distinct"`
	MetaData     interface{} `json:"metadata,omitempty"`
}

type editConversationBody struct {
	Operation string `json:"operation"`
	Property  string `json:"property"`
	Value     []byte `json:"value"`
}

type participantBody struct {
	Operation string `json:"operation"`
	Property  string `json:"property"`
	Value     string `json:"value"`
}

type setParticipants struct {
	Operation string   `json:"operation"`
	Property  string   `json:"property"`
	Value     []string `json:"value"`
}

type metadata struct {
	Operation string      `json:"operation"`
	Property  string      `json:"property"`
	Value     interface{} `json:"value"`
}

// GetConversationID returns the conversation ID from a conversation response object
func (c ConversationResponse) GetConversationID() string {
	return strings.Replace(c.ID, convHead, "", -1)
}

// GetAllConversationsForUser requests all conversations for a specific user
func (l *Layer) GetAllConversationsForUser(userID string, params *QueryParameters) ([]ConversationResponse, error) {
	if userID == "" {
		return []ConversationResponse{}, ErrMissingUserID
	}

	p := Parameters{Path: fmt.Sprintf("users/%s/conversations", userID)}
	resp, err := l.request("GET", &p)
	if err != nil {
		return []ConversationResponse{}, err
	}
	defer resp.Body.Close()

	cr := []ConversationResponse{}
	json.NewDecoder(resp.Body).Decode(&cr)
	return cr, nil
}

// GetConversationForUser request a specific Conversation for a user.
func (l *Layer) GetConversationForUser(userID, convID string, params *QueryParameters) (ConversationResponse, error) {
	if userID == "" {
		return ConversationResponse{}, ErrMissingUserID
	}

	p := Parameters{Path: fmt.Sprintf("users/%s/conversations/%s", userID, convID)}

	resp, err := l.request("GET", &p)
	if err != nil {
		return ConversationResponse{}, err
	}
	defer resp.Body.Close()

	cr := ConversationResponse{}
	json.NewDecoder(resp.Body).Decode(&cr)
	return cr, nil
}

// GetConversation requests the Conversation with the given ID
func (l *Layer) GetConversation(convID string) (ConversationResponse, error) {
	p := Parameters{Path: fmt.Sprintf("conversations/%s", convID)}

	resp, err := l.request("GET", &p)
	if err != nil {
		return ConversationResponse{}, err
	}
	defer resp.Body.Close()

	cr := ConversationResponse{}
	json.NewDecoder(resp.Body).Decode(&cr)
	return cr, nil
}

// CreateConversation creates a conversation between two or more participants
func (l *Layer) CreateConversation(participants []string, distinct bool, metadata interface{}) (ConversationResponse, error) {
	cr := ConversationResponse{}
	if len(participants) == 0 {
		return cr, ErrEmptyParticipants
	}

	c := createConversationBody{Participants: participants, Distinct: distinct, MetaData: metadata}
	body, err := json.Marshal(&c)
	if err != nil {
		return cr, err
	}
	p := Parameters{Path: "conversations", Body: body}

	resp, err := l.request("POST", &p)
	if err != nil {
		return cr, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&cr)
	return cr, nil
}

// AddParticipants adds one or more participants to a conversation
func (l *Layer) AddParticipants(convID string, participants []string) (ok bool, err error) {
	cc := buildParticipantBoy(participants, "add")
	body, err := json.Marshal(&cc)
	if err != nil {
		return false, err
	}
	return l.editConversation(convID, body)
}

// RemoveParticipants removes  one or more participants from a conversation
func (l *Layer) RemoveParticipants(convID string, participants []string) (ok bool, err error) {
	cc := buildParticipantBoy(participants, "remove")
	body, err := json.Marshal(&cc)
	if err != nil {
		return false, err
	}
	return l.editConversation(convID, body)
}

// SetParticipants will replace the entire set of participants with a new list
func (l *Layer) SetParticipants(convID string, participants []string) (ok bool, err error) {
	cc := []setParticipants{setParticipants{Operation: "set", Property: "participants", Value: participants}}
	body, err := json.Marshal(&cc)
	if err != nil {
		return false, err
	}
	return l.editConversation(convID, body)

}

func (l *Layer) editConversation(convID string, body []byte) (bool, error) {
	p := Parameters{Path: fmt.Sprintf("conversations/%s", convID), Body: body}
	resp, err := l.request("PATCH", &p)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return false, fmt.Errorf("Responded with Error Code %d", resp.StatusCode)
	}
	return true, nil

}

func buildParticipantBoy(p []string, operation string) []participantBody {
	cc := []participantBody{}
	for _, v := range p {
		cc = append(cc, participantBody{Operation: operation, Property: "participants", Value: v})
	}

	return cc
}

// DeleteConversation removes a conversation's history
func (l *Layer) DeleteConversation(convID string) (ok bool, err error) {
	p := Parameters{Path: fmt.Sprintf("conversations/%s", convID)}
	resp, err := l.request("DELETE", &p)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return false, fmt.Errorf("Responded with Error Code %d", resp.StatusCode)
	}
	return true, nil
}

// DeleteMetadata removes metadata properties from a conversation
func (l *Layer) DeleteMetadata(convID, property string) (bool, error) {
	cc := []metadata{metadata{Operation: "delete", Property: property}}
	body, err := json.Marshal(&cc)
	if err != nil {
		return false, err
	}
	return l.editConversation(convID, body)
}

// SetMetadata sets metadata properties on a conversation
func (l *Layer) SetMetadata(convID, property string, value interface{}) (bool, error) {
	cc := []metadata{metadata{Operation: "set", Property: property, Value: value}}
	body, err := json.Marshal(&cc)
	if err != nil {
		return false, err
	}
	return l.editConversation(convID, body)
}
