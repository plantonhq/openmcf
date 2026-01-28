package module

import (
	"github.com/pkg/errors"
	kubernetesjenkinsv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesjenkins/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *kubernetesjenkinsv1.KubernetesJenkinsStackInput) error {
	locals := initializeLocals(ctx, stackInput)
	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	//conditionally create namespace resource based on create_namespace flag
	if stackInput.Target.Spec.CreateNamespace {
		_, err = kubernetescorev1.NewNamespace(ctx,
			locals.Namespace,
			&kubernetescorev1.NamespaceArgs{
				Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
					Name:   pulumi.String(locals.Namespace),
					Labels: pulumi.ToStringMap(locals.Labels),
				}),
			},
			pulumi.Provider(kubernetesProvider))
		if err != nil {
			return errors.Wrapf(err, "failed to create %s namespace", locals.Namespace)
		}
	}

	//export name of the namespace
	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))

	//create admin-password secret
	createdAdminPasswordSecret, err := adminCredentials(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create admin password resources")
	}

	//install the jenkins helm-chart
	if err := helmChart(ctx, locals, kubernetesProvider, createdAdminPasswordSecret); err != nil {
		return errors.Wrap(err, "failed to create helm-chart resources")
	}

	//create istio-ingress resources if ingress is enabled.
	if locals.KubernetesJenkins.Spec.Ingress.Enabled {
		if err := ingress(ctx, locals, kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to create ingress resources")
		}
	}

	return nil
}
