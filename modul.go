package clr

import "syscall"

var (
	modOleAut   = syscall.MustLoadDLL("OleAut32.dll")
	modNtDll    = syscall.MustLoadDLL("ntdll.dll")
	modkernel32 = syscall.MustLoadDLL("kernel32.dll")
)
