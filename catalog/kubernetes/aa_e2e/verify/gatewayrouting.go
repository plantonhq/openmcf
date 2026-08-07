package verify

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// GatewayRoutingBehavioralVerifier proves Gateway API routing BEHAVIOR with a
// real implementation installed (istiod implements the `istio` GatewayClass):
//
//  1. the HTTPRoute is Accepted by its parent Gateway,
//  2. the Gateway reaches Programmed (istiod provisioned and configured the
//     gateway deployment — honest on kind because the mesh fixture pins the
//     gateway service type to ClusterIP),
//  3. a live request THROUGH the auto-provisioned gateway (over a
//     port-forward; no cloud LB exists on kind) returns the backend's
//     response for the route's hostname.
//
// This is the customer-grade contract: an Accepted route that never carries
// traffic is exactly the failure mode object-existence checks cannot see.
type GatewayRoutingBehavioralVerifier struct {
	Namespace string
	RouteName string
	// GatewayName is the parent Gateway fixture's name; istiod names the
	// auto-provisioned deployment/service "<gateway>-<class>" (clone-verified
	// GetDefaultName), so the request targets "<GatewayName>-istio".
	GatewayName string
	// Hostname the route matches; sent as the Host header on the live probe.
	Hostname string
}

func (v *GatewayRoutingBehavioralVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] behavioral gateway routing: route %q via gateway %q host %q\n",
		v.RouteName, v.GatewayName, v.Hostname)

	// 1. Route Accepted by its parent. The condition lives per-parent under
	// status.parents, which kubectl wait cannot address — poll the jsonpath.
	if err := pollJSONPath(ctx, kubeconfig, "httproute", v.RouteName, v.Namespace,
		`{.status.parents[0].conditions[?(@.type=="Accepted")].status}`, "True", 3*time.Minute); err != nil {
		return errors.Wrap(err, "route never Accepted by its parent gateway")
	}

	// 2. Gateway Programmed — istiod provisioned the gateway deployment and
	// assigned the (ClusterIP) address.
	if err := kubectlWait(ctx, kubeconfig, "gateway", v.GatewayName, v.Namespace,
		"condition=Programmed", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "gateway %q never Programmed", v.GatewayName)
	}

	// 3. Live request through the gateway: port-forward the auto-provisioned
	// service and expect the backend's 200 for the route's hostname.
	gatewayService := v.GatewayName + "-istio"
	if err := KubectlResourceExists(ctx, kubeconfig, "service", gatewayService, v.Namespace); err != nil {
		return errors.Wrapf(err, "auto-provisioned gateway service %q not found", gatewayService)
	}
	return v.probeThroughGateway(ctx, kubeconfig, gatewayService)
}

func (v *GatewayRoutingBehavioralVerifier) probeThroughGateway(ctx context.Context, kubeconfig, gatewayService string) error {
	const localPort = "18080"

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+gatewayService, localPort+":80", "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "starting port-forward to the gateway service")
	}
	// ONE deferred func, cancel FIRST: with separate defers the LIFO order
	// runs Wait before cancel, and Wait blocks forever on a port-forward
	// that is never told to exit.
	defer func() {
		cancel()
		_ = pf.Wait()
	}()

	// The gateway pod can be Programmed before its Envoy accepted the route
	// config push — retry the probe briefly rather than sleeping blind.
	deadline := time.Now().Add(2 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:"+localPort+"/", nil)
		if err != nil {
			return err
		}
		req.Host = v.Hostname
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("  [verify] live request routed through %q: HTTP %d — traffic flows\n", gatewayService, resp.StatusCode)
				return nil
			}
			lastErr = errors.Errorf("gateway returned HTTP %d for host %s", resp.StatusCode, v.Hostname)
		} else {
			lastErr = err
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Wrapf(lastErr, "request never routed through the gateway (port-forward output: %s)", pfOut.String())
}

func (v *GatewayRoutingBehavioralVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "httproutes.gateway.networking.k8s.io", v.RouteName, v.Namespace)
}

// pollJSONPath polls a kubectl jsonpath expression until it equals want.
func pollJSONPath(ctx context.Context, kubeconfig, kind, name, namespace, jsonPath, want string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		args := []string{"--kubeconfig", kubeconfig, "get", kind, name, "-o", "jsonpath=" + jsonPath}
		if namespace != "" {
			args = append(args, "-n", namespace)
		}
		out, err := exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
		if err == nil && strings.TrimSpace(string(out)) == want {
			return nil
		}
		last = strings.TrimSpace(string(out))
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("%s/%s %s never reached %q (last: %q)", kind, name, jsonPath, want, last)
}
