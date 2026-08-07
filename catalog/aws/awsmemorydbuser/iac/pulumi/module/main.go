package module

import (
	"github.com/pkg/errors"
	awsmemorydbuserv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsmemorydbuser/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the MemoryDB user. The user is the leaf of MemoryDB's
// ACL graph: ACLs reference it by name, and clusters attach the ACLs -- so
// credential material lives here, membership lives on the ACL, and the
// cluster itself never changes when access is granted or revoked.
func Resources(ctx *pulumi.Context, stackInput *awsmemorydbuserv1alpha1.AwsMemorydbUserStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsMemorydbUser.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := user(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create MemoryDB user")
	}

	return nil
}
