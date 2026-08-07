package verify

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// kafkaBinDir is where the Strimzi Kafka images ship the Kafka CLI.
const kafkaBinDir = "/opt/kafka/bin"

// KafkaClusterVerifier checks a Strimzi-operator-managed KRaft Kafka
// cluster to the point it is actually serving: the Kafka resource's Ready
// condition True, every strimzi.io/cluster pod ready, and the bootstrap
// Service present.
//
// When Durability is set (the behavioral-durability scenario), the
// verifier proves replicated durability live: it creates a topic with
// replication factor 3 and min.insync.replicas 2, produces run-unique
// markers with acks=all THROUGH THE BOOTSTRAP SERVICE (the application
// path), DELETES a broker node outright, consumes every marker back
// DURING the outage, waits for full strength, and consumes again. With
// acks=all + min-ISR 2, every acknowledged write is on at least two
// brokers before the producer hears success — losing one broker must lose
// nothing. That is the claim this proof pins down.
type KafkaClusterVerifier struct {
	Namespace   string
	ClusterName string
	// FirstPoolName names the pool whose pods the verifier execs into
	// (`<cluster>-<pool>-<idx>` is Strimzi's pod naming contract).
	FirstPoolName string
	// TotalNodes is the declared node count across all pools.
	TotalNodes int
	// BootstrapPort is the first internal listener's port.
	BootstrapPort int
	Durability    bool
}

func (v *KafkaClusterVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] strimzi kafka cluster %q in namespace %q (nodes=%d)\n",
		v.ClusterName, v.Namespace, v.TotalNodes)

	// Ready flips once the KRaft quorum formed, every node rolled, and
	// the entity operator (when enabled) deployed. Nodes roll SERIALLY,
	// so the budget scales with the node count.
	readyBudget := 10*time.Minute + time.Duration(v.TotalNodes)*3*time.Minute
	if err := kubectlWait(ctx, kubeconfig, "kafka", v.ClusterName, v.Namespace,
		"condition=Ready", readyBudget); err != nil {
		return errors.Wrapf(err, "Kafka %q never reached the Ready condition", v.ClusterName)
	}

	if err := KubectlResourceExists(ctx, kubeconfig, "service",
		v.ClusterName+"-kafka-bootstrap", v.Namespace); err != nil {
		return errors.Wrap(err, "bootstrap service not found")
	}

	if !v.Durability {
		return nil
	}
	return v.proveStreamDurability(ctx, kubeconfig)
}

func (v *KafkaClusterVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "kafka", v.ClusterName, v.Namespace)
}

// proveStreamDurability runs the produce → broker-loss → consume-during-
// outage → recovery → consume-again cycle through the bootstrap Service.
func (v *KafkaClusterVerifier) proveStreamDurability(ctx context.Context, kubeconfig string) error {
	if v.TotalNodes < 3 {
		return errors.New("behavioral durability needs 3 broker nodes — the replicated topic must survive the loss")
	}

	bootstrap := fmt.Sprintf("%s-kafka-bootstrap:%d", v.ClusterName, v.BootstrapPort)
	execPod := fmt.Sprintf("%s-%s-0", v.ClusterName, v.FirstPoolName)
	victim := fmt.Sprintf("%s-%s-1", v.ClusterName, v.FirstPoolName)
	topic := "e2e-durability"
	runID := fmt.Sprintf("run-%d", time.Now().Unix())
	markerCount := 5

	// 1. A topic whose durability contract survives one broker loss:
	//    RF 3, min-ISR 2 — acks=all writes are on two brokers before the
	//    producer hears success.
	createTopic := fmt.Sprintf(
		"%s/kafka-topics.sh --bootstrap-server %s --create --if-not-exists --topic %s --partitions 3 --replication-factor 3 --config min.insync.replicas=2",
		kafkaBinDir, bootstrap, topic)
	if out, err := v.kafkaExec(ctx, kubeconfig, execPod, createTopic); err != nil {
		return errors.Wrapf(err, "failed to create the durability topic: %s", out)
	}

	// 2. Produce the run-unique markers with acks=all through the
	//    bootstrap Service — the application path.
	markers := make([]string, 0, markerCount)
	for i := 0; i < markerCount; i++ {
		markers = append(markers, fmt.Sprintf("%s-marker-%d", runID, i))
	}
	produce := fmt.Sprintf(
		"printf '%s\\n' | %s/kafka-console-producer.sh --bootstrap-server %s --topic %s --request-required-acks all",
		strings.Join(markers, "\\n"), kafkaBinDir, bootstrap, topic)
	if out, err := v.kafkaExec(ctx, kubeconfig, execPod, produce); err != nil {
		return errors.Wrapf(err, "failed to produce the markers: %s", out)
	}
	fmt.Printf("  [verify] %d markers acknowledged (acks=all) through %q — killing broker %q\n",
		markerCount, bootstrap, victim)

	// 3. The disaster: a broker node deleted outright. The controller
	//    quorum keeps 2 of 3; every partition keeps an in-sync replica.
	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"delete", "pod", victim, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
		return errors.Errorf("failed to delete broker node: %v: %s", err, string(out))
	}

	// 4. Every marker must be consumable DURING the outage.
	if err := v.consumeMarkers(ctx, kubeconfig, execPod, bootstrap, topic,
		fmt.Sprintf("%s-outage", runID), markers, 4*time.Minute); err != nil {
		return errors.Wrap(err, "markers not consumable while the broker was down")
	}
	fmt.Printf("  [verify] DURABILITY: %d/%d markers consumed during the broker outage\n",
		markerCount, markerCount)

	// 5. Full recovery: the pod set returns and the cluster reports
	//    Ready again — then the markers must still be there.
	if err := kubectlWait(ctx, kubeconfig, "kafka", v.ClusterName, v.Namespace,
		"condition=Ready", 10*time.Minute); err != nil {
		return errors.Wrap(err, "cluster never returned to Ready after the broker loss")
	}
	if err := v.consumeMarkers(ctx, kubeconfig, execPod, bootstrap, topic,
		fmt.Sprintf("%s-recovered", runID), markers, 3*time.Minute); err != nil {
		return errors.Wrap(err, "markers not consumable after recovery")
	}
	fmt.Printf("  [verify] DURABILITY: markers intact after the broker rejoined — the cluster lost nothing\n")
	return nil
}

