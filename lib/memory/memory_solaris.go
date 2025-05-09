package memory

import (
	"golang.org/x/sys/unix"
	"log"
)

const PHYS_PAGES = 0x1f4

func sysTotalMemory() int {
	memPageSize := unix.Getpagesize()
	// https://man7.org/linux/man-pages/man3/sysconf.3.html
	// _SC_PHYS_PAGES
	memPagesCnt, err := unix.Sysconf(PHYS_PAGES)
	if err != nil {
		log.Panicf("FATAL: error in unix.Sysconf: %s", err)
	}

	return memPageSize * int(memPagesCnt)
}
