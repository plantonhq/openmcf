package verify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// RabbitMqOperatorInstallVerifier checks a RabbitMQ Cluster Operator
// install from its released manifest: the operator Deployment Available,
// the RabbitmqCluster CRD Established, and both admission webhook
// configurations present (their serving certificate is a cert-manager
// Certificate — an Available operator behind failurePolicy: Fail webhooks
// with no issued certificate would still reject every RabbitmqCluster,
// so the webhook wiring is part of the install contract).
//
// Every name below is FIXED by the release manifest (including the
// rabbitmq-system namespace) — the kind installs at most once per
// cluster.
type RabbitMqOperatorInstallVerifier struct{}

const (
	rabbitmqOperatorNamespace  = "rabbitmq-system"
	rabbitmqOperatorDeployment = "rabbitmq-cluster-operator"
	rabbitmqClusterCrd         = "rabbitmqclusters.rabbitmq.com"
)

var rabbitmqOperatorWebhookConfigurations = []struct {
	kind string
	name string
}{
	{"mutatingwebhookconfiguration", "cluster-operator-mutating-webhook-configuration"},
	{"validatingwebhookconfiguration", "cluster-operator-validating-webhook-configuration"},
}

func (v *RabbitMqOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] rabbitmq cluster operator in namespace %q\n", rabbitmqOperatorNamespace)

	if err := kubectlWait(ctx, kubeconfig, "deployment", rabbitmqOperatorDeployment, rabbitmqOperatorNamespace,
		"condition=Available", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "operator deployment %q never became Available", rabbitmqOperatorDeployment)
	}
	if err := kubectlWait(ctx, kubeconfig, "crd", rabbitmqClusterCrd, "",
		"condition=Established", time.Minute); err != nil {
		return errors.Wrapf(err, "CRD %q never became Established", rabbitmqClusterCrd)
	}
	for _, webhook := range rabbitmqOperatorWebhookConfigurations {
		if err := KubectlResourceExists(ctx, kubeconfig, webhook.kind, webhook.name, ""); err != nil {
			return errors.Wrapf(err, "webhook configuration %q not found", webhook.name)
		}
	}
	fmt.Printf("  [verify] operator Available, CRD Established, both admission webhooks present\n")
	return nil
}

func (v *RabbitMqOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The CRD is one document of the applied manifest and deletes with it
	// — unlike the keep-CRD operator kinds, absence IS the destroy
	// contract here (removing the operator removes every RabbitmqCluster;
	// the spec's CRD-lifecycle note carries the warning).
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", rabbitmqOperatorDeployment, rabbitmqOperatorNamespace); err != nil {
		return err
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "crd", rabbitmqClusterCrd, "")
}

// RabbitMqClusterVerifier checks an operator-managed RabbitMQ cluster to
// the point clients can rely on it: the RabbitmqCluster reaches its
// ClusterAvailable and AllReplicasReady conditions, the operator's
// naming contract holds (the `<name>` client Service, the
// `<name>-nodes` headless Service, the `<name>-default-user`
// credentials Secret), and — on every lane — a LIVE message round-trip
// through the management API as the operator-generated admin user
// (declare a queue, publish a marker, get it back: a broker that cannot
// queue and deliver a message is not a broker). The management plugin
// is one of the operator's always-on essentials, so the API doubles as
// the honest liveness surface.
//
// The behavioral-durability scenario (recognized by name) declares a
// QUORUM queue across the cluster's nodes, DELETES a broker pod, proves
// the marker is served DURING the outage, and re-proves it after the
// pod returns.
type RabbitMqClusterVerifier struct {
	Namespace   string
	ClusterName string
	// Replicas is the declared node count — it scales the readiness
	// budget (nodes boot serially on first cluster formation).
	Replicas   int
	Durability bool
}

