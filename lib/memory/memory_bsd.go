//go:build freebsd || openbsd || dragonfly || netbsd

package memory

import (
	"log"
)

// This code has been adopted from https://github.com/pbnjay/memory

func sysTotalMemory() int {
	s, err := sysctlUint64("hw.physmem")
	if err != nil {
		log.Panicf("FATAL: cannot determine system memory: %s", err)
	}
	return int(s)
}
