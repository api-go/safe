//go:build windows
// +build windows

package safe

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// LockMemory 锁定内存页 (Windows 需使用 VirtualLock)
func LockMemory(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	// Windows 内存锁定通常需要处理内存页对齐
	return windows.VirtualLock(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
}

// UnlockMemory 解锁内存页
func UnlockMemory(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return windows.VirtualUnlock(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
}

// DisableCoreDump 禁止进程核心转储
func DisableCoreDump() error {
	windows.SetErrorMode(windows.SEM_NOGPFAULTERRORBOX)
	return nil
}
