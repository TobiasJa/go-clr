//go:build windows
// +build windows

package clr

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// AppDomain is a Windows COM object interface pointer for the .NET AppDomain class.
// The AppDomain object represents an application domain, which is an isolated environment where applications execute.
// This structure only contains a pointer to the AppDomain's virtual function table
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain?view=netframework-4.8
type AppDomain struct {
	vtbl *AppDomainVtbl
}

// AppDomainVtbl is a Virtual Function Table for the AppDomain COM interface
// The Virtual Function Table contains pointers to the COM IUnkown interface
// functions (QueryInterface, AddRef, & Release) as well as the AppDomain object's methods
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain?view=netframework-4.8
type AppDomainVtbl struct {
	QueryInterface            uintptr
	AddRef                    uintptr
	Release                   uintptr
	GetTypeInfoCount          uintptr
	GetTypeInfo               uintptr
	GetIDsOfNames             uintptr
	Invoke                    uintptr
	GetToString               uintptr
	Equals                    uintptr
	GetHashCode               uintptr
	GetType                   uintptr
	InitializeLifetimeService uintptr
	GetLifetimeService        uintptr
	GetEvidence               uintptr
	add_DomainUnload          uintptr
	remove_DomainUnload       uintptr
	add_AssemblyLoad          uintptr
	remove_AssemblyLoad       uintptr
	add_ProcessExit           uintptr
	remove_ProcessExit        uintptr
	add_TypeResolve           uintptr
	remove_TypeResolve        uintptr
	add_ResourceResolve       uintptr
	remove_ResourceResolve    uintptr
	add_AssemblyResolve       uintptr
	remove_AssemblyResolve    uintptr
	add_UnhandledException    uintptr
	remove_UnhandledException uintptr
	DefineDynamicAssembly     uintptr
	DefineDynamicAssembly_2   uintptr
	DefineDynamicAssembly_3   uintptr
	DefineDynamicAssembly_4   uintptr
	DefineDynamicAssembly_5   uintptr
	DefineDynamicAssembly_6   uintptr
	DefineDynamicAssembly_7   uintptr
	DefineDynamicAssembly_8   uintptr
	DefineDynamicAssembly_9   uintptr
	CreateInstance            uintptr
	CreateInstanceFrom        uintptr
	CreateInstance_2          uintptr
	CreateInstanceFrom_2      uintptr
	CreateInstance_3          uintptr
	CreateInstanceFrom_3      uintptr
	Load                      uintptr
	Load_2                    uintptr
	Load_3                    uintptr
	Load_4                    uintptr
	Load_5                    uintptr
	Load_6                    uintptr
	Load_7                    uintptr
	ExecuteAssembly           uintptr
	ExecuteAssembly_2         uintptr
	ExecuteAssembly_3         uintptr
	GetFriendlyName           uintptr
	GetBaseDirectory          uintptr
	GetRelativeSearchPath     uintptr
	GetShadowCopyFiles        uintptr
	GetAssemblies             uintptr
	AppendPrivatePath         uintptr
	ClearPrivatePath          uintptr
	SetShadowCopyPath         uintptr
	ClearShadowCopyPath       uintptr
	SetCachePath              uintptr
	SetData                   uintptr
	GetData                   uintptr
	SetAppDomainPolicy        uintptr
	SetThreadPrincipal        uintptr
	SetPrincipalPolicy        uintptr
	DoCallBack                uintptr
	GetDynamicDirectory       uintptr
}

// GetDefaultAppDomain is a wrapper function that returns an appDomain from an existing ICORRuntimeHost object
func GetDefaultAppDomain(runtimeHost *ICORRuntimeHost) (appDomain *AppDomain, err error) {
	debugPrint("Entering into appdomain.GetAppDomain()...")
	return runtimeHost.GetDefaultDomain()
}

func (obj *AppDomain) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into AppDomain.QueryInterface()...")
	var ppvObject unsafe.Pointer
	err := NewHResultChecker("AppDomain::QueryInterface").CheckHResultSyscallError(syscall.SyscallN(
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

func (obj *AppDomain) AddRef() error {
	debugPrint("Entering into AppDomain.AddRef()...")
	return NewHResultChecker("AppDomain::AddRef").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	))
}

