package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// CertManagerInstallVerifier checks a cert-manager installation: the three
// component Deployments are Available and the core CRDs are Established —
// the exact preconditions every Issuer/Certificate apply depends on.
type CertManagerInstallVerifier struct {
	Namespace     string
	ComponentName string
}

func (v *CertManagerInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] cert-manager installation %q in namespace %q\n", v.ComponentName, v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for cert-manager %q", v.Namespace, v.ComponentName)
	}

	for _, deployment := range []string{"cert-manager", "cert-manager-webhook", "cert-manager-cainjector"} {
		if err := kubectlWait(ctx, kubeconfig, "deployment", deployment, v.Namespace,
			"condition=Available", 2*time.Minute); err != nil {
			return errors.Wrapf(err, "deployment %q not available in namespace %q", deployment, v.Namespace)
		}
	}

	for _, crd := range []string{"certificates.cert-manager.io", "issuers.cert-manager.io", "clusterissuers.cert-manager.io"} {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %q not established", crd)
		}
	}

	return nil
}

func (v *CertManagerInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The namespace (created by the module when create_namespace is set)
	// disappearing is the destroy signal; the CRDs deliberately survive
	// uninstall (crds.keep_on_uninstall defaults true) so certificate data
	// is never cascade-deleted by a release removal.
	return KubectlResourceAbsent(ctx, kubeconfig, "namespace", v.Namespace, "")
}

// IssuerVerifier checks a cert-manager Issuer or ClusterIssuer reaches the
// Ready condition. In-cluster backends (self-signed, CA) become Ready
// without any external dependency, so Ready is a fair kind-cluster
// assertion for them.
type IssuerVerifier struct {
	// "issuer" or "clusterissuer".
	Kind      string
	Name      string
	Namespace string // empty for ClusterIssuer
}

func (v *IssuerVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] %s %q ready\n", v.Kind, v.Name)

	if err := KubectlResourceExists(ctx, kubeconfig, v.Kind, v.Name, v.Namespace); err != nil {
		return err
	}
	return kubectlWait(ctx, kubeconfig, v.Kind, v.Name, v.Namespace, "condition=Ready", 2*time.Minute)
}

func (v *IssuerVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, v.Kind, v.Name, v.Namespace)
}

// CertificateVerifier checks REAL issuance: the Certificate reaches Ready
// and the target TLS Secret materializes with certificate and key data —
// the customer-grade proof that the whole chain (issuer, order, signing,
// secret write-back) actually worked.
type CertificateVerifier struct {
	Name       string
	Namespace  string
	SecretName string
}

func (v *CertificateVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] certificate %q issuance into secret %q\n", v.Name, v.SecretName)

	if err := KubectlResourceExists(ctx, kubeconfig, "certificate", v.Name, v.Namespace); err != nil {
		return err
	}
	if err := kubectlWait(ctx, kubeconfig, "certificate", v.Name, v.Namespace,
		"condition=Ready", 3*time.Minute); err != nil {
		return errors.Wrapf(err, "certificate %q never became Ready", v.Name)
	}

	// The signed material must actually be in the Secret.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "secret", v.SecretName, "-n", v.Namespace,
		"-o", "jsonpath={.data.tls\\.crt}").Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		return errors.Errorf("secret %q has no tls.crt data (issuance did not complete)", v.SecretName)
	}
	out, err = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "secret", v.SecretName, "-n", v.Namespace,
		"-o", "jsonpath={.data.tls\\.key}").Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		return errors.Errorf("secret %q has no tls.key data (issuance did not complete)", v.SecretName)
	}

	return nil
}

func (v *CertificateVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "certificate", v.Name, v.Namespace)
}

// kubectlWait shells out to `kubectl wait` — the same condition-polling the
// kubectl CLI gives users, so a verifier pass means the user's own wait
// would pass too.
func kubectlWait(ctx context.Context, kubeconfig, kind, name, namespace, condition string, timeout time.Duration) error {
	args := []string{"--kubeconfig", kubeconfig, "wait", kind + "/" + name,
		"--for", condition, "--timeout", timeout.String()}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	out, err := exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
	if err != nil {
		return errors.Errorf("kubectl wait %s/%s --for %s failed: %s", kind, name, condition, strings.TrimSpace(string(out)))
	}
	return nil
}
