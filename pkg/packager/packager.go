// Package packager provides a mechanism for packaging CloudFormation templates
// and exporting the resultant artifacts to S3. The functionality of this
// package is similar to that of the "aws cloudformation package" command.
package packager

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/awslabs/goformation/v4"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/mholt/archiver/v3"
)

type Format bool

const (
	FormatJSON Format = false
	FormatYAML Format = true
)

var (
	ErrUnknownFormat = errors.New("unknown format")
)

func isZipFile(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	var buf [512]byte
	if _, err := io.ReadFull(f, buf[:]); err != nil {
		return false, nil
	}

	return http.DetectContentType(buf[:]) == "application/zip", nil
}

type Packager struct {
	Format   Format
	KMSKeyID string
	S3Bucket string
	S3Prefix string

	exportedResources ExportedResources
	svc               s3iface.S3API
}

func New(cfgs ...*aws.Config) (*Packager, error) {
	sess, err := session.NewSession(cfgs...)
	if err != nil {
		return nil, err
	}

	return &Packager{
		svc:               s3.New(sess),
		exportedResources: defaultExportedResources,
	}, nil
}

func NewFromS3(svc s3iface.S3API) *Packager {
	return &Packager{
		Format:            FormatYAML,
		exportedResources: defaultExportedResources,
		svc:               svc,
	}
}

func (p *Packager) PackageTemplateFile(name string) ([]byte, error) {
	template, err := goformation.Open(name)
	if err != nil {
		return nil, err
	}

	if err := p.exportTemplateArtifacts(template); err != nil {
		return nil, err
	}

	return p.marshalTemplate(template)
}

func tempZipFile(sources []string, dir, pattern string) (string, error) {
	name, err := tempName(dir, pattern+"*.zip")
	if err != nil {
		return "", nil
	}

	if err := archiver.Archive(sources, name); err != nil {
		return "", err
	}

	return name, nil
}

// exportTemplateArtifacts exports the template's artifacts by uploading them
// to S3.
func (p *Packager) exportTemplateArtifacts(template *cloudformation.Template) error {
	for _, resource := range template.Resources {
		for _, exportedResource := range p.exportedResources {
			fmt.Fprintln(os.Stderr, exportedResource)
			untyped, ok := exportedResource.GetProperty(resource)
			if !ok {
				continue
			}

			switch path := untyped.(type) {
			case string:
				// Skip if the path is an S3 path.
				if strings.HasPrefix(path, "s3://") {
					continue
				}

				// Upload the file.
				remotePath, err := p.ProcessAndUpload(path, exportedResource.ForceZip)
				if err != nil {
					return err
				}

				fmt.Fprintln(os.Stderr, "Replacing property", remotePath)

				// Replace the original path with the new remote path.
				exportedResource.ReplaceProperty(resource, remotePath)
			default:
				return errors.New("invalid property type")
			}
		}
	}

	return nil
}

func (p *Packager) marshalTemplate(template *cloudformation.Template) ([]byte, error) {
	if p.Format == FormatJSON {
		return template.JSON()
	} else {
		return template.YAML()
	}
}
