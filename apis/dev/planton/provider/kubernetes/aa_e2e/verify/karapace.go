package verify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// KarapaceVerifier checks a Karapace schema registry to the point clients
// can rely on it: the registry Deployment Available, the Service present,
// and — because a registry that cannot persist schemas is worthless — a
// LIVE register/fetch round-trip through the Confluent-compatible REST
// API: POST a schema for a run-unique subject, then read the exact schema
// back. The round-trip exercises the whole chain (HTTP API → leader
// election → the Kafka-backed _schemas topic → the in-memory reader).
type KarapaceVerifier struct {
	Namespace    string
	RegistryName string
	Port         int
	// RestProxy asserts the optional second deployment when the scenario
	// enables the role.
	RestProxy bool
}

func (v *KarapaceVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] karapace schema registry %q in namespace %q\n", v.RegistryName, v.Namespace)

	if err := kubectlWait(ctx, kubeconfig, "deployment", v.RegistryName, v.Namespace,
		"condition=Available", 6*time.Minute); err != nil {
		return errors.Wrapf(err, "registry deployment %q never became Available (does it reach its Kafka cluster?)", v.RegistryName)
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.RegistryName, v.Namespace); err != nil {
		return errors.Wrap(err, "registry service not found")
	}
	if v.RestProxy {
		if err := kubectlWait(ctx, kubeconfig, "deployment", v.RegistryName+"-rest", v.Namespace,
			"condition=Available", 4*time.Minute); err != nil {
			return errors.Wrap(err, "rest-proxy deployment never became Available")
		}
	}
	return v.proveRegisterFetch(ctx, kubeconfig)
}

func (v *KarapaceVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.RegistryName, v.Namespace); err != nil {
		return err
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "service", v.RegistryName, v.Namespace)
}

// proveRegisterFetch drives the SR API over a port-forward: register an
// Avro schema under a run-unique subject, then fetch version 1 back and
// match the schema body.
func (v *KarapaceVerifier) proveRegisterFetch(ctx context.Context, kubeconfig string) error {
	const localPort = "18081"
	port := v.Port
	if port == 0 {
		port = 8081
	}

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.RegistryName, fmt.Sprintf("%s:%d", localPort, port), "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "starting port-forward to the registry service")
	}
	// ONE deferred func, cancel FIRST — Wait blocks forever on a
	// port-forward that is never told to exit.
	defer func() {
		cancel()
		_ = pf.Wait()
	}()

	subject := fmt.Sprintf("e2e-subject-%d", time.Now().Unix())
	registerBody := `{"schema": "{\"type\": \"string\"}"}`
	base := "http://127.0.0.1:" + localPort

	// Registration needs the registry's Kafka reader warmed up and a
	// leader elected — retry rather than sleeping blind.
	deadline := time.Now().Add(4 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			base+"/subjects/"+subject+"/versions", strings.NewReader(registerBody))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(5 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			lastErr = errors.Errorf("register returned HTTP %d: %s", resp.StatusCode, firstLines(string(body), 2))
			time.Sleep(5 * time.Second)
			continue
		}

		// The fetch must return the exact schema we registered.
		fetchReq, err := http.NewRequestWithContext(ctx, http.MethodGet,
			base+"/subjects/"+subject+"/versions/1", nil)
		if err != nil {
			return err
		}
		fetchResp, err := http.DefaultClient.Do(fetchReq)
		if err != nil {
			lastErr = err
			time.Sleep(5 * time.Second)
			continue
		}
		fetchBody, _ := io.ReadAll(fetchResp.Body)
		_ = fetchResp.Body.Close()
		if fetchResp.StatusCode == http.StatusOK && strings.Contains(string(fetchBody), `\"string\"`) {
			fmt.Printf("  [verify] REGISTER/FETCH: schema registered under %q and read back through the SR API\n", subject)
			return nil
		}
		lastErr = errors.Errorf("fetch returned HTTP %d: %s", fetchResp.StatusCode, firstLines(string(fetchBody), 2))
		time.Sleep(5 * time.Second)
	}
	return errors.Wrapf(lastErr, "the schema register/fetch round-trip never succeeded (port-forward output: %s)", pfOut.String())
}
