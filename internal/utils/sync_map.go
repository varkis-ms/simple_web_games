package utils

import (
	"simple_web_games/internal/games"
	"sync"
)

type SyncMap struct {
	mx sync.RWMutex
	m  map[string]games.GameField
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		mx: sync.RWMutex{},
		m:  make(map[string]games.GameField),
	}
}

func (c *SyncMap) Load(key string) (games.GameField, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *SyncMap) Store(key string, value games.GameField) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func (c *SyncMap) Delete(key string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.m, key)
}

func (c *SyncMap) ReadAll() ([]any, []any) {
	c.mx.Lock()
	defer c.mx.Unlock()
	keySlice := make([]any, len(c.m))
	valueSlice := make([]any, len(c.m))
	i := 0
	for k, v := range c.m {
		keySlice[i] = k
		valueSlice[i] = v
		i++
	}
	return keySlice, valueSlice
}
