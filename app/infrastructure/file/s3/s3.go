package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"mime/multipart"
	"os"
	"sync"
	"time"
)

type S3 struct {
	s3session      *session.Session
	baseBucketName string
}

type partUploadResult struct {
	completedPart *s3.CompletedPart
	err           error
}

func NewS3Session() *S3 {
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	s := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_REGION")),
	}))
	bucket := os.Getenv("AWS_BUCKET_NAME")
	return &S3{
		s3session:      s,
		baseBucketName: bucket,
	}
}

// UploadFile ファイルをS3にアップロードする
func (s *S3) UploadFile(folder string, key string, file *os.File) error {
	_, err := s3.New(s.s3session).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.baseBucketName),
		Key:    aws.String(folder + "/" + key),
		Body:   file,
	})
	return err
}

// UploadMultipartFile multipart requestで受け取ったファイルをS3にアップロードする
func (s *S3) UploadMultipartFile(folder string, fh *multipart.FileHeader, file multipart.File) (string, error) {
	// create a new uploader
	uploader := s3manager.NewUploader(s.s3session, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5MB
		u.LeavePartsOnError = true
		u.Concurrency = 5
	})
	// create a new file upload input
	input := &s3manager.UploadInput{
		Bucket: aws.String(s.baseBucketName),
		Key:    aws.String(folder + "/" + fh.Filename),
		Body:   file,
	}
	// Upload the file to S3.
	upload, err := uploader.Upload(input)
	if err != nil {
		log.Printf("failed to upload file, %v", err)
		return "", err
	}

	return upload.Location, nil
}

// UploadBigFile 大容量ファイルをS3にアップロードする
func (s *S3) UploadBigFile(folder string, fh *multipart.FileHeader, file multipart.File) (*string, error) {
	// create a new uploader
	uploader := s3manager.NewUploader(s.s3session, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5MB
		u.LeavePartsOnError = true
		u.Concurrency = 5
	})
	// 大容量データアップロード用
	// ファイルの有効期限を設定
	expiryDate := time.Now().AddDate(0, 0, 30)
	// multipart uploadの準備
	// multipart uploadを紐付けるためのUploadIDを取得
	s3client := s3.New(s.s3session)
	createdResp, err := s3client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket:  aws.String(s.baseBucketName),
		Key:     aws.String(folder + "/" + fh.Filename),
		Expires: &expiryDate,
	})

	var wg sync.WaitGroup
	ch := make(chan partUploadResult)
	var start, currentSize int
	var remaining = int(fh.Size)
	var partNum = 1 // 分割順序
	var completedParts []*s3.CompletedPart
	var partSize = int(uploader.PartSize)
	for start = 0; remaining > 0; start += partSize {
		wg.Add(1)
		if remaining < partSize {
			currentSize = remaining
		} else {
			currentSize = partSize
		}

		go func(partNum int, start int, currentSize int) {
			defer wg.Done()
			// ファイルからバッファに読み込み
			buffer := make([]byte, currentSize)
			_, err := file.Read(buffer)
			if err != nil {
				panic(err)
			}
			// バッファからS3にmultipart upload
			part, err := s3client.UploadPart(&s3.UploadPartInput{
				Body:       bytes.NewReader(buffer),
				Bucket:     createdResp.Bucket,
				Key:        createdResp.Key,
				PartNumber: aws.Int64(int64(partNum)),
				UploadId:   createdResp.UploadId,
				// ContentLength: aws.Int64(int64(currentSize)),
			})
			// アップロードが完了した分割データの情報をチャンネルに送信
			if err != nil {
				ch <- partUploadResult{err: err}
			} else {
				ch <- partUploadResult{completedPart: &s3.CompletedPart{
					ETag:       part.ETag,
					PartNumber: aws.Int64(int64(partNum)),
				}}
			}
		}(partNum, start, currentSize)

		remaining -= currentSize
		fmt.Printf("Uplaodind of part %v started and remaning is %v \n", partNum, remaining)
		partNum++
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// アップロードが完了した分割データの情報を受信
	for partUploadResult := range ch {
		// エラーが発生した場合はアップロードを中断
		if partUploadResult.err != nil {
			// エラーが発生する前までにアップロードした分割データを削除
			_, err = s3client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
				Bucket:   createdResp.Bucket,
				Key:      createdResp.Key,
				UploadId: createdResp.UploadId,
			})
			// これ以上アップロードを続けないようにチャンネルを閉じる
			close(ch)
			remaining = 0
			if err != nil {
				log.Printf("Failed to abort multipart upload, %v", err)
				return nil, err
			}
		}
		// 成功した場合はアップロードが完了した分割データの情報をcompletedPartsに追加
		log.Printf("Uploading completed part %v", *partUploadResult.completedPart.PartNumber)
		completedParts = append(completedParts, partUploadResult.completedPart)
	}

	// 全ての分割データのアップロードが完了したことをS3に通知
	wg.Wait()
	resp, err := s3client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   createdResp.Bucket,
		Key:      createdResp.Key,
		UploadId: createdResp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	return resp.Location, err
}

// DeleteFile ファイルを削除する
func (s *S3) DeleteFile(folder string, key string) error {
	_, err := s3.New(s.s3session).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.baseBucketName),
		Key:    aws.String(folder + "/" + key),
	})
	if err != nil {
		log.Printf("failed to delete file, %v", err)
	}
	return err
}

// GetFileURL バケット内にあるファイルのURLを取得する
func (s *S3) GetFileURL(folder string) string {
	objects, err := s3.New(s.s3session).ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.baseBucketName),
		Prefix: aws.String(folder),
	})
	if err != nil {
		log.Printf("failed to get file url, %v", err)
		return ""
	}
	// フォルダだけ（サイズ0）が取得される場合もある
	var filename string
	for _, object := range objects.Contents {
		log.Printf("object key: %v", *object.Key)
		if *object.Size > 0 {
			filename = *object.Key
			break
		}
	}
	return *s.s3session.Config.Endpoint + "/" + s.baseBucketName + "/" + folder + "/" + filename
}
