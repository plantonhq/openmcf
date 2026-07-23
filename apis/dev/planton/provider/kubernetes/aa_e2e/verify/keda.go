package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// KedaInstallVerifier checks a KEDA installation to the point the
// autoscaling engine is actually serving: operator and metrics API server
// Available, and the cluster-wide v1beta1.external.metrics.k8s.io
// APIService Available (the aggregation layer reaches KEDA — the exact
// dependency every ScaledObject-driven HPA has).
//
// When Behavioral is set (the behavioral-scaling scenario), the verifier
// additionally proves the ENGINE SCALES: it applies a deterministic
// cron-trigger ScaledObject against the scenario's target-Deployment
// fixture, waits for KEDA to materialize the HPA and drive the Deployment
// above its baseline replica count, then deletes the ScaledObject it
// created. The ScaledObject is verifier-owned on purpose: scenario fixtures
// deploy BEFORE the component under test, so a fixture-borne ScaledObject
// would land before the KEDA CRDs exist — the same reason the CR cannot be
// a registry-driven chain entry here.
type KedaInstallVerifier struct {
	Namespace string
	// Behavioral switches on the live scale proof (see above).
	Behavioral bool
	// BehavioralSqs proves the CLOUD pod-identity hop (the aws-sqs-irsa
	// scenario on the aws-eks profile): a verifier-owned driver chain — a
	// scale target, a pod-identity TriggerAuthentication, and an
	// aws-sqs-queue ScaledObject against the batch queue — plus real
	// messages enqueued via the AWS CLI must produce a real 0→N scale-up,
	// meaning the operator read queue depth out of SQS with the assumed
	// IRSA role, keyless. The whole chain is verifier-owned because every
	// CR is served by CRDs this component installs.
	BehavioralSqs bool
}

// kedaSqs* identify the verifier-owned SQS driver chain. Queue coordinates
// come from the batch bootstrap's exported environment.
const (
	kedaSqsNamespace  = "e2e-keda-sqs"
	kedaSqsTargetName = "e2e-keda-sqs-target"
)

// kedaScaleTargetName/Namespace identify the scale-target Deployment fixture
// the behavioral scenario deploys (fixture-scale-target.yaml under the
// kuberneteskeda e2e folder) — fixed so the verifier and the fixture cannot
// drift apart silently. The fixture lives in ITS OWN namespace (the
// scenario's workload namespace), not KEDA's install namespace: ScaledObjects
// are namespace-local to the workloads they scale, so the proof exercises
// exactly the cross-namespace watch posture a real cluster runs.
const (
	kedaScaleTargetName      = "e2e-keda-scale-target"
	kedaScaleTargetNamespace = "e2e-keda-behavioral"
)

func (v *KedaInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] keda installation in namespace %q\n", v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for keda", v.Namespace)
	}

	// The chart names its components fixed (keda-operator /
	// keda-operator-metrics-apiserver) regardless of release name — one
	// installation per cluster.
	for _, deploy := range []string{"keda-operator", "keda-operator-metrics-apiserver"} {
		if err := kubectlWait(ctx, kubeconfig, "deployment", deploy, v.Namespace,
			"condition=Available", 3*time.Minute); err != nil {
			return errors.Wrapf(err, "keda deployment %q not available", deploy)
		}
	}

	if err := kubectlWait(ctx, kubeconfig, "apiservice", "v1beta1.external.metrics.k8s.io", "",
		"condition=Available", 2*time.Minute); err != nil {
		return errors.Wrap(err, "v1beta1.external.metrics.k8s.io APIService not available")
	}

	if v.BehavioralSqs {
		return v.proveSqsScaling(ctx, kubeconfig)
	}
	if !v.Behavioral {
		return nil
	}
	return v.proveScaling(ctx, kubeconfig)
}

