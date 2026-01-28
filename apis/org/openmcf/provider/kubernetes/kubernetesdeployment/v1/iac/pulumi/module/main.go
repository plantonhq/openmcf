package module

import (
	"github.com/pkg/errors"
	kubernetesdeploymentv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesdeployment/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *kubernetesdeploymentv1.KubernetesDeploymentStackInput) error {
	locals, err := initializeLocals(ctx, stackInput)
	if err != nil {
		return errors.Wrap(err, "failed to initialize locals")
	}

	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	//conditionally create namespace resource based on create_namespace flag
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

	//create ConfigMaps from spec before deployment
	_, err = configMaps(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create configmaps")
	}

	//create kubernetes deployment resources
	createdDeployment, err := deployment(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create microservice deployment")
	}

	//create kubernetes service resources
	if err := service(ctx, locals, kubernetesProvider, createdDeployment); err != nil {
		return errors.Wrap(err, "failed to create microservice kubernetes service resource")
	}

	//create kubernetes secret with app secrets
	if err := secret(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create secret")
	}

	//create istio-ingress resources if ingress is enabled.
	if locals.KubernetesDeployment.Spec.Ingress != nil && locals.KubernetesDeployment.Spec.Ingress.Enabled {
		if err := ingress(ctx, locals, kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to create istio ingress resources")
		}
	}

	//create pod disruption budget if enabled
	if err := podDisruptionBudget(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create pod disruption budget")
	}

	return nil
}
