// Package uniq implements a uniqueness filter. It supports atomic
// multi-key put and delete operations. There is currently no lookup
// operation (deemed unnecessary for its current purpose).
//
// TODO(as): Replace sync mutexes with downgradable locks
// and benchmark
package uniq

import "sync"

// New returns a ready-to-use deterministic uniqueness filter
func New() *Filter {
	return &Filter{m: map[string]struct{}{}}
}

// Unique is a non-counting deterministic existence
// filter. It supports atomic multi-key put and delete
// operations.
type Filter struct {
	sync.Mutex
	m map[string]struct{}
	// TODO(as): use sys package
	_ [96]byte // cache line pad
}

// PutAll either commits the entire list of keys to the
// filter, or none at all. It returns false if any of the
// provided keys exist in the filter.
func (f *Filter) PutAll(key ...string) bool {
	f.Lock()
	defer f.Unlock()
	for _, key := range key {
		if _, exist := f.m[key]; exist {
			return false
		}
	}
	for _, key := range key {
		f.m[key] = struct{}{}
	}
	return true
}

// DelAll removes the list of keys from the filter. It does
// not check for existence.
func (f *Filter) DelAll(key ...string) {
	f.Lock()
	defer f.Unlock()
	for _, key := range key {
		delete(f.m, key)
	}
}
