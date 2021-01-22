package hive

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	w := NewWorker()
	loopSize := 10
	count := 0
	callbackCount := 0
	var wg sync.WaitGroup
	for i := 0; i < loopSize; i++ {
		wg.Add(1)
		w.Submit(func() {
			count++
		}, func() {
			defer wg.Done()
			callbackCount++
		})
	}
	wg.Wait()
	assert.Equal(t, loopSize, count)
	assert.Equal(t, loopSize, callbackCount)
}

func TestWorkerClose(t *testing.T) {
	w := NewWorker()
	running := true
	w.Submit(func() {
		time.Sleep(time.Millisecond * 100)
	}, func() {
		running = false
	})
	w.Close()
	time.Sleep(time.Millisecond * 200)
	assert.False(t, running)
}
