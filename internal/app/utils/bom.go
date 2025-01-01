package utils

import (
	"fmt"
	"os"
)

// Constants representing a UTF-8 BOM.
// https://en.wikipedia.org/wiki/Byte_order_mark#UTF-8
const (
	bomByte1 = 0xEF
	bomByte2 = 0xBB
	bomByte3 = 0xBF
	bomSize  = 3
)

// ReadFileWithUTF8BOM reads a file, returning its contents as bytes.
// If a UTF-8 BOM (0xEF, 0xBB, 0xBF) is present, it's removed.
func ReadFileWithUTF8BOM(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	if hasUTF8BOM(content) {
		return content[bomSize:], nil
	}
	return content, nil
}

// hasUTF8BOM checks if the provided data begins with the UTF-8 BOM bytes.
func hasUTF8BOM(data []byte) bool {
	return len(data) >= bomSize &&
		data[0] == bomByte1 &&
		data[1] == bomByte2 &&
		data[2] == bomByte3
}
