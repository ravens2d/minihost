package cache

import "github.com/redis/go-redis/v9"

type Cache interface {
}

type cache struct {
	cache *redis.Client
}

func New() (Cache, error) {
	return &cache{
		cache: redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
	}, nil
}
