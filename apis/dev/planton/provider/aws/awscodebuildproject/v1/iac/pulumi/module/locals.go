package module

import (
	awscodebuildprojectv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscodebuildproject/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsCodeBuildProject *awscodebuildprojectv1.AwsCodeBuildProject

	// ProjectName is the project's cloud name, derived from metadata.name.
	// CodeBuild project names are create-time-immutable (a name change
	// replaces the project), which is why the name is not spec surface.
	// Same basis as the Terraform module.
	ProjectName string

	// Labels are the resource-identity tags, matching the Terraform module
	// key-for-key. Identity tagging is the only tagging surface this module
	// manages; user-defined custom tags are a platform-wide concern, not
	// per-kind spec surface.
	Labels map[string]string
}

func initializeLocals(ctx *pulumi.Context, in *awscodebuildprojectv1.AwsCodeBuildProjectStackInput) *Locals {
	locals := &Locals{}
	locals.AwsCodeBuildProject = in.Target

	locals.ProjectName = locals.AwsCodeBuildProject.Metadata.Name

	locals.Labels = map[string]string{
		"planton.ai/resource":      "true",
		"planton.ai/organization":  locals.AwsCodeBuildProject.Metadata.Org,
		"planton.ai/environment":   locals.AwsCodeBuildProject.Metadata.Env,
		"planton.ai/resource-kind": "AwsCodeBuildProject",
		"planton.ai/resource-id":   locals.AwsCodeBuildProject.Metadata.Id,
	}

	return locals
}
