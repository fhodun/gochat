package server

import (
	"github.com/fhodun/gochat/packets"
)

type hub struct {
	Clients    map[*client]bool
	Register   chan *client
	Unregister chan *client
	Broadcast  chan packets.Message
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.Register:
			for c := range h.Clients { //TODO: repeated code
				select {
				case c.Send <- packets.Message{
					Type: packets.RegisterChatter,
					Text: client.Nickname,
				}:
				default:
					h.Unregister <- client
				}
			}
			h.Clients[client] = true
			var clientsNicknames []string
			for c := range h.Clients {
				clientsNicknames = append(clientsNicknames, c.Nickname)
			}
			client.Send <- packets.Message{
				Type:    packets.ChattersList,
				Content: clientsNicknames,
			}
		case client := <-h.Unregister:
			nickname := client.Nickname
			if _, ok := h.Clients[client]; ok {
				close(client.Send)
				delete(h.Clients, client)
			}
			for c := range h.Clients { //TODO: repeated code
				select {
				case c.Send <- packets.Message{
					Type: packets.UnregisterChatter,
					Text: nickname,
				}:
				default:
					h.Unregister <- c
				}
			}
		case m := <-h.Broadcast:
			for client := range h.Clients { //TODO: repeated code
				select {
				case client.Send <- m:
				default:
					h.Unregister <- client
				}
			}
		}
	}
}
