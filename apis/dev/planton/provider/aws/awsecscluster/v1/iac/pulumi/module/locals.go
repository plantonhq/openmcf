package module

import (
	"strconv"

	awsecsclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsecscluster/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsEcsCluster *awsecsclusterv1.AwsEcsCluster

	// ClusterName is metadata.name -- create-only in AWS (changing it
	// replaces the cluster), and the basis both engines share so a
	// manifest deploys identically on either.
	ClusterName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsecsclusterv1.AwsEcsClusterStackInput) *Locals {
	locals := &Locals{}
	locals.AwsEcsCluster = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.ClusterName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEcsCluster.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
