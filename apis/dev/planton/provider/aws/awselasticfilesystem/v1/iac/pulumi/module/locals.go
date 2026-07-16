package module

import (
	"strconv"

	awselasticfilesystemv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awselasticfilesystem/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	AwsElasticFileSystem *awselasticfilesystemv1.AwsElasticFileSystem
	AwsTags              map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awselasticfilesystemv1.AwsElasticFileSystemStackInput) *Locals {
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
