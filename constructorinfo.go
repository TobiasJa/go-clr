package clr

import (
	"fmt"
	"syscall"
	"unsafe"
)

// from mscorlib.tlh

type ConstructorInfo struct {
	vtbl *ConstructorInfoVtbl
}

// TypeVtbl Discovers the attributes of a method and provides access to method metadata.
// Inheritance: Object -> MemberInfo -> Type
// Object Class: https://docs.microsoft.com/en-us/dotnet/api/system.object?view=net-5.0
type ConstructorInfoVtbl struct {
	QueryInterface               uintptr
	AddRef                       uintptr
	Release                      uintptr
	GetTypeInfoCount             uintptr
	GetTypeInfo                  uintptr
	GetIDsOfNames                uintptr
	Invoke                       uintptr
	ToString                     uintptr
	Equals                       uintptr
	GetHashCode                  uintptr
	GetType                      uintptr
	get_MemberType               uintptr
	get_name                     uintptr
	get_DeclaringType            uintptr
	get_ReflectedType            uintptr
	GetCustomAttributes          uintptr
	GetCustomAttributes_2        uintptr
	IsDefined                    uintptr
	GetParameters                uintptr
	GetMethodImplementationFlags uintptr
	get_MethodHandle             uintptr
	get_Attributes               uintptr
	get_CallingConvention        uintptr
	Invoke_2                     uintptr
	get_IsPublic                 uintptr
	get_IsPrivate                uintptr
	get_IsFamily                 uintptr
	get_IsAssembly               uintptr
	get_IsFamilyAndAssembly      uintptr
	get_IsFamilyOrAssembly       uintptr
	get_IsStatic                 uintptr
	get_IsFinal                  uintptr
	get_IsVirtual                uintptr
	get_IsHideBySig              uintptr
	get_IsAbstract               uintptr
	get_IsSpecialName            uintptr
	get_IsConstructor            uintptr
	Invoke_3                     uintptr
	Invoke_4                     uintptr
	Invoke_5                     uintptr
}

func (obj *ConstructorInfo) Invoke(args *SafeArray) (*Variant, error) {
	argsLen, err := SafeArrayGetArrayLength(args)
	if err != nil {
		return nil, err
	}
	parameters, err := obj.GetParameters()
	if err != nil {
		return nil, err
	}
	parameterLen, err := SafeArrayGetArrayLength(parameters)
	if err != nil {
		return nil, err
	}
	if parameterLen != argsLen {
		return nil, fmt.Errorf("arguments do not match method signature: %d given, %d expected", argsLen, parameterLen)
	}
	var variant *Variant
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Invoke_5,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(&variant)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the ConstructorInfo::Invoke method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the ConstructorInfo::Invoke method returned a non-zero HRESULT: 0x%x", hr)
	}
	return variant, nil
}

func (obj *ConstructorInfo) GetParameters() (*SafeArray, error) {
	debugPrint("Entering into ConstructorInfo.GetParameters()...")
	var safeArray *SafeArray
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetParameters,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(safeArray)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the ConstructorInfo::GetParameters method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the ConstructorInfo::GetParameters method returned a non-zero HRESULT: 0x%x", hr)
	}
	return safeArray, nil
}

func (obj *ConstructorInfo) ToString() (string, error) {
	debugPrint("Entering into ConstructorInfo.ToString()...")
	var object *string
	hr, _, err := syscall.SyscallN(
		obj.vtbl.ToString,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&object)),
	)
	if err != syscall.Errno(0) {
		return "", fmt.Errorf("the ConstructorInfo::ToString method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return "", fmt.Errorf("the ConstructorInfo::ToString method returned a non-zero HRESULT: 0x%x", hr)
	}
	return ReadUnicodeStr(unsafe.Pointer(object)), nil
}
