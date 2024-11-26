//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICLRRuntimeHost struct {
	vtbl *ICLRRuntimeHostVtbl
}

// ICLRRuntimeHostVtbl provides functionality similar to that of the ICorRuntimeHost interface
// provided in the .NET Framework version 1, with the following changes
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimehost-interface
type ICLRRuntimeHostVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	// Start Initializes the CLR into a process.
	Start uintptr
	// Stop Stops the execution of code by the runtime.
	Stop uintptr
	// SetHostControl sets the host control interface. You must call SetHostControl before calling Start.
	SetHostControl uintptr
	// GetCLRControl gets an interface pointer of type ICLRControl that hosts can use to customize
	// aspects of the common language runtime (CLR).
	GetCLRControl uintptr
	// UnloadAppDomain Unloads the AppDomain that corresponds to the specified numeric identifier.
	UnloadAppDomain uintptr
	// ExecuteInAppDomain Specifies the AppDomain in which to execute the specified managed code.
	ExecuteInAppDomain uintptr
	// GetCurrentAppDomainID gets the numeric identifier of the AppDomain that is currently executing.
	GetCurrentAppDomainId uintptr
	// ExecuteApplication used in manifest-based ClickOnce deployment scenarios to specify the application
	// to be activated in a new domain.
	ExecuteApplication uintptr
	// ExecuteInDefaultAppDomain Invokes the specified method of the specified type in the specified assembly.
	ExecuteInDefaultAppDomain uintptr
}

// GetICLRRuntimeHost is a wrapper function that takes an ICLRRuntimeInfo object and
// returns an ICLRRuntimeHost and loads it into the current process
func GetICLRRuntimeHost(runtimeInfo *ICLRRuntimeInfo) (*ICLRRuntimeHost, error) {
	debugPrint("Entering into iclrruntimehost.GetICLRRuntimeHost()...")
	runtimeHost, err := runtimeInfo.GetInterface(CLSID_CLRRuntimeHost, IID_ICLRRuntimeHost)
	if err != nil {
		return nil, err
	}
	err = runtimeHost.(*ICLRRuntimeHost).Start()
	return runtimeHost.(*ICLRRuntimeHost), err
}

func (obj *ICLRRuntimeHost) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into ICLRRuntimeHost.QueryInterface()...")
	var ppvObject unsafe.Pointer
	err := NewHResultChecker("ICLRRuntimeHost::QueryInterface").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.QueryInterface,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&riid)), // A reference to the interface identifier (IID) of the interface being queried for.
			uintptr(unsafe.Pointer(&ppvObject)),
		),
	)
	if err != nil {
		return nil, err
	}
	return ppvObject, nil
}

func (obj *ICLRRuntimeHost) AddRef() (uint32, error) {
	debugPrint("Entering into ICLRRuntimeHost.AddRef()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the ICLRRuntimeHost::AddRef method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

func (obj *ICLRRuntimeHost) Release() (uint32, error) {
	debugPrint("Entering into ICLRRuntimeHost.Release()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the ICLRRuntimeHost::Release method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

// Start Initializes the common language runtime (CLR) into a process.
// HRESULT Start();
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimehost-start-method
func (obj *ICLRRuntimeHost) Start() error {
	debugPrint("Entering into iclrruntimehost.Start()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Start,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		//return fmt.Errorf("the ICLRRuntimeHost::Start method returned an error:\r\n%s", err)
		debugPrint(fmt.Sprintf("the ICLRRuntimeHost::Start method returned an error:\r\n%s", err.Error()))
	}
	if hr != S_OK {
		return fmt.Errorf("the ICLRRuntimeHost::Start method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// ExecuteInDefaultAppDomain Calls the specified method of the specified type in the specified managed assembly.
// HRESULT ExecuteInDefaultAppDomain (
//
//	[in] LPCWSTR pwzAssemblyPath,
//	[in] LPCWSTR pwzTypeName,
//	[in] LPCWSTR pwzMethodName,
//	[in] LPCWSTR pwzArgument,
//
// [out] DWORD *pReturnValue
// );
// An LPCWSTR is a 32-bit pointer to a constant string of 16-bit Unicode characters, which MAY be null-terminated.
// Use syscall.UTF16PtrFromString to turn a string into a LPCWSTR
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimehost-executeindefaultappdomain-method
func (obj *ICLRRuntimeHost) ExecuteInDefaultAppDomain(pwzAssemblyPath, pwzTypeName, pwzMethodName, pwzArgument *uint16) (*uint32, error) {
	var pReturnValue *uint32
	err := NewHResultChecker("ICLRRuntimeHost::ExecuteInDefaultAppDomain").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.ExecuteInDefaultAppDomain,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(pwzAssemblyPath)),
			uintptr(unsafe.Pointer(pwzTypeName)),
			uintptr(unsafe.Pointer(pwzMethodName)),
			uintptr(unsafe.Pointer(pwzArgument)),
			uintptr(unsafe.Pointer(pReturnValue)),
		),
	)
	if err != nil {
		return nil, err
	}
	return pReturnValue, nil
}

// GetCurrentAppDomainID Gets the numeric identifier of the AppDomain that is currently executing.
// HRESULT GetCurrentAppDomainId(
//
//	[out] DWORD* pdwAppDomainId
//
// );
// https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/iclrruntimehost-getcurrentappdomainid-method
func (obj *ICLRRuntimeHost) GetCurrentAppDomainID() (uint32, error) {
	var pdwAppDomainId uint32
	err := NewHResultChecker("ICLRRuntimeHost::GetCurrentAppDomainID").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetCurrentAppDomainId,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&pdwAppDomainId)),
		),
	)
	if err != nil {
		return 0, err
	}
	return pdwAppDomainId, nil
}
