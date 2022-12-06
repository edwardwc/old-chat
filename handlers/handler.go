package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	connectedClients = make(map[string]*websocket.Conn) // client number and connection

	connectedClientsLock sync.Mutex

	connectionLimiterLock sync.Mutex

	ConnectionLimiterCache = ttlcache.New(
		ttlcache.WithTTL[string, *rate.Limiter](time.Second * 5),
	)
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	connectionLimiterLock.Lock()

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	item := ConnectionLimiterCache.Get(ip)
	var limiter *rate.Limiter

	if item == nil {
		limiter = rate.NewLimiter(rate.Limit(5), 10)
		ConnectionLimiterCache.Set(ip, limiter, ttlcache.DefaultTTL)
	} else {
		limiter = item.Value()
	}

	connectionLimiterLock.Unlock()

	if !limiter.Allow() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(429)
		w.Write([]byte(`{ "message": "Too many requests!" }`))
		return
	}
	switch r.Header.Get("upgrade") {
	case "websocket":
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
		w.Write([]byte("Regular request detected"))
	}
}
