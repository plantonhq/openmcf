package module

import (
	"strconv"

	awsfsxwindowsfilesystemv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsfsxwindowsfilesystem/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsFsxWindowsFileSystem *awsfsxwindowsfilesystemv1alpha1.AwsFsxWindowsFileSystem
	AwsTags                 map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsfsxwindowsfilesystemv1alpha1.AwsFsxWindowsFileSystemStackInput) *Locals {
	locals := &Locals{}
	locals.AwsFsxWindowsFileSystem = stackInput.Target

	// Resource-identity tags follow the catalog convention. The Name tag is
	// the resource's metadata.name — FSx has no name argument, so the console
	// name is this tag; the Terraform module pins the same basis, keeping the
	// two engines' physical identity converged.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsFsxWindowsFileSystem.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsFsxWindowsFileSystem.Metadata.Org,
		awstagkeys.Environment:  locals.AwsFsxWindowsFileSystem.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsFsxWindowsFileSystem.String(),
		awstagkeys.ResourceId:   locals.AwsFsxWindowsFileSystem.Metadata.Id,
	}

	return locals
}
