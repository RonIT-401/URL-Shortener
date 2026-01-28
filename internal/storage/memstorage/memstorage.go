package memstorage

import "sync"

type MemStorage struct {
	mu sync.RWMutex
	urls map[string]string
}

func New() *MemStorage {
	return  &MemStorage{urls: make(map[string]string)}
}

// Сохраняем короткий id - длинный id
func (s *MemStorage) Save(id, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[id] = url
	return nil
}

func (s *MemStorage) Get(id string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[id]
	return url, ok, nil
}