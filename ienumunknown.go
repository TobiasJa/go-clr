//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type IEnumUnknown struct {
	vtbl *IEnumUnknownVtbl
}

// IEnumUnknownVtbl Enumerates objects implementing the root COM interface, IUnknown.
// Commonly implemented by a component containing multiple objects. For more information, see IEnumUnknown.
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nn-objidl-ienumunknown
type IEnumUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	// Next Retrieves the specified number of items in the enumeration sequence.
	Next uintptr
	// Skip Skips over the specified number of items in the enumeration sequence.
	Skip uintptr
	// Reset Resets the enumeration sequence to the beginning.
	Reset uintptr
	// Clone Creates a new enumerator that contains the same enumeration state as the current one.
	Clone uintptr
}

func (obj *IEnumUnknown) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into IEnumUnknown.QueryInterface()...")
	var ppvObject unsafe.Pointer
	err := NewHResultChecker("IEnumUnknown::QueryInterface").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)), // A reference to the interface identifier (IID) of the interface being queried for.
		uintptr(unsafe.Pointer(&ppvObject)),
	))
	if err != nil {
		return nil, err
	}
	return ppvObject, nil
}

func (obj *IEnumUnknown) AddRef() (uint32, error) {
	debugPrint("Entering into IEnumUnknown.AddRef()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the IEnumUnknown::AddRef method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

func (obj *IEnumUnknown) Release() (count uint32, err error) {
	debugPrint("Entering into IEnumUnknown.Release()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the IEnumUnknown::Release method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

// Next retrieves the specified number of items in the enumeration sequence.
// HRESULT Next(
//
//	ULONG    celt,
//	IUnknown **rgelt,
//	ULONG    *pceltFetched
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nf-objidl-ienumunknown-next
func (obj *IEnumUnknown) Next(celt uint32, pEnumRuntime unsafe.Pointer, pceltFetched *uint32) (int, error) {
	debugPrint("Entering into ienumunknown.Next()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Next,
		uintptr(unsafe.Pointer(obj)),
		uintptr(celt),
		uintptr(pEnumRuntime),
		uintptr(unsafe.Pointer(pceltFetched)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("there was an error calling the IEnumUnknown::Next method:\r\n%s", err)
	}
	if hr != S_OK && hr != S_FALSE {
		return 0, fmt.Errorf("the IEnumUnknown::Next method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return int(hr), nil
}
