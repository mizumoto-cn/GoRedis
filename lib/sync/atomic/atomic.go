package atomic

import "sync/atomic"

// Here we wrap some atomic functions

// Define Boolean an uint32 for atomic operations
type Boolean uint32

// Get() returns the atomic value
func (b *Boolean) Get() bool {
	return atomic.LoadUint32((*uint32)(b)) == 1
}

// Set() sets the atomic value
func (b *Boolean) Set(value bool) {
	if value {
		atomic.StoreUint32((*uint32)(b), 1)
	} else {
		atomic.StoreUint32((*uint32)(b), 0)
	}
}
