package catcherImpl

import (
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/postgres"
	"github.com/lib/pq"
)

type PostgresCatcher struct {
	pglisten *postgres.PGListener
	config   config.ConfigStruct
}

func (p *PostgresCatcher) serveNewTorrent(*pq.Notification) {

}

func (p *PostgresCatcher) HandleNewTorrent() {
	p.pglisten.BeginListen(p.serveNewTorrent)
}

func NewCatcher(cfg config.ConfigStruct) *PostgresCatcher {
	pglisten, err := postgres.NewListener(cfg)
	if err != nil {
		panic(err)
	}

	return &PostgresCatcher{
		pglisten: pglisten,
		config:   cfg,
	}
}
