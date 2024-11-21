//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// SafeArray represents a safe array
// defined in OAIdl.h
//
//	typedef struct tagSAFEARRAY {
//	  USHORT         cDims;
//	  USHORT         fFeatures;
//	  ULONG          cbElements;
//	  ULONG          cLocks;
//	  PVOID          pvData;
//	  SAFEARRAYBOUND rgsabound[1];
//	} SAFEARRAY;
//
// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-safearray
// https://docs.microsoft.com/en-us/archive/msdn-magazine/2017/march/introducing-the-safearray-data-structure
type SafeArray struct {
	// cDims is the number of dimensions
	cDims uint16
	// fFeatures is the feature flags
	fFeatures uint16
	// cbElements is the size of an array element
	cbElements uint32
	// cLocks is the number of times the array has been locked without a corresponding unlock
	cLocks uint32
	// pvData is the data
	pvData uintptr
	// rgsabout is one bound for each dimension
	rgsabound [1]SafeArrayBound
}

// SafeArrayBound represents the bounds of one dimension of the array
//
//	typedef struct tagSAFEARRAYBOUND {
//	  ULONG cElements;
//	  LONG  lLbound;
//	} SAFEARRAYBOUND, *LPSAFEARRAYBOUND;
//
// https://docs.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-safearraybound
type SafeArrayBound struct {
	// cElements is the number of elements in the dimension
	cElements uint32
	// lLbound is the lowerbound of the dimension
	lLbound int32
}

func SafeArrayCreateBytesWithSize(size uint32) (*SafeArray, error) {
	debugPrint("Entering into safearray.SafeArrayCreateWithSize()...")

	safeArrayBounds := SafeArrayBound{
		cElements: size,
		lLbound:   int32(0),
	}
	return SafeArrayCreate(VT_UI1, uint32(1), &safeArrayBounds)
}

func SafeArrayCopyBytes(safeArray *SafeArray, rawBytes []byte) error {
	chunkSize := 4096
	lenBytes := len(rawBytes)
	destPtr := uintptr(unsafe.Pointer(safeArray.pvData))
	sourcePtr := uintptr(unsafe.Pointer(&rawBytes[0]))
	var err error
	for {
		debugPrint(fmt.Sprintf("Bytes: %d; destPtr: %X; sourcePtr: %X\n", lenBytes, destPtr, sourcePtr))
		chunkBytes := min(chunkSize, lenBytes)
		_, _, err = modkernel32.MustFindProc("RtlMoveMemory").Call(
			destPtr,
			sourcePtr,
			uintptr(chunkBytes),
		)
		if err != syscall.Errno(0) {
			return err
		}
		destPtr = uintptr(destPtr + uintptr(chunkBytes))
		sourcePtr = uintptr(sourcePtr + uintptr(chunkBytes))
		lenBytes = lenBytes - chunkBytes
		if lenBytes == 0 {
			break
		}
	}
	return nil
}

// CreateSafeArray is a wrapper function that takes in a Go byte array and creates a SafeArray containing unsigned bytes
// by making two syscalls and copying raw memory into the correct spot.
func CreateSafeArray(rawBytes []byte) (*SafeArray, error) {
	debugPrint("Entering into safearray.CreateSafeArray()...")
	safeArray, err := SafeArrayCreateBytesWithSize(uint32(len(rawBytes)))
	if err != nil {
		return nil, err
	}
	err = SafeArrayLock(safeArray)
	if err != nil {
		return nil, err
	}
	err = SafeArrayCopyBytes(safeArray, rawBytes)
	if err != nil {
		return nil, err
	}
	err = SafeArrayUnlock(safeArray)
	if err != nil {
		return nil, err
	}
	return safeArray, nil
}

