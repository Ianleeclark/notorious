package catcherImpl

import (
	"encoding/json"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/postgres"
	"github.com/lib/pq"
)

type PostgresCatcher struct {
	pglisten    *postgres.PGListener
	config      config.ConfigStruct
	ircNotifier *IRCNotifier
}

func (p *PostgresCatcher) serveNewTorrent(notify *pq.Notification) {
	if p.ircNotifier != nil {
		p.ircNotifier.newTorrentChan <- deserializeNotification(notify)
	}
}

func (p *PostgresCatcher) HandleNewTorrent() {
	go p.pglisten.BeginListen(p.serveNewTorrent)
}

func NewPostgresCatcher(cfg config.ConfigStruct) *PostgresCatcher {
	pglisten, err := postgres.NewListener(cfg)
	if err != nil {
		panic(err)
	}

	var ircNotifier *IRCNotifier
	if cfg.IRCCfg != nil {
		ircNotifier = SpawnNotifier(cfg)
	} else {
		ircNotifier = nil
	}

	return &PostgresCatcher{
		pglisten:    pglisten,
		config:      cfg,
		ircNotifier: ircNotifier,
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
