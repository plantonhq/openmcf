package verify

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// SolrCloudVerifier checks an operator-managed SolrCloud to the point
// clients can rely on it: the SolrCloud reports every node ready and the
// common Service exists. The behavioral-collection scenario (recognized
// by name) additionally proves the cluster WORKS: create a collection
// through the Collections API, index a run-unique document, and query it
// back — the full SolrCloud chain (overseer, ZooKeeper state, shard
// placement, the query path) exercised live.
type SolrCloudVerifier struct {
	Namespace   string
	ClusterName string
	Replicas    int
	// Security drives authenticated API calls with the operator's
	// bootstrapped admin credentials.
	Security   bool
	Behavioral bool
}

func (v *SolrCloudVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] solrcloud %q in namespace %q\n", v.ClusterName, v.Namespace)

	// The ensemble bootstraps before the first Solr node can join —
	// poll the operator's own readiness accounting.
	if err := v.waitForNodesReady(ctx, kubeconfig, 15*time.Minute); err != nil {
		return err
	}
	commonSvc := v.ClusterName + "-solrcloud-common"
	if err := KubectlResourceExists(ctx, kubeconfig, "service", commonSvc, v.Namespace); err != nil {
		return errors.Wrap(err, "common service not found")
	}
	if !v.Behavioral {
		return nil
	}
	return v.proveCollectionQuery(ctx, kubeconfig)
}

func (v *SolrCloudVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "solrcloud", v.ClusterName, v.Namespace)
}

func (v *SolrCloudVerifier) waitForNodesReady(ctx context.Context, kubeconfig string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastReady string
	for time.Now().Before(deadline) {
		ready, _ := kubectlGetJSONPath(ctx, kubeconfig, "solrcloud", v.ClusterName, v.Namespace, "{.status.readyReplicas}")
		lastReady = ready
		if ready == fmt.Sprintf("%d", v.Replicas) {
			fmt.Printf("  [verify] solrcloud reports %s/%d nodes ready\n", ready, v.Replicas)
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("solrcloud never reported %d ready nodes (last readyReplicas %q)", v.Replicas, lastReady)
}

// adminCredentials reads the operator's security-bootstrap Secret
// (`<name>-solrcloud-security-bootstrap`, one key per bootstrapped user
// — the admin user's password under key "admin"; verified in the
// operator's security bootstrap source).
func (v *SolrCloudVerifier) adminCredentials(ctx context.Context, kubeconfig string) (string, string, error) {
	secretName := v.ClusterName + "-solrcloud-security-bootstrap"
	passB64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.admin}")
	if err != nil {
		return "", "", errors.Wrapf(err, "reading %q admin password", secretName)
	}
	pass, err := base64.StdEncoding.DecodeString(passB64)
	if err != nil {
		return "", "", err
	}
	return "admin", strings.TrimSpace(string(pass)), nil
}

// proveCollectionQuery drives the Solr API over a port-forward to the
// common service: CREATE a collection, index a marker document, query it
// back.
func (v *SolrCloudVerifier) proveCollectionQuery(ctx context.Context, kubeconfig string) error {
	const localPort = "18983"

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.ClusterName+"-solrcloud-common", localPort+":80", "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "starting port-forward to the common service")
	}
	// ONE deferred func, cancel FIRST — Wait blocks forever on a
	// port-forward that is never told to exit.
	defer func() {
		cancel()
		_ = pf.Wait()
	}()

	username, password := "", ""
	if v.Security {
		var err error
		username, password, err = v.adminCredentials(ctx, kubeconfig)
		if err != nil {
			return err
		}
	}

	base := "http://127.0.0.1:" + localPort
	stamp := time.Now().Unix()
	collection := fmt.Sprintf("e2e-coll-%d", stamp)
	marker := fmt.Sprintf("e2e-marker-%d", stamp)

	createURL := fmt.Sprintf("%s/solr/admin/collections?action=CREATE&name=%s&numShards=1&replicationFactor=1",
		base, collection)
	if err := v.request(ctx, http.MethodGet, createURL, "", "", username, password, 5*time.Minute); err != nil {
		return errors.Wrap(err, "creating the e2e collection")
	}
	docBody := fmt.Sprintf(`[{"id":"1","marker_s":"%s"}]`, marker)
	indexURL := fmt.Sprintf("%s/solr/%s/update?commit=true", base, collection)
	if err := v.request(ctx, http.MethodPost, indexURL, docBody, "application/json", username, password, 2*time.Minute); err != nil {
		return errors.Wrap(err, "indexing the marker document")
	}
	queryURL := fmt.Sprintf("%s/solr/%s/select?q=%s", base, collection, url.QueryEscape("marker_s:"+marker))
	deadline := time.Now().Add(2 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		body, err := v.get(ctx, queryURL, username, password)
		if err == nil && strings.Contains(body, marker) {
			fmt.Printf("  [verify] COLLECTION/QUERY: collection %q created, marker %q indexed and queried back\n", collection, marker)
			return nil
		}
		if err != nil {
			lastErr = err
		} else {
			lastErr = errors.Errorf("query returned without the marker: %s", firstLines(body, 2))
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Wrap(lastErr, "the indexed document never came back from the query")
}

func (v *SolrCloudVerifier) request(ctx context.Context, method, reqURL, body, contentType, username, password string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, method, reqURL, strings.NewReader(body))
		if err != nil {
			return err
		}
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		if username != "" {
			req.SetBasicAuth(username, password)
		}
		resp, err := http.DefaultClient.Do(req)
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

func (v *SolrCloudVerifier) get(ctx context.Context, reqURL, username, password string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}
	if username != "" {
		req.SetBasicAuth(username, password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("HTTP %d: %s", resp.StatusCode, firstLines(string(body), 2))
	}
	return string(body), nil
}
