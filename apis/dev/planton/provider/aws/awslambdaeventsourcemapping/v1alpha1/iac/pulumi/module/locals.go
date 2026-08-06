package module

import (
	"strconv"

	awslambdaeventsourcemappingv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awslambdaeventsourcemapping/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsLambdaEventSourceMapping *awslambdaeventsourcemappingv1alpha1.AwsLambdaEventSourceMapping

	// MappingName is metadata.name -- the Planton identity for this node.
	// AWS assigns the runtime UUID separately (exported as uuid).
	MappingName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awslambdaeventsourcemappingv1alpha1.AwsLambdaEventSourceMappingStackInput) *Locals {
	locals := &Locals{}
	locals.AwsLambdaEventSourceMapping = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.MappingName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsLambdaEventSourceMapping.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
