package verify

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// OpenSearchClusterVerifier checks an operator-managed OpenSearch
// cluster to the point clients can rely on it: the OpenSearchCluster
// reaches phase RUNNING with every declared node available, and — on
// every lane — a LIVE index/search round-trip through the HTTP API with
// the operator-bootstrapped admin credentials (a search engine that
// cannot index and find a document is not a search engine).
//
// The behavioral-durability scenario (recognized by name) additionally
// indexes a replicated document, DELETES a node pod, proves the document
// is served DURING the outage, and re-proves it after the cluster
// returns to full strength.
type OpenSearchClusterVerifier struct {
	Namespace   string
	ClusterName string
	// TotalNodes is the sum of every declared pool's replicas — the
	// availableNodes target.
	TotalNodes int
	// FirstPoolName names the pool whose pod the durability proof kills.
	FirstPoolName string
	HttpPort      int
	Durability    bool
	// Dashboards asserts the Dashboards deployment when the scenario
	// enables it.
	Dashboards bool
}

func (v *OpenSearchClusterVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] opensearch cluster %q in namespace %q\n", v.ClusterName, v.Namespace)

	// First boot runs TLS bootstrap, the security-config update job and
	// per-pool StatefulSets — poll the operator's own status rather than
	// guessing at pod names.
	if err := v.waitForClusterRunning(ctx, kubeconfig, 15*time.Minute); err != nil {
		return err
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.ClusterName, v.Namespace); err != nil {
		return errors.Wrap(err, "cluster service not found")
	}
	if v.Dashboards {
		if err := kubectlWait(ctx, kubeconfig, "deployment", v.ClusterName+"-dashboards", v.Namespace,
			"condition=Available", 6*time.Minute); err != nil {
			return errors.Wrap(err, "dashboards deployment never became Available")
		}
	}
	return v.proveIndexSearch(ctx, kubeconfig)
}

func (v *OpenSearchClusterVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "opensearchcluster", v.ClusterName, v.Namespace)
}

// waitForClusterRunning polls the CR status until phase RUNNING and
// availableNodes reaches the declared total.
func (v *OpenSearchClusterVerifier) waitForClusterRunning(ctx context.Context, kubeconfig string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastPhase, lastNodes string
	for time.Now().Before(deadline) {
		phase, _ := kubectlGetJSONPath(ctx, kubeconfig, "opensearchcluster", v.ClusterName, v.Namespace, "{.status.phase}")
		nodes, _ := kubectlGetJSONPath(ctx, kubeconfig, "opensearchcluster", v.ClusterName, v.Namespace, "{.status.availableNodes}")
		lastPhase, lastNodes = phase, nodes
		if phase == "RUNNING" && nodes == fmt.Sprintf("%d", v.TotalNodes) {
			fmt.Printf("  [verify] cluster RUNNING with %s/%d nodes available\n", nodes, v.TotalNodes)
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("cluster never reached RUNNING with %d nodes (last phase %q, availableNodes %q)",
		v.TotalNodes, lastPhase, lastNodes)
}

// adminCredentials reads the operator-bootstrapped admin Secret
// (`<cluster>-admin-password`, fields username/password).
func (v *OpenSearchClusterVerifier) adminCredentials(ctx context.Context, kubeconfig string) (string, string, error) {
	secretName := v.ClusterName + "-admin-password"
	userB64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.username}")
	if err != nil {
		return "", "", errors.Wrapf(err, "reading %q username", secretName)
	}
	passB64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.password}")
	if err != nil {
		return "", "", errors.Wrapf(err, "reading %q password", secretName)
	}
	user, err := base64.StdEncoding.DecodeString(userB64)
	if err != nil {
		return "", "", err
	}
	pass, err := base64.StdEncoding.DecodeString(passB64)
	if err != nil {
		return "", "", err
	}
	username := strings.TrimSpace(string(user))
	if username == "" {
		username = "admin"
	}
	return username, strings.TrimSpace(string(pass)), nil
}

