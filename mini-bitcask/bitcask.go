package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"sync"
)

const entryHeaderSize = 14

const (
	PUT uint16 = iota
	DEL
)

// Entry 表示 Bitcask 中的键值对
type Entry struct {
	Key       []byte
	Value     []byte
	KeySize   uint32
	ValueSize uint32
	Mark      uint16 // 墓碑
	CRC       uint32 // CRC 校验码
}

// 内存中的 key -> 文件偏移量的映射
type KeyDir map[string]int64

// Bitcask 结构，包含内存中的索引和文件
type Bitcask struct {
	mu     sync.Mutex
	file   *os.File
	keyDir KeyDir
}

// NewEntry 初始化一个新的 Entry
func NewEntry(key, value []byte, mark uint16) *Entry {
	return &Entry{
		Key:       key,
		Value:     value,
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Mark:      mark,
		CRC:       0, // CRC 初始值，稍后会计算
	}
}

// CalculateCRC 计算 Entry 的 CRC 校验码
func (e *Entry) CalculateCRC() uint32 {
	crc := crc32.NewIEEE()
	crc.Write(e.Key)
	crc.Write(e.Value)
	return crc.Sum32()
}

// Encode 编码 Entry，返回字节数组
func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, entryHeaderSize+e.KeySize+e.ValueSize)

	// 计算并设置 CRC
	e.CRC = e.CalculateCRC()

	// 编码 KeySize, ValueSize, Mark, 和 CRC
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	binary.BigEndian.PutUint16(buf[8:10], e.Mark)
	binary.BigEndian.PutUint32(buf[10:14], e.CRC)

	// 拷贝 Key 和 Value 到 buffer 中
	copy(buf[entryHeaderSize:entryHeaderSize+e.KeySize], e.Key)
	copy(buf[entryHeaderSize+e.KeySize:], e.Value)

	return buf, nil
}

// Decode 将字节数组解码为 Entry
func Decode(buf []byte) (*Entry, error) {
	if len(buf) < entryHeaderSize {
		return nil, fmt.Errorf("buffer size mismatch")
	}

	ks := binary.BigEndian.Uint32(buf[0:4])
	vs := binary.BigEndian.Uint32(buf[4:8])
	mark := binary.BigEndian.Uint16(buf[8:10])
	crc := binary.BigEndian.Uint32(buf[10:14])

	if len(buf) < int(entryHeaderSize+ks+vs) {
		return nil, fmt.Errorf("buffer size mismatch")
	}

	key := make([]byte, ks)
	value := make([]byte, vs)
	copy(key, buf[entryHeaderSize:entryHeaderSize+ks])
	copy(value, buf[entryHeaderSize+ks:])

	e := &Entry{
		Key:       key,
		Value:     value,
		KeySize:   ks,
		ValueSize: vs,
		Mark:      mark,
		CRC:       crc,
	}

	if e.CalculateCRC() != crc {
		return nil, fmt.Errorf("CRC mismatch")
	}

	return e, nil
}

// 写入 Entry 到文件，并更新内存中的偏移量
func (bc *Bitcask) Put(key, value []byte) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// 创建 Entry
	entry := NewEntry(key, value, PUT)
	data, err := entry.Encode()
	if err != nil {
		return err
	}

	// 写入文件
	offset, err := bc.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if _, err := bc.file.Write(data); err != nil {
		return err
	}

	// 更新内存中的索引
	bc.keyDir[string(key)] = offset

	return nil
}

// 删除 Entry（标记为删除）
func (bc *Bitcask) Delete(key []byte) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// 检查键是否存在
	if _, exists := bc.keyDir[string(key)]; !exists {
		return fmt.Errorf("key not found")
	}

	// 创建标记为删除的 Entry
	entry := NewEntry(key, nil, DEL)
	data, err := entry.Encode()
	if err != nil {
		return err
	}

	// 写入文件
	if _, err := bc.file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if _, err := bc.file.Write(data); err != nil {
		return err
	}

	// 从内存索引中删除
	delete(bc.keyDir, string(key))

	return nil
}

// 读取 Entry
func (bc *Bitcask) Get(key []byte) ([]byte, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	offset, exists := bc.keyDir[string(key)]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	// 从文件读取 Entry
	if _, err := bc.file.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	buf := make([]byte, entryHeaderSize)
	if _, err := io.ReadFull(bc.file, buf); err != nil {
		return nil, err
	}

	ks := binary.BigEndian.Uint32(buf[0:4])
	vs := binary.BigEndian.Uint32(buf[4:8])

	buf = append(buf, make([]byte, ks+vs)...)
	if _, err := io.ReadFull(bc.file, buf[entryHeaderSize:]); err != nil {
		return nil, err
	}

	entry, err := Decode(buf)
	if err != nil {
		return nil, err
	}

	if entry.Mark == DEL {
		return nil, fmt.Errorf("key has been deleted")
	}

	return entry.Value, nil
}

// 替代 os.Rename 的跨设备复制函数
func replaceFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	// 将源文件的内容复制到目标文件
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	// 确保数据写入磁盘
	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("error syncing destination file: %v", err)
	}

	// 删除源文件
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("error removing source file: %v", err)
	}

	return nil
}

