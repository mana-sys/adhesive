package packager

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (p *Packager) ProcessAndUpload(path string, forceZip bool) (string, error) {
	// Make sure the path refers to an actual file.
	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	// If the file is a directory, create a ZIP archive from the
	// directory and upload it. Otherwise, upload the file directly.
	if stat.IsDir() {
		path, err = tempZipFile([]string{path}, "", "artifacts")
		if err != nil {
			return "", err
		}
		defer os.Remove(path)
	}

	// TODO: Implement forceZip

	// Upload the file.
	return p.UploadFileWithDedup(path)
}

// Upload uploads the content of the reader to S3. The object will be named
// using the specified key.
func (p *Packager) Upload(r io.ReadSeeker, key string) (string, error) {
	fmt.Printf("Uploading %s\n", key)

	input := &s3.PutObjectInput{
		Bucket: aws.String(p.S3Bucket),
		Body:   r,
		Key:    aws.String(key),
	}

	if p.KMSKeyID != "" {
		input.ServerSideEncryption = aws.String("aws:kms")
		input.SSEKMSKeyId = aws.String(p.KMSKeyID)
	}

	_, err := p.svc.PutObject(input)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("s3://%s/%s", p.S3Bucket, key), nil
}

// UploadFileWithDedup uploads a file to S3. The S3 key of the resultant object is
// based on the file's MD5 sum.
func (p *Packager) UploadFileWithDedup(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return p.UploadWithDedup(b)
}

// UploadWithDedup
func (p *Packager) UploadWithDedup(data []byte) (string, error) {
	sum := md5.Sum(data)
	remotePath := hex.EncodeToString(sum[:])
	if p.S3Prefix != "" {
		remotePath = p.S3Prefix + "/" + remotePath
	}

	return p.Upload(bytes.NewReader(data), remotePath)
}
