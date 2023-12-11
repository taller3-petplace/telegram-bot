package utils

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

func ReadFileWithPath(configFilePath string, suffixToRemove string) ([]byte, error) {
	_, file, _, _ := runtime.Caller(1)
	filePath := strings.TrimSuffix(file, suffixToRemove)
	filePath += configFilePath

	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}

	configFileBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	return configFileBytes, nil
}

func Contains[T comparable](elements []T, target T) bool {
	for idx := range elements {
		if elements[idx] == target {
			return true
		}
	}

	return false
}
