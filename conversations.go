package layer

import (
	"encoding/json"
	"fmt"
	"time"
)

// ConversationResponse contains fields returned in the JSON response of requests made to the conversation endpoint
type ConversationResponse struct {
	ID                 string      `json:"id,omitempty"`
	URL                string      `json:"url,omitempty"`
	MessagesURL        string      `json:"messages_url,omitempty"`
	Created            time.Time   `json:"created_at,omitempty"`
	MetaData           []byte      `json:"metadata,omitempty"`
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
}

// Sender contains information pertaining to the Sender of the LastMessage within a ConversationResponse
type Sender struct {
	UserID string `json:"user_id,omitempty"`
}

type createConversationBody struct {
	Participants []string `json:"participants"`
	Distinct     bool     `json:"distinct"`
	MetaData     []byte   `json:"metadata,omitempty"`
}

type editConversation struct {
	Operation string `json:"operation"`
	Property  string `json:"property"`
	Value     []byte `json:"value"`
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
func (l *Layer) CreateConversation(participants []string, distinct bool, metadata []byte) (ConversationResponse, error) {
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

// EditConversation adds/removes the participants or sets/deletes metadata properties of a conversation
func (l *Layer) EditConversation(convID, operation, property string, value []byte) (ConversationResponse, error) {
	cr := ConversationResponse{}
	c := editConversation{Operation: "set", Property: property, Value: value}
	body, err := json.Marshal(&c)
	if err != nil {
		return cr, err
	}
	p := Parameters{Path: fmt.Sprintf("conversations/%s", convID), Body: body}
	resp, err := l.request("PATCH", &p)
	if err != nil {
		return cr, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&cr)
	return cr, nil
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
