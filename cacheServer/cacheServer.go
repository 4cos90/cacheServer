package cacheServer

import (
	"time"
)

type messageCache struct {
	receiver string
	message  string
	time     time.Time
}

type cache struct {
	Name    string
	Cache   []messageCache
	Timeout time.Duration
}

type CacheServer interface {
	Get(key string) []messageCache
	Set(receiver string, message string, time time.Time) error
	GetAll() []messageCache
}

func (s *cache) clearOutTimeCache() {
	for startindex := 0; startindex < len(s.Cache); startindex++ {
		if time.Now().Sub(s.Cache[startindex].time) < s.Timeout {
			s.Cache = s.Cache[startindex:]
			break
		} else if startindex == (len(s.Cache) - 1) {
			s.Cache = make([]messageCache, 0)
		}
	}
}

func (s *cache) GetAll() []messageCache {
	s.clearOutTimeCache()
	return s.Cache
}

func (s *cache) Get(key string) []messageCache {
	s.clearOutTimeCache()
	rlt := make([]messageCache, 1)
	newCache := make([]messageCache, 1)
	for i := 0; i < len(s.Cache); i++ {
		if s.Cache[i].receiver == key {
			rlt = append(rlt, s.Cache[i])
		} else {
			newCache = append(newCache, s.Cache[i])
		}
	}
	s.Cache = newCache
	return rlt
}

func (s *cache) Set(receiver string, message string, sendtime time.Time) error {
	s.clearOutTimeCache()
	value := messageCache{
		receiver: receiver,
		message:  message,
		time:     sendtime,
	}
	s.Cache = append(s.Cache, value)
	return nil
}

func NewCache(name string, timeout time.Duration) CacheServer {
	return &cache{
		Name:    name,
		Cache:   make([]messageCache, 1),
		Timeout: timeout,
	}
}
