package uniq

import "sync"

// NewCtr returns a ready-to-use deterministic uniqueness counter
func NewCtr() *Ctr {
	return &Ctr{m: make(map[string]int)}
}

// Ctr is like Filter, but associates a count with each key. The operations
// are Inc and Dec. The Inc operations increments the counter and returns
// the post-increment maximum value of any keys provided in the increment
// operation.
//
// The purpose of this structure is a lazy counter. Call Inc on a set 'k'
// and if the max is higher than some 'n', roll back the transaction by
// running it "in-reverse" with Dec.
type Ctr struct {
	sync.Mutex
	m map[string]int
	// TODO(as): use sys package
	_ [96]byte // cache line pad
}

// Inc increments the given keys by 1. It returns
// the maxima of the provided key set before the
// increment occurs. That is, if the input is ("a", "b")
// and the respective counts are (1, 7), it returns 7 and
// then updates the count to (2, 8).
//
// The transaction is atomic with respect to the input
// keys and all other operations on the counter.
func (f *Ctr) Inc(key ...string) (max int) {
	f.Lock()
	defer f.Unlock()
	for _, key := range key {
		f.m[key]++
	}
	return f.max(key...)
}

// Dec decrements the given keys by 1. It does not return anything.
func (f *Ctr) Dec(key ...string) {
	f.Lock()
	defer f.Unlock()
	for _, key := range key {
		f.m[key]--
	}
	return // return the min?
}

// max returns the maximum value of the keyset. The caller
// must be holding the lock on the counter.
func (f *Ctr) max(key ...string) (max int) {
	for _, key := range key {
		if x := f.m[key]; x > max {
			max = x
		}
	}
	return max
}
