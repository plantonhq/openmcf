package module

import (
	"strconv"

	awsrestapidomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsrestapidomain/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awsrestapidomainv1alpha1.AwsRestApiDomain
	Spec   *awsrestapidomainv1alpha1.AwsRestApiDomainSpec

	// EndpointType is the resolved endpoint type: the spec's choice, or
	// REGIONAL when endpoint_configuration is omitted (the right default
	// for almost every new domain). Certificate fan-in keys off it.
	EndpointType string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, in *awsrestapidomainv1alpha1.AwsRestApiDomainStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	locals.EndpointType = "REGIONAL"
	if locals.Spec.EndpointConfiguration != nil && locals.Spec.EndpointConfiguration.Type != "" {
		locals.EndpointType = locals.Spec.EndpointConfiguration.Type
	}

	metadata := in.Target.Metadata

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsRestApiDomain.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
