package rss

import (
	"fmt"
	"github.com/GrappigPanda/notorious/announce/impl/newTorrentType"
	"github.com/GrappigPanda/notorious/config"
	"github.com/gorilla/feeds"
	"log"
	"time"
)

type RSSNotifier struct {
	feed           *feeds.Feed
	killChan       chan bool
	NewTorrentChan chan dataType.NewTorrent
	config         config.ConfigStruct
}

func SpawnNotifier(cfg config.ConfigStruct) *RSSNotifier {
	killChan := make(chan bool)
	NewTorrentChan := make(chan dataType.NewTorrent)

	rssFeed := &feeds.Feed{
		Title:       "Notorious Tracker Torrent Updates",
		Link:        &feeds.Link{Href: "http://localhost"},
		Description: "An RSS feed to notify of new torrents being added",
		Author:      &feeds.Author{Name: "Notorious Tracker", Email: "CONFIGURABLE"},
		Created:     time.Now(),
	}

	rssNotifier := &RSSNotifier{
		feed:           rssFeed,
		killChan:       killChan,
		NewTorrentChan: NewTorrentChan,
		config:         cfg,
	}

	go rssNotifier.notifyCatcher()

	return rssNotifier
}

func (rss *RSSNotifier) notifyCatcher() {
	select {
	case <-rss.killChan:
		log.Println("Received kill notification in IRCNotifer")
		return
	case newTorrent := <-rss.NewTorrentChan:
		rss.sendNotification(newTorrent)
	}
}

func (rss *RSSNotifier) KillNotifier() error {
	rss.killChan <- true
	return nil
}

func (rss *RSSNotifier) sendNotification(torrent dataType.NewTorrent) error {
	rss.feed.Items = append(
		rss.feed.Items,
		&feeds.Item{
			Title: torrent.Name,
			Link: &feeds.Link{
				Href: fmt.Sprintf("http://localhost?infoHash=%s", torrent.InfoHash),
			},
			Description: "New torrent added",
			Author:      &feeds.Author{Name: "Notorious Tracker", Email: "CONFIGURABLE"},
			Created:     time.Now(),
		},
	)

	_, err := rss.feed.ToRss()

	return err
}

func (rss *RSSNotifier) GetRSS() (string, error) {
	return rss.feed.ToRss()
}
