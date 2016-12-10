package announce

// Handled in this file is declaring a consistent API for consumers of new
// torrent additions.
//
// Essentially, whenever a new torrent is added, we want to announce that
// torrent to disparate services (think IRC, websockets, RSS feed).

// AnnouncerNotifier handles declaring the API for notifying different channels
// of nwe torrent additions.
type AnnouncerNotifier interface {
	// SpawnNotifier handles creating a new notifier. The notifier will live by
	// itself and not need any
	SpawnNotifier() *AnnounceNotifier
	// KillNotifier will handle cleanup and closing of the notifier. Necessary
	// for a clean exit for Notorious.
	KillNotifier() error
	// sendNotification Ought to be private, as `SpawnNotifier` spins up a
	// goroutine to handle sendNotifications.
	sendNotification()
}
