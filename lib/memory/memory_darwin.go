package memory

import (
	"log"
)

// This has been adapted from github.com/pbnjay/memory.
func sysTotalMemory() int {
	s, err := sysctlUint64("hw.memsize")
	if err != nil {
		log.Panicf("FATAL: cannot determine system memory: %s", err)
	}
	return int(s)
}
