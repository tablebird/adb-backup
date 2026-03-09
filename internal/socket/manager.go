package socket

import (
	"sync"
)

var (
	clients   = make(map[Client]bool)
	clientsMu sync.Mutex // 保证多线程安全操作clients
)

func deleteClient(client Client) {
	clientsMu.Lock()
	delete(clients, client)
	clientsMu.Unlock()
}

func addClient(client Client) {
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()
}
