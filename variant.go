//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"
)

type VT uint16

const (
	VT_EMPTY           VT = 0x0
	VT_NULL            VT = 0x1
	VT_I2              VT = 0x2
	VT_I4              VT = 0x3
	VT_R4              VT = 0x4
	VT_R8              VT = 0x5
	VT_CY              VT = 0x6
	VT_DATE            VT = 0x7
	VT_BSTR            VT = 0x8
	VT_DISPATCH        VT = 0x9
	VT_ERROR           VT = 0xa
	VT_BOOL            VT = 0xb
	VT_VARIANT         VT = 0xc
	VT_UNKNOWN         VT = 0xd
	VT_DECIMAL         VT = 0xe
	VT_I1              VT = 0x10
	VT_UI1             VT = 0x11
	VT_UI2             VT = 0x12
	VT_UI4             VT = 0x13
	VT_I8              VT = 0x14
	VT_UI8             VT = 0x15
	VT_INT             VT = 0x16
	VT_UINT            VT = 0x17
	VT_VOID            VT = 0x18
	VT_HRESULT         VT = 0x19
	VT_PTR             VT = 0x1a
	VT_SAFEARRAY       VT = 0x1b
	VT_CARRAY          VT = 0x1c
	VT_USERDEFINED     VT = 0x1d
	VT_LPSTR           VT = 0x1e
	VT_LPWSTR          VT = 0x1f
	VT_RECORD          VT = 0x24
	VT_INT_PTR         VT = 0x25
	VT_UINT_PTR        VT = 0x26
	VT_FILETIME        VT = 0x40
	VT_BLOB            VT = 0x41
	VT_STREAM          VT = 0x42
	VT_STORAGE         VT = 0x43
	VT_STREAMED_OBJECT VT = 0x44
	VT_STORED_OBJECT   VT = 0x45
	VT_BLOB_OBJECT     VT = 0x46
	VT_CF              VT = 0x47
	VT_CLSID           VT = 0x48
	VT_BSTR_BLOB       VT = 0xfff
	VT_VECTOR          VT = 0x1000
	VT_ARRAY           VT = 0x2000
	VT_BYREF           VT = 0x4000
	VT_RESERVED        VT = 0x8000
	VT_ILLEGAL         VT = 0xffff
	VT_ILLEGALMASKED   VT = 0xfff
	VT_TYPEMASK        VT = 0xfff
)

// from https://github.com/go-ole/go-ole/blob/master/variant_amd64.go
// https://docs.microsoft.com/en-us/windows/win32/winauto/variant-structure
// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-variant
// https://docs.microsoft.com/en-us/previous-versions/windows/embedded/ms891678(v=msdn.10)?redirectedfrom=MSDN
// VARIANT Type Constants https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-oaut/3fe7db9f-5803-4dc4-9d14-5425d3f5461f
type Variant struct {
	VT         VT // VARTYPE
	wReserved1 uint16
	wReserved2 uint16
	wReserved3 uint16
	Val        uintptr
	_          [8]byte
}

// ToIUnknown converts Variant to Unknown object.
func (v *Variant) ToIUnknown() *IUnknown {
	if v.VT != VT_UNKNOWN {
		return nil
	}
	return (*IUnknown)(unsafe.Pointer(uintptr(v.Val)))
}

// ToIDispatch converts variant to dispatch object.
func (v *Variant) ToIDispatch() *IDispatch {
	if v.VT != VT_DISPATCH {
		return nil
	}
	return (*IDispatch)(unsafe.Pointer(uintptr(v.Val)))
}

// ToString converts variant to Go string.
func (v *Variant) ToString() string {
	if v.VT != VT_BSTR {
		return ""
	}
	return ReadUnicodeStr(unsafe.Pointer(*(**uint16)(unsafe.Pointer(&v.Val))))
}

// Clear the memory of variant object.
func (v *Variant) Clear() error {
	modOleAuto := syscall.MustLoadDLL("OleAut32.dll")
	procVariantClear := modOleAuto.MustFindProc("VariantClear")
	hr, _, err := procVariantClear.Call(uintptr(unsafe.Pointer(v)))
	if err != syscall.Errno(0) {
		return fmt.Errorf("the Variant::Clear method returned an error:\r\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the Variant::Clear method returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// Value returns variant value based on its type.
//
// Currently supported types: 2- and 4-byte integers, strings, bools.
// Note that 64-bit integers, datetimes, and other types are stored as strings
// and will be returned as strings.
//
// Needs to be further converted, because this returns an interface{}.
func (v *Variant) Value() interface{} {
	switch v.VT {
	case VT_I1:
		return int8(v.Val)
	case VT_UI1:
		return uint8(v.Val)
	case VT_I2:
		return int16(v.Val)
	case VT_UI2:
		return uint16(v.Val)
	case VT_I4:
		return int32(v.Val)
	case VT_UI4:
		return uint32(v.Val)
	case VT_I8:
		return int64(v.Val)
	case VT_UI8:
		return uint64(v.Val)
	case VT_INT:
		return int(v.Val)
	case VT_UINT:
		return uint(v.Val)
	case VT_INT_PTR:
		return uintptr(v.Val) // TODO
	case VT_UINT_PTR:
		return uintptr(v.Val)
	case VT_R4:
		return *(*float32)(unsafe.Pointer(&v.Val))
	case VT_R8:
		return *(*float64)(unsafe.Pointer(&v.Val))
	case VT_BSTR:
		return v.ToString()
	case VT_UNKNOWN:
		return v.ToIUnknown()
	case VT_DISPATCH:
		return v.ToIDispatch()
	case VT_BOOL:
		return (v.Val & 0xffff) != 0
	}
	return nil
}
