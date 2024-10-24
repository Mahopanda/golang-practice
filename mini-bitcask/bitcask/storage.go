package bitcask

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

// Storage 是 Bitcask 的存儲結構
type Storage struct {
	index     map[string]valueLocation // 鍵對應的文件偏移
	dataFiles []*DBFile                // 多個數據文件
	lock      sync.RWMutex             // 用來保護併發讀寫
	basePath  string                   // 存儲文件的基礎路徑
	fileLimit int64                    // 每個文件的大小限制（如 2MB）
}

type valueLocation struct {
	fileIndex int   // 文件索引
	offset    int64 // 文件內偏移
}

// NewStorage 初始化一個新的 Storage，並指定文件大小上限
func NewStorage(basePath string, fileLimit int64) *Storage {
	s := &Storage{
		index:     make(map[string]valueLocation),
		dataFiles: []*DBFile{},
		basePath:  basePath,
		fileLimit: fileLimit,
	}
	return s
}

// GetDataFiles 返回目前所有的數據文件
func (s *Storage) GetDataFiles() []*DBFile {
	return s.dataFiles
}

// LoadFiles 重啟時加載所有數據文件，並重建索引
func (s *Storage) LoadFiles() error {
	files, err := ioutil.ReadDir(s.basePath)
	if err != nil {
		return err
	}

	// 遍歷所有檔案，按順序重新加載
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".data" {
			dbFile, err := NewDBFile(s.basePath, s.fileLimit, file.Name())
			if err != nil {
				return err
			}

			// 加入到 dataFiles 列表中
			s.dataFiles = append(s.dataFiles, dbFile)

			// 重新加載文件中的所有 Entry 並重建索引
			if err := s.loadEntriesFromFile(len(s.dataFiles)-1, dbFile); err != nil {
				return err
			}
		}
	}

	// 如果沒有數據文件，創建一個新的
	if len(s.dataFiles) == 0 {
		return s.AddFile()
	}

	return nil
}

// loadEntriesFromFile 從指定的數據文件中讀取所有 Entry，並重建索引
func (s *Storage) loadEntriesFromFile(fileIndex int, dbFile *DBFile) error {
	var offset int64
	for offset < dbFile.Offset {
		entry, err := dbFile.Read(offset)
		if err != nil {
			return err
		}

		// 打印调试信息
		fmt.Printf("Loaded entry: key=%s, offset=%d\n", entry.Key, offset)

		if entry.Mark == PUT {
			s.index[string(entry.Key)] = valueLocation{fileIndex: fileIndex, offset: offset}
		} else if entry.Mark == DEL {
			delete(s.index, string(entry.Key))
		}

		// 移动到下一个条目
		offset += entry.GetSize()
	}
	return nil
}

// AddFile 新增一個新的數據文件
func (s *Storage) AddFile() error {
	fileName := fmt.Sprintf("datafile-%d.data", len(s.dataFiles)+1)
	dbFile, err := NewDBFile(s.basePath, s.fileLimit, fileName)
	if err != nil {
		return err
	}
	s.dataFiles = append(s.dataFiles, dbFile)
	return nil
}

// Put 寫入鍵值對到存儲
func (s *Storage) Put(key string, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// 建立新的 Entry
	newEntry := NewEntry([]byte(key), value, PUT)
	encodedEntry, err := newEntry.Encode()
	if err != nil {
		return err
	}

	// 確認是否需要分割文件
	lastFile := s.getLastFile()
	if lastFile == nil || lastFile.Offset+int64(len(encodedEntry)) > s.fileLimit {
		if err := s.AddFile(); err != nil {
			return err
		}
		lastFile = s.getLastFile()
	}

	// 將 Entry 寫入最後一個數據文件
	err = lastFile.Write(newEntry)
	if err != nil {
		return err
	}

	// 更新索引，保存該鍵的偏移量
	s.index[key] = valueLocation{fileIndex: len(s.dataFiles) - 1, offset: lastFile.Offset - newEntry.GetSize()}
	return nil
}

// Get 讀取指定鍵的值
func (s *Storage) Get(key string) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// 確認鍵是否存在
	loc, ok := s.index[key]
	if !ok {
		return nil, errors.New("鍵不存在")
	}

	// 根據 loc 找到對應的文件並讀取
	if loc.fileIndex < 0 || loc.fileIndex >= len(s.dataFiles) {
		return nil, errors.New("無效的文件索引")
	}

	file := s.dataFiles[loc.fileIndex]
	entry, err := file.Read(loc.offset)
	if err != nil {
		return nil, err
	}

	return entry.Value, nil
}

// getLastFile 返回最後一個數據文件
func (s *Storage) getLastFile() *DBFile {
	if len(s.dataFiles) == 0 {
		return nil
	}
	return s.dataFiles[len(s.dataFiles)-1]
}
