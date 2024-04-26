package utils

import (
	"fmt"
	"runtime"
)

func GetExpectedAssetName(packageName, version string) string {
	return fmt.Sprintf("%s-%s-%s-%s", packageName, version, runtime.GOOS, runtime.GOARCH)
}
