package module

import (
	"strconv"

	awseksnodegroupv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awseksnodegroup/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsEksNodeGroup *awseksnodegroupv1alpha1.AwsEksNodeGroup

	// NodeGroupName is metadata.name truncated to AWS's 63-character node
	// group limit, deterministically, so the same manifest always yields
	// the same name on both engines.
	NodeGroupName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awseksnodegroupv1alpha1.AwsEksNodeGroupStackInput) *Locals {
	locals := &Locals{}
	locals.AwsEksNodeGroup = stackInput.Target

	locals.NodeGroupName = stackInput.Target.Metadata.Name
	if len(locals.NodeGroupName) > 63 {
		locals.NodeGroupName = locals.NodeGroupName[:63]
	}

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEksNodeGroup.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
