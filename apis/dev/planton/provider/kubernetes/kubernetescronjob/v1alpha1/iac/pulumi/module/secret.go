package module

import (
	"github.com/pkg/errors"
	kubernetesv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/workloadpod"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// secret materializes literal secret env values into ONE workload-scoped
// Kubernetes Secret, collected across the job template's app container, every
// sidecar, and every init container. Secrets referenced from existing
// Kubernetes Secrets (secretRef) are wired directly as env references and
// never pass through here.
func secret(ctx *pulumi.Context, locals *Locals, kubernetesProvider pulumi.ProviderResource, namespaceDeps []pulumi.ResourceOption) error {
	jobTemplate := locals.KubernetesCronJob.Spec.JobTemplate

	secretSourceContainers := make([]*kubernetesv1.WorkloadContainer, 0, 2+len(jobTemplate.Container.Sidecars))
	secretSourceContainers = append(secretSourceContainers, jobTemplate.Container.App)
	secretSourceContainers = append(secretSourceContainers, jobTemplate.Container.Sidecars...)
	if jobTemplate.Pod != nil {
		secretSourceContainers = append(secretSourceContainers, jobTemplate.Pod.InitContainers...)
	}

	dataMap := workloadpod.CollectLiteralEnvSecrets(secretSourceContainers...)
	if len(dataMap) == 0 {
		return nil
	}

	secretArgs := &kubernetescorev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.EnvSecretName),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Type:       pulumi.String("Opaque"),
		StringData: pulumi.ToStringMap(dataMap),
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, namespaceDeps...)
	_, err := kubernetescorev1.NewSecret(ctx,
		locals.EnvSecretName,
		secretArgs,
		opts...)
	if err != nil {
		return errors.Wrap(err, "failed to create env secret resource")
	}

	return nil
}
