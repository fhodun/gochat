package server

type hub struct {
	Clients    map[*client]bool
	Register   chan *client
	Unregister chan *client
	Broadcast  chan []byte
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				close(client.Send)
				delete(h.Clients, client)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
				}
			}
		}
	}
}
