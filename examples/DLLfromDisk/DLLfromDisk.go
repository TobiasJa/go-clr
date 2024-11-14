//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"syscall"

	clr "github.com/tobiasja/go-clr"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println(clr.CLSID_CLRMetaHost)
	fmt.Println(clr.IID_ICLRMetaHost)
	metahost, err := clr.CLRCreateInstance(clr.CLSID_CLRMetaHost, clr.IID_ICLRMetaHost)
	must(err)
	fmt.Println("[+] Got metahost")

	installedRuntimes, err := clr.GetInstalledRuntimes(metahost)
	must(err)
	fmt.Printf("[+] Found installed runtimes: %s\n", installedRuntimes)
	versionString := "v4.0.30319"
	pwzVersion, err := syscall.UTF16PtrFromString(versionString)
	must(err)

	runtimeInfo, err := metahost.GetRuntime(pwzVersion, clr.IID_ICLRRuntimeInfo)
	must(err)
	fmt.Printf("[+] Using runtime: %s\n", versionString)

	isLoadable, err := runtimeInfo.IsLoadable()
	must(err)
	if !isLoadable {
		log.Fatal("[!] IsLoadable returned false. Bailing...")
	}

	runtimeHost, err := runtimeInfo.GetInterface(clr.CLSID_CLRRuntimeHost, clr.IID_ICLRRuntimeHost)
	must(err)

	fmt.Println("[+] Start Runtime Host")
	err = runtimeHost.(*clr.ICLRRuntimeHost).Start()
	must(err)
	fmt.Println("[+] Loaded CLR into this process")

	fmt.Println("[+] Executing assembly...")
	pDLLPath, err := syscall.UTF16PtrFromString("TestDLL.dll")
	must(err)
	pTypeName, err := syscall.UTF16PtrFromString("TestDLL.HelloWorld")
	must(err)
	pMethodName, err := syscall.UTF16PtrFromString("SayHello")
	must(err)
	pArgument, err := syscall.UTF16PtrFromString("foobar")
	must(err)
	ret, err := runtimeHost.(*clr.ICLRRuntimeHost).ExecuteInDefaultAppDomain(
		pDLLPath,
		pTypeName,
		pMethodName,
		pArgument,
	)
	if *ret != 0 {
		err = fmt.Errorf("the ICLRRuntimeHost::ExecuteInDefaultAppDomain method returned a non-zero return value: %d", *ret)
	}
	must(err)
}
