package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/time/rate"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	wsConnectionLimiterLock sync.Mutex

	wsConnectionLimiterCache = ttlcache.New(
		ttlcache.WithTTL[*websocket.Conn, *rate.Limiter](time.Second * 5),
	)
)

func reader(conn *websocket.Conn, client string) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		wsConnectionLimiterLock.Lock()

		item := wsConnectionLimiterCache.Get(conn)
		var limiter *rate.Limiter

		if item == nil {
			limiter = rate.NewLimiter(rate.Limit(10), 20)
			wsConnectionLimiterCache.Set(conn, limiter, ttlcache.DefaultTTL)
		} else {
			limiter = item.Value()
		}

		wsConnectionLimiterLock.Unlock()
		if !limiter.Allow() {
			sendAdminMessage("You are trying to send too many messages :)", conn)
		} else {
			message := string(p)
			if string(message[0]) == "/" {
				switch message[1:len(message)] {
				case "count":
					sendAdminMessage(len(connectedClients), conn)
				case "users":
					users := ""
					iter := 0
					connectedClientsLock.Lock()
					for i, _ := range connectedClients {
						iter++
						if len(connectedClients) == iter {
							users += i
						} else {
							users += (i + ", ")
						}
					}
					connectedClientsLock.Unlock()
					sendAdminMessage(users, conn)
				default:
					if strings.Contains(message[1:len(message)], "disc0nnect") { // hacker like, I know
						connectedClientsLock.Lock()
						newConn := connectedClients[strings.Split(message[1:len(message)], " ")[1]]
						newConn.Close()
						connectedClientsLock.Unlock()
					}
					sendAdminMessage("Command not found", conn)
				}
			} else {
				sendToAllClients(message, client, messageType, conn)
			}
		}
	}
}

type Message struct {
	Sender  string
	Message string
}

func sendAdminMessage(message any, conn *websocket.Conn) {
	if reflect.TypeOf(message).String() == "int" {
		conn.WriteMessage(1, []byte(fmt.Sprintf(`{"Sender":"update","Message":%v}`, message)))
	} else {
		conn.WriteMessage(1, []byte(fmt.Sprintf(`{"Sender":"update","Message":"%v"}`, message)))
	}
}

func sendToAllClients(message string, sender string, msgtype int, conn *websocket.Conn) {
	messagePrepped, _ := json.Marshal(Message{
		Sender:  sender,
		Message: message,
	})
	connectedClientsLock.Lock()
	for _, v := range connectedClients {
		if v == conn {
			if err := v.WriteMessage(msgtype, []byte(fmt.Sprintf(`{"Sender":"%v (you)":"%v"}`, sender, message))); err != nil {
				log.Println(err)
			}
		} else {
			if err := v.WriteMessage(msgtype, messagePrepped); err != nil {
				log.Println(err)
			}
		}
	}
	connectedClientsLock.Unlock()
}
