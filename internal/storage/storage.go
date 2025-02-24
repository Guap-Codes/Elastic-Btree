package storage

import (
	"elastic-btree/internal/tree" // Import the tree package
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	//"sync"
)

// Storage represents the persistent storage layer for the B-tree.
type Storage struct {
	filePath string // Path to the file where the tree is stored
}

// NewStorage creates a new Storage instance with the given file path.
func NewStorage(filePath string) *Storage {
	return &Storage{
		filePath: filePath,
	}
}

// SaveTree serializes the tree and saves it to disk.
func (s *Storage) SaveTree(tree *tree.Tree) error {
	if tree == nil {
		return errors.New("tree is nil")
	}

	// Serialize the tree to JSON
	data, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("failed to serialize tree: %v", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the serialized data to the file
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tree to file: %v", err)
	}

	return nil
}

// LoadTree loads the tree from disk and deserializes it.
func (s *Storage) LoadTree() (*tree.Tree, error) {
	// Read the file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", s.filePath)
		}
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Deserialize the tree
	var tree tree.Tree
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, fmt.Errorf("failed to deserialize tree: %v", err)
	}

	// Reinitialize fields that can't be serialized
	//tree.Lock = sync.RWMutex{}
	tree.Comparator = func(a, b int) int { return a - b } // Default comparator

	tree.RebuildParentPointers()
	return &tree, nil
}

// DeleteTree deletes the tree's storage file.
func (s *Storage) DeleteTree() error {
	if err := os.Remove(s.filePath); err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to delete
		}
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}
