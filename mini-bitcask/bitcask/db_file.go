package bitcask

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// DBFile 用來表示 Bitcask 中的一個數據檔案
type DBFile struct {
	File          *os.File
	Offset        int64
	HeaderBufPool *sync.Pool
	MaxSize       int64 // 文件大小上限，用於觸發分割
}

// NewDBFile 創建一個新的數據檔案
func NewDBFile(path string, maxSize int64, fileName string) (*DBFile, error) {
	fullPath := filepath.Join(path, fileName)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	pool := &sync.Pool{New: func() interface{} {
		return make([]byte, entryHeaderSize) // 預設Header大小
	}}

	fmt.Printf("file stat size: %d\n", stat.Size())
	return &DBFile{
		File:          file,
		Offset:        stat.Size(),
		HeaderBufPool: pool,
		MaxSize:       maxSize,
	}, nil
}

// Write 寫入 Entry 到文件
// Write 寫入 Entry 到文件
func (df *DBFile) Write(e *Entry) (offset int64, err error) {
	enc, err := e.Encode()
	if err != nil {
		fmt.Printf("encode error: %v\n", err)
		return 0, err
	}

	// 寫入數據
	n, err := df.File.WriteAt(enc, df.Offset)
	if err != nil {
		fmt.Printf("write error: %v\n", err)
		return 0, err
	}

	// 確保數據被同步到磁盤
	err = df.File.Sync()
	if err != nil {
		fmt.Printf("sync error: %v\n", err)
		return 0, err
	}

	// 更新文件偏移量
	df.Offset += int64(n) // 用寫入的字節數更新偏移量
	return df.Offset, nil
}

// Read 從 offset 讀取 Entry
// Read 從 offset 讀取 Entry
func (df *DBFile) Read(offset int64) (*Entry, error) {
	fmt.Printf("Reading from offset: %d\n", offset)

	// 读取 Header
	buf := make([]byte, df.MaxSize)
	n, err := df.File.ReadAt(buf, offset)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Read %d bytes for header at offset %d\n", n, offset)

	// 解码 Header
	e, err := Decode(buf)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	// 根據 KeySize 和 ValueSize 來計算 Entry 的總大小
	entrySize := int64(entryHeaderSize + e.KeySize + e.ValueSize)

	// 重新读取完整的 Entry 数据（Header + Key + Value）
	fullBuf := make([]byte, entrySize)
	n, err = df.File.ReadAt(fullBuf, offset)
	if err != nil {
		fmt.Printf("read error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Read %d bytes for full entry at offset %d\n", n, offset)

	// 重新解码完整的 Entry
	e, err = Decode(fullBuf)
	if err != nil {
		fmt.Printf("Decode full entry error: %v\n", err)
		return nil, err
	}

	return e, nil
}

// Split 如果文件大小超過限制，則觸發文件分割
func (df *DBFile) Split() bool {
	return df.Offset >= df.MaxSize
}
