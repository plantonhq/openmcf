package module

import (
	"strconv"

	awshttpapidomainv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awshttpapidomain/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the domain name.
type Locals struct {
	AwsHttpApiDomain *awshttpapidomainv1.AwsHttpApiDomain
	AwsTags          map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awshttpapidomainv1.AwsHttpApiDomainStackInput) *Locals {
	locals := &Locals{}
	locals.AwsHttpApiDomain = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsHttpApiDomain.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsHttpApiDomain.Metadata.Org,
		awstagkeys.Environment:  locals.AwsHttpApiDomain.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsHttpApiDomain.String(),
		awstagkeys.ResourceId:   locals.AwsHttpApiDomain.Metadata.Id,
	}

	return locals
}
