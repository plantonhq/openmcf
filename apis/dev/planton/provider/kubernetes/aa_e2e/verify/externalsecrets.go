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
	// BehavioralAwsSm proves a REAL backend read (the aws-sm-irsa scenario):
	// a verifier-owned ExternalSecret against this store must materialize a
	// Kubernetes Secret whose value equals what the batch bootstrap stored
	// in AWS Secrets Manager — the sync loop crossed the cloud boundary
	// with the IRSA identity, keyless. The ExternalSecret is verifier-owned
	// because it must reference the store under test, which fixtures
	// (deploying first) cannot.
	BehavioralAwsSm bool
}

// esoAwsSm* identify the verifier-owned sync driver. The expected value is
// the bootstrap's known plaintext (bootstrap.sh writes the source secret) —
// asserting the exact value proves the read, not merely a Ready condition.
const (
	esoAwsSmNamespace     = "e2e-eso"
	esoAwsSmTargetSecret  = "e2e-sm-proof"
	esoAwsSmExpectedValue = "e2e-sup3r-s3cret"
)

func (v *SecretStoreVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] %s %q ready\n", v.Kind, v.Name)

	if err := KubectlResourceExists(ctx, kubeconfig, v.Kind, v.Name, v.Namespace); err != nil {
		return err
	}
	// Ready is already a real backend call for cloud arms — ESO validates
	// the connection before granting the condition.
	if err := kubectlWait(ctx, kubeconfig, v.Kind, v.Name, v.Namespace, "condition=Ready", 3*time.Minute); err != nil {
		return err
	}

	if !v.BehavioralAwsSm {
		return nil
	}
	return v.proveAwsSmRead(ctx, kubeconfig)
}

// proveAwsSmRead drives the verifier-owned ExternalSecret and asserts the
// materialized Secret's value round-tripped from AWS Secrets Manager.
func (v *SecretStoreVerifier) proveAwsSmRead(ctx context.Context, kubeconfig string) error {
	smSecretName := os.Getenv("PLANTON_E2E_SM_SECRET_NAME")
	if smSecretName == "" {
		return errors.New("PLANTON_E2E_SM_SECRET_NAME unset — the batch bootstrap exports it")
	}
	fmt.Printf("  [verify] behavioral backend read: %q must sync out of AWS Secrets Manager\n", smSecretName)

	storeKind := "ClusterSecretStore"
	if v.Kind == "secretstore" {
		storeKind = "SecretStore"
	}
	externalSecret := fmt.Sprintf(`apiVersion: external-secrets.io/v1
kind: ExternalSecret
metadata:
  name: e2e-sm-proof
  namespace: %s
spec:
  refreshInterval: 30s
  secretStoreRef:
    kind: %s
    name: %s
  target:
    name: %s
  data:
    - secretKey: password
      remoteRef:
        key: %s
        property: password
`, esoAwsSmNamespace, storeKind, v.Name, esoAwsSmTargetSecret, smSecretName)
	esFile, err := writeTempManifest(externalSecret)
	if err != nil {
		return err
	}
	defer os.Remove(esFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", esFile); err != nil {
		return errors.Wrap(err, "failed to apply ExternalSecret")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "externalsecret", "e2e-sm-proof",
			"-n", esoAwsSmNamespace, "--ignore-not-found")
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "secret", esoAwsSmTargetSecret,
			"-n", esoAwsSmNamespace, "--ignore-not-found")
	}()

	deadline := time.Now().Add(3 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "secret", esoAwsSmTargetSecret, "-n", esoAwsSmNamespace,
			"-o", "go-template={{index .data \"password\" | base64decode}}").CombinedOutput()
		value := strings.TrimSpace(string(out))
		if err == nil && value == esoAwsSmExpectedValue {
			fmt.Printf("  [verify] materialized Secret matches the Secrets Manager value — a real keyless read\n")
			return nil
		}
		last = fmt.Sprintf("value-matches=%v err=%v", value == esoAwsSmExpectedValue, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("the ExternalSecret never materialized the Secrets Manager value (last: %s)", last)
}

func (v *SecretStoreVerifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
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
