package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// sparkOperatorCrds are the two CRDs the chart ships from its crds/
// directory: Helm installs them once and KEEPS them on uninstall (the
// crds/-directory posture), so destroy asserts their survival, not their
// deletion.
var sparkOperatorCrds = []string{
	"sparkapplications.spark.apache.org",
	"sparkclusters.spark.apache.org",
}

// sparkProofLabel marks every CR this verifier creates so leftovers are
// findable (and assertable-absent) across all namespaces after destroy.
const sparkProofLabel = "planton.ai/e2e-proof"

// SparkOperatorVerifier checks an Apache Spark Kubernetes Operator
// install to the kind's actual definition of working: not "the pods are
// up" but A SPARK JOB RUNS TO COMPLETION. The verifier submits a
// verifier-owned SparkApplication (upstream's own SparkPi e2e shape) and
// waits for the operator to drive it to the Succeeded state — the whole
// submit/schedule/run/collect loop proven live, including the workload
// service account the module plants. On the fenced posture it also
// proves the watch fence: a SparkApplication OUTSIDE the workload
// namespaces is never reconciled (status stays empty), the same
// assertion upstream's watched-namespaces e2e makes.
type SparkOperatorVerifier struct {
	// Namespace is the release namespace; the module pins the chart
	// fullname to the resource name, so the Deployment is
	// deployment/<Name> here.
	Namespace string
	Name      string
	// WorkloadNamespaces is the watch fence from the spec. Non-empty
	// means the operator watches ONLY these namespaces (the chart
	// creates them and plants the workload identity there); empty means
	// cluster-wide watch with the workload identity in the release
	// namespace.
	WorkloadNamespaces []string
	// WorkloadServiceAccount is the driver service account the module
	// plants for jobs ("spark" upstream default when unset).
	WorkloadServiceAccount string
}

func (v *SparkOperatorVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] spark-operator %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 5*time.Minute); err != nil {
		return errors.Wrap(err, "the operator deployment never rolled out")
	}
	if err := waitForCrdsEstablished(ctx, kubeconfig, sparkOperatorCrds); err != nil {
		return err
	}
	fmt.Printf("  [verify] operator rolled out and both spark.apache.org CRDs Established\n")

	if err := v.proveJobRun(ctx, kubeconfig); err != nil {
		return err
	}
	if len(v.WorkloadNamespaces) > 0 {
		return v.proveWatchFence(ctx, kubeconfig)
	}
	return nil
}

