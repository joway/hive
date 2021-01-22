package hive

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"sync"
	"testing"
	"time"
)

var poolSize = 100
var loopSize = 100

func benchFunc() {
	time.Sleep(100 * time.Nanosecond)
}

func BenchmarkHive(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := New(WithSize(poolSize))
	defer p.Close()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(loopSize)
		for j := 0; j < loopSize; j++ {
			_ = p.Submit(context.Background(), func() {
				defer wg.Done()
				benchFunc()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkAnts(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(poolSize)
	defer p.Release()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(loopSize)
		for j := 0; j < loopSize; j++ {
			_ = p.Submit(func() {
				defer wg.Done()
				benchFunc()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}
