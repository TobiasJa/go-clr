//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	clr "github.com/tobiasja/go-clr"
)

func main() {
	fmt.Println("[+] Loading DLL from Disk")
	ret, err := clr.ExecuteDLLFromDisk(
		"v4",
		"TestDLL.dll",
		"TestDLL.HelloWorld",
		"SayHello",
		"foobar")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] DLL Return Code: %d\n", ret)

	fmt.Println("[+] Executing EXE from memory")
	exebytes, err := os.ReadFile("helloworld.exe")
	if err != nil {
		log.Fatal(err)
	}
	runtime.KeepAlive(exebytes)

	ret2, err := clr.ExecuteByteArray("v2", exebytes, []string{"test", "test2"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] EXE Return Code: %d\n", ret2)
}
