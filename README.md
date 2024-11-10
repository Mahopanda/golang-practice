# golang-practice

## 簡易鍵值儲存庫
### 簡介
以 Golang 撰寫，實作一個簡單的 Bitcask 鍵值儲存庫

### 功能模組概述

* 採用 Bitcask 結構，利用硬碟文件儲存資料，記憶體索引加速查詢。
* 支援 Put (插入/更新)、Get (查詢) 和 Delete (刪除) 操作，資料依序追加到文件中。
* 支援資料合併功能 (Merge)，減少碎片、節省儲存空間。
* bitcask.go、entry.go 和 keydir.go 模組負責處理資料庫的基本 CRUD 操作和記憶體索引管理。

## B+樹索引
### 簡介
支援多欄位查詢的 B+ 樹索引程式。這些功能被包裝成模組化的套件，可以儲存樹資料到檔案、可建立索引。

### 功能模組概述

* 支援 Insert (插入)、Search (查詢)、Update (更新) 以及 Delete (刪除) 操作。
* 實作的多欄位索引結構，便於快速查詢多個欄位資料。
* 建立簡易 Database 模組，支援多欄位索引，允許根據 ID 和名稱等不同欄位進行查詢。

### SQL Parser 模組
#### 簡介
實作一個 SQL 解析器，支援基本的 SQL 語句解析，包括資料定義語言（DDL）和資料操作語言（DML）。

#### 支援的 SQL 語句
1. 資料定義語言 (DDL)
   * CREATE TABLE - 建立新表格
   * DROP TABLE - 刪除現有表格

2. 資料操作語言 (DML)
   * SELECT - 查詢資料，支援條件篩選
   * INSERT - 插入單筆資料
   * UPDATE - 更新資料，支援條件篩選
   * DELETE - 刪除資料，支援條件篩選
   * BULK INSERT - 批量插入多筆資料

#### 語法特點
* WHERE 子句支援基本的比較運算（=）
* INSERT 支援單筆和批量插入
* CREATE TABLE 支援基本的資料型別定義（INTEGER、STRING）
* 支援標準 SQL 語法，如：
  ```sql
  -- 單筆插入
  INSERT INTO users (name, age) VALUES ("John", 25);
  
  -- 批量插入
  INSERT INTO employees (id, name, age, department) VALUES 
    (1, "Peter Parker", 19, "IT"),
    (2, "Bruce Wayne", 29, "Automobile");
  
  -- 條件查詢
  SELECT name, age FROM users WHERE age = 30;
  ```

#### 技術實現
* 採用遞歸下降解析方法
* 支援詞法分析（Lexical Analysis）和語法分析（Syntax Analysis）
* 使用模組化設計，易於擴展新的 SQL 功能
* 完整的錯誤處理和錯誤提示機制

## 練習 Memory Mapped Files

* 利用 mmap 和 munmap 函數將文件內容映射到內存中，實作文件讀寫。

## Binary Search Tree

以Golang 建立一個基本的 BinarySearchTree (BST) 資料結構。原本一開始是以建立二元樹為出發點，後來調整成平衡樹。
此 BST 支援增刪查改 (CRUD) 操作，並且具有自我平衡的功能，使得樹在節點插入或刪除後能夠保持平衡 (AVL 樹)。
提供了 main.go 中的範例展示了如何操作此 BST 並以圖形化方式顯示樹的結構。

### 功能概述
#### 1. 建立與初始化 BST

* 透過 BinarySearchTree 結構來表示 BST，其中包含節點 (Node) 的資料結構，每個節點包含 Value, Left, Right, 及 Height 用於 AVL 平衡操作。

#### 2. 節點的插入 (AddNode)

* 使用 addNode 方法將值插入 BST 中，且自動維持樹的平衡性。
* 插入後的樹結構將透過 AVL 旋轉 (左旋與右旋) 來保持平衡。

#### 3. 節點的查詢 (Search)

* 使用 search 方法進行節點查詢，支援根據節點值來查找並回傳該節點。

#### 4. 節點的刪除 (DeleteNode)

* 使用 deleteNode 方法刪除指定的節點，並自動平衡 AVL 樹。
* 在節點有兩個子節點時，會找到右子樹的最小值作為替代。

#### 5. 圖形化顯示樹結構

* 透過 printTree 方法來生成樹的視覺化結構，顯示父節點與子節點的連結關係，便於觀察樹的形態。

  
