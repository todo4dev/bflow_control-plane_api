package s3compat

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	adapter "src/application/adapter/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type S3CompatStorageAdapter struct {
	client *awss3.Client
	bucket string
}

var _ adapter.IStorageAdapter = (*S3CompatStorageAdapter)(nil)

func NewS3CompatStorageAdapter(storageConfig *adapter.StorageConfig) *S3CompatStorageAdapter {
	if storageConfig == nil {
		panic("storage/s3compat: storage config is nil")
	}

	parsedEndpoint, err := url.Parse(storageConfig.Endpoint)
	if err != nil || parsedEndpoint.Scheme == "" || parsedEndpoint.Host == "" {
		panic(fmt.Errorf("storage/s3compat: invalid Endpoint: %q", storageConfig.Endpoint))
	}

	awsConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(storageConfig.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			storageConfig.AccessKeyID,
			storageConfig.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		panic(fmt.Errorf("storage/s3compat: failed to load AWS config: %w", err))
	}

	client := awss3.NewFromConfig(awsConfig, func(options *awss3.Options) {
		options.BaseEndpoint = aws.String(storageConfig.Endpoint)
		options.UsePathStyle = true
	})

	return &S3CompatStorageAdapter{
		client: client,
		bucket: storageConfig.Bucket,
	}
}

func (a *S3CompatStorageAdapter) Ping(ctx context.Context) error {
	_, err := a.client.HeadBucket(ctx, &awss3.HeadBucketInput{Bucket: aws.String(a.bucket)})
	return err
}

func (a *S3CompatStorageAdapter) Write(ctx context.Context, input adapter.WriteInput) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if input.Stream == nil {
		return "", fmt.Errorf("storage/s3compat: input stream is nil")
	}
	if strings.TrimSpace(input.FilePath) == "" {
		return "", fmt.Errorf("storage/s3compat: filePath is required")
	}

	key := sanitizeKey(input.FilePath)

	contentType := strings.TrimSpace(input.MimeType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	hasher := sha256.New()
	body := io.TeeReader(input.Stream, hasher)

	_, err := a.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("storage/s3compat: PutObject failed: %w", err)
	}

	digest := hex.EncodeToString(hasher.Sum(nil))

	_, err = a.client.PutObjectTagging(ctx, &awss3.PutObjectTaggingInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
		Tagging: &s3types.Tagging{
			TagSet: []s3types.Tag{
				{Key: aws.String("digest"), Value: aws.String(digest)},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("storage/s3compat: PutObjectTagging failed: %w", err)
	}

	return fmt.Sprintf("/%s?digest=%s", input.FilePath, digest), nil
}

func (a *S3CompatStorageAdapter) Read(ctx context.Context, filePath string, digest string) (*adapter.ReadResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(filePath) == "" {
		return nil, fmt.Errorf("storage/s3compat: filePath is required")
	}

	key := sanitizeKey(filePath)

	headOut, err := a.client.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("storage/s3compat: file not found")
		}
		return nil, fmt.Errorf("storage/s3compat: HeadObject failed: %w", err)
	}

	tagDigest, err := a.getDigestTag(ctx, key)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("storage/s3compat: file not found")
		}
		return nil, fmt.Errorf("storage/s3compat: GetObjectTagging failed: %w", err)
	}
	if tagDigest != digest {
		return nil, fmt.Errorf("storage/s3compat: invalid file hash")
	}

	getOut, err := a.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("storage/s3compat: file not found")
		}
		return nil, fmt.Errorf("storage/s3compat: GetObject failed: %w", err)
	}

	mimeType := "application/octet-stream"
	if headOut.ContentType != nil && *headOut.ContentType != "" {
		mimeType = *headOut.ContentType
	}

	var size int64
	if headOut.ContentLength != nil {
		size = *headOut.ContentLength
	}

	return &adapter.ReadResult{
		MimeType: mimeType,
		Size:     size,
		Stream:   getOut.Body,
	}, nil
}

func (a *S3CompatStorageAdapter) Delete(ctx context.Context, filePath string, digest string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("storage/s3compat: filePath is required")
	}

	key := sanitizeKey(filePath)

	tagDigest, err := a.getDigestTag(ctx, key)
	if err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("storage/s3compat: file not found")
		}
		return fmt.Errorf("storage/s3compat: GetObjectTagging failed: %w", err)
	}
	if tagDigest != digest {
		return fmt.Errorf("storage/s3compat: invalid file hash")
	}

	_, err = a.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("storage/s3compat: file not found")
		}
		return fmt.Errorf("storage/s3compat: DeleteObject failed: %w", err)
	}

	return nil
}

func (a *S3CompatStorageAdapter) getDigestTag(ctx context.Context, key string) (string, error) {
	out, err := a.client.GetObjectTagging(ctx, &awss3.GetObjectTaggingInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}

	for _, tag := range out.TagSet {
		if tag.Key != nil && *tag.Key == "digest" && tag.Value != nil {
			return *tag.Value, nil
		}
	}

	return "", nil
}

func sanitizeKey(filePath string) string {
	return strings.TrimLeft(filePath, "/\\")
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var respErr *smithyhttp.ResponseError
	if errors.As(err, &respErr) {
		return respErr.HTTPStatusCode() == 404
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NoSuchKey", "NotFound":
			return true
		}
	}

	return false
}
