//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type IErrorInfo struct {
	vtbl *IErrorInfoVtbl
}

// IErrorInfoVtbl returns information about an error in addition to the return code.
// It returns the error message, name of the component and GUID of the interface in
// which the error occurred, and the name and topic of the Help file that applies to the error.
// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/ms723041(v=vs.85)
type IErrorInfoVtbl struct {
	// QueryInterface Retrieves pointers to the supported interfaces on an object.
	QueryInterface uintptr
	// AddRef Increments the reference count for an interface pointer to a COM object.
	// You should call this method whenever you make a copy of an interface pointer.
	AddRef uintptr
	// Release Decrements the reference count for an interface on a COM object.
	Release uintptr
	// GetDescription Returns a text description of the error
	GetDescription uintptr
	// GetGUID Returns the GUID of the interface that defined the error.
	GetGUID uintptr
	// GetHelpContext Returns the Help context ID for the error.
	GetHelpContext uintptr
	// GetHelpFile Returns the path of the Help file that describes the error.
	GetHelpFile uintptr
	// GetSource Returns the name of the component that generated the error, such as "ODBC driver-name".
	GetSource uintptr
}

func (obj *IErrorInfo) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into IErrorInfo.QueryInterface()...")
	var ppvObject unsafe.Pointer
	err := NewHResultChecker("IErrorInfo::QueryInterface").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)),
		uintptr(unsafe.Pointer(ppvObject)),
	))
	if err != nil {
		return nil, err
	}
	return ppvObject, nil
}

func (obj *IErrorInfo) AddRef() (uint32, error) {
	debugPrint("Entering into IErrorInfo.AddRef()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the IErrorInfo::AddRef method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

func (obj *IErrorInfo) Release() (uint32, error) {
	debugPrint("Entering into IErrorInfo.Release()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the IErrorInfo::Release method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

// GetDescription Returns a text description of the error.
// HRESULT GetDescription (
//
//	BSTR *pbstrDescription);
//
// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/ms714318(v=vs.85)
func (obj *IErrorInfo) GetDescription() (*string, error) {
	debugPrint("Entering into IErrorInfo.GetDescription()...")
	var pbstrDescription *string
	err := NewHResultChecker("IErrorInfo::QueryInterface").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetDescription,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&pbstrDescription)),
		),
	)
	if err != nil {
		return nil, err
	}
	return pbstrDescription, nil
}

// GetGUID Returns the globally unique identifier (GUID) of the interface that defined the error.
// HRESULT GetGUID(
//
//	GUID *pGUID
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-ierrorinfo-getguid
func (obj *IErrorInfo) GetGUID() (*windows.GUID, error) {
	debugPrint("Entering into ierrorinfo.GetGUID()...")
	var pGUID *windows.GUID
	err := NewHResultChecker("IErrorInfo::QueryInterface").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetGUID,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(pGUID)),
		),
	)
	if err != nil {
		return nil, err
	}
	return pGUID, nil
}

// GetErrorInfo Obtains the error information pointer set by the previous call to SetErrorInfo in the current logical thread.
// HRESULT GetErrorInfo(
//
//	ULONG      dwReserved,
//	IErrorInfo **pperrinfo
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-geterrorinfo
func GetErrorInfo() (*IErrorInfo, error) {
	debugPrint("Entering into ierrorinfo.GetErrorInfo()...")
	var pperrinfo *IErrorInfo
	err := NewHResultChecker("IErrorInfo::QueryInterface").CheckHResultError(
		modOleAut.MustFindProc("GetErrorInfo").Call(0, uintptr(unsafe.Pointer(&pperrinfo))),
	)
	if err != nil {
		return nil, err
	}
	return pperrinfo, nil
}
