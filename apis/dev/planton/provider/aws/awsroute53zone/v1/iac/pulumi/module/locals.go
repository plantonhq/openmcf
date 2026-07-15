package module

import (
	"strconv"
	"strings"

	awsroute53zonev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsroute53zone/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource, the derived
// zone identity, and the identity tag set.
type Locals struct {
	AwsRoute53Zone *awsroute53zonev1.AwsRoute53Zone

	// ZoneName is the zone's domain name — metadata.name IS the domain
	// (create-time immutable in Route 53).
	ZoneName string

	// ResourceName is the Pulumi logical name: the domain with dots replaced
	// so nested resource names stay readable.
	ResourceName string

	AwsTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsroute53zonev1.AwsRoute53ZoneStackInput) *Locals {
	locals := &Locals{}
	locals.AwsRoute53Zone = stackInput.Target

	locals.ZoneName = locals.AwsRoute53Zone.Metadata.Name
	locals.ResourceName = strings.ReplaceAll(locals.ZoneName, ".", "-")

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.ZoneName,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsRoute53Zone.Metadata.Org,
		awstagkeys.Environment:  locals.AwsRoute53Zone.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsRoute53Zone.String(),
		awstagkeys.ResourceId:   locals.AwsRoute53Zone.Metadata.Id,
	}

	return locals
}