// Merge 功能，将有效的 PUT 记录合并到新文件
func (bc *Bitcask) Merge() error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "bitcask-merge")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// 将所有 PUT 的 Entry 写入临时文件
	for _, offset := range bc.keyDir {
		if _, err := bc.file.Seek(offset, io.SeekStart); err != nil {
			return err
		}

		buf := make([]byte, entryHeaderSize)
		if _, err := io.ReadFull(bc.file, buf); err != nil {
			return err
		}

		ks := binary.BigEndian.Uint32(buf[0:4])
		vs := binary.BigEndian.Uint32(buf[4:8])

		buf = append(buf, make([]byte, ks+vs)...)
		if _, err := io.ReadFull(bc.file, buf[entryHeaderSize:]); err != nil {
			return err
		}

		entry, err := Decode(buf)
		if err != nil {
			return err
		}

		// 仅写入标记为 PUT 的 Entry
		if entry.Mark == PUT {
			data, err := entry.Encode()
			if err != nil {
				return err
			}

			if _, err := tempFile.Write(data); err != nil {
				return err
			}
		}
	}

	// 尝试用临时文件替换原始文件
	tempFileName := tempFile.Name()
	bc.file.Close()

	if err := os.Rename(tempFileName, bc.file.Name()); err != nil {
		// 如果跨设备重命名失败，执行手动复制和替换
		if err := replaceFile(tempFileName, bc.file.Name()); err != nil {
			return fmt.Errorf("error replacing file: %v", err)
		}
	}

	// 重建内存索引
	bc.keyDir = KeyDir{}
	newFile, err := os.OpenFile(bc.file.Name(), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	bc.file = newFile
	return bc.buildIndex()
}

// 构建内存中的 key -> offset 索引
func (bc *Bitcask) buildIndex() error {
	if _, err := bc.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	offset := int64(0)
	for {
		buf := make([]byte, entryHeaderSize)
		_, err := io.ReadFull(bc.file, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		ks := binary.BigEndian.Uint32(buf[0:4])
		vs := binary.BigEndian.Uint32(buf[4:8])

		buf = append(buf, make([]byte, ks+vs)...)
		if _, err := io.ReadFull(bc.file, buf[entryHeaderSize:]); err != nil {
			return err
		}

		entry, err := Decode(buf)
		if err != nil {
			return err
		}

		if entry.Mark == PUT {
			bc.keyDir[string(entry.Key)] = offset
		}

		offset += int64(entryHeaderSize + ks + vs)
	}

	return nil
}

// NewBitcask 初始化一个新的 Bitcask 实例
func NewBitcask(filename string) (*Bitcask, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	bc := &Bitcask{
		file:   file,
		keyDir: KeyDir{},
	}

	if err := bc.buildIndex(); err != nil {
		return nil, err
	}

	return bc, nil
}

func main() {
	// 初始化 Bitcask 数据库
	bitcask, err := NewBitcask("bitcask.db")
	if err != nil {
		fmt.Printf("Error initializing Bitcask: %v\n", err)
		return
	}

	// 1. 写入键值对
	fmt.Println("Inserting 'name' -> 'Alice'")
	if err := bitcask.Put([]byte("name"), []byte("Alice")); err != nil {
		fmt.Printf("Error writing to Bitcask: %v\n", err)
	}

	fmt.Println("Inserting 'address' -> 'Taiwan'")
	if err := bitcask.Put([]byte("address"), []byte("Taiwan")); err != nil {
		fmt.Printf("Error writing to Bitcask: %v\n", err)
	}

	// 2. 读取键值对，验证插入
	value, err := bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask: %v\n", err)
	} else {
		fmt.Printf("Value for 'name': %s\n", value) // 应输出 Alice
	}

	// 3. 更新键值对，验证更新
	fmt.Println("Updating 'name' -> 'Bob'")
	if err := bitcask.Put([]byte("name"), []byte("Bob")); err != nil {
		fmt.Printf("Error updating Bitcask: %v\n", err)
	}

	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask after update: %v\n", err)
	} else {
		fmt.Printf("Value for 'name' after update: %s\n", value) // 应输出 Bob
	}

	// 4. 删除键值对，验证删除
	fmt.Println("Deleting 'name'")
	if err := bitcask.Delete([]byte("name")); err != nil {
		fmt.Printf("Error deleting from Bitcask: %v\n", err)
	}

	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Expected error reading deleted 'name': %v\n", err) // 应输出错误信息
	} else {
		fmt.Printf("Unexpected value for deleted 'name': %s\n", value)
	}

	// 5. 再次插入键值对，验证重新插入
	fmt.Println("Inserting 'name' -> 'Charlie'")
	if err := bitcask.Put([]byte("name"), []byte("Charlie")); err != nil {
		fmt.Printf("Error writing to Bitcask: %v\n", err)
	}

	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask after re-insertion: %v\n", err)
	} else {
		fmt.Printf("Value for 'name' after re-insertion: %s\n", value) // 应输出 Charlie
	}

	// 6. 合并文件，验证合并操作
	fmt.Println("Merging entries...")
	if err := bitcask.Merge(); err != nil {
		fmt.Printf("Error during merge: %v\n", err)
	}

	// 7. 读取键值对，验证合并后数据是否仍然有效
	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask after merge: %v\n", err)
	} else {
		fmt.Printf("Value for 'name' after merge: %s\n", value) // 应输出 Charlie
	}

	// 8. 讀取 address
	value, err = bitcask.Get([]byte("address"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask: %v\n", err)
	} else {
		fmt.Printf("Value for 'address': %s\n", value)
	}
}
