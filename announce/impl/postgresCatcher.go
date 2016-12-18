package catcherImpl

import (
	"encoding/json"
	"github.com/GrappigPanda/notorious/announce/impl/irc"
	"github.com/GrappigPanda/notorious/announce/impl/newTorrentType"
	"github.com/GrappigPanda/notorious/announce/impl/rss"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/postgres"
	"github.com/lib/pq"
)

type PostgresCatcher struct {
	pglisten    *postgres.PGListener
	config      config.ConfigStruct
	ircNotifier *irc.IRCNotifier
	rssNotifier *rss.RSSNotifier
}

func (p *PostgresCatcher) serveNewTorrent(notify *pq.Notification) {
	deserializedNotification := deserializeNotification(notify)

	if p.ircNotifier != nil {
		p.ircNotifier.NewTorrentChan <- deserializedNotification
	}

	if p.rssNotifier != nil {
		p.rssNotifier.NewTorrentChan <- deserializedNotification
	}
}

func (p *PostgresCatcher) HandleNewTorrent() {
	go p.pglisten.BeginListen(p.serveNewTorrent)
}

func (p *PostgresCatcher) GetRSSNotifier() *rss.RSSNotifier {
	return p.rssNotifier
}

func NewPostgresCatcher(cfg config.ConfigStruct) *PostgresCatcher {
	pglisten, err := postgres.NewListener(cfg)
	if err != nil {
		panic(err)
	}

	var ircNotifier *irc.IRCNotifier
	if cfg.IRCCfg != nil {
		ircNotifier = irc.SpawnNotifier(cfg)
	} else {
		ircNotifier = nil
	}

	var rssNotifier *rss.RSSNotifier
	if cfg.UseRSS == true {
		rssNotifier = rss.SpawnNotifier(cfg)
	} else {
		ircNotifier = nil
	}

	return &PostgresCatcher{
		pglisten:    pglisten,
		config:      cfg,
		ircNotifier: ircNotifier,
		rssNotifier: rssNotifier,
	}
}

func deserializeNotification(notify *pq.Notification) dataType.NewTorrent {
	var torrent dataType.NewTorrent
	err := json.Unmarshal([]byte(notify.Extra), &torrent)
	if err != nil {
		println(err)
	}

	return torrent
}
