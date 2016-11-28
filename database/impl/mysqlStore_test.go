package sqlStoreImpl

import (
	. "github.com/GrappigPanda/notorious/database"
	"testing"
	"time"
)

var MYSQLSTORE MySQLStore

func TestInitMySQLStore(t *testing.T) {
	MYSQLSTORE = InitMySQLStore()
}

func TestOpenConnection(t *testing.T) {
	_, err := MYSQLSTORE.OpenConnection()
	if err != nil {
		t.Errorf("TestMySQLStore:TestOpenConnection() failed to open a connection with err %v", err)
	}
}

func TestPeerUpdate(t *testing.T) {
	expectedReturn := &PeerStats{
		Downloaded: 6,
		Uploaded:   21,
		Ip:         "127.0.0.1",
	}

	newPeer := &PeerStats{
		Downloaded: 1,
		Uploaded:   1,
		Ip:         "127.0.0.1",
	}

	MYSQLSTORE.dbPool.Save(&newPeer)

	peerUpdate := PeerTrackerDelta{
		Uploaded:   20,
		Downloaded: 5,
		IP:         "127.0.0.1",
		Event:      PEERUPDATE,
	}

	go func() { MYSQLSTORE.UpdateConsumer <- peerUpdate }()
	time.Sleep(5 * time.Millisecond)

	retval := &PeerStats{}
	MYSQLSTORE.dbPool.Where("Ip = ?", peerUpdate.IP).Find(&retval)

	if retval.Downloaded != expectedReturn.Downloaded {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Downloaded,
			retval.Downloaded)
	}

	if retval.Uploaded != expectedReturn.Uploaded {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Uploaded,
			retval.Uploaded)
	}

	if retval.Ip != expectedReturn.Ip {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Ip,
			retval.Ip)
	}
}
