package storage

import (
	"mime/multipart"
	"sync"
)

type Storage interface {
	BulkStore(files []*multipart.FileHeader) <-chan storingResult
	BulkDelete(urls []string) <-chan storingResult
	Close() error
}

type storingResult struct {
	OriginalName string
	PublicUrl    string
	Error        error
}

type storageAsyncFunc[T any] func(T) storingResult

func storageBulkProcess[T any](objs []T, fnc storageAsyncFunc[T]) <-chan storingResult {
	results := make(chan storingResult, len(objs))
	var wg sync.WaitGroup
	for _, file := range objs {
		wg.Add(1)
		go func(t T) {
			defer wg.Done()
			results <- fnc(t)
		}(file)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}
