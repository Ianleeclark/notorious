package peerStore

type PeerStore interface {
    SetKeyIfNotExists(string, string) bool
    SetKV(string, string) 
    RemoveKV(string, string)
    KeyExists(string) bool
    GetKeyVal(string) []string
    GetAllPeers(string) []string
}
