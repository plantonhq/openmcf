package module

import (
	"github.com/pkg/errors"
	kuberneteskeycloakv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kuberneteskeycloak/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *kuberneteskeycloakv1.KubernetesKeycloakStackInput) error {
	//initialize locals
	locals := initializeLocals(ctx, stackInput)

	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to setup kubernetes provider")
	}

	// Conditionally create namespace based on create_namespace flag.
	// The return value is discarded because downstream resources should use
	// pulumi.Provider(kubernetesProvider) instead of pulumi.Parent(namespace)
	// to avoid nil pointer panics when create_namespace is false.
	_, err = namespace(ctx, stackInput, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace")
	}

	// TODO: Keycloak Helm chart deployment
	// When implementing, resources should use:
	// - pulumi.Provider(kubernetesProvider) for Kubernetes resources
	// - pulumi.String(locals.Namespace) for namespace references

	return nil
}
