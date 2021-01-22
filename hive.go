package hive

import (
	"context"
	"errors"
	"github.com/joway/pond"
	"time"
)

const (
	DefaultSize        = 10
	DefaultNonblocking = false
	DefaultMinIdle     = 0
	DefaultMaxIdle     = 10
	DefaultMinIdleTime = time.Minute * 5
)

var (
	ErrOverload = errors.New("pool is overload")
)

type Hive struct {
	workers *pond.Pool

	/**
	Size of the goroutine pool.
	*/
	Size int
	/**
	If true, return ErrOverload when all workers is busy.
	Otherwise, block Submit method.
	*/
	Nonblocking bool
	/**
	The minimum size of the idle goroutines.
	*/
	MinIdle int
	/**
	The maximal size of the idle goroutines.
	*/
	MaxIdle int
	/**
	The minimum time that idle goroutine should be reserved.
	*/
	MinIdleTime time.Duration
}

type Option func(h *Hive)

func WithSize(size int) Option {
	return func(h *Hive) {
		h.Size = size
	}
}

func WithNonblocking(b bool) Option {
	return func(h *Hive) {
		h.Nonblocking = b
	}
}

func WithMinIdle(n int) Option {
	return func(h *Hive) {
		h.MinIdle = n
	}
}

func WithMaxIdle(n int) Option {
	return func(h *Hive) {
		h.MaxIdle = n
	}
}

func WithMinIdleTime(t time.Duration) Option {
	return func(h *Hive) {
		h.MinIdleTime = t
	}
}

func defaultHive() *Hive {
	return &Hive{
		Size:        DefaultSize,
		Nonblocking: DefaultNonblocking,
		MinIdle:     DefaultMinIdle,
		MaxIdle:     DefaultMaxIdle,
		MinIdleTime: DefaultMinIdleTime,
	}
}

func New(options ...Option) (*Hive, error) {
	h := defaultHive()
	for _, optFunc := range options {
		optFunc(h)
	}

	pConfig := pond.NewDefaultConfig()
	pConfig.MaxSize = h.Size
	pConfig.Nonblocking = h.Nonblocking
	pConfig.ObjectCreateFactory = func(ctx context.Context) (interface{}, error) {
		return NewWorker(), nil
	}
	pConfig.ObjectDestroyFactory = func(ctx context.Context, object interface{}) error {
		w := object.(*Worker)
		w.Close()
		return nil
	}
	workers, err := pond.New(pConfig)
	if err != nil {
		return nil, err
	}
	h.workers = workers

	return h, nil
}

func (h *Hive) Submit(ctx context.Context, task Task) error {
	object, err := h.workers.BorrowObject(ctx)
	if err != nil {
		if err == pond.ErrPoolExhausted {
			return ErrOverload
		}
		return err
	}

	worker := object.(*Worker)
	worker.Submit(task, func() {
		_ = h.workers.ReturnObject(ctx, object)
	})
	return nil
}

func (h *Hive) Close() error {
	return h.workers.Close(context.Background())
}
