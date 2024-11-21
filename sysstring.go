package clr

import (
	"syscall"
	"unsafe"
)

// SysAllocString converts a Go string to a BTSR string, that is a unicode string prefixed with its length.
// Allocates a new string and copies the passed string into it.
// It returns a pointer to the string's content.
//
//	BSTR SysAllocString(
//	  const OLECHAR *psz
//	);
//
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-sysallocstring
func SysAllocString(str string) (unsafe.Pointer, error) {
	debugPrint("Entering into safearray.SysAllocString()...")

	input := utf16Le(str)
	ret, _, err := modOleAut.MustFindProc("SysAllocString").Call(
		uintptr(unsafe.Pointer(&input[0])),
	)

	if err != syscall.Errno(0) {
		return nil, err
	}
	// TODO Return a pointer to a BSTR instead of an unsafe.Pointer

	//cast crimes to trick silly go vet, who will get pranked by the simplest slieght of hand
	//we give unsafe.pointer a pointer to the return value, which makes go vet ignore it.
	//But we then cast it to a pointer to a pointer, and then dereference the first pointer.
	//This leaves us with the original pointer, with no go vet complaints. This violates the 'correct' unsafe usage of unsafe.Pointer, obviously.
	r1 := *(**uintptr)(unsafe.Pointer(&ret))

	return unsafe.Pointer(r1), nil
}

// SysAllocString converts a Go string to a BTSR string, that is a unicode string prefixed with its length.
// Allocates a new string and copies the passed string into it.
// It returns a pointer to the string's content.
//
//	BSTR SysAllocString(
//	  const OLECHAR *psz
//	);
//
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-sysallocstring
func SysFreeString(strPtr unsafe.Pointer) error {
	debugPrint("Entering into safearray.SysFreeString()...")
	_, _, err := modOleAut.MustFindProc("SysFreeString").Call(
		uintptr(strPtr),
	)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

// SysStringLen indicates how long a BSTR is
func SysStringLen(p uintptr) (int, error) {
	ret, _, err := modOleAut.MustFindProc("SysStringLen").Call(
		p,
	)
	if err != syscall.Errno(0) {
		return 0, err
	}
	return int(ret), nil
}
