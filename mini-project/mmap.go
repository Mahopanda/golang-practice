package main

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// mmap 函數，將文件映射到內存中
func mmap(file *os.File, length int) ([]byte, error) {
	// 使用 syscall 設置內存映射
	data, err := syscall.Mmap(int(file.Fd()), 0, length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("failed to mmap file: %v", err)
	}
	return data, nil
}

// munmap 函數，用於解除內存映射
func munmap(data []byte) error {
	return syscall.Munmap(data)
}

// 使用 unsafe 來查看數據指針（示例用途，不推薦在生產環境使用）
func modifyDataUnsafe(data []byte, newData string) {
	copy(data, newData)
}

func main() {
	// 打開文件
	file, err := os.OpenFile("example.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// 確保文件有足夠大小
	if _, err := file.Write([]byte("Hello, mmap!")); err != nil {
		fmt.Println("Failed to write initial data:", err)
		return
	}

	// 獲取文件長度
	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file info:", err)
		return
	}
	length := int(stat.Size())

	// 創建 mmap
	data, err := mmap(file, length)
	if err != nil {
		fmt.Println("Failed to mmap file:", err)
		return
	}
	defer munmap(data)

	// 讀取數據
	fmt.Println("File content:", string(data))

	// 修改 mmap 中的內容（將同步到文件中）
	modifyDataUnsafe(data, "Hi, mmap!")
	fmt.Println("Modified content:", string(data))

	// 確保文件保存
	if err := unix.Msync(data, unix.MS_SYNC); err != nil {
		fmt.Println("Failed to sync mmap:", err)
		return
	}

	fmt.Println("Data modified successfully in the mmap and file.")
}
