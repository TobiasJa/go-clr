//go:build windows
// +build windows

package clr

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// AppDomainSetup is a Windows COM object interface pointer for the .NET AppDomain class.
// The AppDomain object represents an application domain, which is an isolated environment where applications execute.
// This structure only contains a pointer to the AppDomain's virtual function table
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain?view=netframework-4.8
type AppDomainSetup struct {
	vtbl *AppDomainSetupVtbl
}

// AppDomainSetupVtbl is a Virtual Function Table for the AppDomain COM interface
// The Virtual Function Table contains pointers to the COM IUnkown interface
// functions (QueryInterface, AddRef, & Release) as well as the AppDomain object's methods
// https://docs.microsoft.com/en-us/dotnet/api/system.appdomain?view=netframework-4.8
type AppDomainSetupVtbl struct {
	QueryInterface            uintptr
	AddRef                    uintptr
	Release                   uintptr
	get_ApplicationBase       uintptr
	put_ApplicationBase       uintptr
	get_ApplicationName       uintptr
	put_ApplicationName       uintptr
	get_CachePath             uintptr
	put_CachePath             uintptr
	get_ConfigurationFile     uintptr
	put_ConfigurationFile     uintptr
	get_DynamicBase           uintptr
	put_DynamicBase           uintptr
	get_LicenseFile           uintptr
	put_LicenseFile           uintptr
	get_PrivateBinPath        uintptr
	put_PrivateBinPath        uintptr
	get_PrivateBinPathProbe   uintptr
	put_PrivateBinPathProbe   uintptr
	get_ShadowCopyDirectories uintptr
	put_ShadowCopyDirectories uintptr
	get_ShadowCopyFiles       uintptr
	put_ShadowCopyFiles       uintptr
}

func (obj *AppDomainSetup) QueryInterface(riid windows.GUID) (unsafe.Pointer, error) {
	debugPrint("Entering into AppDomainSetup.QueryInterface()...")
	var ppvObject unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::QueryInterface").CheckHResultError(syscall.SyscallN(
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

func (obj *AppDomainSetup) AddRef() error {
	debugPrint("Entering into AppDomainSetup.AddRef()...")
	return NewHResultChecker("AppDomainSetup::AddRef").CheckHResultError(syscall.SyscallN(
		obj.vtbl.AddRef,
		uintptr(unsafe.Pointer(obj)),
	))
}

func (obj *AppDomainSetup) Release() error {
	debugPrint("Entering into AppDomainSetup.Release()...")
	return NewHResultChecker("AppDomainSetup::Release").CheckHResultError(syscall.SyscallN(
		obj.vtbl.Release,
		uintptr(unsafe.Pointer(obj)),
	))
}

func (obj *AppDomainSetup) GetApplicationBase() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetApplicationBase()...")
	var basePtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetApplicationBase").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_ApplicationBase,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&basePtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(basePtr)), nil
}

func (obj *AppDomainSetup) PutApplicationBase(base string) error {
	debugPrint("Entering into AppDomainSetup.PutApplicationBase()...")
	basePtr, err := SysAllocString(base)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutApplicationBase SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutApplicationBase").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_ApplicationBase,
		uintptr(unsafe.Pointer(obj)),
		uintptr(basePtr),
	))
}

func (obj *AppDomainSetup) GetApplicationName() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetApplicationName()...")
	var namePtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetApplicationName").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_ApplicationName,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&namePtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(namePtr)), nil
}

func (obj *AppDomainSetup) PutApplicationName(name string) error {
	debugPrint("Entering into AppDomainSetup.PutApplicationName()...")
	namePtr, err := SysAllocString(name)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutApplicationName SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutApplicationName").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_ApplicationName,
		uintptr(unsafe.Pointer(obj)),
		uintptr(namePtr),
	))
}

func (obj *AppDomainSetup) GetCachePath() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetCachePath()...")
	var cachePathPtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetCachePath").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_CachePath,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&cachePathPtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(cachePathPtr)), nil
}

func (obj *AppDomainSetup) PutCachePath(cachePath string) error {
	debugPrint("Entering into AppDomainSetup.PutCachePath()...")
	cachePathPtr, err := SysAllocString(cachePath)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutCachePath SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutCachePath").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_CachePath,
		uintptr(unsafe.Pointer(obj)),
		uintptr(cachePathPtr),
	))
}

func (obj *AppDomainSetup) GetConfigurationFile() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetConfigurationFile()...")
	var configurationFilePtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetConfigurationFile").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_ConfigurationFile,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&configurationFilePtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(configurationFilePtr)), nil
}

func (obj *AppDomainSetup) PutConfigurationFile(configurationFile string) error {
	debugPrint("Entering into AppDomainSetup.PutConfigurationFile()...")
	configurationFilePtr, err := SysAllocString(configurationFile)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutConfigurationFile SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutConfigurationFile").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_ConfigurationFile,
		uintptr(unsafe.Pointer(obj)),
		uintptr(configurationFilePtr),
	))
}

