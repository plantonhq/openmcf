package module

import (
	"strconv"

	awsredshiftserverlessnamespacev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsredshiftserverlessnamespace/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRedshiftServerlessNamespace *awsredshiftserverlessnamespacev1alpha1.AwsRedshiftServerlessNamespace

	// NamespaceName is metadata.name -- create-only in AWS, and the
	// basis both engines share so a manifest deploys identically on
	// either.
	NamespaceName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsredshiftserverlessnamespacev1alpha1.AwsRedshiftServerlessNamespaceStackInput) *Locals {
	locals := &Locals{}
	locals.AwsRedshiftServerlessNamespace = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.NamespaceName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsRedshiftServerlessNamespace.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
