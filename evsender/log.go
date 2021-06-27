package evsender

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var events map[string][]Event
var eventsMux = sync.RWMutex{}

func Log(event Event) {
	eventsMux.Lock()
	defer func() {
		eventsMux.Unlock()
	}()
	if event.CreatedAt == 0 {
		event.CreatedAt = time.Now().Unix()
	}
	events[event.Name] = append(events[event.Name], event)

	if len(events) >= batchLimit {
		go sendEvents(events)
		events = nil
	}
}

func sendEvents(events map[string][]Event) {
	if senderClient == nil {
		return
	}
	_, err := senderClient.SendBatch(events)
	if err != nil {
		log.Error(err)
		return
	}
}
