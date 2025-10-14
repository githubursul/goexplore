//go:build darwin && arm64

package gonative

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>

#ifndef GONATIVE_MACOS_H
#define GONATIVE_MACOS_H

typedef void (*plog_callback)(int, const char*, int);
typedef int (*pfactorial_proc)(int);
typedef void (*pset_log_callback_proc)(plog_callback);

extern void goLog(int level, char* message, int pii);

void log_callback(int level, const char* message, int pii) {
    goLog(level, (char*)message, pii);
}

int call_factorial_proc(pfactorial_proc f, int n) {
    if (f == NULL) {
        return 0;
    }
    return f(n);
}

int call_set_log_callback_proc(pset_log_callback_proc f, plog_callback cb) {
    if (f == NULL) {
        return 0;
    }
    f(cb);
    return 1;
}

#endif
*/
import "C"

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type nativeAPI struct {
	factorial        unsafe.Pointer
	set_log_callback unsafe.Pointer
}

//go:embed libTestNativeLibrary.dylib
var dylibBytes []byte
var dylib unsafe.Pointer
var api *nativeAPI

func initNativeAPI() {
	var err error
	if api == nil {
		dylib, err = loadLibrary()
		if err != nil {
			panic(err)
		}
		if dylib == nil {
			panic("Failed to load a shared library, but no error returned")
		}
	}

	factorialFuncNameC := C.CString("factorial")
	defer C.free(unsafe.Pointer(factorialFuncNameC))
	factorialSymbol := C.dlsym(dylib, factorialFuncNameC)
	if factorialSymbol == nil {
		panic("Failed to find 'factorial' symbol in shared library")
	}

	setLogCallbackFuncNameC := C.CString("set_log_callback")
	defer C.free(unsafe.Pointer(setLogCallbackFuncNameC))
	setLogCallbackSymbol := C.dlsym(dylib, setLogCallbackFuncNameC)

	api = &nativeAPI{
		factorial:        factorialSymbol,
		set_log_callback: setLogCallbackSymbol,
	}
}

func loadLibrary() (unsafe.Pointer, error) {
	// Create a temporary file to write the embedded DLL
	tempDir := os.TempDir()
	dylibPath := filepath.Join(tempDir, "TestNativeLibrary-tmp.dylib")

	// Write the embedded DLL bytes to the temporary file
	err := os.WriteFile(dylibPath, dylibBytes, 0644)
	if err != nil {
		return nil, err
	}

	dylibPathC := C.CString(dylibPath)
	defer C.free(unsafe.Pointer(dylibPathC))
	dylibHandle := C.dlopen(dylibPathC, C.RTLD_NOW)

	if dylibHandle == nil {
		return nil, fmt.Errorf("failed to load library: %s", dylibPath)
	}

	return dylibHandle, nil
}

//export goLog
func goLog(level C.int, message *C.char, pii C.int) {
	fmt.Printf("[Level %d]: %s\n", level, C.GoString(message))
}

func enableLogging() {
	C.call_set_log_callback_proc((C.pset_log_callback_proc)(api.set_log_callback), (C.plog_callback)(C.log_callback))
}

// Exports
func Factorial(n int) int {
	initNativeAPI()
	enableLogging()
	result := int(C.call_factorial_proc((C.pfactorial_proc)(api.factorial), C.int(n)))
	if result == 0 {
		panic("Native DLL call failed")
	}
	return result
}