// proveIndexSearch drives the HTTP API over a port-forward: create an
// index, index a run-unique marker document, and search it back. The
// durability variant uses number_of_replicas=1 and re-proves the search
// through a live node loss.
func (v *OpenSearchClusterVerifier) proveIndexSearch(ctx context.Context, kubeconfig string) error {
	const localPort = "19200"
	port := v.HttpPort
	if port == 0 {
		port = 9200
	}

	// A kubectl port-forward through a Service binds ONE backing pod at
	// establishment and dies with it — the durability proof DELETES a pod,
	// so the tunnel must be re-establishable (caught live: the
	// during-outage search got connection-refused from the dead tunnel,
	// not from the cluster).
	startTunnel := func() (stop func(), err error) {
		pfCtx, cancel := context.WithCancel(ctx)
		pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
			"port-forward", "svc/"+v.ClusterName, fmt.Sprintf("%s:%d", localPort, port), "-n", v.Namespace)
		var pfOut strings.Builder
		pf.Stdout = &pfOut
		pf.Stderr = &pfOut
		if err := pf.Start(); err != nil {
			cancel()
			return nil, errors.Wrap(err, "starting port-forward to the cluster service")
		}
		// ONE stop func, cancel FIRST — Wait blocks forever on a
		// port-forward that is never told to exit.
		return func() {
			cancel()
			_ = pf.Wait()
		}, nil
	}
	stopTunnel, err := startTunnel()
	if err != nil {
		return err
	}
	defer func() { stopTunnel() }()

	username, password, err := v.adminCredentials(ctx, kubeconfig)
	if err != nil {
		return err
	}

	// The scenarios declare the operator's TLS-on-HTTP posture, so the
	// API speaks https with the operator's self-signed certificates.
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	base := "https://127.0.0.1:" + localPort

	marker := fmt.Sprintf("e2e-marker-%d", time.Now().Unix())
	index := fmt.Sprintf("e2e-index-%d", time.Now().Unix())

	replicas := 0
	if v.Durability {
		// A replicated index is the durability proof's substrate: the
		// primary's loss must leave the document served from a replica.
		replicas = 1
	}
	createBody := fmt.Sprintf(`{"settings":{"number_of_shards":1,"number_of_replicas":%d}}`, replicas)
	if err := v.request(ctx, client, http.MethodPut, base+"/"+index, createBody, username, password, 4*time.Minute); err != nil {
		return errors.Wrap(err, "creating the e2e index")
	}
	docBody := fmt.Sprintf(`{"marker":"%s"}`, marker)
	if err := v.request(ctx, client, http.MethodPut, base+"/"+index+"/_doc/1?refresh=true", docBody, username, password, 2*time.Minute); err != nil {
		return errors.Wrap(err, "indexing the marker document")
	}
	if err := v.search(ctx, client, base, index, marker, username, password, 2*time.Minute); err != nil {
		return errors.Wrap(err, "the indexed document never came back from _search")
	}
	fmt.Printf("  [verify] INDEX/SEARCH: marker %q indexed and found through the HTTP API\n", marker)

	if !v.Durability {
		return nil
	}

	// THE durability proof: kill a node, the document must stay served.
	victim := fmt.Sprintf("%s-%s-0", v.ClusterName, v.FirstPoolName)
	fmt.Printf("  [verify] DURABILITY: deleting node pod %q\n", victim)
	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"delete", "pod", victim, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
		return errors.Wrapf(err, "deleting node pod: %s", string(out))
	}
	// The deleted pod may be the tunnel's backing pod — re-establish so
	// connection-refused means the CLUSTER, never the dead tunnel.
	stopTunnel()
	if stopTunnel, err = startTunnel(); err != nil {
		return err
	}
	if err := v.search(ctx, client, base, index, marker, username, password, 4*time.Minute); err != nil {
		return errors.Wrap(err, "the document was NOT served during the node outage")
	}
	fmt.Printf("  [verify] DURABILITY: marker served DURING the node outage\n")

	if err := v.waitForClusterRunning(ctx, kubeconfig, 10*time.Minute); err != nil {
		return errors.Wrap(err, "cluster never returned to full strength after the node loss")
	}
	// Fresh tunnel again: the recovery window churns endpoints.
	stopTunnel()
	if stopTunnel, err = startTunnel(); err != nil {
		return err
	}
	if err := v.search(ctx, client, base, index, marker, username, password, 2*time.Minute); err != nil {
		return errors.Wrap(err, "the document was lost after recovery")
	}
	fmt.Printf("  [verify] DURABILITY: full strength restored and marker re-read\n")
	return nil
}

// request retries an authenticated request until 2xx or the budget
// expires (first requests race the security plugin's own warm-up).
func (v *OpenSearchClusterVerifier) request(ctx context.Context, client *http.Client, method, url, body, username, password string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
		if err != nil {
			return err
		}
		req.SetBasicAuth(username, password)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(5 * time.Second)
			continue
		}
		respBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		lastErr = errors.Errorf("HTTP %d: %s", resp.StatusCode, firstLines(string(respBody), 2))
		time.Sleep(5 * time.Second)
	}
	return lastErr
}

// search retries until the marker document comes back.
func (v *OpenSearchClusterVerifier) search(ctx context.Context, client *http.Client, base, index, marker, username, password string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			base+"/"+index+"/_search?q=marker:"+marker, nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(username, password)
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(5 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusOK && strings.Contains(string(body), marker) {
			return nil
		}
		lastErr = errors.Errorf("search HTTP %d: %s", resp.StatusCode, firstLines(string(body), 2))
		time.Sleep(5 * time.Second)
	}
	return lastErr
}
