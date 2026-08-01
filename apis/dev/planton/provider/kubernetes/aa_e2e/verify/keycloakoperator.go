package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// keycloakCrds are the four CRDs the operator bundle installs. The module
// applies them LAST (so a crashed-then-recovered operator finds them) and
// deletes them FIRST on destroy (draining CR finalizers while the operator
// is still alive to process them).
var keycloakCrds = []string{
	"keycloaks.k8s.keycloak.org",
	"keycloakrealmimports.k8s.keycloak.org",
	"keycloakoidcclients.k8s.keycloak.org",
	"keycloaksamlclients.k8s.keycloak.org",
}

// KeycloakOperatorVerifier checks a Keycloak Operator install to the point
// a KubernetesKeycloak declaration could be applied against it: the
// operator Deployment rolled out (upstream's fixed name
// `keycloak-operator`), its Service present, all four k8s.keycloak.org
// CRDs established — and THE DESIGN INVARIANT proven: installing the
// operator alone deploys NO Keycloak server (the declaration kind owns
// every server). Destroy asserts the release-manifest-bundle posture:
// workloads and CRDs delete with the resource.
type KeycloakOperatorVerifier struct {
	// Namespace is the install namespace from the spec (upstream's
	// resource NAMES are fixed; the namespace is the only placement
	// choice).
	Namespace string
}

func (v *KeycloakOperatorVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] keycloak-operator in namespace %q\n", v.Namespace)

	// JOSDK crash-loops until the CRDs land (the module applies them
	// after the Deployment by design) — the rollout wait absorbs those
	// early restarts.
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/keycloak-operator", v.Namespace, 5*time.Minute); err != nil {
		return errors.Wrap(err, "the operator deployment never rolled out")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", "keycloak-operator", v.Namespace); err != nil {
		return errors.Wrap(err, "the operator service not found")
	}
	for _, crd := range keycloakCrds {
		if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
			return errors.Wrapf(err, "CRD %q was not installed", crd)
		}
	}

	// THE DESIGN INVARIANT: the operator install deploys no Keycloak
	// server — a Keycloak CR appearing here would mean the bundle grew
	// an auto-provision path the two-kind grain does not expect.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "keycloaks.k8s.keycloak.org", "-A", "-o", "name").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "listing Keycloak CRs: %s", firstLines(string(out), 3))
	}
	if strings.TrimSpace(string(out)) != "" {
		return errors.Errorf("a Keycloak CR exists after installing the operator alone (found: %s) — the declaration kind owns every server", strings.TrimSpace(string(out)))
	}
	fmt.Printf("  [verify] INVARIANT: operator rolled out, 4 CRDs established, and NO Keycloak server deployed — the declaration kind owns servers\n")
	return nil
}

func (v *KeycloakOperatorVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", "keycloak-operator", v.Namespace); err != nil {
		return err
	}
	// The release-manifest-bundle posture: CRDs delete WITH the
	// resource. Deletion drains any remaining CR finalizers first, so
	// check with a bounded wait rather than a one-shot race.
	deadline := time.Now().Add(2 * time.Minute)
	for _, crd := range keycloakCrds {
		for {
			err := KubectlResourceAbsent(ctx, kubeconfig, "crd", crd, "")
			if err == nil {
				break
			}
			if time.Now().After(deadline) {
				return errors.Wrapf(err, "CRD %q never finished deleting after destroy (the bundle posture: CRDs delete with the resource)", crd)
			}
			time.Sleep(5 * time.Second)
		}
	}
	fmt.Printf("  [verify] DESTROY: operator workloads and all 4 CRDs gone (the release-manifest-bundle posture)\n")
	return nil
}
