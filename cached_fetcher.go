package main

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

func NewCachedRateFetcher(f RateFetcher, c *cache.Cache) RateFetcher {
	return cachedRateFetcher{f, c}
}

type cachedRateFetcher struct {
	rateFetcher RateFetcher
	cache       *cache.Cache
}

func (f cachedRateFetcher) FetchExchangeRate(date, base, other string) (float64, error) {
	cacheKey := fmt.Sprintf("%s:%s:%s", date, base, other)

	cached, found := f.cache.Get(cacheKey)
	if found {
		return cached.(float64), nil
	}

	rate, err := f.rateFetcher.FetchExchangeRate(date, base, other)
	if err != nil {
		return 0, err
	}

	f.cache.Set(cacheKey, rate, cache.DefaultExpiration)

	return rate, nil
}
