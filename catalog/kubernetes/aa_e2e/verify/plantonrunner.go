package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// PlantonRunnerVerifier checks a Planton runner install to the module's
// contract: the module-created token Secret present, the chart's
// Deployment rendered pinned to exactly ONE replica with the Recreate
// strategy (two live pods under one runner name would revoke each other's
// keys), and the pod template carrying the enrollment env wiring — the
// token via a secretKeyRef (never inline) and the registration name.
//
// DELIBERATELY no pod-readiness assertion: the module installs with no
// Helm wait because the runner's readiness contract is its control-plane
// work queue — the kind lanes run with an obviously-fake token whose join
// is refused, and the pod restarting is the DESIGNED behavior there, not
// a failure. Destroy asserts the workload and the token Secret are gone.
type PlantonRunnerVerifier struct {
	Namespace string
	Name      string
}

func (v *PlantonRunnerVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] planton-runner %q in namespace %q\n", v.Name, v.Namespace)

	// The module-created token Secret (`<name>-token`) — the chart reads
	// it via its existingSecret form; the token never rides rendered
	// chart values.
	tokenSecret := v.Name + "-token"
	if err := KubectlResourceExists(ctx, kubeconfig, "secret", tokenSecret, v.Namespace); err != nil {
		return errors.Wrap(err, "the token Secret not found")
	}

	// fullnameOverride pins the Deployment name to the resource name.
	if err := KubectlResourceExists(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "the runner Deployment was not rendered")
	}

	// The singleton law and the identity-safe rollout strategy, straight
	// off the live object.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", "jsonpath={.spec.replicas}/{.spec.strategy.type}").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading deployment replicas/strategy: %s", firstLines(string(out), 3))
	}
	if got := strings.TrimSpace(string(out)); got != "1/Recreate" {
		return errors.Errorf("the runner Deployment must run exactly one replica with the Recreate strategy (a rolling surge would revoke the live pod's key); got %q", got)
	}

	// The enrollment env contract on the pod template: the token arrives
	// through a secretKeyRef naming the module's Secret, and the runner
	// registers itself under the declared name.
	out, err = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", `jsonpath={.spec.template.spec.containers[0].env[?(@.name=="PLANTON_RUNNER_TOKEN")].valueFrom.secretKeyRef.name}`).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading the token env wiring: %s", firstLines(string(out), 3))
	}
	if got := strings.TrimSpace(string(out)); got != tokenSecret {
		return errors.Errorf("PLANTON_RUNNER_TOKEN must ride a secretKeyRef naming %q (never inline values); got %q", tokenSecret, got)
	}

	out, err = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", `jsonpath={.spec.template.spec.containers[0].env[?(@.name=="PLANTON_RUNNER_NAME")].value}`).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading the runner-name env: %s", firstLines(string(out), 3))
	}
	if got := strings.TrimSpace(string(out)); got == "" {
		return errors.New("PLANTON_RUNNER_NAME is missing from the pod template -- the runner would register under the chart default instead of the declared name")
	}

	fmt.Printf("  [verify] CONTRACT: token Secret %q wired via secretKeyRef; 1 replica, Recreate strategy\n", tokenSecret)
	return nil
}

func (v *PlantonRunnerVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return err
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "secret", v.Name+"-token", v.Namespace); err != nil {
		return err
	}
	fmt.Printf("  [verify] DESTROY: the runner Deployment and its token Secret are gone\n")
	return nil
}

