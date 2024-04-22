package main

import "runtime"

// getExecutableFileExtention get the appropriate executable file type for the machine
func getExecutableFileExtention() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	return os + arch
}
