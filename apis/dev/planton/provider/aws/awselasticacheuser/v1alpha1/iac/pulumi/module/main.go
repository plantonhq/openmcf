package module

import (
	"github.com/pkg/errors"
	awselasticacheuserv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awselasticacheuser/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the ElastiCache RBAC user. The user is a leaf of the
// RBAC graph: user groups reference it by id, and caches reference the
// groups -- so credential material lives here, membership lives on the
// group, and the cache itself never changes when access is granted or
// revoked.
func Resources(ctx *pulumi.Context, stackInput *awselasticacheuserv1alpha1.AwsElasticacheUserStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsElasticacheUser.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := user(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create ElastiCache user")
	}

	return nil
}
