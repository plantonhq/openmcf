package module

import (
	"strconv"

	awslambdav1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awslambda/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsLambda *awslambdav1alpha1.AwsLambda

	// FunctionName is metadata.name -- create-time immutable in AWS,
	// and the basis both engines share so a manifest deploys
	// identically on either.
	FunctionName string

	// LogGroupName is where the function's logs land: the custom group
	// from logging_config, or the AWS-default "/aws/lambda/<name>" that
	// Lambda creates on first invocation.
	LogGroupName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awslambdav1alpha1.AwsLambdaStackInput) *Locals {
	locals := &Locals{}
	locals.AwsLambda = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.FunctionName = metadata.Name

	locals.LogGroupName = "/aws/lambda/" + locals.FunctionName
	if lc := stackInput.Target.Spec.LoggingConfig; lc != nil && lc.LogGroup.GetValue() != "" {
		locals.LogGroupName = lc.LogGroup.GetValue()
	}

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsLambda.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
