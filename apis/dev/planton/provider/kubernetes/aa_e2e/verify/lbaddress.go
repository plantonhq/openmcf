package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Cloud-LB address proofs. Both engines deliberately never block on load
// balancer provisioning at deploy time (the documented never-wait posture:
// composition must not deadlock on cluster capabilities), so the address is
// asserted HERE, by lanes that run only where a cloud LB controller exists
// (the aws-eks profile). These verifiers make the provisioned address itself
// the recorded proof.

// waitForServiceLbAddress polls a Service's .status.loadBalancer.ingress
// until the cloud populates a hostname or IP, returning the address.
func waitForServiceLbAddress(ctx context.Context, kubeconfig, name, namespace string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "service", name, "-n", namespace,
			"-o", "jsonpath={.status.loadBalancer.ingress[0].hostname}{.status.loadBalancer.ingress[0].ip}").CombinedOutput()
		address := strings.TrimSpace(string(out))
		if err == nil && address != "" {
			return address, nil
		}
		last = fmt.Sprintf("address=%q err=%v", address, err)
		time.Sleep(10 * time.Second)
	}
	return "", errors.Errorf("service %s/%s never received a load-balancer address (last: %s)", namespace, name, last)
}

// ServiceLbAddressVerifier proves cloud-LB provisioning for a LoadBalancer
// Service: existence, then a real address in .status.loadBalancer.ingress.
type ServiceLbAddressVerifier struct {
	Namespace string
	Name      string
}

func (v *ServiceLbAddressVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] LoadBalancer service %q must receive a real cloud address\n", v.Name)
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return err
	}
	// NLB provisioning typically lands in 2-3 minutes; 6 covers cold paths.
	address, err := waitForServiceLbAddress(ctx, kubeconfig, v.Name, v.Namespace, 6*time.Minute)
	if err != nil {
		return err
	}
	fmt.Printf("  [verify] cloud LB provisioned: %s\n", address)
	return nil
}

func (v *ServiceLbAddressVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// Deleting the Service is what deprovisions the cloud LB (the in-tree
	// controller finalizes it); absence of the Service is the contract the
	// AWS-side audit then cross-checks by tag.
	return KubectlResourceAbsent(ctx, kubeconfig, "service", v.Name, v.Namespace)
}

// GatewayLbAddressVerifier proves the Gateway API cloud-LB address story: an
// istiod-provisioned gateway service keeps its LoadBalancer default, so the
// Gateway must reach Programmed AND surface a real address in
// .status.addresses — the half of the Gateway story kind lanes pin away.
type GatewayLbAddressVerifier struct {
	Namespace string
	Name      string
}

func (v *GatewayLbAddressVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] gateway %q must be Programmed with a real cloud address\n", v.Name)

	if err := KubectlResourceExists(ctx, kubeconfig, "gateway.gateway.networking.k8s.io", v.Name, v.Namespace); err != nil {
		return err
	}
	// Programmed requires istiod to have provisioned the gateway deployment
	// AND the cloud to have assigned the address — one wait covers both.
	if err := kubectlWait(ctx, kubeconfig, "gateway.gateway.networking.k8s.io", v.Name, v.Namespace,
		"condition=Programmed", 6*time.Minute); err != nil {
		return errors.Wrap(err, "gateway never reached Programmed (istiod provisioning or cloud LB failed)")
	}

	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "gateway.gateway.networking.k8s.io", v.Name, "-n", v.Namespace,
		"-o", "jsonpath={.status.addresses[0].value}").CombinedOutput()
	address := strings.TrimSpace(string(out))
	if err != nil || address == "" {
		return errors.Errorf("gateway %q Programmed but carries no address (out=%q err=%v)", v.Name, address, err)
	}
	fmt.Printf("  [verify] gateway address: %s\n", address)
	return nil
}

func (v *GatewayLbAddressVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "gateway.gateway.networking.k8s.io", v.Name, v.Namespace)
}
