//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ISupportErrorInfo Ensures that error information can be propagated up the call chain correctly.
// Automation objects that use the error handling interfaces must implement ISupportErrorInfo
// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/nn-oaidl-isupporterrorinfo
type ISupportErrorInfo struct {
	vtbl *ISupportErrorInfoVtbl
}

type ISupportErrorInfoVtbl struct {
	// QueryInterface Retrieves pointers to the supported interfaces on an object.
	QueryInterface uintptr
	// AddRef Increments the reference count for an interface pointer to a COM object.
	// You should call this method whenever you make a copy of an interface pointer.
	AddRef uintptr
	// Release Decrements the reference count for an interface on a COM object.
	Release uintptr
	// InterfaceSupportsErrorInfo Indicates whether an interface supports the IErrorInfo interface.
	// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-isupporterrorinfo-interfacesupportserrorinfo
	InterfaceSupportsErrorInfo uintptr
}

// QueryInterface queries a COM object for a pointer to one of its interface;
// identifying the interface by a reference to its interface identifier (IID).
// If the COM object implements the interface, then it returns a pointer to that interface after calling IUnknown::AddRef on it.
// HRESULT QueryInterface(
//
//	REFIID riid,
//	void   **ppvObject
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-queryinterface(refiid_void)
func (obj *ISupportErrorInfo) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into ISupportErrorInfo.QueryInterface()...")
	var ppvObject unsafe.Pointer
	err := NewHResultChecker("ISupportErrorInfo::QueryInterface").CheckHResultSyscallError(syscall.SyscallN(
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

// AddRef Increments the reference count for an interface pointer to a COM object.
// You should call this method whenever you make a copy of an interface pointer
// ULONG AddRef();
// https://docs.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-addref
func (obj *ISupportErrorInfo) AddRef() (uint32, error) {
	debugPrint("Entering into ISupportErrorInfo.AddRef()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the ISupportErrorInfo::AddRef method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

// Release Decrements the reference count for an interface on a COM object.
// ULONG Release();
// https://docs.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-release
func (obj *ISupportErrorInfo) Release() (uint32, error) {
	debugPrint("Entering into ISupportErrorInfo.Release()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the ISupportErrorInfo::Release method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

// InterfaceSupportsErrorInfo
// HRESULT InterfaceSupportsErrorInfo(
//
//	REFIID riid
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-isupporterrorinfo-interfacesupportserrorinfo
func (obj *ISupportErrorInfo) InterfaceSupportsErrorInfo(riid windows.GUID) error {
	debugPrint("Entering into isupporterrorinfo.InterfaceSupportsErrorInfo()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.InterfaceSupportsErrorInfo,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the ISupportErrorInfo::InterfaceSupportsErrorInfo method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the ISupportErrorInfo::InterfaceSupportsErrorInfo method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}