func (v *RabbitMqClusterVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] rabbitmq cluster %q in namespace %q\n", v.ClusterName, v.Namespace)

	// First boot pulls the server image and forms the cluster node by
	// node — the budget scales with the declared node count.
	budget := 10*time.Minute + time.Duration(v.Replicas-1)*5*time.Minute
	if err := kubectlWait(ctx, kubeconfig, "rabbitmqcluster", v.ClusterName, v.Namespace,
		"condition=ClusterAvailable", budget); err != nil {
		return errors.Wrapf(err, "cluster %q never became ClusterAvailable", v.ClusterName)
	}
	if err := kubectlWait(ctx, kubeconfig, "rabbitmqcluster", v.ClusterName, v.Namespace,
		"condition=AllReplicasReady", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "cluster %q never reached AllReplicasReady", v.ClusterName)
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.ClusterName, v.Namespace); err != nil {
		return errors.Wrap(err, "client service not found")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.ClusterName+"-nodes", v.Namespace); err != nil {
		return errors.Wrap(err, "headless service not found")
	}

	username, password, err := v.defaultUserCredentials(ctx, kubeconfig)
	if err != nil {
		return err
	}
	fmt.Printf("  [verify] default-user Secret present (user %q)\n", username)

	return v.proveMessageRoundTrip(ctx, kubeconfig, username, password)
}

func (v *RabbitMqClusterVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "rabbitmqcluster", v.ClusterName, v.Namespace); err != nil {
		return err
	}
	// The operator's deletion finalizer owns the cascade (the CR deletes
	// with BACKGROUND propagation) — assert its children actually went.
	if err := KubectlResourceAbsent(ctx, kubeconfig, "statefulset", v.ClusterName+"-server", v.Namespace); err != nil {
		return err
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "secret", v.ClusterName+"-default-user", v.Namespace)
}

// defaultUserCredentials reads the operator-generated admin credentials
// from the `<name>-default-user` Secret — the same read a real client
// (or a Service Binding consumer) performs, so it proves the outputs
// contract, not just Secret existence.
func (v *RabbitMqClusterVerifier) defaultUserCredentials(ctx context.Context, kubeconfig string) (string, string, error) {
	secretName := v.ClusterName + "-default-user"
	encodedUsername, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.username}")
	if err != nil {
		return "", "", errors.Wrapf(err, "reading username from secret %q", secretName)
	}
	encodedPassword, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.password}")
	if err != nil {
		return "", "", errors.Wrapf(err, "reading password from secret %q", secretName)
	}
	username, err := base64.StdEncoding.DecodeString(encodedUsername)
	if err != nil {
		return "", "", errors.Wrap(err, "decoding username")
	}
	password, err := base64.StdEncoding.DecodeString(encodedPassword)
	if err != nil {
		return "", "", errors.Wrap(err, "decoding password")
	}
	if len(username) == 0 || len(password) == 0 {
		return "", "", errors.Errorf("secret %q carries empty credentials", secretName)
	}
	return string(username), string(password), nil
}

