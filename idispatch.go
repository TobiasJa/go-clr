package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type IDispatch struct {
	vtbl *IDispatchVtbl
}

type IDispatchVtbl struct {
	QueryInterface   uintptr
	AddRef           uintptr
	Release          uintptr
	GetTypeInfoCount uintptr
	GetTypeInfo      uintptr
	GetIDsOfNames    uintptr
	Invoke           uintptr
}

func (obj *IDispatch) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into IDispatch.QueryInterface()...")
	var ppvObject unsafe.Pointer
	hr, _, err := syscall.SyscallN(
		obj.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&riid)),
		uintptr(unsafe.Pointer(ppvObject)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the IDispatch::QueryInterface method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the IDispatch::QueryInterface method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return ppvObject, nil
}

func (obj *IDispatch) AddRef() error {
	debugPrint("Entering into IDispatch.AddRef()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the IDispatch::AddRef method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the IDispatch::AddRef method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func (obj *IDispatch) Release() error {
	debugPrint("Entering into IDispatch.Release()...")
	hr, _, err := syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the IDispatch::Release method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the IDispatch::Release method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func (obj *IDispatch) GetIDsOfName(names []string) ([]int32, error) {
	wnames := make([]*uint16, len(names))
	for i := 0; i < len(names); i++ {
		wnames[i] = windows.StringToUTF16Ptr(names[i])
	}
	dispid := make([]int32, len(names))
	namelen := uint32(len(names))
	hr, _, err := syscall.SyscallN(
		obj.vtbl.GetIDsOfNames,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&IID_NULL)),
		uintptr(unsafe.Pointer(&wnames[0])),
		uintptr(namelen),
		uintptr(GetUserDefaultLCID()),
		uintptr(unsafe.Pointer(&dispid[0])),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the IDispatch::GetIDsOfName method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the IDispatch::GetIDsOfName method method returned a non-zero HRESULT: 0x%x", hr)
	}
	return dispid, nil
}

// func (v *IDispatch) Invoke(dispid int32, dispatch int16, params ...interface{}) (result *Variant, err error) {
// 	result, err = invoke(v, dispid, dispatch, params...)
// 	return
// }

// func (v *IDispatch) GetTypeInfoCount() (c uint32, err error) {
// 	c, err = getTypeInfoCount(v)
// 	return
// }

// func (v *IDispatch) GetTypeInfo() (tinfo *ITypeInfo, err error) {
// 	tinfo, err = getTypeInfo(v)
// 	return
// }

// GetSingleIDOfName is a helper that returns single display ID for IDispatch name.
//
// This replaces the common pattern of attempting to get a single name from the list of available
// IDs. It gives the first ID, if it is available.
func (v *IDispatch) GetSingleIDOfName(name string) (displayID int32, err error) {
	var displayIDs []int32
	displayIDs, err = v.GetIDsOfName([]string{name})
	if err != nil {
		return
	}
	displayID = displayIDs[0]
	return
}

// // InvokeWithOptionalArgs accepts arguments as an array, works like Invoke.
// //
// // Accepts name and will attempt to retrieve Display ID to pass to Invoke.
// //
// // Passing params as an array is a workaround that could be fixed in later versions of Go that
// // prevent passing empty params. During testing it was discovered that this is an acceptable way of
// // getting around not being able to pass params normally.
// func (v *IDispatch) InvokeWithOptionalArgs(name string, dispatch int16, params []interface{}) (result *VARIANT, err error) {
// 	displayID, err := v.GetSingleIDOfName(name)
// 	if err != nil {
// 		return
// 	}

// 	if len(params) < 1 {
// 		result, err = v.Invoke(displayID, dispatch)
// 	} else {
// 		result, err = v.Invoke(displayID, dispatch, params...)
// 	}

// 	return
// }

// // CallMethod invokes named function with arguments on object.
// func (v *IDispatch) CallMethod(name string, params ...interface{}) (*VARIANT, error) {
// 	return v.InvokeWithOptionalArgs(name, DISPATCH_METHOD, params)
// }

// // GetProperty retrieves the property with the name with the ability to pass arguments.
// //
// // Most of the time you will not need to pass arguments as most objects do not allow for this
// // feature. Or at least, should not allow for this feature. Some servers don't follow best practices
// // and this is provided for those edge cases.
// func (v *IDispatch) GetProperty(name string, params ...interface{}) (*VARIANT, error) {
// 	return v.InvokeWithOptionalArgs(name, DISPATCH_PROPERTYGET, params)
// }

// // PutProperty attempts to mutate a property in the object.
// func (v *IDispatch) PutProperty(name string, params ...interface{}) (*VARIANT, error) {
// 	return v.InvokeWithOptionalArgs(name, DISPATCH_PROPERTYPUT, params)
// }
