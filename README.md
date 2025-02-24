# Elastic B-Tree

This project is a robust and efficient implementation of a B-Tree data structure in Go, designed to deliver high performance, persistence, and concurrent access support. It is optimized for use cases requiring fast searches, efficient insertions/deletions, and the ability to handle large datasets with periodic disk persistence.

## Features

- **Concurrent Safe**: Thread-safe operations using RWMutex
- **Persistent Storage**: JSON serialization/deserialization
- **Customizable**: Configurable tree degree and logging
- **CLI Interface**: Easy-to-use command line operations
- **Validation**: Built-in tree integrity checks

## Installation

```bash
git clone https://github.com/Guap-Codes/Elastic-Btree.git

cd elastic-btree

go build -o elastic-btree cmd/main.go
```

## Usage

```bash
# Insert key-value pair
./elastic-btree insert 42 "important value"

# Search for a key
./elastic-btree search 42

# Delete a key
./elastic-btree delete 42

# Save tree to disk
./elastic-btree save

# Print tree structure
./elastic-btree print

# Validate tree integrity
./elastic-btree validate
```

## Configuration

Environment Variables:

 - TREE_DEGREE: Minimum degree of the B-Tree (default: 3)

 - STORAGE_PATH: Path to persistence file (default: data/tree.json)

 - LOG_LEVEL: Logging level (debug/info/warn/error)


## Benchmarks

- Run performance tests:
```bash
go test -bench=. -benchmem -v
```

- Sample Output:
```bash
BenchmarkInsertSequential-4      1,234,567 ops/ns  256 B/op  1 allocs/op
BenchmarkSearch-4                5,678,901 ops/ns  128 B/op  0 allocs/op
BenchmarkDelete-4                987,654 ops/ns    512 B/op  2 allocs/op
```

## License
MIT License - See LICENSE for details