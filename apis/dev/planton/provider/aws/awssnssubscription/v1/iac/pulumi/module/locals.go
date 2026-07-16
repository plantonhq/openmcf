package module

import (
	awssnssubscriptionv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awssnssubscription/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
//
// A subscription has no cloud-side name and is untaggable in AWS — it is
// identified by a server-assigned ARN — so unlike sibling modules there is no
// derived cloud name and no identity-tag map here. metadata.name only drives
// the Pulumi resource name.
type Locals struct {
	Target *awssnssubscriptionv1.AwsSnsSubscription
	Spec   *awssnssubscriptionv1.AwsSnsSubscriptionSpec
}

func initializeLocals(_ *pulumi.Context, in *awssnssubscriptionv1.AwsSnsSubscriptionStackInput) *Locals {
	return &Locals{
		Target: in.Target,
		Spec:   in.Target.Spec,
	}
}
