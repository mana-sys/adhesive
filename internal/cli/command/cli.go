package command

import (
	"os"

	"github.com/aws/aws-sdk-go/service/glue"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mana-sys/adhesive/internal/cli/config"
	log "github.com/sirupsen/logrus"
)

// State represents the state of the Adhesive workflow.
type State struct {
	workflow string
}

// Adhesive represents a running instance of the Adhesive application.
type AdhesiveCli struct {
	State  *State
	Config *config.Config

	FoundConfigFile bool

	cfn  *cloudformation.CloudFormation
	s3   *s3.S3
	glue *glue.Glue
	sess *session.Session
}

func NewAdhesiveCli(path string) (*AdhesiveCli, error) {
	var (
		conf            = config.NewConfig()
		foundConfigFile bool
	)

	if path == "" {
		path = "adhesive.toml"
	}

	// Try reading configuration from adhesive.toml
	err := config.LoadConfigFileInto(conf, path)
	if pathErr, ok := err.(*os.PathError); ok && os.IsNotExist(pathErr) {
		foundConfigFile = true
	} else if err != nil {
		return nil, err
	}

	return &AdhesiveCli{
		Config:          conf,
		FoundConfigFile: foundConfigFile,
	}, nil
}

func (cli *AdhesiveCli) InitializeClients() error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Logger: aws.LoggerFunc(func(args ...interface{}) {
				log.Debug(args...)
			}),
			LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
			Region:   aws.String(cli.Config.Region),
		},
		Profile:           cli.Config.Profile,
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return err
	}

	cli.cfn = cloudformation.New(sess)
	cli.s3 = s3.New(sess)
	cli.glue = glue.New(sess)
	cli.sess = sess
	return nil
}

func (cli *AdhesiveCli) S3() *s3.S3 {
	return cli.s3
}

func (cli *AdhesiveCli) CloudFormation() *cloudformation.CloudFormation {
	return cli.cfn
}

func (cli *AdhesiveCli) Glue() *glue.Glue {
	return cli.glue
}

func (cli *AdhesiveCli) Session() *session.Session {
	return cli.sess
}
