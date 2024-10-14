package s3

import (
	"context"
	"io"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/escalopa/cloudphoto/domain"
	"github.com/google/uuid"
)

type credentials struct{}

func (c *credentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("ACCESS_KEY"),
		SecretAccessKey: os.Getenv("SECRET_KEY"),
	}, nil
}

type Client struct {
	s3Client *s3.Client
}

func New() *Client {
	s3Client := s3.New(s3.Options{
		Credentials:  &credentials{},
		Region:       os.Getenv("REGION"),
		BaseEndpoint: aws.String(os.Getenv("BASE_ENDPOINT")),
	})
	return &Client{s3Client: s3Client}
}

func (c *Client) Upload(ctx context.Context, bucket string, key string, body io.Reader, contentType string) error {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

func (c *Client) Download(ctx context.Context, bucket string, key string) ([]byte, error) {
	out, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}

func (c *Client) ListObjects(ctx context.Context, bucket string, prefix string) (*domain.Folder, error) {
	out, err := c.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(formatPrefix(prefix)),
	})
	if err != nil {
		return nil, err
	}
	return c.buildFolder(out, prefix), nil
}

func (c *Client) GenKey() string {
	key, _ := uuid.NewV7()
	return key.String()
}

func (c *Client) buildFolder(out *s3.ListObjectsV2Output, prefix string) *domain.Folder {
	folder := &domain.Folder{
		Key:  prefix,
		Name: formatBase(prefix),
	}
	// Read files
	for _, item := range out.Contents {
		itemFile := &domain.File{
			Key:  aws.ToString(item.Key),
			Name: formatBase(aws.ToString(item.Key)),
			Size: aws.ToInt64(item.Size),
		}
		folder.Files = append(folder.Files, itemFile)
	}
	// Read folders
	for _, prefix := range out.CommonPrefixes {
		subFolder := &domain.Folder{
			Key:  aws.ToString(prefix.Prefix),
			Name: formatBase(aws.ToString(prefix.Prefix)),
		}
		folder.Folders = append(folder.Folders, subFolder)
	}
	return folder
}

func formatPrefix(prefix string) string {
	if prefix == "" {
		return ""
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return prefix
}

func formatBase(prefix string) string {
	s := path.Base(prefix)
	if s == "." {
		return ""
	}
	return s
}