// proveJobRun is THE JOB PROOF: a verifier-owned SparkApplication driven
// to the Succeeded state by the operator. The CR body is upstream's own
// e2e job (spark-kubernetes-operator tests/e2e/watched-namespaces/
// spark-example.yaml): SparkPi from the in-image example jar at
// local:///opt/spark/examples/jars/spark-examples.jar, one executor,
// image apache/spark:{{SPARK_VERSION}}-scala — the operator itself
// substitutes {{SPARK_VERSION}} from runtimeVersions.sparkVersion
// (SparkAppSubmissionWorker.java:139), pinned to 4.2.0 as upstream's
// examples and e2e do at this operator version.
//
// On the fenced posture the job runs INSIDE the first workload
// namespace — the same run simultaneously proves the fence admits work
// and the per-namespace workload RBAC/service-account actually submits.
func (v *SparkOperatorVerifier) proveJobRun(ctx context.Context, kubeconfig string) error {
	jobNamespace := v.Namespace
	if len(v.WorkloadNamespaces) > 0 {
		jobNamespace = v.WorkloadNamespaces[0]
	}
	serviceAccount := v.WorkloadServiceAccount
	if serviceAccount == "" {
		// The upstream contract: the workload identity is "spark"
		// unless the spec overrides it.
		serviceAccount = "spark"
	}
	proofName := v.Name + "-e2e-proof"

	// ALWAYS sweep the proof CR — a leaked SparkApplication would hold
	// driver resources and poison the destroy-side leftover assertion.
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		_ = kubectlDeleteResource(cleanupCtx, kubeconfig, "sparkapplications.spark.apache.org", proofName, jobNamespace)
		_ = KubectlResourceAbsent(cleanupCtx, kubeconfig, "sparkapplications.spark.apache.org", proofName, jobNamespace)
	}()

	manifest := v.proofSparkApplication(proofName, jobNamespace, serviceAccount)
	if out, err := applyManifestString(ctx, kubeconfig, manifest); err != nil {
		return errors.Errorf("applying the proof SparkApplication: %v: %s", err, firstLines(out, 3))
	}

	// The operator's state machine (spark-operator-api
	// ApplicationStateSummary.java): Submitted → DriverRequested →
	// DriverStarted → DriverReady → RunningHealthy → Succeeded, then the
	// terminal ResourceReleased once pods/services are garbage-collected.
	// The printer-column jsonpath is
	// .status.currentState.currentStateSummary (SparkApplication.java).
	// Poll on a 10-minute budget (image pull + driver + executor + Pi).
	deadline := time.Now().Add(10 * time.Minute)
	lastState := ""
	for time.Now().Before(deadline) {
		out, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "sparkapplications.spark.apache.org", proofName, "-n", jobNamespace,
			"-o", "jsonpath={.status.currentState.currentStateSummary}").Output()
		state := strings.TrimSpace(string(out))
		if state != "" && state != lastState {
			fmt.Printf("  [verify] job proof %s/%s: state %s\n", jobNamespace, proofName, state)
			lastState = state
		}
		switch state {
		case "Succeeded":
			fmt.Printf("  [verify] THE JOB PROOF: SparkApplication %s reached state %s — the operator ran a Spark job to completion\n", proofName, state)
			return nil
		case "ResourceReleased", "TerminatedWithoutReleaseResources":
			// Terminal after EITHER success or failure — the poll can
			// skip past the short-lived Succeeded observation, so read
			// the state history the operator keeps and require Succeeded
			// in it (upstream's own e2e asserts this exact history).
			history, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
				"get", "sparkapplications.spark.apache.org", proofName, "-n", jobNamespace,
				"-o", "jsonpath={.status.stateTransitionHistory}").Output()
			if strings.Contains(string(history), "Succeeded") {
				fmt.Printf("  [verify] THE JOB PROOF: SparkApplication %s terminated (%s) with Succeeded in its state history — the operator ran a Spark job to completion\n", proofName, state)
				return nil
			}
			return errors.Errorf("the proof SparkApplication terminated (%s) WITHOUT a Succeeded state (history: %s)", state, firstLines(string(history), 3))
		case "Failed", "SchedulingFailure", "DriverStartTimedOut", "ExecutorsStartTimedOut", "DriverReadyTimedOut", "DriverEvicted":
			return errors.Errorf("the proof SparkApplication reached failure state %s — the operator cannot run a job to completion", state)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Second):
		}
	}
	return errors.Errorf("the proof SparkApplication never reached Succeeded within 10 minutes (last state %q)", lastState)
}

