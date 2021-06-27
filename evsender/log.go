package evsender

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

var events []interface{}
var eventsMux = sync.RWMutex{}

func Log(event interface{}) {
	eventsMux.Lock()
	defer func() {
		eventsMux.Unlock()
	}()
	events = append(events, event)

	if len(events) >= batchLimit {
		go sendEvents(events)
		events = nil
	}
}

func sendEvents(events []interface{}) {
	if senderClient == nil {
		return
	}
	_, err := senderClient.SendBatch(events)
	if err != nil {
		log.Error(err)
		return
	}
}
