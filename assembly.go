//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// from mscorlib.tlh

type Assembly struct {
	vtbl *AssemblyVtbl
}

// AssemblyVtbl is a COM virtual table of functions for the Assembly Class
// https://docs.microsoft.com/en-us/dotnet/api/system.reflection.assembly?view=netframework-4.8
type AssemblyVtbl struct {
	QueryInterface              uintptr
	AddRef                      uintptr
	Release                     uintptr
	GetTypeInfoCount            uintptr
	GetTypeInfo                 uintptr
	GetIDsOfNames               uintptr
	Invoke                      uintptr
	get_ToString                uintptr
	Equals                      uintptr
	GetHashCode                 uintptr
	GetType                     uintptr
	get_CodeBase                uintptr
	get_EscapedCodeBase         uintptr
	GetName                     uintptr
	GetName_2                   uintptr
	get_FullName                uintptr
	get_EntryPoint              uintptr
	GetType_2                   uintptr
	GetType_3                   uintptr
	GetExportedTypes            uintptr
	GetTypes                    uintptr
	GetManifestResourceStream   uintptr
	GetManifestResourceStream_2 uintptr
	GetFile                     uintptr
	GetFiles                    uintptr
	GetFiles_2                  uintptr
	GetManifestResourceNames    uintptr
	GetManifestResourceInfo     uintptr
	get_Location                uintptr
	get_Evidence                uintptr
	GetCustomAttributes         uintptr
	GetCustomAttributes_2       uintptr
	IsDefined                   uintptr
	GetObjectData               uintptr
	add_ModuleResolve           uintptr
	remove_ModuleResolve        uintptr
	GetType_4                   uintptr
	GetSatelliteAssembly        uintptr
	GetSatelliteAssembly_2      uintptr
	LoadModule                  uintptr
	LoadModule_2                uintptr
	CreateInstance              uintptr
	CreateInstance_2            uintptr
	CreateInstance_3            uintptr
	GetLoadedModules            uintptr
	GetLoadedModules_2          uintptr
	GetModules                  uintptr
	GetModules_2                uintptr
	GetModule                   uintptr
	GetReferencedAssemblies     uintptr
	get_GlobalAssemblyCache     uintptr
}

func (obj *Assembly) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into Assembly.QueryInterface()...")
	var ppvObject unsafe.Pointer
	hr, _, err := syscall.SyscallN(
		obj.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)),
		uintptr(ppvObject),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the Assembly::QueryInterface method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the Assembly::QueryInterface method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return ppvObject, nil
}

func (obj *Assembly) AddRef() error {
	debugPrint("Entering into Assembly.AddRef()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the Assembly::AddRef method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the Assembly::AddRef method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func (obj *Assembly) Release() error {
	debugPrint("Entering into Assembly.Release()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the Assembly::Release method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the Assembly::Release method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// GetEntryPoint returns the assembly's MethodInfo
//
//	 virtual HRESULT __stdcall get_EntryPoint (
//	/*[out,retval]*/ struct _MethodInfo * * pRetVal ) = 0;
//
// https://docs.microsoft.com/en-us/dotnet/api/system.reflection.assembly.entrypoint?view=netframework-4.8#System_Reflection_Assembly_EntryPoint
// https://docs.microsoft.com/en-us/dotnet/api/system.reflection.methodinfo?view=netframework-4.8
func (obj *Assembly) GetEntryPoint() (pRetVal *MethodInfo, err error) {
	debugPrint("Entering into Assembly.GetEntryPoint()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.get_EntryPoint,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pRetVal)),
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the Assembly::GetEntryPoint method returned an error:\r\n%s", err)
		return
	}
	if hr != S_OK {
		err = fmt.Errorf("the Assembly::GetEntryPoint method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	return
}

func (obj *Assembly) GetFullName() (string, error) {
	debugPrint("Entering into Assembly.GetFullName()...")
	var err error
	var pRetValBSTR unsafe.Pointer
	hr, _, err := syscall.SyscallN(
		obj.vtbl.get_FullName,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pRetValBSTR)),
	)
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the Assembly::GetFullName method returned an error:\r\n%s", err)
		return "", err
	}
	if hr != S_OK {
		err = fmt.Errorf("the Assembly::GetFullName method returned a non-zero HRESULT: 0x%x", hr)
		return "", err
	}
	return ReadUnicodeStr(pRetValBSTR), nil
}

func (obj *Assembly) CreateInstance(instance string) (unsafe.Pointer, error) {
	debugPrint("Entering into Assembly.CreateInstance()...")
	var instanceObj unsafe.Pointer
	instancePtr, err := SysAllocString(instance)
	if err != nil {
		return nil, fmt.Errorf("the Assembly::CreateInstance SysAllocString returned error: %v", err)
	}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.CreateInstance,
		uintptr(unsafe.Pointer(obj)),
		uintptr(instancePtr),
		uintptr(instanceObj),
	)
	errFree := SysFreeString(instancePtr)
	if errFree != nil {
		return nil, fmt.Errorf("the Assembly::CreateInstance error free String:\r\n%v", errFree)
	}
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the Assembly::CreateInstance method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the Assembly::CreateInstance method returned a non-zero HRESULT: 0x%x", hr)
	}
	return instanceObj, nil
}

