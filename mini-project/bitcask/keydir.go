package bitcask

import "sync"

type KeyDir struct {
	mu    sync.RWMutex
	index map[string]int64
}

// Key Directory 索引管理
// NewKeyDir 初始化 KeyDir
func NewKeyDir() *KeyDir {
	return &KeyDir{index: make(map[string]int64)}
}

// Get 返回 key 對應的文件偏移量
func (kd *KeyDir) Get(key string) (int64, bool) {
	kd.mu.RLock()
	defer kd.mu.RUnlock()
	offset, ok := kd.index[key]
	return offset, ok
}

// Put 更新 key 的偏移量
func (kd *KeyDir) Put(key string, offset int64) {
	kd.mu.Lock()
	defer kd.mu.Unlock()
	kd.index[key] = offset
}

// Delete 從索引中刪除 key
func (kd *KeyDir) Delete(key string) {
	kd.mu.Lock()
	defer kd.mu.Unlock()
	delete(kd.index, key)
}

// ListKeys 返回所有 key
func (kd *KeyDir) ListKeys() []string {
	kd.mu.RLock()
	defer kd.mu.RUnlock()
	keys := make([]string, 0, len(kd.index))
	for key := range kd.index {
		keys = append(keys, key)
	}
	return keys
}
