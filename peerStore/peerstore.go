package peerStore

// PeerStore Represents the base implementation of how we store peers. We
// currently use Redis for storing peers, but this interface allows for
// extensions to other third-party systems.
type PeerStore interface {
	SetKeyIfNotExists(string, string) bool
	SetKV(string, string)
	RemoveKV(string, string)
	KeyExists(string) bool
	GetKeyVal(string) []string
	GetAllPeers(string) []string
	SetIPMember(string, string) int
	CreateNewTorrentKey(string)
}
