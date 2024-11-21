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

type Type struct {
	vtbl *TypeVtbl
}

// TypeVtbl Discovers the attributes of a method and provides access to method metadata.
// Inheritance: Object -> MemberInfo -> Type
// Object Class: https://docs.microsoft.com/en-us/dotnet/api/system.object?view=net-5.0
type TypeVtbl struct {
	QueryInterface            uintptr
	AddRef                    uintptr
	Release                   uintptr
	GetTypeInfoCount          uintptr
	GetTypeInfo               uintptr
	GetIDsOfNames             uintptr
	Invoke                    uintptr
	ToString                  uintptr
	Equals                    uintptr
	GetHashCode               uintptr
	GetType                   uintptr
	get_MemberType            uintptr
	get_Name                  uintptr
	get_DeclaringType         uintptr
	get_ReflectedType         uintptr
	GetCustomAttributes       uintptr
	GetCustomAttributes_2     uintptr
	IsDefined                 uintptr
	get_Guid                  uintptr
	get_Module                uintptr
	get_Assembly              uintptr
	get_TypeHandle            uintptr
	get_FullName              uintptr
	get_Namespace             uintptr
	get_AssemblyQualifiedName uintptr
	GetArrayRank              uintptr
	get_BaseType              uintptr
	GetConstructors           uintptr
	GetInterface              uintptr
	GetInterfaces             uintptr
	FindInterfaces            uintptr
	GetEvent                  uintptr
	GetEvents                 uintptr
	GetEvents_2               uintptr
	GetNestedTypes            uintptr
	GetNestedType             uintptr
	GetMember                 uintptr
	GetDefaultMembers         uintptr
	FindMembers               uintptr
	GetElementType            uintptr
	IsSubclassOf              uintptr
	IsInstanceOfType          uintptr
	IsAssignableFrom          uintptr
	GetInterfaceMap           uintptr
	GetMethod                 uintptr
	GetMethod_2               uintptr
	GetMethodes               uintptr
	GetField                  uintptr
	GetFields                 uintptr
	GetProperty               uintptr
	GetProperty_2             uintptr
	GetProperties             uintptr
	GetMember_2               uintptr
	GetMembers                uintptr
	InvokeMember              uintptr
	get_UnterlyingSystemType  uintptr
	InvokeMember_2            uintptr
	InvokeMember_3            uintptr
	GetConstructor            uintptr
	GetConstructor_2          uintptr
	GetConstructor_3          uintptr
	GetConstructors_2         uintptr
	get_TypeInitializer       uintptr
	GetMethod_3               uintptr
	GetMethod_4               uintptr
	GetMethod_5               uintptr
	GetMethod_6               uintptr
	GetMethodes_2             uintptr
	GetField_2                uintptr
	GetFields_2               uintptr
	GetInterface_2            uintptr
	GetEvent_2                uintptr
	GetProperty_3             uintptr
	GetProperty_4             uintptr
	GetProperty_5             uintptr
	GetProperty_6             uintptr
	GetProperty_7             uintptr
	GetProperties_2           uintptr
	GetNestedTypes_2          uintptr
	GetNestedType_2           uintptr
	GetMember_3               uintptr
	GetMembers_2              uintptr
	get_Attributes            uintptr
	get_IsNotPublic           uintptr
	get_IsPublic              uintptr
	get_IsNestedPublic        uintptr
	get_IsNestedPrivate       uintptr
	get_IsNestedFamily        uintptr
	get_IsNestedAssembly      uintptr
	get_IsNestedFamANDAssem   uintptr
	get_IsNestedFamORAssem    uintptr
	get_IsAutoLayout          uintptr
	get_IsLayoutSequential    uintptr
	get_IsExplicitLayout      uintptr
	get_IsClass               uintptr
	get_IsInterface           uintptr
	get_IsValueType           uintptr
	get_IsAbstract            uintptr
	get_IsSealed              uintptr
	get_IsEnum                uintptr
	get_IsSpecialName         uintptr
	get_IsImport              uintptr
	get_IsSerializable        uintptr
	get_IsAnsiClass           uintptr
	get_IsUnicodeClass        uintptr
	get_IsAutoClass           uintptr
	get_IsArray               uintptr
	get_IsByRef               uintptr
	get_IsPointer             uintptr
	get_IsPrimitive           uintptr
	get_IsCOMObject           uintptr
	get_HasElementType        uintptr
	get_IsContextful          uintptr
	get_IsMarshalByRef        uintptr
	Equals_2                  uintptr
}

func (obj *Type) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into Type.QueryInterface()...")
	var ppvObject unsafe.Pointer
	hr, _, err := syscall.SyscallN(
		obj.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)), // A reference to the interface identifier (IID) of the interface being queried for.
		uintptr(ppvObject),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the Type::QueryInterface method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the Type::QueryInterface method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return ppvObject, nil
}

