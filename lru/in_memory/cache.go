package in_memory

import (
	"container/list"
	"gitlab.com/artilligence/http-db-saver/domain"
)

const storageElementsCount = 1000

type hashMap map[string]interface{}

type InMemoryLRU struct {
	storage hashMap
	list    list.List
}

func NewInMemoryLRU() domain.LRU {
	return &InMemoryLRU{
		storage: make(hashMap),
	}
}

func (lru *InMemoryLRU) Put(key string, value interface{}) {
	val, ok := lru.storage[key]
	l := lru.list

	// if element exist in cache, move it to beginning
	if ok {
		for e := l.Front(); e != nil; e = e.Next() {
			if val == e.Value {
				l.MoveBefore(e, l.Front())
				return
			}
		}
	}

	// if element not exist in cache and total elements limit reached
	// remove oldest element
	if !ok && l.Len() == storageElementsCount {
		back := l.Back()

		delete(lru.storage, key)
		l.Remove(back)
	}

	lru.storage[key] = value
	if l.Len() == 0 {
		l.PushBack(value)
	} else {
		l.InsertBefore(value, l.Front())
	}
}

func (lru *InMemoryLRU) Get(key string) (interface{}, bool) {
	value, ok := lru.storage[key]

	return value, ok
}


