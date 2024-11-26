package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// from mscorlib.tlh

type PropertyInfo struct {
	vtbl *PropertyInfoVtbl
}

// TypeVtbl Discovers the attributes of a method and provides access to method metadata.
// Inheritance: Object -> MemberInfo -> Type
// Object Class: https://docs.microsoft.com/en-us/dotnet/api/system.object?view=net-5.0
type PropertyInfoVtbl struct {
	QueryInterface        uintptr
	AddRef                uintptr
	Release               uintptr
	GetTypeInfoCount      uintptr
	GetTypeInfo           uintptr
	GetIDsOfNames         uintptr
	Invoke                uintptr
	ToString              uintptr
	Equals                uintptr
	GetHashCode           uintptr
	GetType               uintptr
	get_MemberType        uintptr
	get_name              uintptr
	get_DeclaringType     uintptr
	get_ReflectedType     uintptr
	GetCustomAttributes   uintptr
	GetCustomAttributes_2 uintptr
	IsDefined             uintptr
	get_PropertyType      uintptr
	GetValue              uintptr
	GetValue_2            uintptr
	SetValue              uintptr
	SetValue_2            uintptr
	GetAccessors          uintptr
	GetGetMethod          uintptr
	GetSetMethod          uintptr
	GetIndexParameters    uintptr
	get_Attributes        uintptr
	get_CanRead           uintptr
	get_CanWrite          uintptr
	GetAccessors_2        uintptr
	GetGetMethod_2        uintptr
	GetSetMethod_2        uintptr
	get_IsSpecialName     uintptr
}

func (obj *PropertyInfo) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into PropertyInfo.QueryInterface()...")
	var ppvObject unsafe.Pointer
	err := NewHResultChecker("PropertyInfo::QueryInterface").CheckHResultSyscallError(syscall.SyscallN(
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

func (obj *PropertyInfo) AddRef() (uint32, error) {
	debugPrint("Entering into PropertyInfo.AddRef()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the PropertyInfo::AddRef method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

func (obj *PropertyInfo) Release() (uint32, error) {
	debugPrint("Entering into PropertyInfo.Release()...")
	ret, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("the PropertyInfo::Release method returned an error:\r\n%s", err)
	}
	return *(*uint32)(unsafe.Pointer(*((**uintptr)(unsafe.Pointer(&ret))))), nil
}

func (obj *PropertyInfo) ToString() (string, error) {
	debugPrint("Entering into PropertyInfo.ToString()...")
	var object *string
	err := NewHResultChecker("PropertyInfo::QueryInterface").CheckHResultSyscallError(
		syscall.SyscallN(
			obj.vtbl.ToString,
			uintptr(unsafe.Pointer(obj)),
			uintptr(unsafe.Pointer(&object)),
		),
	)
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(object)), nil
}

// virtual HRESULT __stdcall GetValue (
//
//	/*[in]*/ VARIANT obj,
//	/*[in]*/ SAFEARRAY * index,
//	/*[out,retval]*/ VARIANT * pRetVal ) = 0;
func (obj *PropertyInfo) GetValue(instance Variant, index *SafeArray) (*Variant, error) {
	debugPrint("Entering into PropertyInfo.GetValue()...")
	var retVar *Variant
	var indexPtr uintptr
	removeArray := false
	if index == nil {
		//index, _ = SafeArrayCreateVector(VT_EMPTY, 0, 0)
		indexPtr = uintptr(0)
		//removeArray = true
	} else {
		indexPtr = uintptr(unsafe.Pointer(index))
	}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetValue,
		uintptr(unsafe.Pointer(obj)),
		uintptr(*(*uintptr)(unsafe.Pointer(&instance))),
		indexPtr,
		uintptr(unsafe.Pointer(&retVar)),
	)
	if removeArray {
		SafeArrayDestroy(index)
	}
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the PropertyInfo::GetValue method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the PropertyInfo::GetValue method returned a non-zero HRESULT: 0x%x", hr)
	}
	return retVar, nil
}
