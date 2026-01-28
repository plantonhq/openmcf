package module

import (
	"github.com/pkg/errors"
	kubernetesstrimzikafkaoperatorv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesstrimzikafkaoperator/v1"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// kafkaOperator installs the Strimzi Kafka Operator on the target cluster.
//
// The function:
//  1. Conditionally creates the namespace based on create_namespace flag.
//  2. Deploys the Helm chart (watch‑any‑namespace=true so one install can
//     manage topics/streams across all namespaces).
//  3. Exports the namespace name so other stacks can import it later.
func kafkaOperator(
	ctx *pulumi.Context,
	target *kubernetesstrimzikafkaoperatorv1.KubernetesStrimziKafkaOperator,
	kubernetesProvider *pulumikubernetes.Provider,
) error {
	// Compute locals to get the namespace name
	l := newLocals(&kubernetesstrimzikafkaoperatorv1.KubernetesStrimziKafkaOperatorStackInput{Target: target})

	// ---------------------------------------------------------------------
	// 1. Namespace - conditionally create based on create_namespace flag
	// ---------------------------------------------------------------------
	if target.Spec.CreateNamespace {
		_, err := corev1.NewNamespace(
			ctx,
			l.namespace,
			&corev1.NamespaceArgs{
				Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
					Name:   pulumi.String(l.namespace),
					Labels: l.labels,
				}),
			},
			pulumi.Provider(kubernetesProvider),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create namespace for Kafka operator")
		}
	}

	// ---------------------------------------------------------------------
	// 2. Helm release
	// ---------------------------------------------------------------------
	_, err := helm.NewRelease(
		ctx,
		l.HelmReleaseName,
		&helm.ReleaseArgs{
			Name:            pulumi.String(l.HelmReleaseName),
			Namespace:       pulumi.String(l.namespace),
			Chart:           pulumi.String(vars.HelmChartName),
			Version:         pulumi.String(vars.HelmChartVersion),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180),
			Values: pulumi.Map{
				"watchAnyNamespace": pulumi.Bool(true),
			},
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(vars.HelmChartRepo),
			},
		},
		pulumi.Provider(kubernetesProvider),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create Strimzi Helm release")
	}

	// ---------------------------------------------------------------------
	// 3. Stack output
	// ---------------------------------------------------------------------
	ctx.Export(OpNamespace, pulumi.String(l.namespace))

	return nil
}
