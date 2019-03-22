package domain

type LRU interface {
	Put(key string, value interface{})
	Get(key string) (interface{}, bool)
}