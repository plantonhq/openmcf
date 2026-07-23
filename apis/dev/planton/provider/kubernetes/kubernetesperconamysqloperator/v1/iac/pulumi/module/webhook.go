package module

import (
	admissionregistrationv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/admissionregistration/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// validationWebhook renders the operator's CR-validation webhook —
// MODULE-OWNED in the widened-watch arms, deliberately. Upstream behavior:
// an operator with cluster-scoped RBAC (which widened watch grants)
// registers this ONE fixed-name, cluster-scoped
// ValidatingWebhookConfiguration at startup, pointing at a Service in its
// own namespace with failurePolicy Fail — and NOTHING ever removes it. The
// object cannot ride the operator Deployment's ownerReference (Kubernetes
// never garbage-collects a cluster-scoped dependent of a namespaced
// owner), and an operator that finds the object already present updates
// only the CA bundle, never the service pointer. Left to the operator,
// uninstalling it therefore strands a Fail-closed webhook whose service no
// longer exists — bricking every future PerconaXtraDBCluster admission in
// the cluster.
//
// The module renders the object FIRST (the release depends on it), so the
// operator's startup create hits AlreadyExists and merely refreshes the CA
// bundle into the module-declared shape — and destroy removes the webhook
// with the resource. The CA bundle and the operator-stamped ownerReference
// are the operator's to manage; IgnoreChanges keeps the module from
// fighting them.
//
// Own-namespace installs render nothing: the chart's namespaced Role
// carries no admissionregistration permissions, the operator's own
// registration attempt is Forbidden (which it treats as a soft skip), and
// the webhook simply does not exist — the upstream posture, preserved.
//
// One cluster carries at most ONE widened-watch operator: the webhook name
// is fixed upstream, so a second widened installation would contend for
// the same object (documented in the kind's README).
func validationWebhook(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.Resource,
) (pulumi.Resource, error) {
	if !locals.WatchWidened {
		return nil, nil
	}

	return admissionregistrationv1.NewValidatingWebhookConfiguration(ctx,
		"percona-xtradbcluster-webhook",
		&admissionregistrationv1.ValidatingWebhookConfigurationArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String("percona-xtradbcluster-webhook"),
				Labels: pulumi.ToStringMap(locals.Labels),
				// The object is a deliberately-shared, fixed-name cluster
				// singleton that the OPERATOR also writes (CA bundle,
				// ownerReference) with its own field manager. patchForce
				// tells the provider's server-side apply to take the
				// fields the module declares even when another manager
				// holds them — without it, adopting an object a previous
				// operator instance created fails with an apply conflict.
				Annotations: pulumi.StringMap{
					"pulumi.com/patchForce": pulumi.String("true"),
				},
			},
			Webhooks: admissionregistrationv1.ValidatingWebhookArray{
				&admissionregistrationv1.ValidatingWebhookArgs{
					Name:                    pulumi.String("validationwebhook.pxc.percona.com"),
					AdmissionReviewVersions: pulumi.StringArray{pulumi.String("v1")},
					ClientConfig: &admissionregistrationv1.WebhookClientConfigArgs{
						Service: &admissionregistrationv1.ServiceReferenceArgs{
							Name:      pulumi.String("percona-xtradb-cluster-operator"),
							Namespace: pulumi.String(locals.Namespace),
							Path:      pulumi.String("/validate-percona-xtradbcluster"),
							Port:      pulumi.Int(443),
						},
					},
					SideEffects:   pulumi.String("None"),
					FailurePolicy: pulumi.String("Fail"),
					Rules: admissionregistrationv1.RuleWithOperationsArray{
						&admissionregistrationv1.RuleWithOperationsArgs{
							ApiGroups:   pulumi.StringArray{pulumi.String("pxc.percona.com")},
							ApiVersions: pulumi.StringArray{pulumi.String("*")},
							Resources:   pulumi.StringArray{pulumi.String("perconaxtradbclusters/*")},
							Operations:  pulumi.StringArray{pulumi.String("CREATE"), pulumi.String("UPDATE")},
						},
					},
				},
			},
		},
		append([]pulumi.ResourceOption{
			pulumi.Provider(kubernetesProvider),
			// The operator patches the CA bundle in at startup (it issues
			// the serving certificate); the module never carries
			// certificate material. The ownerReference it stamps is
			// likewise its own.
			pulumi.IgnoreChanges([]string{
				"webhooks[0].clientConfig.caBundle",
				"metadata.ownerReferences",
			}),
		}, dependsOn(dependencies)...)...)
}
