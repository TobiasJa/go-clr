package clr

import (
	"fmt"
	"syscall"
	"unsafe"
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

func (obj *PropertyInfo) ToString() (string, error) {
	debugPrint("Entering into PropertyInfo.ToString()...")
	var object *string
	hr, _, err := syscall.SyscallN(
		obj.vtbl.ToString,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&object)),
	)
	if err != syscall.Errno(0) {
		return "", fmt.Errorf("the PropertyInfo::ToString method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return "", fmt.Errorf("the PropertyInfo::ToString method returned a non-zero HRESULT: 0x%x", hr)
	}
	return ReadUnicodeStr(unsafe.Pointer(object)), nil
}

func (obj *PropertyInfo) GetValue(instance unsafe.Pointer, index *SafeArray) (*Variant, error) {
	debugPrint("Entering into PropertyInfo.GetValue()...")
	var retVar *Variant
	removeArray := false
	if index == nil {
		index, _ = SafeArrayCreateVector(VT_EMPTY, 0, 0)
		removeArray = true
	}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetValue,
		uintptr(unsafe.Pointer(obj)),
		uintptr(instance),
		uintptr(unsafe.Pointer(index)),
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
