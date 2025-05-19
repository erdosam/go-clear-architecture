package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/erdosam/go-clear-architecture/config"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type googleStorage struct {
	log        logger.Interface
	client     *storage.Client
	bucket     *storage.BucketHandle
	isEmulated bool
}

var _ Storage = &googleStorage{}

func New(l logger.Interface, cfg *config.Config) Storage {
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		l.Fatal(err)
	}
	b := c.Bucket(cfg.Google.BucketName)
	e := os.Getenv("STORAGE_EMULATOR_HOST") != ""
	return &googleStorage{l, c, b, e}
}

func (s *googleStorage) BulkStore(files []*multipart.FileHeader) <-chan storingResult {
	return storageBulkProcess(files, func(f *multipart.FileHeader) storingResult {
		ctx := context.Background()
		src, err := f.Open()
		if err != nil {
			return storingResult{
				OriginalName: f.Filename,
				Error:        err,
			}
		}
		defer src.Close()

		obj := s.bucket.Object(f.Filename)
		wrt := obj.NewWriter(ctx)
		if _, err := io.Copy(wrt, src); err != nil {
			_ = wrt.Close()
			return storingResult{
				OriginalName: f.Filename,
				Error:        err,
			}
		}
		if err := wrt.Close(); err != nil {
			return storingResult{
				OriginalName: f.Filename,
				Error:        err,
			}
		}
		if !s.isEmulated {
			if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
				return storingResult{
					OriginalName: f.Filename,
					Error:        err,
				}
			}
		}
		url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucket.BucketName(), obj.ObjectName())
		s.log.Info("File has been stored to %s", url)
		return storingResult{
			OriginalName: f.Filename,
			PublicUrl:    url,
			Error:        nil,
		}
	})
}

func (s *googleStorage) BulkDelete(urls []string) <-chan storingResult {
	return storageBulkProcess(urls, func(url string) storingResult {
		ctx := context.Background()
		path := s.extractObjectPathFromURL(url)
		if path == "" {
			return storingResult{
				OriginalName: path,
				Error:        errors.New("no such file or directory in the storage"),
			}
		}
		if err := s.bucket.Object(path).Delete(ctx); err != nil {
			return storingResult{
				OriginalName: path,
				Error:        errors.New("no such file or directory in the storage"),
			}
		}
		s.log.Info("File %s has been deleted", url)
		return storingResult{
			OriginalName: url,
			PublicUrl:    "",
			Error:        nil,
		}
	})
}

func (s *googleStorage) Close() error {
	return s.client.Close()
}

func (s *googleStorage) extractObjectPathFromURL(url string) string {
	const prefix = "https://storage.googleapis.com/"
	if !strings.HasPrefix(url, prefix) {
		return ""
	}
	parts := strings.SplitN(strings.TrimPrefix(url, prefix), "/", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
