package reaper

import (
	"fmt"
	"github.com/GrappigPanda/notorious/peerStore/redis"
	"sync"
	"testing"
	"time"
)

func TestConvertTimeToUnixTimeStamp(t *testing.T) {
	_, err := convertTimeToUnixTimeStamp("1480308617")
	if err != nil {
		t.Fatalf("Invalid time conversion -- TestConvertTime")
	}
}

/*
 *func TestReapInfoHash(t *testing.T) {
 *    c := redisPeerStore.OpenClient()
 *
 *    infoHash := "TestInfoHash-Reaper"
 *    redisPeerStore.CreateNewTorrentKey(c, infoHash)
 *
 *    keymember := fmt.Sprintf("%s:complete", infoHash)
 *
 *    for i := 0; i <= 3; i++ {
 *        redisPeerStore.SetKeyVal(c, keymember, fmt.Sprintf("127.0.0.%d:5454", i))
 *    }
 *
 *    peersReaped := 0
 *
 *    peerReapCount := make(chan int)
 *    currTime := int64(time.Now().UTC().Unix()) + 99999999
 *
 *    go func() {
 *        select {
 *        case count := <-peerReapCount:
 *            peersReaped += count
 *        }
 *    }()
 *
 *    reapInfoHash(c, "TestInfoHash-Reaper:complete", peerReapCount, currTime)
 *    time.Sleep(50 * time.Millisecond)
 *
 *    if peersReaped != 4 {
 *        t.Fatalf("Expected 4 peers to be reaped, only %v reaped", peersReaped)
 *    }
 *}
 *
 */
func TestReapPeers(t *testing.T) {
	c := redisPeerStore.OpenClient()

	infoHash := "TestInfoHash-reapPeers"

	var wgList sync.WaitGroup
	keyMemberList := make([]string, 6000)
	for i := 0; i < cap(keyMemberList)-1; i++ {
		wgList.Add(1)
		go func(wg *sync.WaitGroup) {
			keyMemberList[i] = fmt.Sprintf("%s:complete", infoHash)
			wg.Done()
		}(&wgList)
	}
	wgList.Wait()

	var wg sync.WaitGroup
	for i := 0; i <= 150000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			redisPeerStore.SetKeyVal(c, fmt.Sprintf("%s:complete", infoHash), fmt.Sprintf("127.0.0.%d:5454", i))
			wg.Done()
		}(&wg)
	}
	wg.Wait()

	currTime := int64(time.Now().UTC().Unix()) + 99999999
	reapedPeers := reapPeers(currTime)
	time.Sleep(1 * time.Second)

	if reapedPeers != 150000 {
		t.Fatalf("Expected to reap 150000 peers, reaped %v", reapedPeers)
	}
}

func TestStartReapingSchedule(t *testing.T) {

}
