package bitcask

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a temporary directory for testing
func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "bitcask_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// Helper function to cleanup temp directory
func cleanupTempDir(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Failed to cleanup temp dir: %v", err)
	}
}

func TestStorageFileSplitting(t *testing.T) {
	//測試文件分割與寫入功能
	// 初始化測試環境
	dir := createTempDir(t)
	defer cleanupTempDir(t, dir)

	// 設定文件大小上限為 100 bytes 來測試文件分割
	fileLimit := int64(100)
	storage := NewStorage(dir, fileLimit)

	// 寫入一些鍵值對
	err := storage.Put("key1", []byte("value1-longer-value-to-test-file-splitting"))
	assert.NoError(t, err, "Put should succeed")

	err = storage.Put("key2", []byte("another-long-value-to-trigger-file-splitting"))
	assert.NoError(t, err, "Put should succeed")

	// 應該已經有兩個檔案，因為每次寫入都應該超過 100 bytes
	assert.Equal(t, 2, len(storage.GetDataFiles()), "Should have two data files")
}

func TestStoragePutAndGet(t *testing.T) {
	//測試數據寫入與讀取
	// 初始化測試環境
	dir := createTempDir(t)
	defer cleanupTempDir(t, dir)

	storage := NewStorage(dir, 1024) // 使用 1KB 文件上限

	// 寫入鍵值對
	err := storage.Put("username", []byte("john_doe"))
	assert.NoError(t, err, "Put should succeed")

	// 讀取鍵值對
	value, err := storage.Get("username")
	assert.NoError(t, err, "Get should succeed")
	assert.Equal(t, []byte("john_doe"), value, "Value should match")
}

// func TestStorageReload(t *testing.T) {
// 	//測試重啟後能夠重新載入數據
// 	// 初始化測試環境
// 	dir := createTempDir(t)
// 	defer cleanupTempDir(t, dir)

// 	storage := NewStorage(dir, 1024) // 使用 1KB 文件上限

// 	// 寫入鍵值對
// 	err := storage.Put("username", []byte("john_doe"))
// 	assert.NoError(t, err, "Put should succeed")

// 	// 重啟存儲，重新加載檔案
// 	storage = NewStorage(dir, 1024)
// 	err = storage.LoadFiles()
// 	assert.NoError(t, err, "LoadFiles should succeed")

// 	// 讀取之前寫入的鍵值對
// 	value, err := storage.Get("username")
// 	assert.NoError(t, err, "Get should succeed after reload")
// 	assert.Equal(t, []byte("john_doe"), value, "Value should match after reload")
// }

// func TestStorageDelete(t *testing.T) {
// 	//測試刪除功能
// 	// 初始化測試環境
// 	dir := createTempDir(t)
// 	defer cleanupTempDir(t, dir)

// 	storage := NewStorage(dir, 1024) // 使用 1KB 文件上限

// 	// 寫入鍵值對
// 	err := storage.Put("username", []byte("john_doe"))
// 	assert.NoError(t, err, "Put should succeed")

// 	// 刪除該鍵
// 	err = storage.Put("username", []byte("TOMBSTONE")) // 使用 Tombstone 表示刪除
// 	assert.NoError(t, err, "Delete should succeed")

// 	// 讀取該鍵，應該返回不存在
// 	_, err = storage.Get("username")
// 	assert.Error(t, err, "Get should fail for deleted key")
// }

// func TestStorageMultiFileLoad(t *testing.T) {
// 	//測試多檔案加載與跨檔案讀取
// 	// 初始化測試環境
// 	dir := createTempDir(t)
// 	defer cleanupTempDir(t, dir)

// 	// 設定文件大小上限為 100 bytes，讓它分割為多個文件
// 	storage := NewStorage(dir, 100)

// 	// 寫入多個鍵值對，強制文件分割
// 	err := storage.Put("key1", []byte("value1"))
// 	assert.NoError(t, err, "Put should succeed")

// 	err = storage.Put("key2", []byte("value2"))
// 	assert.NoError(t, err, "Put should succeed")

// 	err = storage.Put("key3", []byte("value3"))
// 	assert.NoError(t, err, "Put should succeed")

// 	// 應該會有多個文件
// 	assert.Equal(t, 2, len(storage.GetDataFiles()), "Should have two data files")
// 	assert.GreaterOrEqual(t, len(storage.GetDataFiles()), 2, "Should have at least two data files")

// 	// 重新加載並讀取數據
// 	storage = NewStorage(dir, 100)
// 	err = storage.LoadFiles()
// 	assert.NoError(t, err, "LoadFiles should succeed")

// 	// 驗證數據正確讀取
// 	value, err := storage.Get("key1")
// 	assert.NoError(t, err, "Get should succeed for key1")
// 	assert.Equal(t, []byte("value1"), value, "Value for key1 should match")

// 	value, err = storage.Get("key2")
// 	assert.NoError(t, err, "Get should succeed for key2")
// 	assert.Equal(t, []byte("value2"), value, "Value for key2 should match")

// 	value, err = storage.Get("key3")
// 	assert.NoError(t, err, "Get should succeed for key3")
// 	assert.Equal(t, []byte("value3"), value, "Value for key3 should match")
// }
