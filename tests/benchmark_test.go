package tree_test

import (
	"io"
	"math/rand"
	"testing"
	"time"
	"elastic-btree/internal/storage"
	"elastic-btree/internal/tree"
	"elastic-btree/pkg/logger"
)

const (
	benchmarkDegree = 100
	numPreloadKeys  = 100000
)

func newTestTree() *tree.Tree {
	return tree.NewTree(benchmarkDegree, logger.New(logger.Error, io.Discard))
}

func BenchmarkInsertSequential(b *testing.B) {
	t := newTestTree()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		t.Insert(i, struct{}{})
	}
}

func BenchmarkInsertRandom(b *testing.B) {
	t := newTestTree()
	rand.Seed(time.Now().UnixNano())
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Insert(rand.Intn(b.N*2), struct{}{})
	}
}

func BenchmarkSearch(b *testing.B) {
	t := newTestTree()
	for i := 0; i < numPreloadKeys; i++ {
		t.Insert(i, struct{}{})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Search(i % numPreloadKeys)
	}
}

func BenchmarkDelete(b *testing.B) {
	t := newTestTree()
	keys := make([]int, numPreloadKeys)
	for i := 0; i < numPreloadKeys; i++ {
		keys[i] = i
		t.Insert(i, struct{}{})
	}
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i < numPreloadKeys {
			t.Delete(keys[i])
		}
	}
}

func BenchmarkBulkInsertAndSave(b *testing.B) {
	t := newTestTree()
	storage := storage.NewStorage("benchmark.db")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Insert(i, struct{}{})
		if i%1000 == 0 {
			storage.SaveTree(t)
		}
	}
}