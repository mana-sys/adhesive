package packager

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/mock"
)

type mockS3 struct {
	mock.Mock
	s3iface.S3API
}

func (m *mockS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.PutObjectOutput), nil
}

func TestPackager_Upload(t *testing.T) {
	m := &mockS3{}
	p := NewFromS3(m)
	p.S3Bucket = "bucket"
	r := bytes.NewReader([]byte("content"))
	m.On("PutObject", &s3.PutObjectInput{
		Bucket: aws.String(p.S3Bucket),
		Body:   r,
		Key:    aws.String("key"),
	}).Return(&s3.PutObjectOutput{})

	key, err := p.Upload(r, "key")
	assert.Nil(t, err)
	assert.Equal(t, "s3://bucket/key", key)

	m.AssertExpectations(t)
}

func TestPackager_UploadWithDedup(t *testing.T) {
	m := &mockS3{}
	p := NewFromS3(m)
	p.S3Bucket = "bucket"
	data := []byte("content")
	hash := md5.Sum(data)
	key := hex.EncodeToString(hash[:])
	m.On("PutObject", &s3.PutObjectInput{
		Bucket: aws.String(p.S3Bucket),
		Body:   bytes.NewReader(data),
		Key:    aws.String(key),
	}).Return(&s3.PutObjectOutput{})

	remoteKey, err := p.UploadWithDedup(data)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("s3://%s/%s", p.S3Bucket, key), remoteKey)
}
