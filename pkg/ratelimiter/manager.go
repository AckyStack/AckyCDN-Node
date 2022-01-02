package ratelimiter

import (
	"github.com/kelindar/binary"
	"sync"
	"time"

	"github.com/gofiber/storage/memory"
)

type item struct {
	currHits int
	prevHits int
	exp      uint64
}

type manager struct {
	pool   sync.Pool
	memory *memory.Storage
}

func newManager() *manager {
	// Create new storage handler
	manager := &manager{
		pool: sync.Pool{
			New: func() interface{} {
				return new(item)
			},
		},
	}
	// Fallback too memory storage
	manager.memory = memory.New()
	return manager
}

// acquire returns an *entry from the sync.Pool
func (m *manager) acquire() *item {
	return m.pool.Get().(*item)
}

// release and reset *entry to sync.Pool
func (m *manager) release(e *item) {
	e.prevHits = 0
	e.currHits = 0
	e.exp = 0
	m.pool.Put(e)
}

// get data from storage or memory
func (m *manager) get(key string) (it *item) {
	dataBin, _ := m.memory.Get(key)
	if dataBin == nil {
		it = m.acquire()
		return it
	}
	data := &item{}
	binary.Unmarshal(dataBin, data)
	return data
}

// get raw data from storage or memory
func (m *manager) getRaw(key string) (raw []byte) {
	raw, _ = m.memory.Get(key)
	return raw
}

// set data to storage or memory
func (m *manager) set(key string, it *item, exp time.Duration) {
	data, _ := binary.Marshal(it)
	m.memory.Set(key, data, exp)
}

// set data to storage or memory
func (m *manager) setRaw(key string, raw []byte, exp time.Duration) {
	m.memory.Set(key, raw, exp)
}

// delete data from storage or memory
func (m *manager) delete(key string) {
	m.memory.Delete(key)
}
