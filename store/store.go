package store

import "sync"

type Store interface {
	Keys() []string
	Put(k, v string)
	Get(k string) string
	Delete(k string)
}

func New() (Store, error) {
	return &memoryStore{
		container: map[string]string{},
	}, nil
}

type memoryStore struct {
	sync.RWMutex

	container map[string]string
}

func (ms *memoryStore) Keys() []string {
	ms.RLock()
	defer ms.RUnlock()
	keys := []string{}
	for k, _ := range ms.container {
		keys = append(keys, k)
	}
	return keys
}

func (ms *memoryStore) Put(k, v string) {
	ms.Lock()
	defer ms.Unlock()
	ms.container[k] = v
}

func (ms *memoryStore) Get(k string) string {
	ms.RLock()
	defer ms.RUnlock()
	return ms.container[k]
}

func (ms *memoryStore) Delete(k string) {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.container, k)
}
