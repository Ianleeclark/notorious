package catcherImpl

import (
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/mysql"
	"github.com/GrappigPanda/notorious/database/schemas"
)

type MysqlCatcher struct {
	mysqlListen chan schemas.WhiteTorrent
	config      config.ConfigStruct
}

func (m *MysqlCatcher) serveNewTorrent(notify *pq.Notification) {
}

func (m *MysqlCatcher) HandleNewTorrent() {
}

func NewCatcher(cfg config.ConfigStruct) *MysqlCatcher {
	mysqlListen := make(chan schemas.WhiteTorrent)

	return &MysqlCatcher{
		mysqlListen: mysqlListen,
		config:      cfg,
	}
}
