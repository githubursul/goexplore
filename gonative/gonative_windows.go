//go:build windows && amd64

package gonative

import (
	_ "embed"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

type nativeAPI struct {
	factorial *windows.LazyProc
}

//go:embed TestNativeLibrary.dll
var dllBytes []byte
var api *nativeAPI

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

	api = &nativeAPI{
		factorial: factorialProc,
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

// Exports

func Factorial(n int) int {
	initNativeAPI()

	// Call the factorial function from the DLL
	result, _, lastErr := api.factorial.Call(uintptr(n))
	if result == 0 {
		panic("Native DLL call failed, lastErr: " + lastErr.Error())
	}

	return int(result)
}
