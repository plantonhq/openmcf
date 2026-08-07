package module

import (
	"strconv"

	awssesemailidentityv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssesemailidentity/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	AwsSesEmailIdentity *awssesemailidentityv1alpha1.AwsSesEmailIdentity
	AwsTags             map[string]string
	// Dkim is emitted only when the manifest configures signing.
	Dkim *awssesemailidentityv1alpha1.AwsSesEmailIdentityDkimSigning
	// Byodkim is selected by the key/selector pair.
	Byodkim bool
	// Policies maps policy name -> spec entry for per-name satellites.
	Policies map[string]*awssesemailidentityv1alpha1.AwsSesEmailIdentityPolicy
}

func initializeLocals(_ *pulumi.Context, stackInput *awssesemailidentityv1alpha1.AwsSesEmailIdentityStackInput) *Locals {
	locals := &Locals{}
	locals.AwsSesEmailIdentity = stackInput.Target

	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsSesEmailIdentity.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsSesEmailIdentity.Metadata.Org,
		awstagkeys.Environment:  locals.AwsSesEmailIdentity.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsSesEmailIdentity.String(),
		awstagkeys.ResourceId:   locals.AwsSesEmailIdentity.Metadata.Id,
	}

	spec := locals.AwsSesEmailIdentity.Spec
	locals.Dkim = spec.DkimSigning
	locals.Byodkim = locals.Dkim != nil && locals.Dkim.GetDomainSigningPrivateKey() != ""

	locals.Policies = make(map[string]*awssesemailidentityv1alpha1.AwsSesEmailIdentityPolicy)
	for _, policy := range spec.GetPolicies() {
		locals.Policies[policy.Name] = policy
	}

	return locals
}
