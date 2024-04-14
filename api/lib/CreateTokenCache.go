package lib

import(
	"time"
	"github.com/allegro/bigcache"
)

type TokenCache struct {
	token *bigcache.BigCache
}

func CreateTokenCache() (*TokenCache, error) {
	cache, cacheError := bigcache.NewBigCache(bigcache.DefaultConfig(1 * time.Hour))
	if cacheError != nil {
		return nil, cacheError
	}

	return &TokenCache{
		token: cache,
	}, nil
} 