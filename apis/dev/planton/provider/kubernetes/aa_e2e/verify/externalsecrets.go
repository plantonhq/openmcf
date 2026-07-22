package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// externalSecretTargetName resolves the materialized Secret's name from an
// ExternalSecret manifest: spec.target.name when set, else metadata.name —
// the same defaulting both IaC modules apply.
func externalSecretTargetName(manifestPath string) (string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read manifest %s", manifestPath)
	}
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return "", errors.Wrapf(err, "failed to parse manifest YAML %s", manifestPath)
	}
	if spec, ok := raw["spec"].(map[string]interface{}); ok {
		if target, ok := spec["target"].(map[string]interface{}); ok {
			if name, ok := target["name"].(string); ok && name != "" {
				return name, nil
			}
		}
	}
	if metadata, ok := raw["metadata"].(map[string]interface{}); ok {
		if name, ok := metadata["name"].(string); ok && name != "" {
			return name, nil
		}
	}
	return "", errors.Errorf("manifest %s carries neither spec.target.name nor metadata.name", manifestPath)
}

// ExternalSecretsOperatorInstallVerifier checks an External Secrets Operator
// installation: the three component Deployments (controller, webhook,
// cert-controller) are Available and the core CRDs are Established — the
// exact preconditions every SecretStore/ExternalSecret apply depends on.
type ExternalSecretsOperatorInstallVerifier struct {
	Namespace     string
	ComponentName string
}

func (v *ExternalSecretsOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] external-secrets operator installation %q in namespace %q\n", v.ComponentName, v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for external-secrets operator %q", v.Namespace, v.ComponentName)
	}

	// The release name is fixed to "external-secrets" (one installation per
	// cluster), so the three Deployment names are deterministic.
	for _, deployment := range []string{"external-secrets", "external-secrets-webhook", "external-secrets-cert-controller"} {
		if err := kubectlWait(ctx, kubeconfig, "deployment", deployment, v.Namespace,
			"condition=Available", 2*time.Minute); err != nil {
			return errors.Wrapf(err, "deployment %q not available in namespace %q", deployment, v.Namespace)
		}
	}

	for _, crd := range []string{"externalsecrets.external-secrets.io", "secretstores.external-secrets.io", "clustersecretstores.external-secrets.io"} {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %q not established", crd)
		}
	}

	return nil
}

func (v *ExternalSecretsOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The namespace (created by the module when create_namespace is set)
	// disappearing is the destroy signal; the CRDs deliberately survive
	// uninstall (crds.keep_on_uninstall defaults true) so SecretStore/
	// ExternalSecret data is never cascade-deleted by a release removal.
	return KubectlResourceAbsent(ctx, kubeconfig, "namespace", v.Namespace, "")
}

// SecretStoreVerifier checks an External Secrets Operator SecretStore or
// ClusterSecretStore reaches the Ready condition. In-cluster backends (the
// fake provider) become Ready without any external dependency, so Ready is a
// fair kind-cluster assertion for them.
type SecretStoreVerifier struct {
	// "secretstore" or "clustersecretstore" (fully-qualified resources are
	// unambiguous either way because ESO owns both names).
	Kind      string
	Name      string
	Namespace string // empty for ClusterSecretStore
}

func (v *SecretStoreVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] %s %q ready\n", v.Kind, v.Name)

	if err := KubectlResourceExists(ctx, kubeconfig, v.Kind, v.Name, v.Namespace); err != nil {
		return err
	}
	return kubectlWait(ctx, kubeconfig, v.Kind, v.Name, v.Namespace, "condition=Ready", 2*time.Minute)
}

func (v *SecretStoreVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, v.Kind, v.Name, v.Namespace)
}

// ExternalSecretSyncVerifier checks a REAL sync: the ExternalSecret reaches
// Ready (reason SecretSynced) and the target Secret materializes with data —
// the customer-grade proof that the whole loop (store connection, backend
// read, Secret write-back) actually worked.
type ExternalSecretSyncVerifier struct {
	Name       string
	Namespace  string
	SecretName string
}

func (v *ExternalSecretSyncVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] externalsecret %q sync into secret %q\n", v.Name, v.SecretName)

	if err := KubectlResourceExists(ctx, kubeconfig, "externalsecret", v.Name, v.Namespace); err != nil {
		return err
	}
	if err := kubectlWait(ctx, kubeconfig, "externalsecret", v.Name, v.Namespace,
		"condition=Ready", 2*time.Minute); err != nil {
		return errors.Wrapf(err, "externalsecret %q never became Ready", v.Name)
	}

	// The synced material must actually be in the Secret.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "secret", v.SecretName, "-n", v.Namespace,
		"-o", "jsonpath={.data}").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "materialized secret %q not readable: %s", v.SecretName, string(out))
	}
	if strings.TrimSpace(string(out)) == "" || strings.TrimSpace(string(out)) == "{}" {
		return errors.Errorf("materialized secret %q exists but carries no data", v.SecretName)
	}

	return nil
}

func (v *ExternalSecretSyncVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "externalsecret", v.Name, v.Namespace)
}
