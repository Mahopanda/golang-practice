package bitcask

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
)

type Bitcask struct {
	mu     sync.Mutex
	file   *os.File
	keyDir *KeyDir
}

func NewBitcask(filename string) (*Bitcask, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	bc := &Bitcask{
		file:   file,
		keyDir: NewKeyDir(),
	}

	if err := bc.buildIndex(); err != nil {
		return nil, err
	}

	return bc, nil
}

func (bc *Bitcask) Put(key, value []byte) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	entry := NewEntry(key, value, PUT)
	data, err := entry.Encode()
	if err != nil {
		return err
	}

	offset, err := bc.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if _, err := bc.file.Write(data); err != nil {
		return err
	}

	bc.keyDir.Put(string(key), offset)
	return nil
}

func (bc *Bitcask) Get(key []byte) ([]byte, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	offset, exists := bc.keyDir.Get(string(key))
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

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

func (bc *Bitcask) Delete(key []byte) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if _, exists := bc.keyDir.Get(string(key)); !exists {
		return fmt.Errorf("key not found")
	}

	entry := NewEntry(key, nil, DEL)
	data, err := entry.Encode()
	if err != nil {
		return err
	}

	if _, err := bc.file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if _, err := bc.file.Write(data); err != nil {
		return err
	}

	bc.keyDir.Delete(string(key))
	return nil
}

func (bc *Bitcask) Merge() error {
	// 合併邏輯和上面的 Merge 邏輯基本保持一致
	return nil
}

// buildIndex 構建內存索引
func (bc *Bitcask) buildIndex() error {
	// 和之前的 buildIndex 基本保持一致
	return nil
}
