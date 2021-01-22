package hive

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
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
	var finished atomic.Value
	w.Submit(func() {
		time.Sleep(time.Millisecond * 100)
	}, func() {
		finished.Store(true)
	})
	w.Close()
	time.Sleep(time.Millisecond * 200)
	assert.True(t, finished.Load().(bool))
}
