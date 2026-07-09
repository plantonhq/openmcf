package module

import (
	awscodepipelinev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscodepipeline/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsCodePipeline *awscodepipelinev1.AwsCodePipeline

	// PipelineName is the pipeline's cloud name, derived from metadata.name.
	// Pipeline names are create-time-immutable (a name change replaces the
	// pipeline), which is why the name is not spec surface. Same basis as
	// the Terraform module.
	PipelineName string

	// Labels are the resource-identity tags, matching the Terraform module
	// key-for-key. Identity tagging is the only tagging surface this module
	// manages; user-defined custom tags are a platform-wide concern, not
	// per-kind spec surface.
	Labels map[string]string
}

func initializeLocals(ctx *pulumi.Context, in *awscodepipelinev1.AwsCodePipelineStackInput) *Locals {
	locals := &Locals{}
	locals.AwsCodePipeline = in.Target

	locals.PipelineName = locals.AwsCodePipeline.Metadata.Name

	locals.Labels = map[string]string{
		"planton.ai/resource":      "true",
		"planton.ai/organization":  locals.AwsCodePipeline.Metadata.Org,
		"planton.ai/environment":   locals.AwsCodePipeline.Metadata.Env,
		"planton.ai/resource-kind": "AwsCodePipeline",
		"planton.ai/resource-id":   locals.AwsCodePipeline.Metadata.Id,
	}

	return locals
}
