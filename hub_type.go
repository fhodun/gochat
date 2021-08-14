package main

type HubStruct struct {
	clients    map[*ClientStruct]bool
	register   chan *ClientStruct
	unregister chan *ClientStruct
	broadcast  chan []byte
}

func (h *HubStruct) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
				}
			}
		}
	}
}
