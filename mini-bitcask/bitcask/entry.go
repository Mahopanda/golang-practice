package bitcask

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
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
	CRC       uint32 //  CRC 校驗碼
}

// NewEntry 初始化一個新的 Entry
func NewEntry(key, value []byte, mark uint16) *Entry {
	return &Entry{
		Key:       key,
		Value:     value,
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Mark:      mark,
		CRC:       0, // CRC 初始值，稍後會計算
	}
}

// GetSize 返回 Entry 編碼後的大小
func (e *Entry) GetSize() int64 {
	return int64(entryHeaderSize + e.KeySize + e.ValueSize)
}

// CalculateCRC 計算 Entry 的 CRC 校驗碼
func (e *Entry) CalculateCRC() uint32 {
	// 使用簡單的 CRC32 校驗數據完整性
	crc := crc32.NewIEEE()
	crc.Write(e.Key)
	crc.Write(e.Value)
	e.CRC = crc.Sum32()
	return e.CRC
}

// Encode 編碼 Entry，返回字節數組
func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, entryHeaderSize+e.KeySize+e.ValueSize)

	// 编码 KeySize 和 ValueSize
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	binary.BigEndian.PutUint16(buf[8:10], e.Mark)
	binary.BigEndian.PutUint32(buf[10:14], e.CalculateCRC())

	// 拷贝 Key 和 Value 到 buffer 中
	copy(buf[entryHeaderSize:entryHeaderSize+e.KeySize], e.Key)
	copy(buf[entryHeaderSize+e.KeySize:], e.Value)
	crc := e.CalculateCRC()
	// 打印编码后的字节数组
	fmt.Printf("Encoded Entry: %v\n", buf)
	fmt.Printf("Encoded Entry: KeySize=%d, ValueSize=%d, CRC=%d\n", buf[entryHeaderSize:entryHeaderSize+e.KeySize], buf[entryHeaderSize+e.KeySize:], crc)

	return buf, nil
}

// Decode 将字节数组解码为 Entry
func Decode(buf []byte) (*Entry, error) {
	// 检查缓冲区长度是否足够解析出 Header
	if len(buf) < entryHeaderSize {
		fmt.Printf("buffer size mismatch: expected %d bytes, got %d\n", entryHeaderSize, len(buf))
		return nil, fmt.Errorf("buffer size mismatch: expected %d bytes, got %d", entryHeaderSize, len(buf))
	}

	// 解码 KeySize 和 ValueSize
	ks := binary.BigEndian.Uint32(buf[0:4])
	vs := binary.BigEndian.Uint32(buf[4:8])
	mark := binary.BigEndian.Uint16(buf[8:10])
	crc := binary.BigEndian.Uint32(buf[10:14])

	// 打印解析出的 KeySize 和 ValueSize
	fmt.Printf("Decoding Entry: KeySize=%d, ValueSize=%d, CRC=%d\n", ks, vs, crc)

	// 检查缓冲区长度是否足够包含 Key 和 Value
	if len(buf) < int(entryHeaderSize+ks+vs) {
		return nil, fmt.Errorf("buffer size mismatch: expected %d bytes, got %d", entryHeaderSize+ks+vs, len(buf))
	}

	// 提取 Key 和 Value
	key := make([]byte, ks)
	value := make([]byte, vs)

	// copy(key, buf[entryHeaderSize:entryHeaderSize+ks])
	// copy(value, buf[entryHeaderSize+ks:])

	// 打印解码后的 Key 和 Value
	fmt.Printf("Decoded Key: %s, Decoded Value: %s\n", string(key), string(value))

	// 创建 Entry 对象
	e := &Entry{

		KeySize:   ks,
		ValueSize: vs,
		Mark:      mark,
	}

	// 验证 CRC 校验码
	if e.CalculateCRC() != crc {
		return nil, fmt.Errorf("CRC mismatch: expected %d, got %d", crc, e.CalculateCRC())
	}

	return e, nil
}
