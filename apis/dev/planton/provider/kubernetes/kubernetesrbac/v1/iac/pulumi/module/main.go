package module

import (
	"github.com/pkg/errors"
	kubernetesrbacv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesrbac/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the main entry point for the Pulumi module.
// It deploys one RBAC grant: a role (created or existing) plus, when subjects are
// present, a binding that points every subject at that role.
func Resources(ctx *pulumi.Context, stackInput *kubernetesrbacv1.KubernetesRbacStackInput) error {
	// Initialize locals with derived values (scope, role name/kind, binding name/kind)
	locals, err := initializeLocals(ctx, stackInput)
	if err != nil {
		return errors.Wrap(err, "failed to initialize locals")
	}

	// Create Kubernetes provider from credentials
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(
		ctx,
		stackInput.ProviderConfig,
		"kubernetes",
	)
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	// Create the Role/ClusterRole when the spec defines one. When binding to an
	// existing role, createdRole stays nil and only the binding is deployed.
	createdRole, err := createRole(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create role")
	}

	// Create the RoleBinding/ClusterRoleBinding when there are subjects to bind
	if err := createBinding(ctx, locals, kubernetesProvider, createdRole); err != nil {
		return errors.Wrap(err, "failed to create binding")
	}

	// Export outputs
	if err := exportOutputs(ctx, locals); err != nil {
		return errors.Wrap(err, "failed to export outputs")
	}

	return nil
}
