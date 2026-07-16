package module

import (
	"strconv"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"

	awscloudfrontv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscloudfront/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds resolved input values and tag metadata for the Pulumi stack.
type Locals struct {
	AwsCloudFront *awscloudfrontv1.AwsCloudFront
	AwsTags       map[string]string
}

// initializeLocals prepares a Locals object by resolving stack input and metadata-derived tags.
func initializeLocals(_ *pulumi.Context, stackInput *awscloudfrontv1.AwsCloudFrontStackInput) *Locals {
	locals := &Locals{}
	locals.AwsCloudFront = stackInput.Target

	metadata := stackInput.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key.
	// CloudFront distributions have no AWS name -- metadata.name drives the
	// Name tag and consumers address the distribution through its
	// ID/ARN/domain outputs.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsCloudFront.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