// proveWatchFence is THE FENCE PROOF: a SparkApplication in the RELEASE
// namespace — OUTSIDE the workload fence — must be IGNORED by the
// operator: its status stays entirely empty (never reconciled). This is
// the same assertion upstream's watched-namespaces e2e makes (apply
// outside the watch set, sleep 60s, status is null); a status appearing
// here would mean the fence leaks and the multi-tenant posture is a lie.
func (v *SparkOperatorVerifier) proveWatchFence(ctx context.Context, kubeconfig string) error {
	fenceName := v.Name + "-e2e-fence"

	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		_ = kubectlDeleteResource(cleanupCtx, kubeconfig, "sparkapplications.spark.apache.org", fenceName, v.Namespace)
		_ = KubectlResourceAbsent(cleanupCtx, kubeconfig, "sparkapplications.spark.apache.org", fenceName, v.Namespace)
	}()

	manifest := v.proofSparkApplication(fenceName, v.Namespace, "spark")
	if out, err := applyManifestString(ctx, kubeconfig, manifest); err != nil {
		return errors.Errorf("applying the fence-probe SparkApplication: %v: %s", err, firstLines(out, 3))
	}

	// 60 seconds is several reconcile intervals — long enough that a
	// leaking watch would have stamped SOMETHING into status.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(60 * time.Second):
	}

	status, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "sparkapplications.spark.apache.org", fenceName, "-n", v.Namespace,
		"-o", "jsonpath={.status}").Output()
	if err != nil {
		return errors.Wrap(err, "reading the fence-probe status")
	}
	if strings.TrimSpace(string(status)) != "" {
		return errors.Errorf("the operator reconciled a SparkApplication OUTSIDE the workload fence (status: %s) — the watch fence leaks", firstLines(string(status), 2))
	}
	fmt.Printf("  [verify] THE FENCE PROOF: SparkApplication outside the workload namespaces stayed unreconciled (empty status) after 60s — the watch fence holds\n")
	return nil
}

// proofSparkApplication renders the verifier-owned CR body (upstream's
// SparkPi e2e shape, GA spark.apache.org/v1).
func (v *SparkOperatorVerifier) proofSparkApplication(name, namespace, serviceAccount string) string {
	return fmt.Sprintf(`apiVersion: spark.apache.org/v1
kind: SparkApplication
metadata:
  name: %s
  namespace: %s
  labels:
    %s: "true"
spec:
  mainClass: "org.apache.spark.examples.SparkPi"
  jars: "local:///opt/spark/examples/jars/spark-examples.jar"
  sparkConf:
    spark.executor.instances: "1"
    spark.kubernetes.container.image: "apache/spark:{{SPARK_VERSION}}-scala"
    spark.kubernetes.authenticate.driver.serviceAccountName: "%s"
  runtimeVersions:
    sparkVersion: "4.2.0"
`, name, namespace, sparkProofLabel, serviceAccount)
}

// VerifyAbsent asserts the destroy posture: the operator Deployment gone,
// both CRDs SURVIVING (the crds/-directory keep — a designed outcome,
// asserted, not tolerated), the WORKLOAD resources surviving (the
// chart annotates every workload resource `helm.sh/resource-policy:
// keep` at the pin — protective by design: a day-2 fence-list edit
// must never cascade-delete a namespace of running Spark jobs, and
// operator uninstall must not abort workloads), and zero verifier
// residue anywhere. After asserting the keeps, the lane SWEEPS them
// (confirmed) — shared-cluster hygiene, the prove-then-clean shape.
func (v *SparkOperatorVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return err
	}
	for _, crd := range sparkOperatorCrds {
		if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
			return errors.Wrapf(err, "CRD %q was DELETED on destroy — the crds/-directory keep posture regressed", crd)
		}
	}
	// The CRDs survive, so the API is still queryable: any proof CR the
	// verifier leaked would linger forever with no operator to GC it.
	leftovers, err := listCustomResourcesAllNamespaces(ctx, kubeconfig, "sparkapplications.spark.apache.org", sparkProofLabel+"=true")
	if err != nil {
		return errors.Wrap(err, "listing leftover proof SparkApplications")
	}
	if leftovers != "" {
		return errors.Errorf("proof SparkApplications survived the verifier sweep: %s", leftovers)
	}

	// THE WORKLOAD KEEP, per arm (the module's one-decision design):
	// the cluster-wide arm keeps the workload ClusterRole (its binding
	// and SA live in the module-deleted release namespace); the fenced
	// arm keeps each workload namespace with its `spark` SA and
	// Role/Binding inside. Assert the designed keep, then sweep for
	// lane hygiene.
	sweep := func(kind, name, namespace string) error {
		if err := kubectlDeleteResource(ctx, kubeconfig, kind, name, namespace); err != nil {
			return err
		}
		return KubectlResourceAbsent(ctx, kubeconfig, kind, name, namespace)
	}
	if len(v.WorkloadNamespaces) == 0 {
		workloadClusterRole := v.Name + "-workload-clusterrole"
		if err := KubectlResourceExists(ctx, kubeconfig, "clusterrole", workloadClusterRole, ""); err != nil {
			return errors.Wrapf(err, "the workload ClusterRole %q was DELETED on destroy — the chart's resource-policy keep regressed", workloadClusterRole)
		}
		if err := sweep("clusterrole", workloadClusterRole, ""); err != nil {
			return errors.Wrap(err, "sweeping the kept workload ClusterRole")
		}
	}
	for _, workloadNs := range v.WorkloadNamespaces {
		if err := KubectlResourceExists(ctx, kubeconfig, "namespace", workloadNs, ""); err != nil {
			return errors.Wrapf(err, "workload namespace %q was DELETED on destroy — the chart's resource-policy keep regressed", workloadNs)
		}
		if err := sweep("namespace", workloadNs, ""); err != nil {
			return errors.Wrapf(err, "sweeping the kept workload namespace %q", workloadNs)
		}
	}
	fmt.Printf("  [verify] DESTROY: operator deployment gone, both CRDs RETAINED by design, workload keep posture ASSERTED (and swept for lane hygiene), no proof CRs left anywhere\n")
	return nil
}

