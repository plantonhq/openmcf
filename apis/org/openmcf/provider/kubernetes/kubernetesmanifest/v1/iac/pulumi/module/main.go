package module

import (
	"github.com/pkg/errors"
	kubernetesmanifestv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesmanifest/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *kubernetesmanifestv1.KubernetesManifestStackInput) error {
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
	var namespaceResource *kubernetescorev1.Namespace
	if stackInput.Target.Spec.CreateNamespace {
		namespaceResource, err = kubernetescorev1.NewNamespace(ctx,
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

	// Apply the manifest YAML using yamlv2.ConfigGroup
	// yamlv2 provides better CRD ordering and await behavior
	if err := applyManifest(ctx, locals, kubernetesProvider, namespaceResource); err != nil {
		return errors.Wrap(err, "failed to apply manifest")
	}

	return nil
}

// applyManifest applies the raw Kubernetes manifest YAML using yamlv2.ConfigGroup
func applyManifest(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider *pulumikubernetes.Provider,
	namespaceResource *kubernetescorev1.Namespace) error {

	opts := []pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider),
	}

	// If namespace was created, depend on it to ensure proper ordering
	if namespaceResource != nil {
		opts = append(opts, pulumi.DependsOn([]pulumi.Resource{namespaceResource}))
	}

	// Use yamlv2.ConfigGroup which handles multi-document YAML and CRD ordering
	_, err := yamlv2.NewConfigGroup(ctx, "manifest", &yamlv2.ConfigGroupArgs{
		Yaml: pulumi.StringPtr(locals.ManifestYAML),
	}, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to create config group from manifest YAML")
	}

	return nil
}
