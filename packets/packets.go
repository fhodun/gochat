package packets

const (
	NewMessage        uint = 1
	ChattersList      uint = 2
	RegisterChatter   uint = 3
	UnregisterChatter uint = 4
)

type Message struct {
	Type    uint     `json:"type"`
	Author  string   `json:"author"`
	Text    string   `json:"text"`
	Content []string `json:"content"`
}
