package module

import (
	"strconv"

	awskmskeyv1alpha1 "github.com/plantonhq/planton/catalog/aws/awskmskey/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsKmsKey *awskmskeyv1alpha1.AwsKmsKey
	AwsTags   map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awskmskeyv1alpha1.AwsKmsKeyStackInput) *Locals {
	locals := &Locals{}
	locals.AwsKmsKey = stackInput.Target

	metadata := stackInput.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsKmsKey.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
