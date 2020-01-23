# Adhesive

Adhesive is a tool to facilitate local development for AWS Glue. Based on the
[AWS SAM CLI](https://github.com/awslabs/aws-sam-cli), Adhesive provides a 
local Glue execution environment for writing and testing your Glue
scripts.

## Getting Started

## CLI Command Reference

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
