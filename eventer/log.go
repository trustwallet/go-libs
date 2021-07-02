package eventer

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var events []Event
var eventsMux = sync.RWMutex{}

func Log(event Event) {
	eventsMux.Lock()
	defer func() {
		eventsMux.Unlock()
	}()
	if event.CreatedAt == 0 {
		event.CreatedAt = time.Now().Unix()
	}
	events = append(events, event)

	if len(events) >= batchLimit {
		go sendEvents(events)
		events = nil
	}
}

func sendEvents(events []Event) {
	if senderClient == nil {
		return
	}
	_, err := senderClient.SendBatch(events)
	if err != nil {
		log.Error(err)
		return
	}
}
