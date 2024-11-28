//go:build windows

package systeminit

import (
	"fmt"
	"github.com/iwdgo/fileattributes"
)

func hidePath(filepath string) error {
	if !IsPathExists(filepath) {
		return fmt.Errorf("资源路径不存在")
	}

	return fileattributes.SetFileAttributes(filepath, fileattributes.FILE_ATTRIBUTE_HIDDEN)
}
