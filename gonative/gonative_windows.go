//go:build windows && amd64

package gonative

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type nativeAPI struct {
	factorial        *windows.LazyProc
	set_log_callback *windows.LazyProc
}

//go:embed TestNativeLibrary.dll
var dllBytes []byte
var api *nativeAPI
var loggingCallback uintptr

func initNativeAPI() {

	var dll *windows.LazyDLL
	var err error
	if api == nil {
		dll, err = loadLibrary()
		if err != nil {
			panic(err)
		}
		if dll == nil {
			panic("Failed to load DLL, but no error returned")
		}
	}

	factorialProc := dll.NewProc("factorial")
	err = factorialProc.Find()
	if err != nil {
		panic("Failed to find 'factorial' procedure in DLL: " + err.Error())
	}

	setLogCallbackProc := dll.NewProc("set_log_callback")
	err = setLogCallbackProc.Find()
	if err != nil {
		panic("Failed to find 'set_log_callback' procedure in DLL: " + err.Error())
	}

	api = &nativeAPI{
		factorial:        factorialProc,
		set_log_callback: setLogCallbackProc,
	}
}

func loadLibrary() (*windows.LazyDLL, error) {
	// Create a temporary file to write the embedded DLL
	tempDir := os.TempDir()
	dllPath := filepath.Join(tempDir, "TestNativeLibrary-tmp.dll")

	// Write the embedded DLL bytes to the temporary file
	err := os.WriteFile(dllPath, dllBytes, 0644)
	if err != nil {
		return nil, err
	}

	// Lazily load the DLL from the temporary file
	// Use absolute path to counter dll preloading attacks
	dll := windows.NewLazyDLL(dllPath)

	err = dll.Load()
	if err != nil {
		return nil, err
	}

	return dll, nil
}

func utf8CStringToString(p *byte) string {
	if p == nil {
		return ""
	}
	// Find the length up to the first 0 byte.
	var n int
	for *(*byte)(unsafe.Add(unsafe.Pointer(p), n)) != 0 {
		n++
	}
	// Convert to a slice and then to string.
	bytes := unsafe.Slice(p, n)
	return string(bytes)
}

func enableLogging() {
	loggingCallback = syscall.NewCallback(func(logLevel int, message *byte, pii int) uintptr {
		logMessage := utf8CStringToString(message)
		fmt.Printf("[Level %d] %s\n", logLevel, logMessage)
		return 0
	})

	api.set_log_callback.Call(loggingCallback)
}

// Exports

func Factorial(n int) int {
	initNativeAPI()
	enableLogging()

	// Call the factorial function from the DLL
	result, _, lastErr := api.factorial.Call(uintptr(n))
	if result == 0 {
		panic("Native DLL call failed, lastErr: " + lastErr.Error())
	}

	return int(result)
}