// SafeArrayAccessData increments the lock count of an array, and retrieves a pointer to the array data
// HRESULT SafeArrayAccessData(
//
//	SAFEARRAY  *psa,
//	void HUGEP **ppvData
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearrayaccessdata
func SafeArrayAccessData(psa *SafeArray) (*uintptr, error) {
	debugPrint("Entering into safearray.SafeArrayAccessData()...")
	var ppvData *uintptr
	hr, _, err := modOleAut.MustFindProc("SafeArrayAccessData").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&ppvData)),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the oleaut32!SafeArrayAccessData function returned a non-zero HRESULT: 0x%x", hr)
	}
	return ppvData, nil
}

// SafeArrayAllocData allocates SafeArray.
//
// HRESULT SafeArrayAllocData(
//
//	[in] SAFEARRAY *psa
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearrayallocdata
func SafeArrayAllocData(psa *SafeArray) error {
	debugPrint("Entering into safearray.SafeArrayAllocData()...")
	_, _, err := modOleAut.MustFindProc("SafeArrayAllocData").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return err
	}
	// if hr != S_OK {
	// 	return fmt.Errorf("the safearray.SafeArrayAllocData() function return 0x%x", hr)
	// }
	return nil
}

// safeArrayAllocDescriptor allocates SafeArray.
//
// HRESULT SafeArrayAllocDescriptor(
//
//	[in]  UINT      cDims,
//	[out] SAFEARRAY **ppsaOut
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearrayallocdescriptor
func SafeArrayAllocDescriptor(cDims uint32) (*SafeArray, error) {
	debugPrint("Entering into safearray.SafeArrayAllocDescriptor()...")
	var safeArray *SafeArray
	hr, _, err := modOleAut.MustFindProc("SafeArrayAllocDescriptor").Call(
		uintptr(cDims),
		uintptr(unsafe.Pointer(&safeArray)),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the safearray.SafeArrayAllocData() function return 0x%x", hr)
	}
	return safeArray, nil
}

// safeArrayAllocDescriptorEx allocates SafeArray.
//
// HRESULT SafeArrayAllocDescriptorEx(
//
//	[in]  VARTYPE   vt,
//	[in]  UINT      cDims,
//	[out] SAFEARRAY **ppsaOut
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearrayallocdescriptorex
func SafeArrayAllocDescriptorEx(vt VT, cDims uint32) (*SafeArray, error) {
	debugPrint("Entering into safearray.SafeArrayAllocDescriptorEx()...")
	var ppsaOut *SafeArray
	hr, _, err := modOleAut.MustFindProc("SafeArrayAllocDescriptorEx").Call(
		uintptr(vt),
		uintptr(cDims),
		uintptr(unsafe.Pointer(&ppsaOut)),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the safearray.SafeArrayAllocDescriptorEx() function return 0x%x", hr)
	}
	return ppsaOut, nil
}

// SafeArrayCopy returns copy of SafeArray.
//
// HRESULT SafeArrayCopy(
//
//	[in]  SAFEARRAY *psa,
//	[out] SAFEARRAY **ppsaOut
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraycopy
func SafeArrayCopy(psa *SafeArray) (*SafeArray, error) {
	debugPrint("Entering into safearray.SafeArrayCopy()...")
	var ppsaOut *SafeArray
	hr, _, err := modOleAut.MustFindProc("SafeArrayCopy").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&ppsaOut)),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the safearray.SafeArrayCopy() function return 0x%x", hr)
	}
	return ppsaOut, nil
}

// SafeArrayCopyData duplicates SafeArray into another SafeArray object.
//
// HRESULT SafeArrayCopyData(
//
//	[in] SAFEARRAY *psaSource,
//	[in] SAFEARRAY *psaTarget
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraycopydata
func SafeArrayCopyData(psaSource *SafeArray, psaTarget *SafeArray) (err error) {
	debugPrint("Entering into safearray.safeArrayCopyData()...")
	hr, _, err := modOleAut.MustFindProc("safeArrayCopyData").Call(
		uintptr(unsafe.Pointer(psaSource)),
		uintptr(unsafe.Pointer(psaTarget)),
	)
	if err != syscall.Errno(0) {
		return err
	}
	if hr != S_OK {
		return fmt.Errorf("the safearray.SafeArrayCopyData() function return 0x%x", hr)
	}
	return nil
}

