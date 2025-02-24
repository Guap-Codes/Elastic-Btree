package main

import (
	"elastic-btree/internal/storage"
	"elastic-btree/internal/tree"
	"elastic-btree/pkg/config"
	"elastic-btree/pkg/logger"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel, os.Stderr)

	// Create tree and storage
	//currentTree := tree.NewTree(cfg.TreeDegree, log)
	storage := storage.NewStorage(cfg.StoragePath)

	// Load tree from disk (if it exists)
	currentTree, err := storage.LoadTree()
	if err != nil {
		log.Infof("No existing tree found, creating a new one")
		currentTree = tree.NewTree(cfg.TreeDegree, log)
	} else {
		log.Infof("Tree loaded from disk")
		// Re-inject dependencies that weren't serialized.
    	currentTree.Logger = log
	}

	if len(os.Args) < 2 {
		printUsage(log)
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "insert":
		handleInsert(currentTree, log, storage)
	case "delete":
		handleDelete(currentTree, log, storage)
	case "search":
		handleSearch(currentTree, log)
	case "save":
		handleSave(currentTree, storage, log)
	case "load":
		currentTree = handleLoad(storage, cfg, log)
	case "print":
		currentTree.PrintTreeStructure()
	case "validate":
		handleValidate(currentTree, log)
	default:
		log.Errorf("Unknown command: %s", command)
		printUsage(log)
		os.Exit(1)
	}
}

func printUsage(log *logger.Logger) {
	log.Infof("Usage: ./main <command> [arguments]")
	log.Infof("Commands:")
	log.Infof("  insert <key> <value> - Insert a key-value pair")
	log.Infof("  delete <key>         - Delete a key")
	log.Infof("  search <key>         - Search for a key")
	log.Infof("  save                 - Save tree to disk")
	log.Infof("  load                 - Load tree from disk")
	log.Infof("  print                - Print tree structure")
	log.Infof("  validate             - Validate tree properties")
}

func handleInsert(t *tree.Tree, log *logger.Logger, s *storage.Storage) {
	if len(os.Args) < 4 {
		log.Errorf("Insert command requires key and value")
		log.Infof("Example: ./main insert 42 \"example value\"")
		os.Exit(1)
	}

	key, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Errorf("Invalid key: %v", err)
		log.Infof("Key must be an integer")
		os.Exit(1)
	}

	value := os.Args[3]
	t.Insert(key, value)
	log.Infof("Inserted key %d with value: %s", key, value)

	// Automatically save after insertion.
    if err := s.SaveTree(t); err != nil {
        log.Errorf("Save failed: %v", err)
        os.Exit(1)
    }
    log.Infof("Tree saved successfully")
}

func handleDelete(t *tree.Tree, log *logger.Logger, s *storage.Storage) {
	if len(os.Args) < 3 {
		log.Errorf("Delete command requires a key")
		os.Exit(1)
	}

	key, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Errorf("Invalid key: %v", err)
		os.Exit(1)
	}

	t.Delete(key)
	log.Infof("Deleted key %d", key)

	 // Automatically save after deletion.
    if err := s.SaveTree(t); err != nil {
        log.Errorf("Save failed: %v", err)
        os.Exit(1)
    }
    log.Infof("Tree saved successfully")
}

func handleSearch(t *tree.Tree, log *logger.Logger) {
	if len(os.Args) < 3 {
		log.Errorf("Search command requires a key")
		os.Exit(1)
	}

	key, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Errorf("Invalid key: %v", err)
		os.Exit(1)
	}

	if value, found := t.Search(key); found {
		log.Infof("Found key %d: %v", key, value)
	} else {
		log.Infof("Key %d not found", key)
	}
}

func handleSave(t *tree.Tree, s *storage.Storage, log *logger.Logger) {
	if err := s.SaveTree(t); err != nil {
		log.Errorf("Save failed: %v", err)
		os.Exit(1)
	}
	log.Infof("Tree saved successfully")
}

func handleLoad(s *storage.Storage, cfg *config.Config, log *logger.Logger) *tree.Tree {
	loadedTree, err := s.LoadTree()
	if err != nil {
		log.Errorf("Load failed: %v", err)
		os.Exit(1)
	}

	// Re-inject dependencies that can't be serialized
//	loadedTree.Lock = sync.RWMutex{}
	loadedTree.Logger = log // Requires logger field to be exported (adjust in tree.go)

	log.Infof("Tree loaded successfully")
	return loadedTree
}

func handleValidate(t *tree.Tree, log *logger.Logger) {
	if valid := t.ValidateTree(); valid {
		log.Infof("Tree validation successful")
	} else {
		log.Errorf("Tree validation failed")
	}
}
