package bplustree

import (
	"encoding/gob"
	"os"
)

// SaveTree serializes the B+ tree to a file.
func (tree *BPlusTree) SaveTree(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(tree)
}

// LoadTree deserializes a B+ tree from a file.
func LoadTree(filename string) (*BPlusTree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var tree BPlusTree
	if err := decoder.Decode(&tree); err != nil {
		return nil, err
	}
	return &tree, nil
}
