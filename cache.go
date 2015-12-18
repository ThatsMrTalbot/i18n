package i18n

import (
	"sync"

	"golang.org/x/text/language"
)

// Cache stores translations in a map for fast access
type Cache struct {
	lock  sync.RWMutex
	cache map[string]map[string]*Translation
}

// Clear translation cache
func (cache *Cache) Clear() {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	cache.cache = make(map[string]map[string]*Translation)
}

// Add translation to cache
func (cache *Cache) Add(translation *Translation) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.cache == nil {
		cache.cache = make(map[string]map[string]*Translation)
	}

	l := translation.Lang.String()

	if _, ok := cache.cache[l]; !ok {
		cache.cache[l] = make(map[string]*Translation)
	}

	cache.cache[l][translation.Key] = translation
}

// Get translation from cache
func (cache *Cache) Get(lang language.Tag, key string) *Translation {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	if cache.cache == nil {
		return nil
	}

	l := lang.String()

	if _, ok := cache.cache[l]; !ok {
		return nil
	}

	if _, ok := cache.cache[l][key]; !ok {
		return nil
	}

	return cache.cache[l][key]
}

// Delete translation from cache
func (cache *Cache) Delete(translation *Translation) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.cache == nil {
		return
	}

	l := translation.Lang.String()

	if _, ok := cache.cache[l]; !ok {
		return
	}

	if _, ok := cache.cache[l][translation.Key]; !ok {
		return
	}

	delete(cache.cache[l], translation.Key)
}
