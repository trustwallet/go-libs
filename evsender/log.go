package evsender

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

var events map[string][]Event
var eventsMux = sync.RWMutex{}

func Log(event Event) {
	eventsMux.Lock()
	defer func() {
		eventsMux.Unlock()
	}()
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
