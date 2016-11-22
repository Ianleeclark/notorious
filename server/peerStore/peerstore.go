package peerStore

type PeerStore interface {
    SetKeyIfNotExists(string, string) bool
    SetKV(string, string) bool
    RemoveKV(string, string)
    KeyExists(string) bool
    GetKeyVal(string) []string
}
