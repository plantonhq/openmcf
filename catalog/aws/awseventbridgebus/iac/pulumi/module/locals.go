package module

import (
	"strconv"

	awseventbridgebusv1alpha1 "github.com/plantonhq/planton/catalog/aws/awseventbridgebus/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target  *awseventbridgebusv1alpha1.AwsEventBridgeBus
	Spec    *awseventbridgebusv1alpha1.AwsEventBridgeBusSpec
	AwsTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, in *awseventbridgebusv1alpha1.AwsEventBridgeBusStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.Target.Metadata.Org,
		awstagkeys.Environment:  locals.Target.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEventBridgeBus.String(),
		awstagkeys.ResourceId:   locals.Target.Metadata.Id,
	}

	return locals
}
