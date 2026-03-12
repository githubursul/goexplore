//go:build darwin && static && !shared

package gonative

/*
#cgo pkg-config: TestNativeLibrary-Static
#include <TestNativeLibrary.h>
*/
import "C"

import (
	"fmt"
)

// Exports
func Factorial(n int) int {

	fmt.Println("I am a static library version of Factorial via pkg-config")
	result := int(C.factorial(C.int(n)))
	return result
}
