//go:build darwin && arm64 && !static && shared

package gonative

/*
#cgo pkg-config: TestNativeLibrary-Shared
#cgo LDFLAGS: -Wl,-rpath,@executable_path -Wl,-rpath,@executable_path/lib -Wl,-rpath,@loader_path -Wl,-rpath,/opt/homebrew/lib
#include <TestNativeLibrary.h>
*/
import "C"

import (
	"fmt"
)

// Exports
func Factorial(n int) int {

	fmt.Println("I am a shared library version of Factorial via pkg-config")
	result := int(C.factorial(C.int(n)))
	return result
}
