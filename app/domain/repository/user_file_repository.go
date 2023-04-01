package repository

import (
	"context"
	"mime/multipart"
)

type UserFileRepository interface {
	GetURLByUser(ctx context.Context, id int) string
	Create(ctx context.Context, id int, f multipart.File, fh *multipart.FileHeader) (string, error)
	Update(ctx context.Context, id int, f multipart.File, fh *multipart.FileHeader) (string, error)
	Delete(ctx context.Context, id int) error
}
