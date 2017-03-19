package ws

import (
	"github.com/GrappigPanda/notorious/announce/impl/newTorrentType"
	"github.com/GrappigPanda/notorious/config"
	"github.com/gorilla/websocket"
	"github.com/trevex/golem"
	"log"
	"net/http"
	"sync"
)

type WSNotifier struct {
	killChan       chan bool
	connections    []*websocket.Conn
	NewTorrentChan chan dataType.NewTorrent
	config         config.ConfigStruct
	roomManager    *golem.RoomManager
	sync.Mutex
}

func SpawnNotifier(config config.ConfigStruct) *WSNotifier {
	killChan := make(chan bool)
	NewTorrentChan := make(chan dataType.NewTorrent)

	websocketNotify := WSNotifier{
		NewTorrentChan: NewTorrentChan,
		killChan:       killChan,
		connections:    []*websocket.Conn{},
		roomManager:    golem.NewRoomManager(),
		config:         config,
	}

	go websocketNotify.notifyCatcher()

	return &websocketNotify
}

func (ws *WSNotifier) notifyCatcher() {
	select {
	case <-ws.killChan:
		log.Println("Received kill notification in IRCNotifer")
		return
	case newTorrent := <-ws.NewTorrentChan:
		ws.sendNotification(newTorrent)
	}
}

func (ws *WSNotifier) KillNotifier() error {
	ws.killChan <- true
	return nil
}

func (ws *WSNotifier) WSHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming websocket connection.")

	log.Println("Added new websocket connection.")
}

func (ws *WSNotifier) sendNotification(newTorrent dataType.NewTorrent) {
	for _, c := range ws.connections {
		websocket.WriteJSON(c, newTorrent)
	}
}