func (obj *Type) AddRef() error {
	debugPrint("Entering into Type.AddRef()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the Type::AddRef method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the Type::AddRef method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func (obj *Type) Release() error {
	debugPrint("Entering into Type.Release()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the Type::Release method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the Type::Release method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// GetString returns a string that represents the current object
// a string version of the method's signature
// public virtual string ToString ();
// https://docs.microsoft.com/en-us/dotnet/api/system.object.tostring?view=net-5.0#System_Object_ToString
func (obj *Type) ToString() (string, error) {
	debugPrint("Entering into Type.ToString()...")
	var object *string
	hr, _, err := syscall.SyscallN(
		obj.vtbl.ToString,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&object)),
	)
	if err != syscall.Errno(0) {
		return "", fmt.Errorf("the Type::ToString method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return "", fmt.Errorf("the Type::ToString method returned a non-zero HRESULT: 0x%x", hr)
	}
	return ReadUnicodeStr(unsafe.Pointer(object)), nil
}

func (obj *Type) GetConstructor(parameterVariant []Variant) (*ConstructorInfo, error) {
	debugPrint("Entering into Type.GetConstructor()...")
	var constructorInfo *ConstructorInfo
	safeArray, err := variantsToSafeArray(parameterVariant)
	if err != nil {
		return nil, err
	}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetConstructor_3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(safeArray)),
		uintptr(unsafe.Pointer(&constructorInfo)),
	)
	errSafeArray := SafeArrayDestroy(safeArray)
	if errSafeArray != nil {
		fmt.Printf("the Type::GetConstructor SafeArrayDestroy method returned an error:\r\n%s", errSafeArray)
	}
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the Type::GetConstructor method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the Type::GetConstructor method returned a non-zero HRESULT: 0x%x", hr)
	}
	return constructorInfo, nil
}

func (obj *Type) GetConstructorBySignatur(signature string) (*ConstructorInfo, error) {
	constructorInfos, err := obj.GetConstructors()
	if err != nil {
		return nil, err
	}
	for _, constructorInfo := range constructorInfos {
		constructorInfoStr, err := constructorInfo.ToString()
		if err != nil {
			continue
		}
		if constructorInfoStr == signature {
			return constructorInfo, nil
		}
	}
	return nil, fmt.Errorf("could not find a constructor with the given signature: %s", signature)
}

func (obj *Type) GetConstructors() ([]*ConstructorInfo, error) {
	debugPrint("Entering into Type.GetConstructors()...")
	constructorInfos := []*ConstructorInfo{}
	// safeArray, err := SafeArrayCreateVector(VT_UNKNOWN, 0, 255)
	// if err != nil {
	// 	return constructorInfos, err
	// }
	var safeArray *SafeArray
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetConstructors_2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&safeArray)),
	)
	if err != syscall.Errno(0) {
		return constructorInfos, fmt.Errorf("the Type::GetConstructor method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return constructorInfos, fmt.Errorf("the Type::GetConstructor method returned a non-zero HRESULT: 0x%x", hr)
	}
	return safeArrayToConstructorInfos(safeArray)
}

func (obj *Type) GetMethod(methodStr string) (*MethodInfo, error) {
	debugPrint("Entering into Type.GetMethod()...")
	var methodInfoPtr *MethodInfo
	methodStrPtr, err := SysAllocString(methodStr)
	if err != nil {
		return nil, fmt.Errorf("the Type::GetMethod SysAllocString returned error: %v", err)
	}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetMethod_6,
		uintptr(unsafe.Pointer(obj)),
		uintptr(methodStrPtr),
		uintptr(unsafe.Pointer(&methodInfoPtr)),
	)
	errFree := SysFreeString(methodStrPtr)
	if errFree != nil {
		return nil, fmt.Errorf("the Type::GetMethod error free String:\r\n%v", errFree)
	}
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the Type::GetMethod method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the Type::GetMethod method returned a non-zero HRESULT: 0x%x", hr)
	}
	return methodInfoPtr, nil
}

func (obj *Type) GetMethodBySignature(signature string) (*MethodInfo, error) {
	methodInfos, err := obj.GetMethodes()
	if err != nil {
		return nil, err
	}
	for _, methodInfo := range methodInfos {
		methodInfoStr, err := methodInfo.ToString()
		if err != nil {
			continue
		}
		if methodInfoStr == signature {
			return methodInfo, nil
		}
	}
	return nil, nil
}

func (obj *Type) GetMethodes() ([]*MethodInfo, error) {
	debugPrint("Entering into Type.GetMethodes()...")
	var safeArray *SafeArray
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetMethodes_2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&safeArray)))
	if err != syscall.Errno(0) {
		return []*MethodInfo{}, fmt.Errorf("the Type.GetMethodes method retured an error:\r\n%s", err)
	}
	if hr != S_OK {
		return []*MethodInfo{}, fmt.Errorf("the Type.GetMethodes method returned a non-zero HRESULT: 0x%x", hr)
	}
	return safeArrayToMethodes(safeArray)
}

