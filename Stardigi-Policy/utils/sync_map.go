package utils

import (
	"sync"
	"sync/atomic"
)

type SyncMap struct {
	RwMutex       *sync.RWMutex
	ConcurrentMap map[string]interface{}
}

func NewSyncMap() *SyncMap {

	rwMutex := new(sync.RWMutex)
	concurrentMap := make(map[string]interface{}, 0)

	return &SyncMap{
		RwMutex:       rwMutex,
		ConcurrentMap: concurrentMap,
	}
}

func (syncMap *SyncMap) Set(key string, value interface{}) {

	syncMap.RwMutex.Lock()
	syncMap.ConcurrentMap[key] = value
	syncMap.RwMutex.Unlock()
}

func (syncMap *SyncMap) Get(key string) (value interface{}, ok bool) {

	syncMap.RwMutex.RLock()
	value, ok = syncMap.ConcurrentMap[key]
	syncMap.RwMutex.RUnlock()
	return
}

func (syncMap *SyncMap) Del(key string) {

	syncMap.RwMutex.Lock()
	delete(syncMap.ConcurrentMap, key)
	syncMap.RwMutex.Unlock()
}

/**
 * 老化value为false的字段
 */
func (syncMap *SyncMap) Age() {

	for k, v := range syncMap.ConcurrentMap {
		if !v.(bool) {
			delete(syncMap.ConcurrentMap, k)
		}
	}
}

/**
 * bkdr哈希算法，对字符串的每个字符转成unit64位，并和hash*seed相加得到
 * 哈希算法是根据提供的key生成固定的为一个较短的整型
 */
func BkdrHash(key string) uint64 {

	var seed uint64 = 13131
	var hashCode uint64

	for _, k := range key {
		kUint := uint64(k)
		hashCode = atomic.AddUint64(&kUint, hashCode*seed)
	}

	return hashCode
}
