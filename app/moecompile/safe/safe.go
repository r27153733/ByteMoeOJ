//go:build linux

package safe

import (
	"os"
	"time"

	"golang.org/x/sys/unix"
)

const (
	MB            = 1024 * 1024
	compileChroot = false
)

func SetLimits() {
	cpuLimit := 50

	var lim unix.Rlimit
	lim.Max = uint64(cpuLimit)
	lim.Cur = uint64(cpuLimit)
	err := unix.Setrlimit(unix.RLIMIT_CPU, &lim)
	if err != nil {
		panic(err)
	}
	time.AfterFunc(time.Second*50, func() {
		os.Exit(-1)
	})

	lim.Max = 500 * MB
	lim.Cur = 500 * MB
	err = unix.Setrlimit(unix.RLIMIT_FSIZE, &lim)
	if err != nil {
		panic(err)
	}

	// 4GB
	lim.Max = MB << 12
	lim.Cur = MB << 12

	err = unix.Setrlimit(unix.RLIMIT_AS, &lim)
	if err != nil {
		panic(err)
	}
}

func SetUID(uid int) {
	if uid == -1 {
		return
	}
	// Set user and group ID
	err := unix.Setgid(uid)
	if err != nil {
		panic(err)
	}
	err = unix.Setuid(uid)
	if err != nil {
		panic(err)
	}
	err = unix.Setresuid(uid, uid, uid)
	if err != nil {
		panic(err)
	}
}
