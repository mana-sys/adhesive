package adhesive

//
//import (
//	"encoding/hex"
//	"errors"
//	"fmt"
//	"math/rand"
//	"strings"
//	"time"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/service/cloudformation"
//	"github.com/aws/aws-sdk-go/service/s3"
//	cfn "github.com/awslabs/goformation/v4/cloudformation"
//	s3cfn "github.com/awslabs/goformation/v4/cloudformation/s3"
//)
//
//const (
//	ManagedBucketStackName = "adhesive-managed-default-bucket"
//	ManagedBucketPrefix    = "adhesive-managed-default-source-bucket"
//)
//
//var ErrNoManagedBucket = errors.New("couldn't find a managed bucket")
//
//func (a *Adhesive) getOrCreateManagedBucket() (string, error) {
//	a.logger.Println("searching for a managed bucket")
//	bucket, err := a.findManagedBucket()
//	if err == ErrNoManagedBucket {
//		name := fmt.Sprintf("%s-%s", ManagedBucketPrefix, getRandomSuffix())
//		a.logger.Infof("could not find a managed bucket")
//		a.logger.Printf("creating a managed bucket: %s\n", name)
//		if err := a.createManagedBucket(name); err != nil {
//			return "", err
//		}
//
//		return name, nil
//	}
//	if err != nil {
//		return "", err
//	}
//
//	a.logger.Printf("found a managed bucket: %s\n", bucket)
//	return bucket, err
//}
//
//// findManagedBucket tries to find a managed artifacts bucket. If multiple
//// buckets are found, the first bucket is chosen. If no matching buckets are
//// found, then ErrNoManagedBucket is returned.
//func (a *Adhesive) findManagedBucket() (name string, err error) {
//	out, err := a.s3.ListBuckets(&s3.ListBucketsInput{})
//	if err != nil {
//		return "", err
//	}
//
//	for _, bucket := range out.Buckets {
//		if strings.HasPrefix(*bucket.Name, ManagedBucketPrefix) {
//			return *bucket.Name, nil
//		}
//	}
//
//	return "", ErrNoManagedBucket
//}
//
//func getRandomSuffix() string {
//	var b [6]byte
//	rand.Read(b[:])
//	return hex.EncodeToString(b[:])
//}
//
//func (a *Adhesive) createManagedBucket(name string) error {
//	template := generateManagedBucketTemplate(name)
//	j, err := template.JSON()
//	if err != nil {
//		return err
//	}
//
//	out, err := a.cfn.CreateStack(&cloudformation.CreateStackInput{
//		StackName:    aws.String(ManagedBucketStackName),
//		TemplateBody: aws.String(string(j)),
//	})
//	if err != nil {
//		return err
//	}
//
//	if err := a.monitorStack(*out.StackId, ManagedBucketStackName); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func generateManagedBucketTemplate(name string) *cfn.Template {
//	template := cfn.NewTemplate()
//	template.Resources["AdhesiveDeploymentBucket"] = &s3cfn.Bucket{
//		BucketName: name,
//		BucketEncryption: &s3cfn.Bucket_BucketEncryption{
//			ServerSideEncryptionConfiguration: []s3cfn.Bucket_ServerSideEncryptionRule{
//				{
//					ServerSideEncryptionByDefault: &s3cfn.Bucket_ServerSideEncryptionByDefault{
//						SSEAlgorithm: "AES256",
//					},
//				},
//			},
//		},
//	}
//
//	return template
//}
//
//func init() {
//	rand.Seed(time.Now().UnixNano())
//}
