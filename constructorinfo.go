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
	err = NewHResultChecker("ConstructorInfo::Invoke").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.Invoke_5,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(&variant)),
	))
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (obj *ConstructorInfo) GetParameters() (*SafeArray, error) {
	debugPrint("Entering into ConstructorInfo.GetParameters()...")
	var safeArray *SafeArray
	err := NewHResultChecker("ConstructorInfo::GetParameters").CheckHResultSyscallError(syscall.SyscallN(
		obj.vtbl.GetParameters,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(safeArray)),
	))
	if err != nil {
		return nil, err
	}
	return safeArray, nil
}

func (obj *ConstructorInfo) ToString() (string, error) {
	debugPrint("Entering into ConstructorInfo.ToString()...")
	var object *string
	err := NewHResultChecker("ConstructorInfo::GetParameters").CheckHResultSyscallError(
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
