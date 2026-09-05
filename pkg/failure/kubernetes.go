package failure

import (
	"fmt"
	"strings"
)

// signatureKubernetesForbidden is the fragment of KubernetesForbidden's
// observation Explain looks for before speaking.
const signatureKubernetesForbidden = "the identity this deploy runs as"

// signatureKubernetesFieldNotInDefinition is the fragment of
// KubernetesFieldNotInDefinition's observation Explain looks for before
// speaking.
const signatureKubernetesFieldNotInDefinition = "as a field its installed definition does not know"

// CustomResourceDefinitionsResource is the resource name the API server uses
// for CRDs in its authorization texts.
const CustomResourceDefinitionsResource = "customresourcedefinitions"

// definitionOwners names, for a custom resource kind the catalog declares,
// the catalog resource that installs its definition and the field on it that
// picks the definition's version. A kind absent here gets the generic next
// step (upgrade whatever installs the definition).
//
// Kept to kinds whose definition another catalog resource owns; a module that
// applies its own chart's definitions never hits this failure, because its
// definitions and its declaration come from the same release.
var definitionOwners = map[string]string{
	"PlantonPlatform": "KubernetesPlantonOperator (spec.chart_version)",
}

// KubernetesFieldNotInDefinition: the API server rejected a field of a custom
// resource because the kind's installed definition (its CRD) predates the
// field. The catalog release the declaration comes from knows the field; the
// definition on the cluster does not — the definition was installed from an
// older release of the chart that owns it, and nothing upgrades that chart
// implicitly. Under server-side apply the server refuses outright (an unknown
// field is never pruned silently), but its sentence names only the mechanism.
// kind is the API server's Kind; field is the dotted path without a leading
// dot (spec.ingress.gatewayRef); raw is the server's text, kept whole.
func KubernetesFieldNotInDefinition(kind, field, raw string) *Failure {
	owner, known := definitionOwners[kind]
	next := fmt.Sprintf("upgrade the resource that installs the %s definition to the current catalog default or newer, apply it, then re-apply this resource", kind)
	if known {
		next = fmt.Sprintf("set %s to the current catalog default or newer, apply it, then re-apply this resource", owner)
	}
	return &Failure{
		Observed: fmt.Sprintf("the API server rejected %s on %s %s (the API server answered: %s)", field, kind, signatureKubernetesFieldNotInDefinition, raw),
		Meaning:  fmt.Sprintf("the definition of %s installed on this cluster predates %s: the catalog release this declaration comes from knows the field, the chart that installed the definition is older than it, and nothing upgrades that chart implicitly", kind, field),
		NextStep: next,
	}
}

// KubernetesForbidden: the API server refused a verb to the identity the
// deploy authenticates as. The server's own text names the identity, the verb,
// the resource, and the scope, which is the whole observation; what it never
// says is why the module needs that right and where the module lists every
// right it needs. scope is the server's phrase ("at the cluster scope" or
// `in the namespace "x"`); raw is the server's text, kept whole.
//
// A denied write on CustomResourceDefinitions gets the fuller answer: a module
// on the catalog's derive branch applies its chart's CRDs itself, outside the
// Helm release, so it needs cluster-scoped CRD rights a namespace-admin
// identity never has, and the user has a way around it that keeps the deploy
// namespaced.
func KubernetesForbidden(user, verb, resource, group, scope, raw string) *Failure {
	qualified := resource
	if group != "" {
		qualified = resource + "." + group
	}
	observed := fmt.Sprintf("%s (%s) may not %s %s %s (the API server answered: %s)", signatureKubernetesForbidden, user, verb, qualified, scope, raw)
	if resource == CustomResourceDefinitionsResource {
		return &Failure{
			Observed: observed,
			Meaning:  "the module applies the chart's CustomResourceDefinitions itself, outside the Helm release, so it needs cluster-scoped rights on CRDs that a namespace-admin identity does not have",
			NextStep: fmt.Sprintf("grant the rules in the module's iac/permissions.yaml (customresourcedefinitions: get, list, create, update, patch, and delete when crds.keep_on_uninstall is false) to %s, or run the deploy with an identity that has them; or set spec.crds.install to false and have a cluster administrator apply the CRDs (`helm template --include-crds` renders them)", user),
		}
	}
	return &Failure{
		Observed: observed,
		Meaning:  fmt.Sprintf("the module creates or manages %s and the identity running the deploy was never granted %s on them %s", qualified, verb, strings.TrimSpace(scope)),
		NextStep: fmt.Sprintf("grant %s on %s %s to %s (the module's iac/permissions.yaml lists every rule it needs), or run the deploy with an identity that has it", verb, qualified, strings.TrimSpace(scope), user),
	}
}
