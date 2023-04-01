package file

import (
	"context"
	"github.com/jpdel518/go-ent/domain/repository"
	"github.com/jpdel518/go-ent/infrastructure/file/s3"
	"mime/multipart"
	"strconv"
	"strings"
)

type userFileRepository struct {
	session *s3.S3
}

func NewUserFileRepository(session *s3.S3) repository.UserFileRepository {
	return &userFileRepository{
		session: session,
	}
}

func (ur userFileRepository) GetURLByUser(ctx context.Context, id int) string {
	return ur.session.GetFileURL("user/avatar/" + strconv.Itoa(id))
}

func (ur userFileRepository) Create(ctx context.Context, id int, f multipart.File, fh *multipart.FileHeader) (string, error) {
	return ur.session.UploadMultipartFile("user/avatar/"+strconv.Itoa(id), fh, f)
}

func (ur userFileRepository) Update(ctx context.Context, id int, f multipart.File, fh *multipart.FileHeader) (string, error) {
	// Get existing avatar image file name
	existingFile := ur.session.GetFileURL("user/avatar/" + strconv.Itoa(id))
	url := strings.Split(existingFile, "/")
	// Upload new avatar image file
	file, err := ur.session.UploadMultipartFile("user/avatar/"+strconv.Itoa(id), fh, f)
	if err != nil {
		return "", err
	}
	// Delete old avatar image file
	// If the file name is the same, it will not be deleted
	filename := url[len(url)-1]
	if fh.Filename != filename {
		err = ur.session.DeleteFile("user/avatar/"+strconv.Itoa(id), filename)
		if err != nil {
			return file, err
		}
	}
	return file, nil
}

func (ur userFileRepository) Delete(ctx context.Context, id int) error {
	existsFile := ur.session.GetFileURL("user/avatar/" + strconv.Itoa(id))
	filename := strings.SplitAfter(existsFile, "user/avatar/"+strconv.Itoa(id)+"/")
	return ur.session.DeleteFile("user/avatar/"+strconv.Itoa(id), filename[1])
}
