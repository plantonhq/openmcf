package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// TektonPipelinesNamespace is where Tekton Pipelines installs its
	// cluster-wide configuration. Tekton's own release manifest pins this
	// name, and the runner binary hard-codes the same constant when its
	// readiness probe reads the events-sink configuration -- it is upstream's
	// contract, not a choice this operator makes.
	TektonPipelinesNamespace = "tekton-pipelines"

	// tektonConfigDefaultsName is Tekton's cluster-wide defaults ConfigMap in
	// TektonPipelinesNamespace; tektonCloudEventsSinkKey is the key naming
	// the URL Tekton posts pipeline CloudEvents to.
	tektonConfigDefaultsName = "config-defaults"
	tektonCloudEventsSinkKey = "default-cloud-events-sink"

	// TektonEventsSinkFieldManager owns EXACTLY the sink key in
	// config-defaults. It is deliberately distinct from SSAFieldManager: the
	// Tekton install apply (which ships config-defaults with only an
	// _example key) and the sink write must never share a field-ownership
	// identity, or a future Tekton version bump re-apply would claw the sink
	// key back out and builds would silently degrade to the reconciliation
	// safety net.
	TektonEventsSinkFieldManager = "planton-operator-events-sink"

	// tektonCloudEventPath is the runner webhook's CloudEvents route. MUST
	// match the runner binary's tektonwebhook.CloudEventPath -- two processes
	// agreeing on a path by convention, exactly like the worker health-check
	// key. A sink URL without this path posts to a 404 and every build
	// silently falls back to the reconciliation safety net.
	tektonCloudEventPath = "/service-hub/tekton/cloud-event"

	// TektonControllerDeploymentName / TektonWebhookDeploymentName are the
	// two Deployments that gate build processing: the controller reconciles
	// PipelineRuns and the webhook admits them. The events and resolvers
	// deployments are deliberately not readiness gates -- builds run without
	// them.
	TektonControllerDeploymentName = "tekton-pipelines-controller"
	TektonWebhookDeploymentName    = "tekton-pipelines-webhook"

	// TektonPipelineRunCRDName is the CRD whose presence means "Tekton
	// Pipelines is installed on this cluster" -- the detect half of the
	// detect-or-install prerequisite idiom.
	TektonPipelineRunCRDName = "pipelineruns.tekton.dev"
)

// TektonEventsSinkURL returns the URL Tekton's cluster-wide CloudEvents sink
// must point at: the runner Service's webhook port (Service port 80, so the
// URL needs no explicit port) plus the runner's CloudEvents route. Composed
// from deterministic names so the sink can be written before the runner
// Service exists.
func TektonEventsSinkURL(crName, namespace string) string {
	return fmt.Sprintf("http://%s%s", RunnerServiceFQDN(crName, namespace), tektonCloudEventPath)
}

// TektonConfigDefaultsSinkFragment builds the minimal ConfigMap apply
// configuration that sets ONLY the CloudEvents sink key on Tekton's
// config-defaults. Applied with TektonEventsSinkFieldManager so the operator
// owns exactly this one key -- Tekton's install owns the rest of the object.
// No owner reference: config-defaults is Tekton's object in Tekton's
// namespace, and a cross-namespace owner reference is invalid anyway.
func TektonConfigDefaultsSinkFragment(sinkURL string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tektonConfigDefaultsName,
			Namespace: TektonPipelinesNamespace,
		},
		Data: map[string]string{
			tektonCloudEventsSinkKey: sinkURL,
		},
	}
}

// TektonEventsSinkReadRoleName returns the Role name granting the runner read
// access to Tekton's sink configuration:
// "{namespace}-{crName}-runner-events-sink-read". The pair lives in Tekton's
// own fixed namespace -- shared by every platform on the cluster -- so the
// name carries the platform's namespace: same-named platforms must never
// share a binding whose subject each reconcile would force-apply to its own
// runner.
func TektonEventsSinkReadRoleName(namespace, crName string) string {
	return fmt.Sprintf("%s-%s-runner-events-sink-read", namespace, crName)
}

// TektonEventsSinkReadRole grants the runner `get` on exactly the
// config-defaults ConfigMap in Tekton's install namespace -- without it the
// readiness probe's events-sink check permanently reports "could not
// determine" instead of a verdict. Scoped by resourceNames to the single
// object the probe reads. No owner reference: the Role lives in Tekton's
// namespace, not the CR's, and cross-namespace owner references are invalid.
func TektonEventsSinkReadRole(crName, runnerNamespace string) *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "Role"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TektonEventsSinkReadRoleName(runnerNamespace, crName),
			Namespace: TektonPipelinesNamespace,
			Labels:    runnerLabels(crName),
		},
		Rules: []rbacv1.PolicyRule{{
			APIGroups:     []string{""},
			Resources:     []string{"configmaps"},
			ResourceNames: []string{tektonConfigDefaultsName},
			Verbs:         []string{"get"},
		}},
	}
}

// TektonEventsSinkReadRoleBinding binds the events-sink-read Role to the
// runner's dedicated ServiceAccount in the CR namespace.
func TektonEventsSinkReadRoleBinding(crName, runnerNamespace string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "RoleBinding"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TektonEventsSinkReadRoleName(runnerNamespace, crName),
			Namespace: TektonPipelinesNamespace,
			Labels:    runnerLabels(crName),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     TektonEventsSinkReadRoleName(runnerNamespace, crName),
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      RunnerServiceAccountName(crName),
			Namespace: runnerNamespace,
		}},
	}
}