func (obj *AppDomain) Release() error {
	debugPrint("Entering into AppDomain.Release()...")
	return NewHResultChecker("AppDomain::Release").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	))
}

// ToString Obtains a string representation that includes the friendly name of the application domain and any context policies.
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain.tostring?view=net-5.0#System_AppDomain_ToString
func (obj *AppDomain) ToString() (string, error) {
	debugPrint("Entering into AppDomain.ToString()...")
	var pDomain *string
	err := NewHResultChecker("AppDomain::GetToString").CheckHResultSyscallError(
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

// GetHashCode serves as the default hash function.
// https://docs.microsoft.com/en-us/dotnet/api/system.object.gethashcode?view=netframework-4.8#System_Object_GetHashCode
func (obj *AppDomain) GetHashCode() (int32, error) {
	debugPrint("Entering into AppDomain.GetHashCode()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.GetHashCode,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the AppDomain.GetHashCode function returned an error:\r\n%s", err)
	}
	// Unable to avoid misuse of unsafe.Pointer because the Windows API call returns the safeArray pointer in the "ret" value. This is a go vet false positive
	return int32(ret), nil
}

// virtual HRESULT __stdcall GetType (
// /*[out,retval]*/ struct _Type * * pRetVal ) = 0;
func (obj *AppDomain) GetType() (*Type, error) {
	debugPrint("Entering into AppDomain.GetType()...")
	var typePtr *Type
	err := NewHResultChecker("AppDomain::GetType").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetType,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&typePtr)),
		),
	)
	if err != nil {
		return nil, err
	}
	return typePtr, nil
}

// _EvidencePtr GetEvidence ( );
func (obj *AppDomain) GetEvidence() (*IUnknown, error) {
	var evidence *IUnknown
	err := NewHResultChecker("AppDomain::GetEvidence").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetEvidence,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&evidence)),
		),
	)
	if err != nil {
		return nil, err
	}
	return evidence, nil
}

// Load_2 takes an assemblystring (name) and checks each Assembly in the appdomain for a prefix match (case insensitive). If a match is found, it is returned.
func (obj *AppDomain) Load_2(assemblyString string) (*Assembly, error) {
	// var err error
	// var pAssembly *Assembly
	// str, _ := SysAllocString(assemblyString)
	// ret, _, _ := syscall.SyscallN(
	// 	obj.vtbl.Load_2,
	// 	uintptr(unsafe.Pointer(obj)),
	// 	uintptr(unsafe.Pointer(str)),
	// 	uintptr(unsafe.Pointer(&pAssembly)))
	// if ret != 0 {
	// 	err = fmt.Errorf("bad load 2: %x", ret)
	// }
	// return pAssembly, err
	//load_2 isn't working nicely, and appears to want a fully qualified assembly name - so let's do the (dumb?) thing and use ListAssemblies to return one that has the right prefix
	debugPrint("Entering into appdomain.Load_2()...")
	asmb, err := obj.GetAssemblies()
	if err != nil {
		return nil, err
	}
	for i := range asmb {
		if name, err := asmb[i].GetFullName(); err == nil {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(assemblyString)) {
				return asmb[i], nil
			}
		}
	}
	return nil, fmt.Errorf("could not find assembly with name %s", assemblyString)
}

// Load_3 Loads an Assembly into this application domain.
// virtual HRESULT __stdcall Load_3 (
// /*[in]*/ SAFEARRAY * rawAssembly,
// /*[out,retval]*/ struct _Assembly * * pRetVal ) = 0;
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain.load?view=net-5.0
func (obj *AppDomain) Load_3(rawAssembly *SafeArray) (*Assembly, error) {
	debugPrint("Entering into appdomain.Load_3()...")
	var assembly *Assembly
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Load_3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(rawAssembly)),
		uintptr(unsafe.Pointer(&assembly)),
	)

	if err != syscall.Errno(0) && err != syscall.Errno(1150) {
		return nil, err
	}

	if hr != S_OK {
		return nil, fmt.Errorf("the appdomain.Load_3 function returned a non-zero HRESULT: 0x%x", hr)
	}

	return assembly, nil
}

