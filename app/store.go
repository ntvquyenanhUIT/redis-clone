package main

import "sync"

type Store struct {
	data map[string]string
	mu   *sync.Mutex
}

func NewDataStore() (error, *Store) {
	return nil, &Store{
		data: make(map[string]string),
		mu:   &sync.Mutex{},
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]

	if !ok {
		return "", false
	} else {
		return val, true
	}
}
