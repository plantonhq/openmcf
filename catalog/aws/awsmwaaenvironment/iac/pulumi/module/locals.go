package module

import (
	"strconv"

	awsmwaaenvironmentv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsmwaaenvironment/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsMwaaEnvironment *awsmwaaenvironmentv1alpha1.AwsMwaaEnvironment

	// EnvironmentName is metadata.name -- create-only in AWS (ForceNew), and
	// the basis both engines share so a manifest deploys identically on
	// either.
	EnvironmentName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsmwaaenvironmentv1alpha1.AwsMwaaEnvironmentStackInput) *Locals {
	locals := &Locals{}
	locals.AwsMwaaEnvironment = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.EnvironmentName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsMwaaEnvironment.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
