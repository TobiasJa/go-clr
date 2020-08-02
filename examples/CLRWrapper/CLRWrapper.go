// +build windows

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"

	clr "github.com/ropnop/go-clr"
)

func main() {
	/*fmt.Println("[+] Loading DLL from Disk")
	ret, err := clr.ExecuteDLLFromDisk(
		"TestDLL.dll",
		"TestDLL.HelloWorld",
		"SayHello",
		"foobar")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] DLL Return Code: %d\n", ret)
	*/
	fmt.Println("[+] Executing EXE from memory")
	exebytes, err := ioutil.ReadFile(`C:\Users\ayoul3\Documents\go\reflect-pe\res\managed.exe`)
	if err != nil {
		log.Fatal(err)
	}
	runtime.KeepAlive(exebytes)

	ret2, err := clr.ExecuteByteArray(exebytes, []string{"test", "test2"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] EXE Return Code: %d\n", ret2)
}
