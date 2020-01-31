package command

import (
	"fmt"
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

func NewAdhesiveCli() *AdhesiveCli {
	return &AdhesiveCli{}
}

func (cli *AdhesiveCli) ReloadConfigFile(path string) error {
	var (
		conf             = config.NewConfig()
		failFileNotFound = true
		foundConfigFile  bool
	)

	if path == "" {
		path = "adhesive.toml"
		failFileNotFound = false
	}

	// Try reading configuration from adhesive.toml.
	err := config.LoadConfigFileInto(conf, path)
	if pathErr, ok := err.(*os.PathError); ok && os.IsNotExist(pathErr) {
		// If a configuration file was explicitly specified, then fail if it
		// can't be read.
		if failFileNotFound {
			return err
		}

		log.Debugf("Unable to read default configuration file adhesive.toml. " +
			"Continuing with defaults.")
	} else if err != nil {
		fmt.Println("Failing")
		return err
	} else {
		foundConfigFile = true
	}

	cli.FoundConfigFile = foundConfigFile
	cli.Config = conf
	return nil
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
