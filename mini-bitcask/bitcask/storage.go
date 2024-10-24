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
		fmt.Printf("encode error: %v\n", err)
		return err
	}

	// 确认是否需要分割文件
	lastFile := s.getLastFile()
	if lastFile == nil || lastFile.Offset+int64(len(encodedEntry)) > s.fileLimit {
		if err := s.AddFile(); err != nil {
			fmt.Printf("add file error: %v\n", err)
			return err
		}
		lastFile = s.getLastFile()
	}
	fmt.Printf("lastFile: %v\n", lastFile.Offset)
	// 保存写入前的文件偏移量
	// fileOffsetBeforeWrite := lastFile.Offset

	// 将 Entry 写入最后一个数据文件
	offset, err := lastFile.Write(newEntry)
	if err != nil {
		fmt.Printf("write error: %v\n", err)
		return err
	}
	fmt.Printf("dataFiles: %v\n", s.dataFiles)
	fmt.Printf("dataFiles length: %v\n", len(s.dataFiles))
	// 更新索引，保存写入前的偏移量，而不是写入后的偏移量
	s.index[string(key)] = valueLocation{fileIndex: len(s.dataFiles) - 1, offset: offset}

	fmt.Printf("File offset after writing key %s: %d\n", key, lastFile.Offset)

	return nil
}

// Get 讀取指定鍵的值
func (s *Storage) Get(key string) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	fmt.Printf("get key: %s\n", key)
	// 确认键是否存在
	loc, ok := s.index[key]
	fmt.Printf("loc: %v\n", loc)
	if !ok {
		fmt.Printf("key not found\n")
		return nil, errors.New("键不存在")
	}
	fmt.Printf("loc fileIndex: %v\n", loc.fileIndex)
	// 打印调试信息，查看从索引中获取的偏移量
	fmt.Printf("Offset retrieved from index for key %s: %d\n", key, loc.offset)

	// 根据 loc 找到对应的文件并读取
	if loc.fileIndex < 0 || loc.fileIndex >= len(s.dataFiles) {
		return nil, errors.New("无效的文件索引")
	}

	file := s.dataFiles[loc.fileIndex]
	if file == nil {
		return nil, errors.New("文件不存在")
	}
	entry, err := file.Read(loc.offset)
	if err != nil {
		return nil, err
	}
	fmt.Printf("get entry: %v\n", entry)
	return entry.Value, nil
}

// getLastFile 返回最後一個數據文件
func (s *Storage) getLastFile() *DBFile {
	if len(s.dataFiles) == 0 {
		return nil
	}
	return s.dataFiles[len(s.dataFiles)-1]
}
