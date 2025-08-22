package main

import (
	"strconv"
	"sync"
	"time"
)

type item struct {
	data string
	ttl  time.Time
}

type Store struct {
	items map[string]item
	mu    *sync.Mutex
}

func NewDataStore() (error, *Store) {
	return nil, &Store{
		items: make(map[string]item),
		mu:    &sync.Mutex{},
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.items[key]; ok {
		entry.data = value
		entry.ttl = time.Now().Add(2 * time.Hour)
		s.items[key] = entry
	} else {
		newItem := item{
			data: value,
			ttl:  time.Now().Add(2 * time.Hour),
		}
		s.items[key] = newItem
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.items[key]

	if !ok {
		return "", false
	} else {
		if val.ttl.Before(time.Now()) {
			return "", false
		} else {
			return val.data, true
		}
	}
}

func (s *Store) SetWithTimeOut(key, value, expiredTime string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// First, convert the string to an integer.
	ms, err := strconv.ParseInt(expiredTime, 10, 64)
	if err != nil {
		// For now, we'll ignore the error, but in a real-world scenario,
		// you'd want to return an error to the client.
		return
	}

	// Next, create a time.Duration from the milliseconds.
	duration := time.Duration(ms) * time.Millisecond
	ttl := time.Now().Add(duration)

	if entry, ok := s.items[key]; ok {
		entry.data = value
		entry.ttl = ttl
		s.items[key] = entry
	} else {
		newItem := item{
			data: value,
			ttl:  ttl,
		}
		s.items[key] = newItem
	}
}
