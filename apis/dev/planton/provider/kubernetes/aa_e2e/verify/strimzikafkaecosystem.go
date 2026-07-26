package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// KafkaConnectVerifier checks a Strimzi-operator-managed Kafka Connect
// cluster to the point it is actually serving: the KafkaConnect resource's
// Ready condition True (workers rolled, the REST API answered the
// operator) and the Connect REST API Service present.
type KafkaConnectVerifier struct {
	Namespace   string
	ConnectName string
}

func (v *KafkaConnectVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] kafka connect cluster %q in namespace %q\n", v.ConnectName, v.Namespace)

	// Ready flips once the workers rolled AND the operator reached the
	// Connect REST API — worker startup includes plugin scanning, so the
	// budget is generous.
	if err := kubectlWait(ctx, kubeconfig, "kafkaconnect", v.ConnectName, v.Namespace,
		"condition=Ready", 12*time.Minute); err != nil {
		return errors.Wrapf(err, "KafkaConnect %q never reached the Ready condition", v.ConnectName)
	}

	if err := KubectlResourceExists(ctx, kubeconfig, "service",
		v.ConnectName+"-connect-api", v.Namespace); err != nil {
		return errors.Wrap(err, "connect REST API service not found")
	}
	return nil
}

func (v *KafkaConnectVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "kafkaconnect", v.ConnectName, v.Namespace)
}

// KafkaConnectorVerifier checks a declared connector all the way to data:
// the KafkaConnector resource Ready (the operator created it on the
// Connect cluster and its tasks came up).
//
// When DataFlow is set (the behavioral-dataflow scenario), the verifier
// proves the pipe end-to-end: the scenario's MirrorSourceConnector
// mirrors the fixture cluster's source topic back to the same cluster
// under the aliased name, so the verifier PRODUCES run-unique markers on
// the source topic and CONSUMES them from the mirrored topic —
// content-matched records crossing a real pipe. A Ready connector that
// moves no data is exactly the failure mode object-existence checks
// cannot see.
type KafkaConnectorVerifier struct {
	Namespace     string
	ConnectorName string
	DataFlow      bool
	// KafkaClusterName names the fixture cluster; SourceTopic and
	// MirroredTopic are the pipe's two ends (paired with the scenario's
	// connector config).
	KafkaClusterName string
	SourceTopic      string
	MirroredTopic    string
}

func (v *KafkaConnectorVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] kafka connector %q in namespace %q\n", v.ConnectorName, v.Namespace)

	if err := kubectlWait(ctx, kubeconfig, "kafkaconnector", v.ConnectorName, v.Namespace,
		"condition=Ready", 6*time.Minute); err != nil {
		return errors.Wrapf(err, "KafkaConnector %q never reached the Ready condition (is it in the Connect cluster's namespace with the right cluster label, and is the class on the workers?)", v.ConnectorName)
	}

	if !v.DataFlow {
		return nil
	}
	return v.proveDataFlow(ctx, kubeconfig)
}

func (v *KafkaConnectorVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "kafkaconnector", v.ConnectorName, v.Namespace)
}

// proveDataFlow produces markers on the source topic and consumes them
// from the connector's mirrored topic — records must cross the pipe.
func (v *KafkaConnectorVerifier) proveDataFlow(ctx context.Context, kubeconfig string) error {
	if v.KafkaClusterName == "" || v.SourceTopic == "" || v.MirroredTopic == "" {
		return errors.New("the behavioral-dataflow scenario needs the fixture cluster name and the pipe's topic pair")
	}

	pod, err := kafkaBrokerPod(ctx, kubeconfig, v.Namespace, v.KafkaClusterName)
	if err != nil {
		return err
	}

	// 1. The source topic (single-broker fixture: RF 1). The connector
	// discovers it on its topic-refresh interval, so creation may race
	// the connector start — the consume loop below absorbs that.
	createTopic := fmt.Sprintf(
		"%s/kafka-topics.sh --bootstrap-server localhost:9092 --create --if-not-exists --topic %s --partitions 1 --replication-factor 1",
		kafkaBinDir, v.SourceTopic)
	if out, err := v.connectorExec(ctx, kubeconfig, pod, createTopic); err != nil {
		return errors.Wrapf(err, "failed to create the source topic: %s", out)
	}

	// 2. Run-unique markers land on the SOURCE end of the pipe.
	runID := fmt.Sprintf("pipe-%d", time.Now().Unix())
	markers := []string{runID + "-0", runID + "-1", runID + "-2"}
	produce := fmt.Sprintf(
		"printf '%s\\n' | %s/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic %s --request-required-acks all",
		strings.Join(markers, "\\n"), kafkaBinDir, v.SourceTopic)
	if out, err := v.connectorExec(ctx, kubeconfig, pod, produce); err != nil {
		return errors.Wrapf(err, "failed to produce the markers: %s", out)
	}

	// 3. The same markers must come out of the MIRRORED end.
	consume := fmt.Sprintf(
		"%s/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic %s --group e2e-dataflow-%d --from-beginning --max-messages %d --timeout-ms 60000",
		kafkaBinDir, v.MirroredTopic, time.Now().Unix(), len(markers))
	deadline := time.Now().Add(6 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		out, err := v.connectorExec(ctx, kubeconfig, pod, consume)
		if err == nil && containsAllMarkers(out, markers) {
			fmt.Printf("  [verify] DATA-FLOW: %d/%d markers crossed the pipe %q → %q through connector %q\n",
				len(markers), len(markers), v.SourceTopic, v.MirroredTopic, v.ConnectorName)
			return nil
		}
		lastErr = errors.Errorf("mirrored consume attempt: err=%v out=%q", err, firstLines(string(out), 3))
		time.Sleep(15 * time.Second)
	}
	return errors.Wrap(lastErr, "the connector never delivered the markers to the mirrored topic")
}

