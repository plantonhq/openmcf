package module

import (
	"strconv"

	awsmskserverlessclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmskserverlesscluster/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsMskServerlessCluster *awsmskserverlessclusterv1.AwsMskServerlessCluster

	// ClusterName is metadata.name -- create-only in AWS (max 64 chars), and
	// the basis both engines share so a manifest deploys identically on
	// either.
	ClusterName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsmskserverlessclusterv1.AwsMskServerlessClusterStackInput) *Locals {
	locals := &Locals{}
	locals.AwsMskServerlessCluster = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.ClusterName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key. Tags are
	// the ONLY mutable surface on a serverless MSK cluster.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsMskServerlessCluster.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
