package w32

import (
	"syscall"
	"unsafe"
)

const (
	MutexName = "SingleInstance_GoFTP"
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex      = kernel32.NewProc("CreateMutexW")
	ERROR_ALREADY_EXISTS = 183
)

func CreateMutex(name string) (uintptr, error) {
	name16, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	ret, _, err := procCreateMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(name16)),
	)

	sysErr := int(err.(syscall.Errno))
	if sysErr == 0 {
		return ret, nil
	}

	return ret, err
}