func (obj *Type) GetProperty(propertyStr string) (*PropertyInfo, error) {
	debugPrint("Entering into Type.GetProperty()...")
	var propertyInfoPtr *PropertyInfo
	propertyStrPtr, err := SysAllocString(propertyStr)
	if err != nil {
		return nil, fmt.Errorf("the Type::GetProperty SysAllocString returned error: %v", err)
	}
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetProperty_7,
		uintptr(unsafe.Pointer(obj)),
		uintptr(propertyStrPtr),
		uintptr(unsafe.Pointer(&propertyInfoPtr)),
	)
	errFree := SysFreeString(propertyStrPtr)
	if errFree != nil {
		return nil, fmt.Errorf("the Type::GetProperty error free String:\r\n%v", errFree)
	}
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the Type::GetProperty method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the Type::GetProperty method returned a non-zero HRESULT: 0x%x", hr)
	}
	return propertyInfoPtr, nil
}

func (obj *Type) GetProperties() ([]*PropertyInfo, error) {
	debugPrint("Entering into Type.GetProperties()...")
	var safeArray *SafeArray
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetProperties_2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&safeArray)))
	if err != syscall.Errno(0) {
		return []*PropertyInfo{}, fmt.Errorf("the Type.GetProperties method retured an error:\r\n%s", err)
	}
	if hr != S_OK {
		return []*PropertyInfo{}, fmt.Errorf("the Type.GetProperties method returned a non-zero HRESULT: 0x%x", hr)
	}
	return safeArrayToProperties(safeArray)
}

func variantsToSafeArray(variants []Variant) (*SafeArray, error) {
	safeArray, err := SafeArrayCreateVector(VT_VARIANT, 0, uint32(len(variants)))
	if err != nil {
		return nil, err
	}
	for i, variant := range variants {
		SafeArrayPutElement(safeArray, int32(i), unsafe.Pointer(&variant))
	}
	return safeArray, nil
}

func safeArrayToConstructorInfos(safeArray *SafeArray) ([]*ConstructorInfo, error) {
	debugPrint("Entering into type.safeArrayToConstructors()...")
	//get dimensions of array (should be 1 for this context always)
	constructorInfos := []*ConstructorInfo{}
	d, err := SafeArrayGetDim(safeArray)
	if err != nil {
		return constructorInfos, err
	}
	if *d != 1 {
		return constructorInfos, fmt.Errorf("expected dimension of 1, got %d", d)
	}

	lbound, err := SafeArrayGetLBound(safeArray, *d)
	if err != nil {
		return constructorInfos, err
	}

	ubound, err := SafeArrayGetUBound(safeArray, *d)
	if err != nil {
		return constructorInfos, err
	}
	for i := lbound; i <= ubound; i++ {
		pConstructorInfo, err := SafeArrayGetElement(safeArray, i)
		if err != nil {
			continue
		}
		constructorInfos = append(constructorInfos, (*ConstructorInfo)(pConstructorInfo))
	}
	err = SafeArrayDestroy(safeArray)
	if err != nil {
		return constructorInfos, err
	}
	return constructorInfos, nil
}

func safeArrayToMethodes(safeArray *SafeArray) ([]*MethodInfo, error) {
	debugPrint("Entering into assembly.safeArrayToMethods()...")
	//get dimensions of array (should be 1 for this context always)
	methodInfos := []*MethodInfo{}
	d, err := SafeArrayGetDim(safeArray)
	if err != nil {
		return methodInfos, err
	}
	if *d != 1 {
		return methodInfos, fmt.Errorf("expected dimension of 1, got %d", d)
	}

	lbound, err := SafeArrayGetLBound(safeArray, *d)
	if err != nil {
		return methodInfos, err
	}

	ubound, err := SafeArrayGetUBound(safeArray, *d)
	if err != nil {
		return methodInfos, err
	}
	for i := lbound; i <= ubound; i++ {
		pMethodInfo, err := SafeArrayGetElement(safeArray, i)
		if err != nil {
			continue
		}
		methodInfos = append(methodInfos, (*MethodInfo)(pMethodInfo))
	}
	err = SafeArrayDestroy(safeArray)
	if err != nil {
		return methodInfos, err
	}
	return methodInfos, nil
}

func safeArrayToProperties(safeArray *SafeArray) ([]*PropertyInfo, error) {
	debugPrint("Entering into assembly.safeArrayToProperties()...")
	//get dimensions of array (should be 1 for this context always)
	propertyInfos := []*PropertyInfo{}
	d, err := SafeArrayGetDim(safeArray)
	if err != nil {
		return propertyInfos, err
	}
	if *d != 1 {
		return propertyInfos, fmt.Errorf("expected dimension of 1, got %d", d)
	}

	lbound, err := SafeArrayGetLBound(safeArray, *d)
	if err != nil {
		return propertyInfos, err
	}

	ubound, err := SafeArrayGetUBound(safeArray, *d)
	if err != nil {
		return propertyInfos, err
	}
	for i := lbound; i <= ubound; i++ {
		pPropertyInfo, err := SafeArrayGetElement(safeArray, i)
		if err != nil {
			continue
		}
		propertyInfos = append(propertyInfos, (*PropertyInfo)(pPropertyInfo))
	}
	err = SafeArrayDestroy(safeArray)
	if err != nil {
		return propertyInfos, err
	}
	return propertyInfos, nil
}