// proveSqsScaling drives the cloud-identity scale loop: target + pod-identity
// TriggerAuthentication + SQS ScaledObject, real messages in, real replicas
// out, queue purged so scale-in state never leaks into the destroy phase.
func (v *KedaInstallVerifier) proveSqsScaling(ctx context.Context, kubeconfig string) error {
	queueURL := os.Getenv("PLANTON_E2E_SQS_QUEUE_URL")
	awsRegion := os.Getenv("PLANTON_E2E_AWS_REGION")
	if queueURL == "" || awsRegion == "" {
		return errors.New("PLANTON_E2E_SQS_QUEUE_URL / PLANTON_E2E_AWS_REGION unset — the batch bootstrap exports them")
	}
	fmt.Printf("  [verify] behavioral SQS scaling: queue depth must drive a real scale-up (keyless)\n")

	if err := v.kubectl(ctx, kubeconfig, "create", "namespace", kedaSqsNamespace); err != nil {
		return errors.Wrap(err, "failed to create sqs driver namespace")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "namespace", kedaSqsNamespace,
			"--ignore-not-found", "--wait=false")
	}()

	driver := fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
spec:
  replicas: 0
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
        - name: worker
          image: registry.k8s.io/pause:3.10
          resources:
            requests:
              cpu: 25m
              memory: 16Mi
---
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: e2e-sqs-podidentity
  namespace: %s
spec:
  podIdentity:
    provider: aws
---
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: %s
  namespace: %s
spec:
  scaleTargetRef:
    name: %s
  minReplicaCount: 0
  maxReplicaCount: 2
  pollingInterval: 10
  cooldownPeriod: 60
  triggers:
    - type: aws-sqs-queue
      authenticationRef:
        name: e2e-sqs-podidentity
      metadata:
        queueURL: %s
        awsRegion: %s
        queueLength: "2"
`, kedaSqsTargetName, kedaSqsNamespace, kedaSqsTargetName, kedaSqsTargetName,
		kedaSqsNamespace, kedaSqsTargetName, kedaSqsNamespace, kedaSqsTargetName,
		queueURL, awsRegion)
	driverFile, err := writeTempManifest(driver)
	if err != nil {
		return err
	}
	defer os.Remove(driverFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", driverFile); err != nil {
		return errors.Wrap(err, "failed to apply SQS driver chain")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "scaledobject", kedaSqsTargetName,
			"-n", kedaSqsNamespace, "--ignore-not-found")
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "triggerauthentication", "e2e-sqs-podidentity",
			"-n", kedaSqsNamespace, "--ignore-not-found")
	}()

	// Real messages: enough to demand both replicas at queueLength=2.
	for i := 0; i < 5; i++ {
		if out, err := exec.CommandContext(ctx, "aws", "sqs", "send-message",
			"--queue-url", queueURL, "--message-body", fmt.Sprintf("e2e-proof-%d", i)).CombinedOutput(); err != nil {
			return errors.Errorf("failed to enqueue proof message: %v: %s", err, string(out))
		}
	}
	// Purge whatever happens: queue depth must never leak into later runs.
	defer func() {
		_ = exec.Command("aws", "sqs", "purge-queue", "--queue-url", queueURL).Run()
	}()

	deadline := time.Now().Add(4 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "deployment", kedaSqsTargetName, "-n", kedaSqsNamespace,
			"-o", "jsonpath={.status.readyReplicas}").CombinedOutput()
		replicas := strings.TrimSpace(string(out))
		if err == nil && (replicas == "2" || replicas == "1") {
			fmt.Printf("  [verify] target scaled 0→%s off real SQS depth — the keyless cloud hop is proven\n", replicas)
			return nil
		}
		last = fmt.Sprintf("readyReplicas=%q err=%v", replicas, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("KEDA never scaled the target off the SQS queue (last: %s)", last)
}

func (v *KedaInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	for _, deploy := range []string{"keda-operator", "keda-operator-metrics-apiserver"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", deploy, v.Namespace); err != nil {
			return err
		}
	}
	// The APIService is chart-owned and must go with the release — a
	// dangling external-metrics backend breaks every future HPA that
	// consults it. The CRDs intentionally SURVIVE uninstall by default
	// (keep_on_uninstall), so their presence is not asserted either way.
	return KubectlResourceAbsent(ctx, kubeconfig, "apiservice", "v1beta1.external.metrics.k8s.io", "")
}

// proveScaling runs the enforce-and-release cycle for autoscaling: apply a
// cron ScaledObject whose window is effectively always-active (Jan 1 through
// Dec 31 — deterministic on any wall clock, unlike a same-day window that
// would flake at midnight boundaries), assert KEDA scales the target
// Deployment from 1 to 2 replicas, then delete the ScaledObject so the
// component's own destroy phase finds the cluster exactly as its fixtures
// left it.
func (v *KedaInstallVerifier) proveScaling(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] behavioral scaling: cron ScaledObject must scale %q above baseline\n", kedaScaleTargetName)

	// The target fixture (a plain Deployment) deploys with the scenario's
	// manifest-path fixtures — gate on it before asserting scaling.
	if err := kubectlWait(ctx, kubeconfig, "deployment", kedaScaleTargetName, kedaScaleTargetNamespace,
		"condition=Available", 2*time.Minute); err != nil {
		return errors.Wrapf(err, "scale-target fixture %q not available", kedaScaleTargetName)
	}

	scaledObject := fmt.Sprintf(`apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: %s
  namespace: %s
