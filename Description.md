# Project Description

A High-Performance, Persistent B-Tree Implementation in Go with Concurrent Access Support.

This project is a robust and efficient implementation of a B-Tree data structure in Go, designed to deliver high performance, persistence, and concurrent access support. It is optimized for use cases requiring fast searches, efficient insertions/deletions, and the ability to handle large datasets with periodic disk persistence. Whether you're building a database, a filesystem, or an in-memory cache, this B-Tree implementation provides the reliability and performance you need.

## Key Features

   -  High Performance

        Fast Searches: Logarithmic time complexity ensures lightning-fast lookups, with benchmarks showing search operations completing in ~1,036 nanoseconds.

        Efficient Inserts/Deletes: Insertions and deletions are optimized to handle both sequential and random data patterns, with latencies in the micro-to-millisecond range.

        Bulk Operations: Supports bulk inserts with periodic persistence, making it ideal for batch processing and large-scale data ingestion.

  -  Persistence

        Disk-Backed Storage: The B-Tree can be serialized and saved to disk, ensuring data durability and the ability to reload the tree after a restart.

        Periodic Saves: Benchmarks demonstrate efficient handling of bulk operations with periodic saves, making it suitable for applications requiring frequent snapshots or checkpoints.

  -  Concurrent Access Support

        Thread-Safe Design: Built-in support for concurrent access using fine-grained locking (sync.RWMutex), allowing multiple readers or a single writer at a time.

        Optimized for Multi-Threading: The implementation minimizes contention, enabling high throughput in multi-threaded environments.

  -  Memory Efficiency

        Low Overhead: Minimal per-operation allocations (2–5) and efficient memory usage ensure low garbage collection pressure.

        Configurable Degree: The B-Tree's Degree parameter allows you to balance memory usage and performance based on your specific workload.

 -  Balanced and Reliable

        Self-Balancing: The B-Tree automatically maintains balance during insertions and deletions, ensuring consistent performance.

        Invariant Checks: Built-in invariant checks validate the tree's structure during operations, preventing corruption and ensuring reliability.

  -  Extensible and Customizable

        Custom Comparators: Supports custom key comparison functions, enabling flexibility in how keys are ordered and compared.

        Logging and Debugging: Integrated logging (Logger) provides detailed insights into tree operations, making it easier to debug and optimize.

## Use Cases

This B-Tree implementation is ideal for a wide range of applications, including:

  -  Databases: As an index structure for fast lookups and efficient data management.

  -  Filesystems: For organizing and retrieving file metadata or block pointers.

  -  In-Memory Caches: As a high-performance key-value store for caching frequently accessed data.

  -  Batch Processing: For handling large-scale data ingestion and periodic persistence.

## Benchmark Results

The implementation has been rigorously tested and benchmarked, demonstrating its performance and efficiency:

  -  Insert (Sequential): ~113,623 ns/op, 95 B/op, 2 allocs/op.

  -  Insert (Random): ~22,078 ns/op, 1,944 B/op, 4 allocs/op.

  -  Search: ~1,036 ns/op, 0 B/op, 0 allocs/op.

  -  Delete: ~18,093 ns/op, 788 B/op, 5 allocs/op.

  -  Bulk Insert and Save: ~133,912 ns/op, 3,717 B/op, 2 allocs/op.

## Why Choose This Implementation?

  -  Performance: Optimized for both read-heavy and write-heavy workloads, with low-latency operations and minimal memory overhead.

  -  Reliability: Built-in invariant checks and persistence ensure data integrity and durability.

  -  Concurrency: Designed for multi-threaded environments, with fine-grained locking to maximize throughput.

  -  Flexibility: Customizable and extensible, making it adaptable to a wide range of use cases.

## Get Started

To start using this B-Tree implementation in your Go projects, simply import the package and follow the documentation. Whether you're building a database, a filesystem, or a custom data structure, this B-Tree provides the performance and reliability you need.