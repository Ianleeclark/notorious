package postgres

import (
	"database/sql"
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/lib/pq"
	"time"
)

/// postgres_libpq differs from postgres.go because this file uses the lib/pq
//API instead of the gorm API. Mostly this is used for listening to `pg_notify`
//events.
///

type callbackFunction func(*pq.Notification)

type PGListener struct {
	listener   *pq.Listener
	connstring string
	conn       *sql.DB
	killListen chan bool
}

func openLibPQConnection(connstring string) (*sql.DB, error) {
	return sql.Open("postgres", connstring)
}

func reportErr(event pq.ListenerEventType, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func NewListener(c config.ConfigStruct) (*PGListener, error) {
	connstring := formatConnectString(c)

	conn, err := openLibPQConnection(connstring)
	if err != nil {
		return nil, fmt.Errorf("Error encountered when opening a connection for a new listener: %v", err.Error())
	}

	listener := pq.NewListener(connstring, 10*time.Second, time.Minute, reportErr)

	err = listener.Listen("new_torrent_added")
	if err != nil {
		return nil, fmt.Errorf("Error opening a listen handle: %v", err.Error())
	}

	killListen := make(chan bool)

	listenObj := PGListener{
		listener:   listener,
		connstring: connstring,
		conn:       conn,
		killListen: killListen,
	}

	return &listenObj, nil
}

// DO we want a callback or something else?
func (pg *PGListener) BeginListen(callback callbackFunction) {
	select {
	case <-pg.killListen:
		return
	case notification := <-pg.listener.Notify:
		println("calling callback")
		callback(notification)
		println("callback called")
	}
}
