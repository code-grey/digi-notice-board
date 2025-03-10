package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	
	"github.com/gorilla/websocket"
	"digi-notice-board/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
	return true
	},
}

type Client struct {
	Conn *websocket.Conn 
}

var (
	wsMutex sync.Mutex
	clients = make(map[*Client]bool)
	Broadcast = make(chan models.Announcement)
)
func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	
	client := &Client{Conn: wsConn}
	wsMutex.Lock()
	clients[client] = true
	wsMutex.Unlock()
	
	log.Println("New WebSocket client connected")
	
	go readMessages(client)
	
}

func readMessages(client *Client) {
	defer func(){
		wsMutex.Lock()
		delete(clients, client)
		wsMutex.Unlock()
		client.Conn.Close()
		log.Println("WebSocket client disconnected")
	}()
	
	for {
        _, _, err := client.Conn.ReadMessage()
        if err != nil {
            log.Println("WebSocket read error:", err)
            break
        }
    } 
}

func StartBroadcast() {
    for {
        announcement := <-Broadcast
        msg, err := json.Marshal(announcement)
        if err != nil {
            log.Println("Error marshaling announcement:", err)
            continue
        }

        wsMutex.Lock()
        for client := range clients {
            err := client.Conn.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Println("Error writing message to client:", err)
                client.Conn.Close()
                delete(clients, client)
            }
        }
        wsMutex.Unlock()
    }
}