func (v *KafkaConnectorVerifier) connectorExec(ctx context.Context, kubeconfig, pod, command string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", pod, "-n", v.Namespace, "-c", "kafka", "--",
		"bash", "-c", command).CombinedOutput()
	return string(out), err
}

// KafkaMirrorMaker2Verifier checks a MirrorMaker 2 deployment: the
// KafkaMirrorMaker2 resource Ready (the engine rolled and every mirror's
// connectors came up).
//
// When Migration is set (the behavioral-migration scenario), the verifier
// proves the migration story live: it creates a topic on the SOURCE
// fixture cluster, produces run-unique markers there, and consumes them
// back from the TARGET fixture cluster on the mirrored topic name
// ("<sourceAlias>.<topic>" under the default replication policy) — records
// crossing clusters is the whole point of the kind.
type KafkaMirrorMaker2Verifier struct {
	Namespace       string
	MirrorMakerName string
	Migration       bool
	// SourceCluster / TargetCluster name the two Kafka fixture clusters;
	// SourceAlias is the mirror's declared source alias (prefixes the
	// mirrored topic name under the default replication policy).
	SourceCluster string
	TargetCluster string
	SourceAlias   string
}

func (v *KafkaMirrorMaker2Verifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] kafka mirrormaker2 %q in namespace %q\n", v.MirrorMakerName, v.Namespace)

	if err := kubectlWait(ctx, kubeconfig, "kafkamirrormaker2", v.MirrorMakerName, v.Namespace,
		"condition=Ready", 12*time.Minute); err != nil {
		return errors.Wrapf(err, "KafkaMirrorMaker2 %q never reached the Ready condition", v.MirrorMakerName)
	}

	if !v.Migration {
		return nil
	}
	return v.proveMigration(ctx, kubeconfig)
}

func (v *KafkaMirrorMaker2Verifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "kafkamirrormaker2", v.MirrorMakerName, v.Namespace)
}

// proveMigration produces on the source cluster and consumes from the
// target cluster — the records must cross.
func (v *KafkaMirrorMaker2Verifier) proveMigration(ctx context.Context, kubeconfig string) error {
	if v.SourceCluster == "" || v.TargetCluster == "" || v.SourceAlias == "" {
		return errors.New("the behavioral-migration scenario needs the source/target fixture cluster names and the source alias")
	}

	sourcePod, err := kafkaBrokerPod(ctx, kubeconfig, v.Namespace, v.SourceCluster)
	if err != nil {
		return errors.Wrap(err, "source cluster broker pod")
	}
	targetPod, err := kafkaBrokerPod(ctx, kubeconfig, v.Namespace, v.TargetCluster)
	if err != nil {
		return errors.Wrap(err, "target cluster broker pod")
	}

	topic := "e2e-migrate"
	mirroredTopic := fmt.Sprintf("%s.%s", v.SourceAlias, topic)
	runID := fmt.Sprintf("migrate-%d", time.Now().Unix())
	markers := []string{runID + "-0", runID + "-1", runID + "-2", runID + "-3", runID + "-4"}

	// 1. The source topic (single-node fixture: RF 1).
	createTopic := fmt.Sprintf(
		"%s/kafka-topics.sh --bootstrap-server localhost:9092 --create --if-not-exists --topic %s --partitions 3 --replication-factor 1",
		kafkaBinDir, topic)
	if out, err := v.mmExec(ctx, kubeconfig, sourcePod, createTopic); err != nil {
		return errors.Wrapf(err, "failed to create the source topic: %s", out)
	}

	// 2. Markers land on the SOURCE cluster — the cluster being migrated
	// away from.
	produce := fmt.Sprintf(
		"printf '%s\\n' | %s/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic %s --request-required-acks all",
		strings.Join(markers, "\\n"), kafkaBinDir, topic)
	if out, err := v.mmExec(ctx, kubeconfig, sourcePod, produce); err != nil {
		return errors.Wrapf(err, "failed to produce the markers on the source: %s", out)
	}
	fmt.Printf("  [verify] %d markers produced on source cluster %q — waiting for the mirror to carry them to %q\n",
		len(markers), v.SourceCluster, v.TargetCluster)

	// 3. The same markers must arrive on the TARGET cluster under the
	// mirrored topic name. MirrorSourceConnector discovers new topics on
	// its refresh interval, so the first attempt can race topic creation
	// — poll patiently.
	consume := fmt.Sprintf(
		"%s/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic %s --group e2e-migrate-check-%d --from-beginning --max-messages %d --timeout-ms 60000",
		kafkaBinDir, mirroredTopic, time.Now().Unix(), len(markers))
	deadline := time.Now().Add(8 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		out, err := v.mmExec(ctx, kubeconfig, targetPod, consume)
		if err == nil && containsAllMarkers(out, markers) {
			fmt.Printf("  [verify] MIGRATION: %d/%d markers consumed from target topic %q — records crossed clusters\n",
				len(markers), len(markers), mirroredTopic)
			return nil
		}
		lastErr = errors.Errorf("target consume attempt: err=%v out=%q", err, firstLines(out, 3))
		time.Sleep(15 * time.Second)
	}
	return errors.Wrap(lastErr, "the mirror never carried the markers to the target cluster")
}