func (obj *AppDomainSetup) GetDynamicBase() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetDynamicBase()...")
	var dynamicBasePtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetDynamicBase").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_DynamicBase,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&dynamicBasePtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(dynamicBasePtr)), nil
}

func (obj *AppDomainSetup) PutDynamicBase(dynamicBase string) error {
	debugPrint("Entering into AppDomainSetup.PutDynamicBase()...")
	dynamicBasePtr, err := SysAllocString(dynamicBase)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutDynamicBase SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutDynamicBase").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_DynamicBase,
		uintptr(unsafe.Pointer(obj)),
		uintptr(dynamicBasePtr),
	))
}

func (obj *AppDomainSetup) GetLicenseFile() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetLicenseFile()...")
	var licenseFilePtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetLicenseFile").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_LicenseFile,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&licenseFilePtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(licenseFilePtr)), nil
}

func (obj *AppDomainSetup) PutLicenseFile(licenseFile string) error {
	debugPrint("Entering into AppDomainSetup.PutLicenseFile()...")
	licenseFilePtr, err := SysAllocString(licenseFile)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutLicenseFile SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutLicenseFile").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_LicenseFile,
		uintptr(unsafe.Pointer(obj)),
		uintptr(licenseFilePtr),
	))
}

func (obj *AppDomainSetup) GetPrivateBinPath() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetPrivateBinPath()...")
	var privateBinPathPtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetPrivateBinPath").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_PrivateBinPath,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&privateBinPathPtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(privateBinPathPtr)), nil
}

func (obj *AppDomainSetup) PutPrivateBinPath(privateBinPath string) error {
	debugPrint("Entering into AppDomainSetup.PutPrivateBinPath()...")
	privateBinPathPtr, err := SysAllocString(privateBinPath)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutPrivateBinPath SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutPrivateBinPath").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_PrivateBinPath,
		uintptr(unsafe.Pointer(obj)),
		uintptr(privateBinPathPtr),
	))
}

func (obj *AppDomainSetup) GetPrivateBinPathProbe() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetPrivateBinPathProbe()...")
	var privateBinPathProbePtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetPrivateBinPathProbe").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_PrivateBinPathProbe,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&privateBinPathProbePtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(privateBinPathProbePtr)), nil
}

func (obj *AppDomainSetup) PutPrivateBinPathProbe(privateBinPathProbe string) error {
	debugPrint("Entering into AppDomainSetup.PutPrivateBinPathProbe()...")
	privateBinPathProbePtr, err := SysAllocString(privateBinPathProbe)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutPrivateBinPathProbe SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutPrivateBinPathProbe").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_PrivateBinPathProbe,
		uintptr(unsafe.Pointer(obj)),
		uintptr(privateBinPathProbePtr),
	))
}

func (obj *AppDomainSetup) GetShadowCopyDirectories() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetShadowCopyDirectories()...")
	var shadowCopyDirectoriesPtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetShadowCopyDirectories").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_ShadowCopyDirectories,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&shadowCopyDirectoriesPtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(shadowCopyDirectoriesPtr)), nil
}

func (obj *AppDomainSetup) PutShadowCopyDirectories(shadowCopyDirectories string) error {
	debugPrint("Entering into AppDomainSetup.PutShadowCopyDirectories()...")
	shadowCopyDirectoriesPtr, err := SysAllocString(shadowCopyDirectories)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutShadowCopyDirectories SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutShadowCopyDirectories").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_ShadowCopyDirectories,
		uintptr(unsafe.Pointer(obj)),
		uintptr(shadowCopyDirectoriesPtr),
	))
}

func (obj *AppDomainSetup) GetShadowCopyFiles() (string, error) {
	debugPrint("Entering into AppDomainSetup.GetShadowCopyFiles()...")
	var shadowCopyFilesPtr unsafe.Pointer
	err := NewHResultChecker("AppDomainSetup::GetShadowCopyFiles").CheckHResultError(syscall.SyscallN(
		obj.vtbl.get_ShadowCopyFiles,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&shadowCopyFilesPtr)),
	))
	if err != nil {
		return "", err
	}
	return ReadUnicodeStr(unsafe.Pointer(shadowCopyFilesPtr)), nil
}

func (obj *AppDomainSetup) PutShadowCopyFiles(shadowCopyFiles string) error {
	debugPrint("Entering into AppDomainSetup.PutShadowCopyFiles()...")
	shadowCopyFilesPtr, err := SysAllocString(shadowCopyFiles)
	if err != nil {
		return fmt.Errorf("the AppDomainSetup::PutShadowCopyFiles SysAllocString returned error: %v", err)
	}
	return NewHResultChecker("AppDomainSetup::PutShadowCopyFiles").CheckHResultError(syscall.SyscallN(
		obj.vtbl.put_ShadowCopyFiles,
		uintptr(unsafe.Pointer(obj)),
		uintptr(shadowCopyFilesPtr),
	))
}
