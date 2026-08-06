package module

import (
	"strconv"

	awsfsxlustrefilesystemv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsfsxlustrefilesystem/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsFsxLustreFileSystem *awsfsxlustrefilesystemv1alpha1.AwsFsxLustreFileSystem
	AwsTags                map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsfsxlustrefilesystemv1alpha1.AwsFsxLustreFileSystemStackInput) *Locals {
	locals := &Locals{}
	locals.AwsFsxLustreFileSystem = stackInput.Target

	// Resource-identity tags follow the catalog convention. The Name tag is
	// the resource's metadata.name — FSx has no name argument, so the console
	// name is this tag; the Terraform module pins the same basis, keeping the
	// two engines' physical identity converged.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsFsxLustreFileSystem.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsFsxLustreFileSystem.Metadata.Org,
		awstagkeys.Environment:  locals.AwsFsxLustreFileSystem.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsFsxLustreFileSystem.String(),
		awstagkeys.ResourceId:   locals.AwsFsxLustreFileSystem.Metadata.Id,
	}

	return locals
}
