package sql

import "sync"

type paramsResult [][]any

func (p paramsResult) release() {
	releaseResultParams(&p)
}

var paramsResultPool = sync.Pool{
	New: func() any { r := make(paramsResult, 0, 10); return &r },
}

func acquireResultParams() paramsResult {
	return *paramsResultPool.Get().(*paramsResult)
}

func releaseResultParams(p *paramsResult) {
	if p != nil {
		newParam := (*p)[:0]
		paramsResultPool.Put(&newParam)
	}
}
