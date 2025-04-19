package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var CachingInMemory *cache.Cache

func InitCache() {
	CachingInMemory = cache.New(2*time.Minute, 10*time.Minute)
}