// SafeArrayCreate creates a new array descriptor, allocates and initializes the data for the array, and returns a pointer to the new array descriptor.
// SAFEARRAY * SafeArrayCreate(
//
//	VARTYPE        vt,
//	UINT           cDims,
//	SAFEARRAYBOUND *rgsabound
//
// );
// Varient types: https://docs.microsoft.com/en-us/windows/win32/api/wtypes/ne-wtypes-varenum
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraycreate
func SafeArrayCreate(vt VT, cDims uint32, rgsabound *SafeArrayBound) (*SafeArray, error) {
	debugPrint("Entering into safearray.SafeArrayCreate()...")
	ret, _, err := modOleAut.MustFindProc("SafeArrayCreate").Call(
		uintptr(vt),
		uintptr(cDims),
		uintptr(unsafe.Pointer(rgsabound)),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	return (*SafeArray)(unsafe.Pointer(ret)), nil
}

// SafeArrayCreateEx creates a new array descriptor, allocates and initializes the data for the array, and returns a pointer to the new array descriptor.
// SAFEARRAY * SafeArrayCreateEx(
//
//	[in] VARTYPE        vt,
//	[in] UINT           cDims,
//	[in] SAFEARRAYBOUND *rgsabound,
//	[in] PVOID          pvExtra
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraycreateex
func SafeArrayCreateEx(vt VT, cDims uint32, rgsabound *SafeArrayBound, pvExtra uintptr) (*SafeArray, error) {
	ret, _, err := modOleAut.MustFindProc("SafeArrayCreateEx").Call(
		uintptr(vt),
		uintptr(cDims),
		uintptr(unsafe.Pointer(rgsabound)),
		pvExtra,
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	if ret != S_OK {
		return nil, fmt.Errorf("the safearray.SafeArrayCreateEx() function return 0x%x and the SafeArray was not created", ret)
	}
	return (*SafeArray)(unsafe.Pointer(ret)), nil
}

// SafeArrayCreateVector creates SafeArray.
// SAFEARRAY * SafeArrayCreateVector(
//
//		[in] VARTYPE vt,
//		[in] LONG    lLbound,
//		[in] ULONG   cElements
//	 );
//
// AKA: SafeArrayCreateVector in Windows API.
func SafeArrayCreateVector(vt VT, lLbound int32, cElements uint32) (*SafeArray, error) {
	ret, _, err := modOleAut.MustFindProc("SafeArrayCreateVector").Call(
		uintptr(vt),
		uintptr(lLbound),
		uintptr(cElements),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	return (*SafeArray)(unsafe.Pointer(ret)), nil
}

// safeArrayCreateVectorEx creates SafeArray.
// SAFEARRAY * SafeArrayCreateVectorEx(
//
//	[in] VARTYPE vt,
//	[in] LONG    lLbound,
//	[in] ULONG   cElements,
//	[in] PVOID   pvExtra
//
// );
// AKA: SafeArrayCreateVectorEx in Windows API.
func SafeArrayCreateVectorEx(vt VT, lLbound int32, cElements uint32, pvExtra uintptr) (*SafeArray, error) {
	ret, _, err := modOleAut.MustFindProc("SafeArrayCreateVectorEx").Call(
		uintptr(vt),
		uintptr(lLbound),
		uintptr(cElements),
		pvExtra,
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	if ret != S_OK {
		return nil, fmt.Errorf("the safearray.SafeArrayCreateVectorEx() function return 0x%x and the SafeArray was not created", ret)
	}
	return (*SafeArray)(unsafe.Pointer(ret)), nil
}

// SafeArrayDestroy Destroys an existing array descriptor and all of the data in the array.
// If objects are stored in the array, Release is called on each object in the array.
// HRESULT SafeArrayDestroy(
//
//	SAFEARRAY *psa
//
// );
func SafeArrayDestroy(psa *SafeArray) error {
	debugPrint("Entering into safearray.SafeArrayDestroy()...")
	hr, _, err := modOleAut.MustFindProc("SafeArrayDestroy").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the oleaut32!SafeArrayDestroy function call returned an error:\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the oleaut32!SafeArrayDestroy function returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArrayDestroyData Destroys an existing array descriptor and all of the data in the array.
// If objects are stored in the array, Release is called on each object in the array.
// HRESULT SafeArrayDestroyData(
//
//	[in] SAFEARRAY *psa
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraydestroydata
func SafeArrayDestroyData(psa *SafeArray) error {
	debugPrint("Entering into safearray.SafeArrayDestroyData()...")
	hr, _, err := modOleAut.MustFindProc("SafeArrayDestroyData").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the oleaut32!SafeArrayDestroyData function call returned an error:\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the oleaut32!SafeArrayDestroyData function returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArrayDestroyDescriptor Destroys an existing array descriptor and all of the data in the array.
// If objects are stored in the array, Release is called on each object in the array.
// HRESULT SafeArrayDestroyDescriptor(
//
//	[in] SAFEARRAY *psa
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraydestroydescriptor
func SafeArrayDestroyDescriptor(psa *SafeArray) error {
	debugPrint("Entering into safearray.SafeArrayDestroyDescriptor()...")
	hr, _, err := modOleAut.MustFindProc("SafeArrayDestroyDescriptor").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return fmt.Errorf("the oleaut32!SafeArrayDestroyDescriptor function call returned an error:\n%s", err)
	}
	if hr != S_OK {
		return fmt.Errorf("the oleaut32!SafeArrayDestroyDescriptor function returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArrayGetDim returns the dimensions of a safearray
// UINT SafeArrayGetDim(
//
//	[in] SAFEARRAY *psa
//
// );
func SafeArrayGetDim(psa *SafeArray) (*uint32, error) {
	debugPrint("Entering into safearray.SafeArrayGetDim()...")
	ret, _, err := modOleAut.MustFindProc("SafeArrayGetDim").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the oleaut32!SafeArrayGetDim function call returned an error:\n%s", err)
	}
	return (*uint32)(unsafe.Pointer(&ret)), nil
}

// SafeArrayGetElement gets an element from the array at the given index
// HRESULT SafeArrayGetElement(
//
//		[in]  SAFEARRAY *psa,
//		[in]  LONG      *rgIndices,
//		[out] void      *pv
//	  );
func SafeArrayGetElement(psa *SafeArray, rgIndices uint32) (unsafe.Pointer, error) {
	debugPrint("Entering into safearray.SafeArrayGetElement()...")
	var pv unsafe.Pointer
	hr, _, err := modOleAut.MustFindProc("SafeArrayGetElement").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&rgIndices)),
		uintptr(unsafe.Pointer(&pv)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the oleaut32!SafeArrayGetElement function call returned an error:\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the oleaut32!SafeArrayGetElement function returned a non-zero HRESULT: 0x%x", hr)
	}
	return pv, nil
}

// SafeArrayGetElementString retrieves element at given index and converts to string.
func SafeArrayGetElementString(psa *SafeArray, rgIndices uint32) (string, error) {
	element, err := SafeArrayGetElement(psa, rgIndices)
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(element), nil
}

// SafeArrayGetElemsize returns the element size of the safearray in bytes
// UINT SafeArrayGetElemsize(
//
//	[in] SAFEARRAY *psa
//	);
func SafeArrayGetElemsize(psa *SafeArray) (*uint32, error) {
	debugPrint("Entering into safearray.SafeArrayGetElemsize()...")
	ret, _, err := modOleAut.MustFindProc("SafeArrayGetElemsize").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the oleaut32!SafeArrayGetElemsize function call returned an error:\n%s", err)
	}
	return (*uint32)(unsafe.Pointer(&ret)), nil
}

// SafeArrayGetIID is the InterfaceID of the elements in the SafeArray.
//
// HRESULT SafeArrayGetIID(
//
//	[in]  SAFEARRAY *psa,
//	[out] GUID      *pguid
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraygetiid
func SafeArrayGetIID(safearray *SafeArray) (*windows.GUID, error) {
	debugPrint("Entering into safearray.SafeArrayGetIID()...")
	var guid *windows.GUID
	hr, _, err := modOleAut.MustFindProc("SafeArrayGetIID").Call(
		uintptr(unsafe.Pointer(safearray)),
		uintptr(unsafe.Pointer(&guid)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("the oleaut32!SafeArrayGetIID function call returned an error:\n%s", err)
	}
	if hr != S_OK {
		return nil, fmt.Errorf("the oleaut32!SafeArrayGetIID function returned a non-zero HRESULT: 0x%x", hr)
	}
	return guid, nil
}

// SafeArrayGetLBound gets the lower bound for any dimension of the specified safe array
// HRESULT SafeArrayGetLBound(
//
//	SAFEARRAY *psa,
//	UINT      nDim,
//	LONG      *plLbound
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraygetlbound
func SafeArrayGetLBound(psa *SafeArray, nDim uint32) (uint32, error) {
	debugPrint("Entering into safearray.SafeArrayGetLBound()...")
	var plLbound uint32
	hr, _, err := modOleAut.MustFindProc("SafeArrayGetLBound").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(nDim),
		uintptr(unsafe.Pointer(&plLbound)),
	)
	if err != syscall.Errno(0) {
		return 0, err
	}
	if hr != S_OK {
		return 0, fmt.Errorf("the oleaut32!SafeArrayGetLBound function returned a non-zero HRESULT: 0x%x", hr)
	}
	return plLbound, nil
}

// SafeArrayGetRecordInfo accesses IRecordInfo info for custom types.
//
// HRESULT SafeArrayGetRecordInfo(
//
//	[in]  SAFEARRAY   *psa,
//	[out] IRecordInfo **prinfo
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraygetrecordinfo
func SafeArrayGetRecordInfo(psa *SafeArray) (interface{}, error) {
	debugPrint("Entering into safearray.SafeArrayGetRecordInfo()...")
	var prinfo interface{}
	hr, _, err := modOleAut.MustFindProc("SafeArrayGetRecordInfo").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&prinfo)),
	)
	if err != syscall.Errno(0) {
		return 0, err
	}
	if hr != S_OK {
		return 0, fmt.Errorf("the oleaut32!SafeArrayGetRecordInfo function returned a non-zero HRESULT: 0x%x", hr)
	}
	return prinfo, nil
}

// SafeArrayGetUBound gets the upper bound for any dimension of the specified safe array
// HRESULT SafeArrayGetUBound(
//
//	SAFEARRAY *psa,
//	UINT      nDim,
//	LONG      *plUbound
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraygetubound
func SafeArrayGetUBound(psa *SafeArray, nDim uint32) (uint32, error) {
	debugPrint("Entering into safearray.SafeArrayGetUBound()...")
	var plUbound uint32
	hr, _, err := modOleAut.MustFindProc("SafeArrayGetUBound").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(nDim),
		uintptr(unsafe.Pointer(&plUbound)),
	)
	if err != syscall.Errno(0) {
		return 0, err
	}
	if hr != S_OK {
		return 0, fmt.Errorf("the oleaut32!SafeArrayGetUBound function returned a non-zero HRESULT: 0x%x", hr)
	}
	return plUbound, nil
}