func (v *KafkaMirrorMaker2Verifier) mmExec(ctx context.Context, kubeconfig, pod, command string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", pod, "-n", v.Namespace, "-c", "kafka", "--",
		"bash", "-c", command).CombinedOutput()
	return string(out), err
}

// firstLines keeps error output readable — verifier retries otherwise
// accumulate full consumer stack traces.
func firstLines(out string, n int) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, " | ")
}

// PAIRED CONSTANTS — these mirror the Kafka-ecosystem E2E fixture
// manifests and the behavioral-dataflow scenario. Changing a fixture's
// metadata.name or the scenario's source file without updating these
// silently breaks the behavioral proofs (the KEDA fixture-namespace
// lesson class).
const (
	// The standard Kafka fixture cluster every ecosystem lane chains
	// (kuberneteskafkatopic/kafkauser/connect/... prerequisites).
	kafkaEcosystemFixtureCluster = "e2e-kaf-fix"
	// The SECOND Kafka fixture the MirrorMaker 2 migration scenario adds
	// as a manifest-path fixture — the cluster being "migrated from".
	kafkaEcosystemSourceCluster = "e2e-kaf-src"
	// The behavioral-dataflow scenario's pipe: the MirrorSourceConnector
	// mirrors connectorDataFlowSourceTopic to the same cluster under the
	// "src." alias prefix.
	connectorDataFlowSourceTopic   = "e2e-pipe-in"
	connectorDataFlowMirroredTopic = "src.e2e-pipe-in"
)

// mirrorMaker2FirstSourceAlias reads the first mirror's source alias from
// a KubernetesKafkaMirrorMaker2 manifest spec.
func mirrorMaker2FirstSourceAlias(spec map[string]interface{}) string {
	if spec == nil {
		return ""
	}
	raw, ok := spec["mirrors"].([]interface{})
	if !ok || len(raw) == 0 {
		return ""
	}
	mirror, ok := raw[0].(map[string]interface{})
	if !ok {
		return ""
	}
	source, ok := mirror["source"].(map[string]interface{})
	if !ok {
		return ""
	}
	alias, _ := source["alias"].(string)
	return alias
}

// karapaceRestProxyEnabled reads the rest-proxy toggle from a
// KubernetesKarapace manifest spec.
func karapaceRestProxyEnabled(spec map[string]interface{}) bool {
	if spec == nil {
		return false
	}
	restProxy, ok := spec["restProxy"].(map[string]interface{})
	if !ok {
		if restProxy, ok = spec["rest_proxy"].(map[string]interface{}); !ok {
			return false
		}
	}
	enabled, _ := restProxy["enabled"].(bool)
	return enabled
}

// kafkaUiFirstClusterName reads the first declared cluster's name from a
// KubernetesKafkaUi manifest spec.
func kafkaUiFirstClusterName(spec map[string]interface{}) string {
	if spec == nil {
		return ""
	}
	raw, ok := spec["clusters"].([]interface{})
	if !ok || len(raw) == 0 {
		return ""
	}
	cluster, ok := raw[0].(map[string]interface{})
	if !ok {
		return ""
	}
	name, _ := cluster["name"].(string)
	return name
}
