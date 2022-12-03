package utils

import (
	"fmt"
	"strings"
)

func GetPartDir(filename string) string {
	return strings.SplitN(filename, ".", 2)[0]
}

func GetPartFilename(filename string, partNum int) string {
	partDir := GetPartDir(filename)
	return fmt.Sprintf("%s/%s-%d", partDir, filename, partNum)
}
