package layer

// GetAllConversationsForUser requests all conversations for a specific user
func GetAllConversationsForUser() {}

// GetConversationForUser request a specific Conversation for a user.
func GetConversationForUser() {}

// GetConversation requests the Conversation with the given ID
func GetConversation() {}

// CreateConversation creates a converstaion between two or more users
func CreateConversation() {}

// EditConversation edits the properties of a conversation
func EditConversation() {}

// DeleteConversation removes a conversation's history
func DeleteConversation() {}

// SendMessage creates a new message in a conversation
func SendMessage() {}

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

// SendAnnouncement messages are sent to all users of the application or to a list of users.
// These Messages will arrive outside of the context of a conversation
func SendAnnouncement() {}

// AddUserToBlockList adds one member to a user's block list
func AddUserToBlockList() {}

// GetUserBlockList Returns an array of all blocked users for the specified
func GetUserBlockList() {}

// UnblockUser Removes a blocked user from the Block List of the specified user
func UnblockUser() {}

// BulkBlock supports bulk operations on a user's blocklist
func BulkBlock() {}
