package catcherImpl

import (
	"encoding/json"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/postgres"
	"github.com/lib/pq"
)

type PostgresCatcher struct {
	pglisten   *postgres.PGListener
	config     config.ConfigStruct
	ircNotifer IRCNotifier
}

func (p *PostgresCatcher) serveNewTorrent(notify *pq.Notification) {
	p.ircNotifer.newTorrentChan <- deserializeNotification(notify)
}

func (p *PostgresCatcher) HandleNewTorrent() {
	go p.pglisten.BeginListen(p.serveNewTorrent)
}

func NewPostgresCatcher(cfg config.ConfigStruct) *PostgresCatcher {
	pglisten, err := postgres.NewListener(cfg)
	if err != nil {
		panic(err)
	}

	return &PostgresCatcher{
		pglisten:   pglisten,
		config:     cfg,
		ircNotifer: *SpawnNotifier(cfg),
	}
}

func deserializeNotification(notify *pq.Notification) newTorrent {
	var torrent newTorrent
	err := json.Unmarshal([]byte(notify.Extra), &torrent)
	if err != nil {
		println(err)
	}

	return torrent
}