// GetFriendlyName returns the friendlyname of the appdomain
//
// _bstr_t GetFriendlyName ( );
func (obj *AppDomain) GetFriendlyName() (string, error) {
	debugPrint("Entering into AppDomain.GetFriendlyName()...")
	var bstrFriendlyname unsafe.Pointer
	err := NewHResultChecker("AppDomain::GetFriendlyName").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetFriendlyName,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&bstrFriendlyname)),
		),
	)
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(bstrFriendlyname)), nil
}

// GetBaseDirectory returns the base directory of the appdomain
//
// _bstr_t GetBaseDirectory ( );
func (obj *AppDomain) GetBaseDirectory() (string, error) {
	debugPrint("Entering into AppDomain.GetBaseDirectory()...")
	var bstrGetBaseDirectory unsafe.Pointer
	err := NewHResultChecker("AppDomain::GetBaseDirectory").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetBaseDirectory,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&bstrGetBaseDirectory)),
		),
	)
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(bstrGetBaseDirectory)), nil
}

// GetRelativeSearchPath returns the relative search path of the appdomain
//
// _bstr_t GetRelativeSearchPath ( );
func (obj *AppDomain) GetRelativeSearchPath() (string, error) {
	debugPrint("Entering into AppDomain.GetRelativeSearchPath()...")
	var bstrRelativeSearchPath unsafe.Pointer
	err := NewHResultChecker("AppDomain::GetRelativeSearchPath").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetRelativeSearchPath,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&bstrRelativeSearchPath)),
		),
	)
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(bstrRelativeSearchPath)), nil
}

// SAFEARRAY * GetAssemblies ( );
func (obj *AppDomain) GetAssemblies() ([]*Assembly, error) {
	debugPrint("Entering into appdomain.GetAssemblies()...")
	var safeArray *SafeArray
	err := NewHResultChecker("AppDomain::GetAssemblies").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.GetAssemblies,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&safeArray)),
		),
	)
	if err != nil {
		return []*Assembly{}, err
	}
	return safeArrayToAssemblies(safeArray)
}

// HRESULT AppendPrivatePath (
//
//	_bstr_t Path );
func (obj *AppDomain) AppendPrivatePath(path string) error {
	debugPrint("Entering into AppDomain.AppendPrivatePath()...")
	pathPtr, err := SysAllocString(path)
	if err != nil {
		return fmt.Errorf("the AppDomain::AppendPrivatePath SysAllocString returned error: %v", err)
	}
	err = NewHResultChecker("AppDomain::AppendPrivatePath").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.CreateInstance,
			uintptr(unsafe.Pointer(obj)),
			uintptr(pathPtr),
		),
	)
	errFree := SysFreeString(pathPtr)
	if errFree != nil {
		return fmt.Errorf("the AppDomain::AppendPrivatePath error free String:\r\n%v", errFree)
	}
	return err
}

// HRESULT ClearPrivatePath ( );
func (obj *AppDomain) ClearPrivatePath() error {
	debugPrint("Entering into AppDomain.ClearPrivatePath()...")
	return NewHResultChecker("AppDomain::ClearPrivatePath").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.CreateInstance,
			uintptr(unsafe.Pointer(obj)),
		),
	)
}

func safeArrayToAssemblies(safeArray *SafeArray) ([]*Assembly, error) {
	debugPrint("Entering into appdomain.safeArrayToAssemblies()...")
	//get dimensions of array (should be 1 for this context always)
	assemblies := []*Assembly{}
	d, err := SafeArrayGetDim(safeArray)
	if err != nil {
		return assemblies, err
	}
	if d == nil || *d != 1 {
		return assemblies, fmt.Errorf("expected dimension of 1, got %d", d)
	}

	lbound, err := SafeArrayGetLBound(safeArray, *d)
	if err != nil {
		return assemblies, err
	}

	ubound, err := SafeArrayGetUBound(safeArray, *d)
	if err != nil {
		return assemblies, err
	}
	for i := lbound; i <= ubound; i++ {
		pApp, err := SafeArrayGetElement(safeArray, i)
		if err != nil {
			continue
		}
		assemblies = append(assemblies, (*Assembly)(pApp))
	}
	err = SafeArrayDestroy(safeArray)
	if err != nil {
		return assemblies, err
	}
	return assemblies, nil
}