// waitForCrdsEstablished waits for each CRD's Established condition —
// existence alone admits a CRD the API server has not yet served.
func waitForCrdsEstablished(ctx context.Context, kubeconfig string, crds []string) error {
	for _, crd := range crds {
		if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"wait", "--for=condition=Established", "crd/"+crd, "--timeout=120s").CombinedOutput(); err != nil {
			return errors.Wrapf(err, "CRD %q never became Established: %s", crd, firstLines(string(out), 3))
		}
	}
	return nil
}

// listCustomResourcesAllNamespaces lists CRs of one kind across all
// namespaces (optionally label-filtered), returning the trimmed
// `-o name` output — empty string means none exist.
func listCustomResourcesAllNamespaces(ctx context.Context, kubeconfig, kind, labelSelector string) (string, error) {
	args := []string{"--kubeconfig", kubeconfig, "get", kind, "-A", "-o", "name"}
	if labelSelector != "" {
		args = append(args, "-l", labelSelector)
	}
	out, err := exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
	if err != nil {
		return "", errors.Errorf("kubectl get %s -A: %v: %s", kind, err, firstLines(string(out), 3))
	}
	return strings.TrimSpace(string(out)), nil
}

// ensureNamespace creates a namespace idempotently (apply tolerates
// AlreadyExists) for proofs that must run in a verifier-chosen namespace.
func ensureNamespace(ctx context.Context, kubeconfig, name string) error {
	manifest := fmt.Sprintf("apiVersion: v1\nkind: Namespace\nmetadata:\n  name: %s\n", name)
	if out, err := applyManifestString(ctx, kubeconfig, manifest); err != nil {
		return errors.Errorf("creating namespace %s: %v: %s", name, err, firstLines(out, 2))
	}
	return nil
}

// sparkWorkloadNamespaces reads the watch fence from spec.workload
// (both manifest key forms tolerated).
func sparkWorkloadNamespaces(spec map[string]interface{}) []string {
	if workload, ok := spec["workload"].(map[string]interface{}); ok {
		return specStringList(workload, "namespaces")
	}
	return nil
}

// sparkWorkloadServiceAccount reads the workload service-account name,
// defaulting to the upstream "spark" contract the spec documents.
func sparkWorkloadServiceAccount(spec map[string]interface{}) string {
	if workload, ok := spec["workload"].(map[string]interface{}); ok {
		for _, key := range []string{"service_account", "serviceAccount"} {
			if name, ok := workload[key].(string); ok && name != "" {
				return name
			}
		}
	}
	return "spark"
}