// consumeMarkers reads the topic from the beginning under a fresh consumer
// group and asserts every marker line came back.
func (v *KafkaClusterVerifier) consumeMarkers(ctx context.Context, kubeconfig, execPod, bootstrap, topic, group string, markers []string, budget time.Duration) error {
	consume := fmt.Sprintf(
		"%s/kafka-console-consumer.sh --bootstrap-server %s --topic %s --group %s --from-beginning --max-messages %d --timeout-ms 120000",
		kafkaBinDir, bootstrap, topic, group, len(markers))

	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		out, err := v.kafkaExec(ctx, kubeconfig, execPod, consume)
		if err == nil && containsAllMarkers(out, markers) {
			return nil
		}
		lastErr = errors.Errorf("consume attempt: err=%v out=%q", err, strings.TrimSpace(out))
		time.Sleep(10 * time.Second)
	}
	return lastErr
}

func containsAllMarkers(out string, markers []string) bool {
	for _, marker := range markers {
		if !strings.Contains(out, marker) {
			return false
		}
	}
	return true
}

// kafkaExec runs a shell command in a Kafka node's kafka container.
func (v *KafkaClusterVerifier) kafkaExec(ctx context.Context, kubeconfig, pod, command string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", pod, "-n", v.Namespace, "-c", "kafka", "--",
		"bash", "-c", command).CombinedOutput()
	return string(out), err
}

// kafkaFirstPool reads the first node pool's name and the total declared
// node count from a KubernetesKafka manifest spec.
func kafkaFirstPool(spec map[string]interface{}) (string, int) {
	if spec == nil {
		return "", 0
	}
	raw, ok := spec["nodePools"]
	if !ok {
		if raw, ok = spec["node_pools"]; !ok {
			return "", 0
		}
	}
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return "", 0
	}
	total := 0
	firstName := ""
	for i, entry := range list {
		pool, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		if i == 0 {
			firstName, _ = pool["name"].(string)
		}
		switch replicas := pool["replicas"].(type) {
		case int:
			total += replicas
		case int64:
			total += int(replicas)
		case float64:
			total += int(replicas)
		}
	}
	return firstName, total
}

// kafkaFirstInternalPort reads the first internal-type listener's port —
// the bootstrap port the verifiers speak through.
func kafkaFirstInternalPort(spec map[string]interface{}) int {
	if spec == nil {
		return 9092
	}
	raw, ok := spec["listeners"]
	if !ok {
		return 9092
	}
	list, ok := raw.([]interface{})
	if !ok {
		return 9092
	}
	for _, entry := range list {
		listener, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		listenerType, _ := listener["type"].(string)
		if listenerType != "" && listenerType != "internal" {
			continue
		}
		switch port := listener["port"].(type) {
		case int:
			return port
		case int64:
			return int(port)
		case float64:
			return int(port)
		}
	}
	return 9092
}

// KafkaTopicVerifier checks a declared topic all the way to Kafka's own
// metadata: the KafkaTopic resource Ready (the topic operator reconciled
// it) AND kafka-topics.sh --describe on the target cluster reporting the
// declared partition count.
type KafkaTopicVerifier struct {
	Namespace   string
	CrName      string
	TopicName   string
	ClusterName string
	// Partitions is the declared count (0 = not declared; the describe
	// assertion then only checks existence).
	Partitions int
}

