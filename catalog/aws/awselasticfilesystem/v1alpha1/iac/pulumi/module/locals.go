package module

import (
	"strconv"

	awselasticfilesystemv1alpha1 "github.com/plantonhq/planton/catalog/aws/awselasticfilesystem/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	AwsElasticFileSystem *awselasticfilesystemv1alpha1.AwsElasticFileSystem
	AwsTags              map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awselasticfilesystemv1alpha1.AwsElasticFileSystemStackInput) *Locals {
	locals := &Locals{}
	locals.AwsElasticFileSystem = stackInput.Target

	// Resource-identity tags follow the catalog convention. The Name tag is the
	// resource's metadata.name — the same basis the Terraform module uses for
	// its creation_token, keeping the two engines' physical identity converged.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsElasticFileSystem.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsElasticFileSystem.Metadata.Org,
		awstagkeys.Environment:  locals.AwsElasticFileSystem.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsElasticFileSystem.String(),
		awstagkeys.ResourceId:   locals.AwsElasticFileSystem.Metadata.Id,
	}

	return locals
}
