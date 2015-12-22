package i18n

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net/http"
	"sync"
	"unsafe"

	"golang.org/x/text/language"
)

func shardGenerate(r *http.Request) string {
	uintVal := uint64(uintptr(unsafe.Pointer(r)))
	hasher := sha1.New()
	binary.Write(hasher, binary.LittleEndian, uintVal)
	return fmt.Sprintf("%x", hasher.Sum(nil))[0:2]
}

type requestShard struct {
	lock sync.RWMutex
	data map[*http.Request]language.Tag
}

type requestLanguageMap struct {
	lock sync.RWMutex
	data map[string]*requestShard
}

func newRequestLanguageMap() *requestLanguageMap {
	return &requestLanguageMap{
		data: make(map[string]*requestShard),
	}
}

func (rMap *requestLanguageMap) getShard(request *http.Request) *requestShard {
	key := shardGenerate(request)
	rMap.lock.RLock()
	shard, ok := rMap.data[key]
	rMap.lock.RUnlock()

	if !ok || shard == nil {
		rMap.lock.Lock()
		shard = &requestShard{
			data: make(map[*http.Request]language.Tag),
		}

		rMap.data[key] = shard
		rMap.lock.Unlock()
	}

	return shard
}

func (rMap *requestLanguageMap) Add(request *http.Request, tag language.Tag) {
	shard := rMap.getShard(request)
	shard.lock.Lock()
	defer shard.lock.Unlock()
	shard.data[request] = tag
}

func (rMap *requestLanguageMap) Delete(request *http.Request) {
	shard := rMap.getShard(request)
	shard.lock.Lock()
	defer shard.lock.Unlock()
	delete(shard.data, request)
}

func (rMap *requestLanguageMap) Get(request *http.Request) (language.Tag, bool) {
	shard := rMap.getShard(request)
	shard.lock.Lock()
	defer shard.lock.Unlock()
	lang, ok := shard.data[request]
	return lang, ok
}