func (v *KafkaTopicVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] kafka topic %q (CR %q) on cluster %q\n", v.TopicName, v.CrName, v.ClusterName)

	if err := kubectlWait(ctx, kubeconfig, "kafkatopic", v.CrName, v.Namespace,
		"condition=Ready", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "KafkaTopic %q never reached the Ready condition (is it in the cluster's own namespace with the right cluster label?)", v.CrName)
	}

	// Ready says the topic operator reconciled; the describe proves the
	// topic exists IN KAFKA with the declared shape.
	pod, err := kafkaBrokerPod(ctx, kubeconfig, v.Namespace, v.ClusterName)
	if err != nil {
		return err
	}
	describe := fmt.Sprintf("%s/kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic %s",
		kafkaBinDir, v.TopicName)
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", pod, "-n", v.Namespace, "-c", "kafka", "--",
		"bash", "-c", describe).CombinedOutput()
	if err != nil {
		return errors.Errorf("kafka-topics describe failed: %v: %s", err, string(out))
	}
	if v.Partitions > 0 && !strings.Contains(string(out), fmt.Sprintf("PartitionCount: %d", v.Partitions)) {
		return errors.Errorf("topic %q exists but not with %d partitions:\n%s", v.TopicName, v.Partitions, string(out))
	}
	fmt.Printf("  [verify] RECONCILED: %s\n", strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0])
	return nil
}

func (v *KafkaTopicVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "kafkatopic", v.CrName, v.Namespace)
}

// KafkaUserVerifier checks a declared user to the point clients can use
// it: the KafkaUser resource Ready and the operator-generated credentials
// Secret present (with the SCRAM password key).
//
// When AuthProof is set (the behavioral-auth scenario), the verifier
// proves the credential live: it builds a SASL_PLAINTEXT/SCRAM-SHA-512
// client config from the generated Secret, produces run-unique markers to
// an ACL-allowed topic through the cluster's SCRAM listener, and consumes
// them back under an ACL-allowed group — authorization and authentication
// exercised on the wire.
type KafkaUserVerifier struct {
	Namespace   string
	Username    string
	ClusterName string
	AuthProof   bool
	// ScramPort is the cluster's scram-sha-512 listener port.
	ScramPort int
	// AclTopic / AclGroup are the topic and consumer group the user's
	// ACLs allow (derived from the manifest's authorization block).
	AclTopic string
	AclGroup string
}

func (v *KafkaUserVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] kafka user %q on cluster %q\n", v.Username, v.ClusterName)

	if err := kubectlWait(ctx, kubeconfig, "kafkauser", v.Username, v.Namespace,
		"condition=Ready", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "KafkaUser %q never reached the Ready condition (is it in the cluster's own namespace with the right cluster label?)", v.Username)
	}

	// The user operator generates the credentials Secret named after the
	// user (scram: password + sasl.jaas.config keys).
	if err := KubectlResourceExists(ctx, kubeconfig, "secret", v.Username, v.Namespace); err != nil {
		return errors.Wrap(err, "credentials secret not found")
	}

	if !v.AuthProof {
		return nil
	}
	return v.proveScramAuth(ctx, kubeconfig)
}

func (v *KafkaUserVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "kafkauser", v.Username, v.Namespace); err != nil {
		return err
	}
	// The user operator garbage-collects the credentials Secret with the
	// user — a lingering Secret would be a leaked credential.
	return KubectlResourceAbsent(ctx, kubeconfig, "secret", v.Username, v.Namespace)
}

