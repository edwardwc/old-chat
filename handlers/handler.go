package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	connectedClients = make(map[string]*websocket.Conn) // client number and connection

	connectedClientsLock sync.Mutex
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header.Get("Upgrade"))
	switch r.Header.Get("Upgrade") {
	case "websocket":
		fmt.Println("Websocket")
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// upgrade this connection to a WebSocket
		// connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		name := ""

		for {
			name = generateName()
			if checkName(name) == true {
				break
			}
		}

		connectedClientsLock.Lock()
		connectedClients[name] = ws
		connectedClientsLock.Unlock()
		sendToAllClients(fmt.Sprintf("New client, %v, is listening for a total of %v clients listening.", name, len(connectedClients)), "system", 1, nil)

		reader(ws, name)

		sendToAllClients(fmt.Sprintf("Client %v has disconnected. There are now %v clients listening.", name, len(connectedClients)-1), "system", 1, nil) // ensure this comes from system on frontend
		connectedClientsLock.Lock()
		delete(connectedClients, name)
		connectedClientsLock.Unlock()
	default:
		if r.URL.Path == "/" {
			fmt.Println("wow lol")
			http.ServeFile(w, r, "/root/chat/index.html")
			return
		}
		fmt.Println("not websocket")
		w.WriteHeader(404)
	}
}
