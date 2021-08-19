package packets

const (
	// New message
	NewMessage uint = 1
	// List of chat members
	ChattersList uint = 2
	// Register new chat member
	RegisterChatter uint = 3
	// Unregister chat member
	UnregisterChatter uint = 4
)

type Message struct {
	// Message type that can be specified by const variables
	Type uint `json:"type"`
	// Author nickname
	Author string `json:"author"`
	// Raw text content
	Text string `json:"text"`
	// Additional content stored in array
	Content []string `json:"content"`
}
