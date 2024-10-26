package main

import (
	"fmt"

	"github.com/Mahopanda/mini-project/bitcask"
)

func main() {
	// 初始化 Bitcask 資料庫
	bitcask, err := bitcask.NewBitcask("bitcask.db")
	if err != nil {
		fmt.Printf("Error initializing Bitcask: %v\n", err)
		return
	}

	// 1. 寫入鍵值對
	fmt.Println("Inserting 'name' -> 'Alice'")
	if err := bitcask.Put([]byte("name"), []byte("Alice")); err != nil {
		fmt.Printf("Error writing to Bitcask: %v\n", err)
	}

	fmt.Println("Inserting 'address' -> 'Taiwan'")
	if err := bitcask.Put([]byte("address"), []byte("Taiwan")); err != nil {
		fmt.Printf("Error writing to Bitcask: %v\n", err)
	}

	// 2. 讀取鍵值對，驗證插入
	value, err := bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask: %v\n", err)
	} else {
		fmt.Printf("Value for 'name': %s\n", value) // 應輸出 Alice
	}

	// 3. 更新鍵值對，驗證更新
	fmt.Println("Updating 'name' -> 'Bob'")
	if err := bitcask.Put([]byte("name"), []byte("Bob")); err != nil {
		fmt.Printf("Error updating Bitcask: %v\n", err)
	}

	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask after update: %v\n", err)
	} else {
		fmt.Printf("Value for 'name' after update: %s\n", value) // 應輸出 Bob
	}

	// 4. 刪除鍵值對，驗證刪除
	fmt.Println("Deleting 'name'")
	if err := bitcask.Delete([]byte("name")); err != nil {
		fmt.Printf("Error deleting from Bitcask: %v\n", err)
	}

	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Expected error reading deleted 'name': %v\n", err) // 應輸出錯誤訊息
	} else {
		fmt.Printf("Unexpected value for deleted 'name': %s\n", value)
	}

	// 5. 再次插入鍵值對，驗證重新插入
	fmt.Println("Inserting 'name' -> 'Charlie'")
	if err := bitcask.Put([]byte("name"), []byte("Charlie")); err != nil {
		fmt.Printf("Error writing to Bitcask: %v\n", err)
	}

	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask after re-insertion: %v\n", err)
	} else {
		fmt.Printf("Value for 'name' after re-insertion: %s\n", value) // 應輸出 Charlie
	}

	// 6. 合併文件，驗證合併操作
	fmt.Println("Merging entries...")
	if err := bitcask.Merge(); err != nil {
		fmt.Printf("Error during merge: %v\n", err)
	}

	// 7. 讀取鍵值對，驗證合併後資料是否仍然有效
	value, err = bitcask.Get([]byte("name"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask after merge: %v\n", err)
	} else {
		fmt.Printf("Value for 'name' after merge: %s\n", value) // 應輸出 Charlie
	}

	// 8. 讀取 address
	value, err = bitcask.Get([]byte("address"))
	if err != nil {
		fmt.Printf("Error reading from Bitcask: %v\n", err)
	} else {
		fmt.Printf("Value for 'address': %s\n", value)
	}
}
