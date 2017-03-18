package websocket

import (
	"github.com/GrappigPanda/notorious/announce/impl/newTorrentType"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebsocketNotifier struct {
	killChan       chan bool
	connections    []*conn
	NewTorrentChan chan dataType.NewTorrent
	config         config.ConfigStruct
}

func SpawnNotifier(config config.ConfigStruct) *WebsocketNotifier {
	killChan := make(chan bool)
	NewTorrentChan := make(chan dataType.NewTorrent)

	websocketNotify := WebsocketNotifier{
		NewTorrentChan: NewTorrentChan,
		killChan:       killChan,
		connections:    make([]*conn),
		config:         config,
	}

	go websocketNotify.notifyCatcher()

	return &WebsocketNotify
}

func (ws *WebsocketNotifier) notifyCatcher() {
	select {
	case <-ws.killChan:
		log.Println("Received kill notification in IRCNotifer")
		return
	case newTorrent := <-ws.NewTorrentChan:
		ws.sendNotification(newTorrent)
	}
}

func (ws *WebsocketNotifier) KillNotifier() error {
	ws.killChan <- true
	return nil
}

func (ws *WebsocketNotifier) wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	// TODO(ian): Do we need a read buffer here?
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection.", http.StatusBadRequest)
	}
}

func (ws *WebsocketNotifier) sendNotification(newTorrent dataType.NewTorrent) {
	for c := range ws.connections {
		WriteJSON(c, newTorrent)
	}
}
