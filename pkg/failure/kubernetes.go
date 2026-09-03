package failure

import (
	"fmt"
	"strings"
)

// signatureKubernetesForbidden is the fragment of KubernetesForbidden's
// observation Explain looks for before speaking.
const signatureKubernetesForbidden = "the identity this deploy runs as"

// CustomResourceDefinitionsResource is the resource name the API server uses
// for CRDs in its authorization texts.
const CustomResourceDefinitionsResource = "customresourcedefinitions"

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
