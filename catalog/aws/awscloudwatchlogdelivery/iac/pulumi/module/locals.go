package module

import (
	"strconv"

	awscloudwatchlogdeliveryv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscloudwatchlogdelivery/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awscloudwatchlogdeliveryv1alpha1.AwsCloudwatchLogDelivery
	Spec   *awscloudwatchlogdeliveryv1alpha1.AwsCloudwatchLogDeliverySpec

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, in *awscloudwatchlogdeliveryv1alpha1.AwsCloudwatchLogDeliveryStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	metadata := in.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key
	// (applied to the taggable four: source, destinations, deliveries,
	// and the cross-account destination; the two policy resources are
	// untaggable at AWS). NOTE: AWS does not return tags on Get for the
	// vended family - the provider tracks them via the tagging API.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsCloudwatchLogDelivery.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
