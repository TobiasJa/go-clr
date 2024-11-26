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
	GetToString                 uintptr
	Equals                      uintptr
	GetHashCode                 uintptr
	GetType                     uintptr
	get_CodeBase                uintptr
	get_EscapedCodeBase         uintptr
	GetName                     uintptr
	GetName_2                   uintptr
	GetFullName                 uintptr
	GetEntryPoint               uintptr
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
	err := NewHResultChecker("Assembly::QueryInterface").CheckHResultSyscallError(syscall.SyscallN(
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

func (obj *Assembly) AddRef() error {
	debugPrint("Entering into Assembly.AddRef()...")
	return NewHResultChecker("Assembly::AddRef").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	))
}

func (obj *Assembly) Release() error {
	debugPrint("Entering into Assembly.Release()...")
	return NewHResultChecker("Assembly::Release").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	))
}

// ToString Obtains a string representation that includes the friendly name of the application domain and any context policies.
// https://docs.microsoft.com/en-us/dotnet/api/system.assembly.tostring?view=net-5.0#System_AppDomain_ToString
func (obj *Assembly) ToString() (string, error) {
	debugPrint("Entering into Assembly.ToString()...")
	var pDomain *string
	err := NewHResultChecker("Assembly::GetToString").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetToString,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&pDomain)),
		),
	)

	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(pDomain)), nil
}

// GetEntryPoint returns the assembly's MethodInfo
//
//	 virtual HRESULT __stdcall get_EntryPoint (
//	/*[out,retval]*/ struct _MethodInfo * * pRetVal ) = 0;
//
// https://docs.microsoft.com/en-us/dotnet/api/system.reflection.assembly.entrypoint?view=netframework-4.8#System_Reflection_Assembly_EntryPoint
// https://docs.microsoft.com/en-us/dotnet/api/system.reflection.methodinfo?view=netframework-4.8
func (obj *Assembly) GetEntryPoint() (*MethodInfo, error) {
	debugPrint("Entering into Assembly.GetEntryPoint()...")
	var pRetVal *MethodInfo
	err := NewHResultChecker("Assembly::GetEntryPoint").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetEntryPoint,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&pRetVal)),
		),
	)
	if err != nil {
		return nil, err
	}
	return pRetVal, nil
}

// _bstr_t GetFullName ( );
func (obj *Assembly) GetFullName() (string, error) {
	debugPrint("Entering into Assembly.GetFullName()...")
	var pRetValBSTR unsafe.Pointer
	err := NewHResultChecker("Assembly::GetFullName").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetFullName,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&pRetValBSTR)),
		),
	)
	if err != nil {
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

func (obj *Assembly) GetType_2(typeStr string) (*Type, error) {
	debugPrint("Entering into Assembly.GetType_2()...")
	var typePtr *Type
	typeStrPtr, err := SysAllocString(typeStr)
	if err != nil {
		return nil, fmt.Errorf("the Assembly::GetType_2 SysAllocString returned error: %v", err)
	}
	err = NewHResultChecker("Assembly::GetType_2").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetType_2,
			uintptr(unsafe.Pointer(obj)),
			uintptr(typeStrPtr),
			uintptr(unsafe.Pointer(&typePtr)),
		),
	)
	errFree := SysFreeString(typeStrPtr)
	if errFree != nil {
		return nil, fmt.Errorf("the Assembly::GetType_2 error free String:\r\n%v", errFree)
	}
	if err != nil {
		return nil, err
	}
	return typePtr, nil
}

func (obj *Assembly) GetTypes() ([]*Type, error) {
	debugPrint("Entering into Assembly.GetTypes()...")
	var safeArray *SafeArray
	types := []*Type{}
	err := NewHResultChecker("Assembly::GetTypes").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetTypes,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&safeArray)),
		),
	)
	if err != nil {
		return types, err
	}
	return safeArrayToTypes(safeArray)
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