// SafeArrayGetVartype gets the VARTYPE stored in the specified safe array
// HRESULT SafeArrayGetVartype(
//
//	SAFEARRAY *psa,
//	VARTYPE   *pvt
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraygetvartype
func SafeArrayGetVartype(psa *SafeArray) (uint16, error) {
	debugPrint("Entering into safearray.SafeArrayGetVartype()...")
	var vt uint16
	hr, _, err := modOleAut.MustFindProc("SafeArrayGetVartype").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&vt)),
	)
	if err != syscall.Errno(0) {
		return 0, err
	}
	if hr != S_OK {
		return 0, fmt.Errorf("the OleAut32!SafeArrayGetVartype function returned a non-zero HRESULT: 0x%x", hr)
	}
	return vt, nil
}

// SafeArrayLock increments the lock count of an array, and places a pointer to the array data in pvData of the array descriptor
// HRESULT SafeArrayLock(
//
//	SAFEARRAY *psa
//
// );
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraylock
func SafeArrayLock(psa *SafeArray) error {
	debugPrint("Entering into safearray.SafeArrayLock()...")
	hr, _, err := modOleAut.MustFindProc("SafeArrayLock").Call(uintptr(unsafe.Pointer(psa)))
	if err != syscall.Errno(0) {
		return err
	}
	if hr != S_OK {
		return fmt.Errorf("the OleAut32!SafeArrayLock function returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArrayPutElement pushes an element to the safe array at a given index
//
//	 HRESULT SafeArrayPutElement(
//		  SAFEARRAY *psa,
//		  LONG      *rgIndices,
//		  void      *pv
//	 );
//
// https://docs.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearrayputelement
func SafeArrayPutElement(psa *SafeArray, rgIndices int32, pv unsafe.Pointer) error {
	debugPrint("Entering into safearray.SafeArrayPutElement()...")
	hr, _, err := modOleAut.MustFindProc("SafeArrayPutElement").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&rgIndices)),
		uintptr(pv),
	)
	if err != syscall.Errno(0) {
		return err
	}
	if hr != S_OK {
		return fmt.Errorf("the OleAut32!SafeArrayPutElement call returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArraySetRecordInfo mutates IRecordInfo info for custom types.
//
// HRESULT SafeArraySetRecordInfo(
//
//	[in] SAFEARRAY   *psa,
//	[in] IRecordInfo *prinfo
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearraysetrecordinfo
func SafeArraySetRecordInfo(psa *SafeArray, prinfo interface{}) error {
	debugPrint("Entering into safearray.SafeArraySetRecordInfo()...")
	hr, _, err := modOleAut.MustFindProc("SafeArraySetRecordInfo").Call(
		uintptr(unsafe.Pointer(psa)),
		uintptr(unsafe.Pointer(&prinfo)),
	)
	if err != syscall.Errno(0) {
		return err
	}
	if hr != S_OK {
		return fmt.Errorf("the OleAut32!SafeArraySetRecordInfo call returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArrayUnaccessData releases raw array.
//
// HRESULT SafeArrayUnaccessData(
//
//	[in] SAFEARRAY *psa
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearrayunaccessdata
func SafeArrayUnaccessData(psa *SafeArray) error {
	debugPrint("Entering into safearray.SafeArrayUnaccessData()...")
	hr, _, err := modOleAut.MustFindProc("SafeArrayUnaccessData").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return err
	}
	if hr != S_OK {
		return fmt.Errorf("the OleAut32!SafeArrayUnaccessData call returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

// SafeArrayUnloack releases raw array.
//
// HRESULT SafeArrayUnlock(
//
//	[in] SAFEARRAY *psa
//
// );
// https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-safearrayunlock
func SafeArrayUnlock(psa *SafeArray) error {
	debugPrint("Entering into safearray.SafeArrayUnlock()...")
	hr, _, err := modOleAut.MustFindProc("SafeArrayUnlock").Call(
		uintptr(unsafe.Pointer(psa)),
	)
	if err != syscall.Errno(0) {
		return err
	}
	if hr != S_OK {
		return fmt.Errorf("the OleAut32!SafeArrayUnlock call returned a non-zero HRESULT: 0x%x", hr)
	}
	return nil
}

func SafeArrayGetArrayLength(safeArray *SafeArray) (uint32, error) {
	upper, err := SafeArrayGetUBound(safeArray, 1)
	if err != nil {
		return 0, err
	}
	lower, err := SafeArrayGetLBound(safeArray, 1)
	if err != nil {
		return 0, err
	}
	if (upper - lower) == 0 {
		return 0, nil
	}
	return (upper - lower + 1), nil
}
