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
	get_ToString              uintptr
	Equals                    uintptr
	GetHashCode               uintptr
	GetType                   uintptr
	InitializeLifetimeService uintptr
	GetLifetimeService        uintptr
	get_Evidence              uintptr
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
	get_FriendlyName          uintptr
	get_BaseDirectory         uintptr
	get_RelativeSearchPath    uintptr
	get_ShadowCopyFiles       uintptr
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
	get_DynamicDirectory      uintptr
}

// GetAppDomain is a wrapper function that returns an appDomain from an existing ICORRuntimeHost object
func GetAppDomain(runtimeHost *ICORRuntimeHost) (appDomain *AppDomain, err error) {
	debugPrint("Entering into appdomain.GetAppDomain()...")
	iu, err := runtimeHost.GetDefaultDomain()
	if err != nil {
		return
	}
	err = iu.QueryInterface(IID_AppDomain, unsafe.Pointer(&appDomain))
	return
}

func (obj *AppDomain) QueryInterface(riid *windows.GUID, ppvObject *uintptr) uintptr {
	debugPrint("Entering into appdomain.QueryInterface()...")
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)),
	)
	return ret
}

func (obj *AppDomain) AddRef() uintptr {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	return ret
}

func (obj *AppDomain) Release() uintptr {
	ret, _, _ := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	return ret
}

// GetHashCode serves as the default hash function.
// https://docs.microsoft.com/en-us/dotnet/api/system.object.gethashcode?view=netframework-4.8#System_Object_GetHashCode
func (obj *AppDomain) GetHashCode() (int32, error) {
	debugPrint("Entering into appdomain.GetHashCode()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.GetHashCode,
		uintptr(unsafe.Pointer(obj)),
		0,
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the appdomain.GetHashCode function returned an error:\r\n%s", err)
	}
	// Unable to avoid misuse of unsafe.Pointer because the Windows API call returns the safeArray pointer in the "ret" value. This is a go vet false positive
	return int32(ret), nil
}

// GetFriendlyName returns the friendlyname of the appdomain
func (obj *AppDomain) GetFriendlyName() (name string, err error) {
	debugPrint("Entering into appdomain.GetFriendlyName()...")
	var bstrFriendlyname unsafe.Pointer
	hr, _, err := syscall.SyscallN(
		obj.vtbl.get_FriendlyName,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&bstrFriendlyname)),
		0,
	)
	if err != syscall.Errno(0) {
		return "", fmt.Errorf("the appdomain.GetFriendlyName function returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		err = fmt.Errorf("the appdomain.GetFriendlyName function returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	return ReadUnicodeStr(unsafe.Pointer(bstrFriendlyname)), nil
}

// Load_3 Loads an Assembly into this application domain.
// virtual HRESULT __stdcall Load_3 (
// /*[in]*/ SAFEARRAY * rawAssembly,
// /*[out,retval]*/ struct _Assembly * * pRetVal ) = 0;
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain.load?view=net-5.0
func (obj *AppDomain) Load_3(rawAssembly *SafeArray) (assembly *Assembly, err error) {
	debugPrint("Entering into appdomain.Load_3()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Load_3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(rawAssembly)),
		uintptr(unsafe.Pointer(&assembly)),
	)

	if err != syscall.Errno(0) {
		if err != syscall.Errno(1150) {
			return
		}
	}

	if hr != S_OK {
		err = fmt.Errorf("the appdomain.Load_3 function returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil

	return
}

// ToString Obtains a string representation that includes the friendly name of the application domain and any context policies.
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain.tostring?view=net-5.0#System_AppDomain_ToString
func (obj *AppDomain) ToString() (domain string, err error) {
	debugPrint("Entering into appdomain.ToString()...")
	var pDomain *string
	hr, _, err := syscall.SyscallN(
		obj.vtbl.get_ToString,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&pDomain)),
	)

	if err != syscall.Errno(0) {
		err = fmt.Errorf("the AppDomain.ToString method retured an error:\r\n%s", err)
		return
	}
	if hr != S_OK {
		err = fmt.Errorf("the AppDomain.ToString method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	err = nil
	domain = ReadUnicodeStr(unsafe.Pointer(pDomain))
	return
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
	asmb, err := obj.ListAssemblies()
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

func (obj *AppDomain) GetAssemblies() (safeArray *SafeArray, err error) {
	debugPrint("Entering into appdomain.GetAssemblies()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetAssemblies,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&safeArray)))
	if err != syscall.Errno(0) {
		err = fmt.Errorf("the AppDomain.GetAssemblies method retured an error:\r\n%s", err)
		return
	}
	if hr != S_OK {
		err = fmt.Errorf("the AppDomain.GetAssemblies method returned a non-zero HRESULT: 0x%x", hr)
		return
	}
	return safeArray, nil
}

func (obj *AppDomain) ListAssemblies() (assemblies []*Assembly, err error) {
	debugPrint("Entering into appdomain.ListAssemblies()...")

	safeArray, err := obj.GetAssemblies()
	if err != nil {
		return
	}
	//get dimensions of array (should be 1 for this context always)
	d, err := SafeArrayGetDim(safeArray)
	if err != nil {
		return
	}
	if d != 1 {
		return nil, fmt.Errorf("expected dimension of 1, got %d", d)
	}

	lbound, err := SafeArrayGetLBound(safeArray, d)
	if err != nil {
		return
	}

	ubound, err := SafeArrayGetUBound(safeArray, d)
	if err != nil {
		return
	}
	arrlen := ubound - lbound
	//avoids allocs (lol, overkill)
	assemblies = make([]*Assembly, 0, arrlen)
	for i := lbound; i <= ubound; i++ {
		pApp, err := SafeArrayGetElement(safeArray, i)
		if err != nil {
			return nil, err
		}
		app := (*Assembly)(pApp)
		assemblies = append(assemblies, app)
	}
	return
}