// proveScramAuth produces and consumes as the user through the SCRAM
// listener, with the client config built from the generated Secret.
func (v *KafkaUserVerifier) proveScramAuth(ctx context.Context, kubeconfig string) error {
	if v.AclTopic == "" || v.AclGroup == "" {
		return errors.New("the behavioral-auth scenario needs a topic ACL and a group ACL to drive the proof")
	}

	// The generated sasl.jaas.config carries the user's SCRAM credential
	// verbatim — the exact line Kafka clients need.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "secret", v.Username, "-n", v.Namespace,
		"-o", "jsonpath={.data.sasl\\.jaas\\.config}").CombinedOutput()
	if err != nil {
		return errors.Errorf("failed to read sasl.jaas.config from the secret: %v: %s", err, string(out))
	}
	jaasConfig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(out)))
	if err != nil {
		return errors.Wrap(err, "failed to decode sasl.jaas.config")
	}

	clientProperties := fmt.Sprintf(
		"security.protocol=SASL_PLAINTEXT\nsasl.mechanism=SCRAM-SHA-512\nsasl.jaas.config=%s\n",
		strings.TrimSpace(string(jaasConfig)))
	encodedProperties := base64.StdEncoding.EncodeToString([]byte(clientProperties))

	pod, err := kafkaBrokerPod(ctx, kubeconfig, v.Namespace, v.ClusterName)
	if err != nil {
		return err
	}
	bootstrap := fmt.Sprintf("%s-kafka-bootstrap:%d", v.ClusterName, v.ScramPort)
	runID := fmt.Sprintf("auth-%d", time.Now().Unix())
	markers := []string{runID + "-0", runID + "-1", runID + "-2"}

	// Stage the client config in the pod, then produce AS THE USER.
	produce := fmt.Sprintf(
		"echo %s | base64 -d > /tmp/e2e-client.properties && printf '%s\\n' | %s/kafka-console-producer.sh --bootstrap-server %s --topic %s --request-required-acks all --producer.config /tmp/e2e-client.properties",
		encodedProperties, strings.Join(markers, "\\n"), kafkaBinDir, bootstrap, v.AclTopic)
	if out, err := v.userExec(ctx, kubeconfig, pod, produce); err != nil {
		return errors.Wrapf(err, "authenticated produce failed: %s", out)
	}

	// Consume the markers back under the ACL-allowed group.
	consume := fmt.Sprintf(
		"%s/kafka-console-consumer.sh --bootstrap-server %s --topic %s --group %s --from-beginning --max-messages %d --timeout-ms 120000 --consumer.config /tmp/e2e-client.properties",
		kafkaBinDir, bootstrap, v.AclTopic, v.AclGroup, len(markers))
	deadline := time.Now().Add(3 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		out, err := v.userExec(ctx, kubeconfig, pod, consume)
		if err == nil && containsAllMarkers(out, markers) {
			fmt.Printf("  [verify] SCRAM-AUTH: produced+consumed %d markers as %q through the scram listener\n",
				len(markers), v.Username)
			return nil
		}
		lastErr = errors.Errorf("authenticated consume attempt: err=%v out=%q", err, strings.TrimSpace(out))
		time.Sleep(10 * time.Second)
	}
	return lastErr
}

func (v *KafkaUserVerifier) userExec(ctx context.Context, kubeconfig, pod, command string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", pod, "-n", v.Namespace, "-c", "kafka", "--",
		"bash", "-c", command).CombinedOutput()
	return string(out), err
}

// kafkaBrokerPod finds one ready broker pod of the given cluster (the
// strimzi.io/broker-role label is stamped by the operator on broker-role
// nodes).
func kafkaBrokerPod(ctx context.Context, kubeconfig, namespace, clusterName string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "pods", "-n", namespace,
		"-l", fmt.Sprintf("strimzi.io/cluster=%s,strimzi.io/broker-role=true", clusterName),
		"-o", "jsonpath={.items[0].metadata.name}").CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return "", errors.Errorf("no broker pod found for cluster %q in %q: %v: %s",
			clusterName, namespace, err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// kafkaClusterRef extracts the target cluster name from a manifest's
// kafkaCluster value-or-ref field: the literal `value` when present, the
// referenced resource's `name` otherwise (a KubernetesKafka's cluster_name
// output IS its metadata.name, so the reference name is the cluster name).
func kafkaClusterRef(spec map[string]interface{}) string {
	if spec == nil {
		return ""
	}
	field, ok := spec["kafkaCluster"].(map[string]interface{})
	if !ok {
		// Scenario manifests are written in the proto's snake_case form.
		if field, ok = spec["kafka_cluster"].(map[string]interface{}); !ok {
			return ""
		}
	}
	if value, ok := field["value"].(string); ok && value != "" {
		return value
	}
	if valueFrom, ok := field["valueFrom"].(map[string]interface{}); ok {
		name, _ := valueFrom["name"].(string)
		return name
	}
	return ""
}

// kafkaUserAcl reads the first topic ACL's name and the first group ACL's
// name from a KubernetesKafkaUser manifest spec — the resources the
// behavioral-auth proof drives through.
func kafkaUserAcl(spec map[string]interface{}) (topic string, group string) {
	if spec == nil {
		return "", ""
	}
	authorization, ok := spec["authorization"].(map[string]interface{})
	if !ok {
		return "", ""
	}
	list, ok := authorization["acls"].([]interface{})
	if !ok {
		return "", ""
	}
	for _, entry := range list {
		acl, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := acl["resource"].(map[string]interface{})
		if !ok {
			continue
		}
		resourceType, _ := resource["type"].(string)
		name, _ := resource["name"].(string)
		if resourceType == "topic" && topic == "" {
			topic = name
		}
		if resourceType == "group" && group == "" {
			group = name
		}
	}
	return topic, group
}
