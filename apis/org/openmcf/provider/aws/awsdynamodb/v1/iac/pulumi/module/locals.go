package module

import (
	"strconv"

	awsdynamodbv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/aws/awsdynamodb/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	Target  *awsdynamodbv1.AwsDynamodb
	Spec    *awsdynamodbv1.AwsDynamodbSpec
	AwsTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, in *awsdynamodbv1.AwsDynamodbStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.Target.Metadata.Org,
		awstagkeys.Environment:  locals.Target.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsDynamodb.String(),
		awstagkeys.ResourceId:   locals.Target.Metadata.Id,
	}
	return locals
}
