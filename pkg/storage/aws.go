package storage

import (
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"mime/multipart"
)

type awsStorage struct {
	log logger.Interface
}

var _ Storage = &awsStorage{}

// IF you want to go using aws
//func New(l logger.Interface) Storage {
//	return &awsStorage{l}
//}

func (s *awsStorage) BulkStore(files []*multipart.FileHeader) <-chan storingResult {
	return storageBulkProcess(files, func(f *multipart.FileHeader) storingResult {
		// TODO store to aws storage
		storedPath := "https://dummy.cloud.stora.ge/" + f.Filename
		s.log.Info("File %s has been stored", f.Filename)
		return storingResult{
			OriginalName: f.Filename,
			PublicUrl:    storedPath,
			Error:        nil,
		}
	})
}

func (s *awsStorage) BulkDelete(urls []string) <-chan storingResult {
	return storageBulkProcess(urls, func(url string) storingResult {
		s.log.Info("File %s has been deleted", url)
		return storingResult{
			OriginalName: url,
			PublicUrl:    "",
			Error:        nil,
		}
	})
}

func (s *awsStorage) Close() error {
	return nil
}
