package layer

// Notification is the payload of a push notification in layer
type Notification struct {
	Title      string                  `json:"title,omitempty"`
	Text       string                  `json:"text,omitempty"`
	Sound      string                  `json:"sound,omitempty"`
	Recipients map[string]Notification `json:"recipients,omitempty"`
}

// SendAnnouncement messages are sent to all users of the application or to a list of users.
// These Messages will arrive outside of the context of a conversation
func SendAnnouncement() {}
