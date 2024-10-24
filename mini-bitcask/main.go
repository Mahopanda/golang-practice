package bitcask

import (
	// "encoding/binary"
	// "errors"
	"os"
	"path/filepath"
	"sync"
)

const (
	KEY_VAL_HEADER_LEN = 4
	MERGE_FILE_EXT     = "merge"
)

type KeyDir map[string](ValuePos)

type ValuePos struct {
	Offset uint64
	Length uint32
}

type MiniBitcask struct {
	log    *Log
	keydir KeyDir
	mu     sync.Mutex
}

type Log struct {
	file *os.File
	path string
	mu   sync.Mutex
}


func NewLog(path string) (*Log, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Log{
		file: file,
		path: path,
	}, nil
}
