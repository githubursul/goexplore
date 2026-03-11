//go:build darwin && static

package gonative

/*
#cgo pkg-config: TestNativeLibrary
#include <TestNativeLibrary.h>
*/
import "C"

import (
	"fmt"
)

// Exports
func Factorial(n int) int {

	fmt.Println("I am a static library version of Factorial")
	result := int(C.factorial(C.int(n)))
	return result
}
