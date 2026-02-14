//go:build darwin && static

package gonative

/*
#cgo LDFLAGS: -L${NATIVE_RUNTIME_DIR} -lTestNativeLibraryStatic
#cgo CFLAGS: -I${NATIVE_RUNTIME_DIR}

#include <TestNativeLibrary.h>
*/

import "C"

// Exports
func Factorial(n int) int {

	result := int(C.factorial(C.int(n)))
	return result
}
