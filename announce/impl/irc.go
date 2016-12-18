package catcherImpl

import (
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/belak/irc"
	"log"
	"net"
)

type IRCNotifier struct {
	killChan       chan bool
	client         *irc.Client
	newTorrentChan chan newTorrent
	config         config.ConfigStruct
}

func SpawnNotifier(cfg config.ConfigStruct) *IRCNotifier {
	killChan := make(chan bool)
	newTorrentChan := make(chan newTorrent)
	client := createIRCHandler(cfg)

	ircNotify := &IRCNotifier{
		killChan:       killChan,
		client:         client,
		newTorrentChan: newTorrentChan,
		config:         cfg,
	}

	go ircNotify.notifyCatcher()

	return ircNotify
}

func (irc *IRCNotifier) KillNotifier() error {
	irc.killChan <- true
	return nil
}

func (irc *IRCNotifier) notifyCatcher() {
	select {
	case <-irc.killChan:
		log.Println("Received kill notification in IRCNotifer")
		return
	case newTorrent := <-irc.newTorrentChan:
		log.Println("New Torrent added.")
		irc.sendNotification(newTorrent)
	}
}

func (irc *IRCNotifier) sendNotification(torrent newTorrent) {
	irc.client.Write(fmt.Sprintf("PRIVMSG %s %s\r\n", (*irc.config.IRCCfg).Chan, formatIRCMessage(torrent)))
}

func formatIRCMessage(torrent newTorrent) string {
	return fmt.Sprintf(
		"New_torrent_added._Name:_%s_InfoHash:_%s",
		torrent.Name,
		torrent.InfoHash,
	)
}

func createIRCHandler(cfg config.ConfigStruct) *irc.Client {
	conn, err := net.Dial("tcp", createServerPort(cfg))
	if err != nil {
		panic("Failed to connect to remote IRC server")
	}

	ircConfig := irc.ClientConfig{
		Nick:    (*cfg.IRCCfg).Nick,
		Pass:    (*cfg.IRCCfg).Pass,
		User:    (*cfg.IRCCfg).Nick,
		Name:    (*cfg.IRCCfg).Nick,
		Handler: irc.HandlerFunc(ircHandler),
	}

	client := irc.NewClient(conn, ircConfig)
	go func() {
		err = client.Run()
		if err != nil {
			log.Panicf("%v", err)
			panic("Failed to start IRC client.")
		}
	}()

	return client
}

func ircHandler(c *irc.Client, m *irc.Message) {
	if m.Command == "001" {
		c.Write("JOIN #12765")
	}
}

func createServerPort(cfg config.ConfigStruct) string {
	return fmt.Sprintf(
		"%s:%d",
		(*cfg.IRCCfg).Server,
		(*cfg.IRCCfg).Port,
	)
}
