package kast

import "sync"

var (
	poolUsedKey = sync.Pool{
		New: func() any {
			return make(map[string]struct{})
		},
	}
)

func getUsedKeyMap() map[string]struct{} {
	return poolUsedKey.Get().(map[string]struct{})
}

func putUsedKeyMap(m map[string]struct{}) {
	clear(m)
	poolUsedKey.Put(m)
}
