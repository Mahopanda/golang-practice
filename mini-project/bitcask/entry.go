package bitcask

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const entryHeaderSize = 14

type EntryType uint16

const (
	PUT EntryType = iota
	DEL
)

type Entry struct {
	Key       []byte
	Value     []byte
	KeySize   uint32
	ValueSize uint32
	Mark      EntryType // 墓碑，用於標記是否已刪除
	CRC       uint32    // CRC 校驗碼
}

// NewEntry 初始化並返回一個新的 Entry
func NewEntry(key, value []byte, mark EntryType) *Entry {
	return &Entry{
		Key:       key,
		Value:     value,
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Mark:      mark,
		CRC:       0,
	}
}

// CalculateCRC 計算並返回 Entry 的 CRC 校驗碼
func (e *Entry) CalculateCRC() uint32 {
	crc := crc32.NewIEEE()
	crc.Write(e.Key)
	crc.Write(e.Value)
	return crc.Sum32()
}

// Encode 將 Entry 編碼為字節數組
func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, entryHeaderSize+e.KeySize+e.ValueSize)

	e.CRC = e.CalculateCRC()
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	binary.BigEndian.PutUint16(buf[8:10], uint16(e.Mark))
	binary.BigEndian.PutUint32(buf[10:14], e.CRC)

	copy(buf[entryHeaderSize:entryHeaderSize+e.KeySize], e.Key)
	copy(buf[entryHeaderSize+e.KeySize:], e.Value)

	return buf, nil
}

// Decode 將字節數組解碼為 Entry
func Decode(buf []byte) (*Entry, error) {
	if len(buf) < entryHeaderSize {
		return nil, errors.New("buffer too small to decode Entry")
	}

	ks := binary.BigEndian.Uint32(buf[0:4])
	vs := binary.BigEndian.Uint32(buf[4:8])
	mark := binary.BigEndian.Uint16(buf[8:10])
	crc := binary.BigEndian.Uint32(buf[10:14])

	if len(buf) < int(entryHeaderSize+ks+vs) {
		return nil, errors.New("buffer size mismatch")
	}

	key := buf[entryHeaderSize : entryHeaderSize+ks]
	value := buf[entryHeaderSize+ks : entryHeaderSize+ks+vs]

	entry := &Entry{
		Key:       key,
		Value:     value,
		KeySize:   ks,
		ValueSize: vs,
		Mark:      EntryType(mark),
		CRC:       crc,
	}

	if entry.CalculateCRC() != crc {
		return nil, errors.New("CRC mismatch")
	}

	return entry, nil
}
