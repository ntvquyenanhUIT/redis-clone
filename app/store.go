package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type RedisObject struct {
	value interface{}
	ttl   time.Time
}

type Store struct {
	items map[string]RedisObject
	mu    *sync.Mutex
}

func NewDataStore() (error, *Store) {
	return nil, &Store{
		items: make(map[string]RedisObject),
		mu:    &sync.Mutex{},
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.items[key]; ok {
		entry.value = value
		entry.ttl = time.Now().Add(2 * time.Hour)
		s.items[key] = entry
	} else {
		newItem := RedisObject{
			value: value,
			ttl:   time.Now().Add(2 * time.Hour),
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
			if strVal, ok := val.value.(string); ok {
				return strVal, true
			}
			return "", false
		}
	}
}

func (s *Store) SetWithTimeOut(key, value, expiredTime string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// First, convert the string to an integer.
	ms, err := strconv.ParseInt(expiredTime, 10, 64)
	if err != nil {
		return
	}

	duration := time.Duration(ms) * time.Millisecond
	ttl := time.Now().Add(duration)

	if entry, ok := s.items[key]; ok {
		entry.value = value
		entry.ttl = ttl
		s.items[key] = entry
	} else {
		newItem := RedisObject{
			value: value,
			ttl:   ttl,
		}
		s.items[key] = newItem
	}
}

func (s *Store) RPush(key, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	obj, exists := s.items[key]

	if !exists {
		newList := NewDoublyLinkedList()
		newList.RPush(value)

		s.items[key] = RedisObject{
			value: newList,
		}
		return newList.len, nil
	}

	list, ok := obj.value.(*DoublyLinkedList)
	if !ok {
		// The key exists but holds something else (like a string).
		// This is a protocol error, just like in real Redis.
		return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	list.RPush(value)
	return list.len, nil
}
