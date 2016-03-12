package server


// TODO(ian): Finish crafting a response.
type AnnounceResponseFailure struct {
	failure string
}

type AnnounceResponse struct {
	interval int // Interval in seconds a client should wait |.| messages
	trackerId string
	complete	uint
	incomplete	uint
	peers PeerList
}
