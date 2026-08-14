package module

import (
	"strconv"

	awscloudtrailv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscloudtrail/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awscloudtrailv1alpha1.AwsCloudTrail
	Spec   *awscloudtrailv1alpha1.AwsCloudTrailSpec

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, in *awscloudtrailv1alpha1.AwsCloudTrailStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	metadata := in.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsCloudTrail.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
