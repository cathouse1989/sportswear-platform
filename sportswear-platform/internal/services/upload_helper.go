package services

import (
	"io"
	"os"
)

// mkdirAll 创建目录（递归）
func mkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

// saveFile 保存文件到本地磁盘
func saveFile(file io.Reader, destPath string) error {
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, file)
	return err
}
