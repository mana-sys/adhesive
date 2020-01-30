# Adhesive

Adhesive is a tool to facilitate local development for AWS Glue. Based on the
[AWS SAM CLI](https://github.com/awslabs/aws-sam-cli), Adhesive provides a 
local Glue execution environment for writing and testing your Glue
scripts.

## Table of Contents

- [Installation](#Installation)
- [Getting Started](#getting-started)
- [CLI Command Reference](#cli-command-reference)
  * [deploy](#adhesive-deploy)
  * [local](#adhesive-local)
  * [package](#adhesive-package)
  * [remove](#adhesive-remove)

## Installation
To install Adhesive, follow the approrpriate instructions for your platform.

### Prerequisites
The following prerequisites must be installed for Adhesive to work correctly:
- [Docker](https://docs.docker.com/install/)

### MacOS (Homebrew)
```shell script
brew tap mana-sys/adhesive
brew install adhesive
```

## <a name="getting-started"></a>Getting Started

## <a name="cli-command-reference"></a>CLI Command Reference

### `adhesive deploy`
The `adhesive deploy` command deploys your Glue jobs using CloudFormation.

Like the AWS SAM CLI, `adhesive deploy` comes with a guided mode, which walks
you through setting the parameters required for your deployment. This mode also
offers the option for setting default arguments in the `adhesive.toml` 
configuration file. After these default arguments are set, subsequent 
deployments may be done by simply executing `adhesive deploy` again.

| Option | Description |
| --- | --- |
| `-g`, `--guided` | Allow Adhesive to guide you through the deployment |
| `--stack-name` | The name of the CloudFormation stack being deployed to |
| `--template-file` |  The path to your CloudFormation template (default "template.yml") |

### `adhesive local`
The `adhesive local` command is a top level command for subcommands that 
allow for executing Glue jobs and test suites locally. It is comprised
of the following subcommands:
- `pip` - Install Python dependencies in the local environment
- `pyspark` - Run `pyspark` in the local execution environment
- `pytest` - Run test suites for Glue scripts written in Python
- `spark-submit` - Run a Glue job

### `adhesive package`
The `adhesive package` command packages the Glue jobs in the specified
AWS CloudFormation template. It uploads your scripts and their dependencies
to Amazon S3 and outputs a copy of the original template, with references
to local artifacts replaced with their corresponding Amazon S3 locations.

**Usage**:
```
adhesive package [flags]
```

| Option | Description |
| --- | --- |
| `--template-file` | The path where your AWS CloudFormation template is located. |
| `--s3-bucket` | The S3 bucket where artifacts will be uploaded. |
| `--s3-prefix` | The prefix added to the names of the artifacts uploaded to the S3 bucket. |

### `adhesive remove`
The `adhesive remove` command removes the current deployment of your Glue jobs.

| Option | Description |
| --- | --- |
| `--stack-name` | The name of the CloudFormation stack to remove |