// proveMessageRoundTrip drives the management API over a port-forward:
// declare a durable queue, publish a run-unique marker, and get it back.
// The durability variant declares a QUORUM queue and re-proves the get
// through a live broker loss.
func (v *RabbitMqClusterVerifier) proveMessageRoundTrip(ctx context.Context, kubeconfig, username, password string) error {
	const localPort = "25672"

	// A kubectl port-forward through a Service binds ONE backing pod at
	// establishment and dies with it — the durability proof DELETES a
	// pod, so the tunnel must be re-establishable (otherwise
	// connection-refused means the dead tunnel, not the broker).
	startTunnel := func() (stop func(), err error) {
		pfCtx, cancel := context.WithCancel(ctx)
		pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
			"port-forward", "svc/"+v.ClusterName, localPort+":15672", "-n", v.Namespace)
		var pfOut strings.Builder
		pf.Stdout = &pfOut
		pf.Stderr = &pfOut
		if err := pf.Start(); err != nil {
			cancel()
			return nil, errors.Wrap(err, "starting port-forward to the client service")
		}
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

	client := &http.Client{Timeout: 30 * time.Second}
	base := "http://127.0.0.1:" + localPort
	queue := "e2e-marker"
	marker := fmt.Sprintf("e2e-marker-%d", time.Now().Unix())

	// The durability substrate is a quorum queue (Raft-replicated across
	// the cluster's nodes); the plain lanes prove the same round-trip on
	// a classic durable queue.
	queueArguments := map[string]interface{}{}
	if v.Durability {
		queueArguments["x-queue-type"] = "quorum"
	}
	declareBody := map[string]interface{}{
		"durable":   true,
		"arguments": queueArguments,
	}
	if err := v.apiCall(ctx, client, username, password, http.MethodPut,
		base+"/api/queues/%2F/"+queue, declareBody, 4*time.Minute, ""); err != nil {
		return errors.Wrap(err, "declaring the queue")
	}
	// The queue is verifier-owned — remove it so a failed later assertion
	// never blocks a re-run (deleting the cluster removes it anyway; this
	// keeps re-runs against a surviving cluster clean).
	defer func() {
		_ = v.apiCall(context.Background(), client, username, password, http.MethodDelete,
			base+"/api/queues/%2F/"+queue, nil, 30*time.Second, "")
	}()

	publishBody := map[string]interface{}{
		"properties":       map[string]interface{}{"delivery_mode": 2},
		"routing_key":      queue,
		"payload":          marker,
		"payload_encoding": "string",
	}
	if err := v.apiCall(ctx, client, username, password, http.MethodPost,
		base+"/api/exchanges/%2F/amq.default/publish", publishBody, 2*time.Minute, `"routed":true`); err != nil {
		return errors.Wrap(err, "publishing the marker")
	}
	if err := v.getMarker(ctx, client, base, username, password, queue, marker, 2*time.Minute); err != nil {
		return errors.Wrap(err, "the published marker never came back")
	}
	fmt.Printf("  [verify] MESSAGE ROUND-TRIP: marker %q published and consumed as user %q\n", marker, username)

	if !v.Durability {
		return nil
	}

	// THE durability proof: kill the first broker, the quorum queue's
	// surviving members must keep serving the marker.
	victim := v.ClusterName + "-server-0"
	fmt.Printf("  [verify] DURABILITY: deleting broker pod %q\n", victim)
	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"delete", "pod", victim, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
		return errors.Wrapf(err, "deleting broker pod: %s", string(out))
	}
	// The deleted pod may be the tunnel's backing pod — re-establish so
	// a refused connection means the BROKER, never the dead tunnel.
	stopTunnel()
	if stopTunnel, err = startTunnel(); err != nil {
		return err
	}
	if err := v.getMarker(ctx, client, base, username, password, queue, marker, 4*time.Minute); err != nil {
		return errors.Wrap(err, "the marker was NOT served during the broker outage")
	}
	fmt.Printf("  [verify] DURABILITY: marker served DURING the broker outage\n")

	if err := kubectlWait(ctx, kubeconfig, "pod", victim, v.Namespace,
		"condition=Ready", 8*time.Minute); err != nil {
		return errors.Wrap(err, "the deleted broker never returned to Ready")
	}
	stopTunnel()
	if stopTunnel, err = startTunnel(); err != nil {
		return err
	}
	if err := v.getMarker(ctx, client, base, username, password, queue, marker, 2*time.Minute); err != nil {
		return errors.Wrap(err, "the marker was lost after the broker recovered")
	}
	fmt.Printf("  [verify] DURABILITY: broker recovered and marker re-read\n")
	return nil
}

// apiCall issues one management-API request with basic auth, retrying
// until 2xx (and, when expect is non-empty, a body containing it) or the
// budget expires (the first calls race management-plugin startup).
func (v *RabbitMqClusterVerifier) apiCall(ctx context.Context, client *http.Client,
	username, password, method, url string, body interface{}, budget time.Duration, expect string,
) error {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		var payload io.Reader
		if body != nil {
			encoded, err := json.Marshal(body)
			if err != nil {
				return err
			}
			payload = strings.NewReader(string(encoded))
		}
		req, err := http.NewRequestWithContext(ctx, method, url, payload)
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
		responseBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 &&
			(expect == "" || strings.Contains(string(responseBody), expect)) {
			return nil
		}
		lastErr = errors.Errorf("%s %s: HTTP %d: %s", method, url, resp.StatusCode, firstLines(string(responseBody), 2))
		time.Sleep(5 * time.Second)
	}
	return lastErr
}

// getMarker retries the management API's queue-get until the marker
// payload comes back. ack_requeue_true keeps the message on the queue so
// the durability proof can re-read it across the broker loss and after
// recovery.
func (v *RabbitMqClusterVerifier) getMarker(ctx context.Context, client *http.Client,
	base, username, password, queue, marker string, budget time.Duration,
) error {
	getBody := map[string]interface{}{
		"count":    1,
		"ackmode":  "ack_requeue_true",
		"encoding": "auto",
	}
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		err := v.apiCall(ctx, client, username, password, http.MethodPost,
			base+"/api/queues/%2F/"+queue+"/get", getBody, 15*time.Second, marker)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(5 * time.Second)
	}
	return lastErr
}
