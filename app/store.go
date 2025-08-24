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
	items         map[string]RedisObject
	mu            *sync.Mutex
	waitingClient map[string][]chan string
}

func NewDataStore() (error, *Store) {
	return nil, &Store{
		items:         make(map[string]RedisObject),
		mu:            &sync.Mutex{},
		waitingClient: make(map[string][]chan string),
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

	if len(s.waitingClient[key]) > 0 {
		s.waitingClient[key][0] <- value
		close(s.waitingClient[key][0])
		s.waitingClient[key] = s.waitingClient[key][1:]
		return 1, nil
	}

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

func (s *Store) LPush(key, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.waitingClient[key]) > 0 {
		s.waitingClient[key][0] <- value
		close(s.waitingClient[key][0])
		s.waitingClient[key] = s.waitingClient[key][1:]
		// The value is consumed by the waiting client, not stored in the list.
		return 1, nil
	}

	obj, exists := s.items[key]

	if !exists {
		newList := NewDoublyLinkedList()
		newList.LPush(value)

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

	list.LPush(value)
	return list.len, nil
}

func (s *Store) LRange(key string, start, end int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, exists := s.items[key]
	if !exists {
		return []string{}, nil
	}

	list, ok := obj.value.(*DoublyLinkedList)

	if !ok {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	result := list.LRange(start, end)
	return result, nil
}

func (s *Store) LLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, exists := s.items[key]

	if !exists {
		return 0, nil
	}

	list, ok := obj.value.(*DoublyLinkedList)

	if !ok {
		return -1, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	result := list.Len()
	return result, nil

}

func (s *Store) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, exists := s.items[key]

	if !exists {
		return "", false, nil
	}

	list, ok := obj.value.(*DoublyLinkedList)

	if !ok {
		return "", false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	val, hasDeleted := list.LPop()

	return val, hasDeleted, nil

}

func (s *Store) BLPop(key string, timeout int) (string, error) {

	// try pop first, if pop successfully, just return as usual
	// else, create a sepearate channel for it and add to the map
	s.mu.Lock()

	obj, exists := s.items[key]

	if exists {
		list, ok := obj.value.(*DoublyLinkedList)

		if !ok {
			s.mu.Unlock()
			return "", fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}

		val, hasPopped := list.LPop()

		if hasPopped {
			s.mu.Unlock()
			return val, nil
		}
	}

	// no list, or no item to pop
	waiter := make(chan string, 1)
	s.waitingClient[key] = append(s.waitingClient[key], waiter)
	s.mu.Unlock()

	var timer <-chan time.Time
	if timeout > 0 {
		timer = time.After(time.Duration(timeout) * time.Second)
	}

	select {
	case val := <-waiter:
		return val, nil
	case <-timer:
		s.mu.Lock()
		queue := s.waitingClient[key]
		newQueue := make([]chan string, 0)
		for _, ch := range queue {
			if ch != waiter {
				newQueue = append(newQueue, ch)
			}
		}
		s.waitingClient[key] = newQueue
		s.mu.Unlock()

		return "", nil
	}

}
