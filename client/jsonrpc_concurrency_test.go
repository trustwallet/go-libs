package client

import (
	"sync"
	"testing"
)

func TestConcurrentRPCRequestIDsAreUnique(t *testing.T) {
	const workers, requestsPerWorker = 32, 256
	start := make(chan struct{})
	ids := make(chan int64, workers*requestsPerWorker)
	var workersDone sync.WaitGroup
	workersDone.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func() {
			defer workersDone.Done()
			batch := make(RpcRequests, requestsPerWorker)
			for i := range batch {
				batch[i] = &RpcRequest{Method: "test"}
			}
			<-start
			for _, request := range batch.fillDefaultValues() {
				ids <- request.Id
			}
		}()
	}
	close(start)
	workersDone.Wait()
	close(ids)

	seen := make(map[int64]struct{}, workers*requestsPerWorker)
	for id := range ids {
		seen[id] = struct{}{}
	}
	if len(seen) != workers*requestsPerWorker {
		t.Fatalf("got %d unique request IDs, want %d", len(seen), workers*requestsPerWorker)
	}
}
