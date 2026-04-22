//go:build !windows
// +build !windows

package safe

import (
	"golang.org/x/sys/unix"
)

// LockMemory 锁定内存页，防止其被置换到硬盘 (Swap)
func LockMemory(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return unix.Mlock(buf)
}

// UnlockMemory 解锁内存页
func UnlockMemory(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return unix.Munlock(buf)
}

// DisableCoreDump 禁止进程核心转储
func DisableCoreDump() error {
	var rlimit unix.Rlimit
	rlimit.Cur = 0
	rlimit.Max = 0
	return unix.Setrlimit(unix.RLIMIT_CORE, &rlimit)
}
