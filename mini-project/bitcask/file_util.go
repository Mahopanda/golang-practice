package bitcask

import (
	"fmt"
	"io"
	"os"
)

// 文件操作工具函數
// ReplaceFile 替代 os.Rename 的跨設備複製函數
func ReplaceFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("error syncing destination file: %v", err)
	}

	return os.Remove(src)
}
