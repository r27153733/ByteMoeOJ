//go:build linux

package main

import (
	"github.com/r27153733/ByteMoeOJ/app/moecompile/safe"
)

func init() {
	safe.SetLimits()
	safe.SetUID(*safeUID)
}
