package module

import (
	"strconv"

	awsopensearchdomainv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsopensearchdomain/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awsopensearchdomainv1alpha1.AwsOpenSearchDomain
	Spec   *awsopensearchdomainv1alpha1.AwsOpenSearchDomainSpec
	// DomainName is metadata.name -- create-only in AWS, constrained to
	// ^[a-z][0-9a-z\-]{2,27}$ (3-28 chars), and the basis both engines share so
	// a manifest deploys identically on either.
	DomainName string
	AwsTags    map[string]string
}

func initializeLocals(ctx *pulumi.Context, in *awsopensearchdomainv1alpha1.AwsOpenSearchDomainStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec
	locals.DomainName = in.Target.Metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.Target.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.Target.Metadata.Org,
		awstagkeys.Environment:  locals.Target.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsOpenSearchDomain.String(),
		awstagkeys.ResourceId:   locals.Target.Metadata.Id,
	}

	return locals
}
