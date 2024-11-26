package clr

import "syscall"

var (
	modOleAut   = syscall.MustLoadDLL("OleAut32.dll")
	modkernel32 = syscall.MustLoadDLL("kernel32.dll")
)
