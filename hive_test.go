package hive

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestHive(t *testing.T) {
	h, err := New(
		WithSize(10),
	)
	assert.NoError(t, err)

	var last int32 = 0
	var wg sync.WaitGroup
	for i := int32(0); i < int32(h.Size); i++ {
		n := i
		wg.Add(1)
		err := h.Submit(context.Background(), func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * 100)
			atomic.StoreInt32(&last, n)
		})
		assert.NoError(t, err)
	}
	wg.Add(1)
	err = h.Submit(context.Background(), func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 200)
		atomic.StoreInt32(&last, -1)
	})
	assert.NoError(t, err)
	wg.Wait()
	assert.Equal(t, int32(-1), last)
}

func TestHiveWithNonblocking(t *testing.T) {
	h, err := New(
		WithSize(10),
		WithNonblocking(true),
	)
	assert.NoError(t, err)

	for i := int32(0); i < int32(h.Size); i++ {
		err := h.Submit(context.Background(), func() {
			time.Sleep(time.Millisecond * 100)
		})
		assert.NoError(t, err)
	}
	err = h.Submit(context.Background(), func() {
	})
	assert.Equal(t, ErrOverload, err)
}
