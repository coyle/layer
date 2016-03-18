package layer

// Notification is the payload of a push notification in layer
type Notification struct {
	Title      string                  `json:"title,omitempty"`
	Text       string                  `json:"text,omitempty"`
	Sound      string                  `json:"sound,omitempty"`
	Recipients map[string]Notification `json:"recipients,omitempty"`
}

// SetUsersBadge sets an external unread count for a particular user
func SetUsersBadge() {}

// GetUsersBadge reads the badge for a particular user
func GetUsersBadge() {}
