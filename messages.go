package layer

// SendMessage creates a new message in a conversation
func SendMessage() {}

// UploadRichContent allows messages whose body is larger than 2KB to be sent. Must be called prior to sending
// message with rich content
func UploadRichContent() {}

// SendMessageWithRichContent will send a Message that includes the Rich Content part of a messgae
// once the Rich Content upload has completed from a send message request
func SendMessageWithRichContent() {}

// GetMessagesForUser requests all messages in a conversation from a specific user's perspective
func GetMessagesForUser() {}

// GetAllMessages requests all messages in a conversation from the System's perspective
func GetAllMessages() {}

// GetMessageForUser requests a single message from a conversation from a specific user's perspective
func GetMessageForUser() {}

// GetMessage request a single message from a conversation from the System's perspective
func GetMessage() {}

// DeleteMessage causes the message to be destroyed for all recipients.
func DeleteMessage() {}
