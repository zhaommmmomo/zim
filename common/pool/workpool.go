package pool

import (
	"github.com/panjf2000/ants/v2"
	"github.com/zhaommmmomo/zim/common/log"
)

type WorkPool struct {
	pool *ants.Pool
}

const (
	defaultWorkPoolSize = 1024
)

func DefaultWorkPool() *WorkPool {
	pool, err := ants.NewPool(defaultWorkPoolSize)
	if err != nil {
		log.Error("create work pool fail.", log.Int("workPoolSize", defaultWorkPoolSize), log.Err(err))
		return nil
	}
	return &WorkPool{pool: pool}
}

func NewWorkPool(size int) *WorkPool {
	if size <= 0 {
		size = defaultWorkPoolSize
	}
	pool, err := ants.NewPool(size)
	if err != nil {
		log.Error("create work pool fail.", log.Int("workPoolSize", size), log.Err(err))
		return nil
	}
	return &WorkPool{pool: pool}
}

func (wp *WorkPool) Submit(f func()) error {
	return wp.pool.Submit(f)
}