// VerifyRuntimeFailureCause pins the fake-token pod's designed failure to the
// refused control-plane join and nothing else: the container's status must
// prove the image PULLED and the process RAN (a non-zero terminated state or
// a restart, never ImagePullBackOff/ErrImagePull), and the pod's own log must
// carry the runner's join-step error ("joining as runner ...") -- the line the
// binary emits only after it read the token from the Secret and REACHED the
// control plane. A "dialing control-plane" line instead means the failure is
// network, not the token, and the phase fails with that evidence.
//
// Polling is bounded: a fresh pod needs a restart or two before its status
// and logs attest the cause.
func (v *PlantonRunnerVerifier) VerifyRuntimeFailureCause(ctx context.Context, kubeconfig, cause string) error {
	if cause != "refused-join" {
		return errors.Errorf("unsupported runtime failure cause %q for the runner (supported: refused-join)", cause)
	}

	deadline := time.Now().Add(3 * time.Minute)
	var lastState string
	for {
		podName, stateErr := v.crashedRunnerPod(ctx, kubeconfig)
		if stateErr == nil && podName != "" {
			// The container ran and exited; now pin WHY from its own log.
			logOut, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
				"logs", podName, "-n", v.Namespace, "--tail", "50").CombinedOutput()
			logs := string(logOut)
			if !strings.Contains(logs, "joining as runner") {
				// The crash may predate log capture on this restart; try the
				// previous container instance before deciding.
				prevOut, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
					"logs", podName, "-n", v.Namespace, "--previous", "--tail", "50").CombinedOutput()
				logs = string(prevOut)
			}
			switch {
			case strings.Contains(logs, "joining as runner"):
				fmt.Printf("  [verify] CAUSE: image pulled, container ran, and the log carries the join-step error -- the ONLY failure is the refused join\n")
				return nil
			case strings.Contains(logs, "dialing control-plane"):
				return errors.Errorf("the runner failed DIALING the control plane (network), not joining with the token -- pod %s logs: %s",
					podName, firstLines(logs, 5))
			}
			lastState = fmt.Sprintf("pod %s crashed but its log does not yet carry the join line: %s", podName, firstLines(logs, 3))
		} else if stateErr != nil {
			lastState = stateErr.Error()
		}

		if time.Now().After(deadline) {
			return errors.Errorf("the runner pod never attested the refused join within the window; last state: %s", lastState)
		}
		time.Sleep(10 * time.Second)
	}
}

// crashedRunnerPod finds the runner Deployment's pod and returns its name once
// its container status proves the image pulled and the process ran-and-exited.
// An image-pull failure state fails IMMEDIATELY -- that is a registry problem
// wearing a crash costume, never the designed refused join.
func (v *PlantonRunnerVerifier) crashedRunnerPod(ctx context.Context, kubeconfig string) (string, error) {
	// The chart's Deployment selects its pods by the Helm instance label
	// (fullnameOverride pins names, not labels, so read the selector from
	// the live Deployment rather than assuming label conventions).
	selOut, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", "jsonpath={.spec.selector.matchLabels}").CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "reading the Deployment selector: %s", firstLines(string(selOut), 3))
	}
	var labels map[string]string
	if err := json.Unmarshal(selOut, &labels); err != nil {
		return "", errors.Wrapf(err, "parsing the Deployment selector %q", string(selOut))
	}
	var selector []string
	for k, val := range labels {
		selector = append(selector, k+"="+val)
	}
	sort.Strings(selector)

	podsOut, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "pods", "-n", v.Namespace, "-l", strings.Join(selector, ","),
		"-o", "json").CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "listing the runner pods: %s", firstLines(string(podsOut), 3))
	}

	var podList struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Status struct {
				ContainerStatuses []struct {
					ImageID      string `json:"imageID"`
					RestartCount int    `json:"restartCount"`
					State        struct {
						Waiting *struct {
							Reason string `json:"reason"`
						} `json:"waiting"`
						Terminated *struct {
							ExitCode int `json:"exitCode"`
						} `json:"terminated"`
					} `json:"state"`
					LastState struct {
						Terminated *struct {
							ExitCode int `json:"exitCode"`
						} `json:"terminated"`
					} `json:"lastState"`
				} `json:"containerStatuses"`
			} `json:"status"`
		} `json:"items"`
	}
	if err := json.Unmarshal(podsOut, &podList); err != nil {
		return "", errors.Wrap(err, "parsing the runner pod list")
	}
	if len(podList.Items) == 0 {
		return "", errors.New("no runner pod scheduled yet")
	}

	pod := podList.Items[0]
	if len(pod.Status.ContainerStatuses) == 0 {
		return "", errors.Errorf("pod %s has no container status yet", pod.Metadata.Name)
	}
	cs := pod.Status.ContainerStatuses[0]

	if cs.State.Waiting != nil {
		switch cs.State.Waiting.Reason {
		case "ErrImagePull", "ImagePullBackOff":
			return "", errors.Errorf("pod %s cannot PULL the image (%s) -- a registry problem, not the designed refused join",
				pod.Metadata.Name, cs.State.Waiting.Reason)
		}
	}
	if cs.ImageID == "" {
		return "", errors.Errorf("pod %s has not pulled the image yet", pod.Metadata.Name)
	}

	ranAndExited := cs.RestartCount >= 1 ||
		(cs.State.Terminated != nil && cs.State.Terminated.ExitCode != 0) ||
		(cs.LastState.Terminated != nil && cs.LastState.Terminated.ExitCode != 0)
	if !ranAndExited {
		return "", errors.Errorf("pod %s has not exited yet (restarts=%d)", pod.Metadata.Name, cs.RestartCount)
	}
	return pod.Metadata.Name, nil
}
