package jobs

import (
	"sync/atomic"
)

var readOnlyMode int32 // 0 = normal, 1 = read-only

func IsReadOnly() bool {
	return atomic.LoadInt32(&readOnlyMode) == 1
}

func SetReadOnly(readonly bool) {
	if readonly {
		atomic.StoreInt32(&readOnlyMode, 1)
	} else {
		atomic.StoreInt32(&readOnlyMode, 0)
	}
}
