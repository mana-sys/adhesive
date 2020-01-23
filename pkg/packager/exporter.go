package packager

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/glue"
)

type ExportedResources []ExportedResource

type ExportedResource struct {
	GetProperty     func(cloudformation.Resource) (interface{}, bool)
	ReplaceProperty func(resource cloudformation.Resource, path string)
	ForceZip        bool
}

var (
	// Command.ScriptLocation property of AWS::Glue::Job.
	glueJobCommandScriptLocationResource = ExportedResource{
		GetProperty: func(resource cloudformation.Resource) (interface{}, bool) {
			if job, ok := resource.(*glue.Job); ok {
				return job.Command.ScriptLocation, true
			}

			return nil, false
		},
		ReplaceProperty: func(resource cloudformation.Resource, path string) {
			resource.(*glue.Job).Command.ScriptLocation = path
		},
	}

	// DefaultArguments."--extra-py-files" property of AWS::Glue::Job
	glueJobDefaultArgumentsExtraPyFilesResource = ExportedResource{
		GetProperty: func(resource cloudformation.Resource) (interface{}, bool) {
			if job, ok := resource.(*glue.Job); ok {
				if args, ok := job.DefaultArguments.(map[string]interface{}); ok {
					if v, ok := args["--extra-py-files"]; ok {
						return v, true
					}
				}
			}

			return nil, false
		},
		ReplaceProperty: func(resource cloudformation.Resource, path string) {
			resource.(*glue.Job).DefaultArguments.(map[string]interface{})["--extra-py-files"] = path
		},
	}

	// DefaultArguments."--extra-files" property of AWS::Glue::Job
	glueJobDefaultArgumentsExtraFilesResource = ExportedResource{
		GetProperty: func(resource cloudformation.Resource) (interface{}, bool) {
			if job, ok := resource.(*glue.Job); ok {
				if args, ok := job.DefaultArguments.(map[string]interface{}); ok {
					if v, ok := args["--extra-files"]; ok {
						return v, true
					}
				}
			}

			return nil, false
		},
		ReplaceProperty: func(resource cloudformation.Resource, path string) {
			resource.(*glue.Job).DefaultArguments.(map[string]interface{})["--extra-files"] = path
		},
	}
)

var defaultExportedResources = ExportedResources{
	glueJobCommandScriptLocationResource,
	glueJobDefaultArgumentsExtraFilesResource,
	glueJobDefaultArgumentsExtraPyFilesResource,
}