func (obj *Assembly) GetType(typeStr string) (*Type, error) {
	debugPrint("Entering into Assembly.GetType()...")
	var typePtr *Type
	typeStrPtr, err := SysAllocString(typeStr)
	if err != nil {
		return nil, fmt.Errorf("the Assembly::GetType SysAllocString returned error: %v", err)
	}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetType_2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(typeStrPtr),
		uintptr(unsafe.Pointer(&typePtr)),
	)
	errFree := SysFreeString(typeStrPtr)
	if errFree != nil {
		return nil, fmt.Errorf("the Assembly::GetType error free String:\r\n%v", errFree)
	}
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the Assembly::GetType method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the Assembly::GetType method returned a non-zero HRESULT: 0x%x", hr)
	}
	return typePtr, nil
}

func (obj *Assembly) GetTypes() ([]*Type, error) {
	debugPrint("Entering into Assembly.GetTypes()...")
	var err error
	var safeArray *SafeArray
	types := []*Type{}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetTypes,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&safeArray)),
	)
	if err != syscall.Errno(0) {
		return types, fmt.Errorf("the Assembly::GetTypes method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return types, fmt.Errorf("the Assembly::GetTypes method returned a non-zero HRESULT: 0x%x", hr)
	}
	return safeArrayToTypes(safeArray)
}

// ToString Obtains a string representation that includes the friendly name of the application domain and any context policies.
// https://docs.microsoft.com/en-us/dotnet/api/system.assembly.tostring?view=net-5.0#System_AppDomain_ToString
func (obj *Assembly) ToString() (string, error) {
	debugPrint("Entering into Assembly.ToString()...")
	var pDomain *string
	hr, _, err := syscall.SyscallN(
		obj.vtbl.get_ToString,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pDomain)),
	)

	if err != syscall.Errno(0) {
		return "", fmt.Errorf("the Assembly.ToString method retured an error:\r\n%s", err)
	}
	if hr != S_OK {
		return "", fmt.Errorf("the Assembly.ToString method returned a non-zero HRESULT: 0x%x", hr)
	}
	return ReadUnicodeStr(unsafe.Pointer(pDomain)), nil
}

func safeArrayToTypes(safeArray *SafeArray) ([]*Type, error) {
	debugPrint("Entering into assembly.SafeArrayToTypes()...")
	//get dimensions of array (should be 1 for this context always)
	types := []*Type{}
	d, err := SafeArrayGetDim(safeArray)
	if err != nil {
		return types, err
	}
	if *d != 1 {
		return types, fmt.Errorf("expected dimension of 1, got %d", d)
	}

	lbound, err := SafeArrayGetLBound(safeArray, *d)
	if err != nil {
		return types, err
	}

	ubound, err := SafeArrayGetUBound(safeArray, *d)
	if err != nil {
		return types, err
	}
	for i := lbound; i <= ubound; i++ {
		pType, err := SafeArrayGetElement(safeArray, i)
		if err != nil {
			continue
		}
		types = append(types, (*Type)(pType))
	}
	err = SafeArrayDestroy(safeArray)
	if err != nil {
		return types, err
	}
	return types, nil
}