spec:
  scaleTargetRef:
    name: %s
  minReplicaCount: 1
  maxReplicaCount: 2
  triggers:
    - type: cron
      metadata:
        timezone: Etc/UTC
        start: 0 0 1 1 *
        end: 59 23 31 12 *
        desiredReplicas: "2"
`, kedaScaleTargetName, kedaScaleTargetNamespace, kedaScaleTargetName)

	tmpFile, err := writeTempManifest(scaledObject)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	applyCmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "apply", "-f", tmpFile)
	if out, err := applyCmd.CombinedOutput(); err != nil {
		return errors.Errorf("failed to apply ScaledObject: %v: %s", err, string(out))
	}
	// The ScaledObject is verifier-owned: delete it whatever happens below,
	// so a failed assertion cannot leak a CR that would block later lanes.
	defer func() {
		deleteCmd := exec.Command("kubectl", "--kubeconfig", kubeconfig,
			"delete", "scaledobject", kedaScaleTargetName, "-n", kedaScaleTargetNamespace, "--ignore-not-found")
		_ = deleteCmd.Run()
	}()

	// KEDA reconciles the ScaledObject into an HPA and the cron trigger is
	// active immediately — poll the target's replica count above baseline.
	deadline := time.Now().Add(3 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "deployment", kedaScaleTargetName, "-n", kedaScaleTargetNamespace,
			"-o", "jsonpath={.status.readyReplicas}")
		out, err := cmd.CombinedOutput()
		replicas := strings.TrimSpace(string(out))
		if err == nil && replicas == "2" {
			fmt.Printf("  [verify] target scaled to 2 ready replicas — KEDA drove a real scale-up\n")
			return nil
		}
		last = fmt.Sprintf("readyReplicas=%q err=%v", replicas, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("KEDA never scaled the target above baseline (last: %s)", last)
}

func (v *KedaInstallVerifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
}

// writeTempManifest writes a manifest document to a temp file for kubectl
// apply -f (KEDA's CR is verifier-owned, so it never lives in the repo as a
// deployable fixture — see the KedaInstallVerifier doc comment).
func writeTempManifest(content string) (string, error) {
	dir, err := os.MkdirTemp("", "planton-e2e-verify-*")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir for verifier manifest")
	}
	path := filepath.Join(dir, "manifest.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return "", errors.Wrap(err, "failed to write verifier manifest")
	}
	return path, nil
}
