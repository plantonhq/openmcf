package module

import (
	"github.com/pkg/errors"
	kubernetesstatefulsetv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesstatefulset/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *kubernetesstatefulsetv1.KubernetesStatefulSetStackInput) error {
	locals, err := initializeLocals(ctx, stackInput)
	if err != nil {
		return errors.Wrap(err, "failed to initialize locals")
	}

	// Create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	// Conditionally create namespace resource based on create_namespace flag
	if stackInput.Target.Spec.CreateNamespace {
		_, err = kubernetescorev1.NewNamespace(ctx,
			locals.Namespace,
			&kubernetescorev1.NamespaceArgs{
				Metadata: kubernetesmetav1.ObjectMetaPtrInput(
					&kubernetesmetav1.ObjectMetaArgs{
						Name:   pulumi.String(locals.Namespace),
						Labels: pulumi.ToStringMap(locals.Labels),
					}),
			}, pulumi.Provider(kubernetesProvider))
		if err != nil {
			return errors.Wrapf(err, "failed to create %s namespace", locals.Namespace)
		}
	}

	// Create ConfigMaps from spec before StatefulSet
	_, err = configMaps(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create configmaps")
	}

	// Create the headless service for stable network identity (required for StatefulSet)
	createdHeadlessService, err := headlessService(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create headless service")
	}

	// Create the StatefulSet
	createdStatefulSet, err := statefulSet(ctx, locals, kubernetesProvider, createdHeadlessService)
	if err != nil {
		return errors.Wrap(err, "failed to create stateful set")
	}

	// Create ClusterIP service for client access (if ports are defined)
	if err := clientService(ctx, locals, kubernetesProvider, createdStatefulSet); err != nil {
		return errors.Wrap(err, "failed to create client service")
	}

	// Create kubernetes secret with app secrets
	if err := secret(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create secret")
	}

	// Create istio-ingress resources if ingress is enabled
	if locals.KubernetesStatefulSet.Spec.Ingress != nil && locals.KubernetesStatefulSet.Spec.Ingress.Enabled {
		if err := ingress(ctx, locals, kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to create istio ingress resources")
		}
	}

	// Create pod disruption budget if enabled
	if err := podDisruptionBudget(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create pod disruption budget")
	}

	return nil
}
