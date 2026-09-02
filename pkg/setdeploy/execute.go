package setdeploy

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/manifestgraph"
	"github.com/plantonhq/planton/pkg/outputs"
	"github.com/plantonhq/planton/pkg/protobufyaml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NodeDeployer hands one preflighted node to its IaC engine and returns the
// captured outputs. It is an interface so the execution loop's own behavior —
// ordering, resolution, output feeding, failure honesty — is testable without
// an engine, and so each CLI wires the engines its build carries.
type NodeDeployer interface {
	Deploy(node NodePlan, manifestPath string) (*outputs.CaptureResult, error)
}

// Events is the execution loop's voice — the library never prints. Every
// callback carries the node's plan so renderers can speak in the user's own
// identifiers.
type Events interface {
	// NodeStarting fires before a node's handoff; position is 1-based within
	// the deploy order.
	NodeStarting(position, total int, node NodePlan)
	// NodeSucceeded fires after a node's apply and output capture.
	NodeSucceeded(node NodePlan, captured *outputs.CaptureResult)
	// NodeWarning fires for non-fatal facts the user must see (a sensitive
	// output entering a downstream manifest).
	NodeWarning(node NodePlan, message string)
	// NodeFailed fires once, for the node that stopped the run.
	NodeFailed(node NodePlan, err error)
}

// NodeStatus is one node's outcome in a Result.
type NodeStatus string

const (
	NodeStatusSucceeded    NodeStatus = "succeeded"
	NodeStatusFailed       NodeStatus = "failed"
	NodeStatusNeverStarted NodeStatus = "never-started"
)

// Result is the execution loop's honest summary: which nodes succeeded (their
// outputs captured), which one failed (the engine's error verbatim), and
// which never started. State backends hold everything either way — re-running
// the same command re-applies completed nodes as no-ops and continues.
type Result struct {
	// Statuses is indexed like Plan.Set.Nodes.
	Statuses []NodeStatus
	// Outputs holds each succeeded node's captured outputs by identity.
	Outputs map[manifestgraph.Identity]*outputs.CaptureResult
	// FailedErr is the failing node's engine error, verbatim. Nil when every
	// node succeeded.
	FailedErr error
}

// Succeeded reports whether every node deployed.
func (r *Result) Succeeded() bool {
	return r.FailedErr == nil
}

// Execute walks the plan's deploy order sequentially: resolve the node's
// references from the outputs captured so far, hand off, capture, fold, next.
// It refuses to run a refused plan — the wall's verdict is not advisory.
//
// On a node failure the loop STOPS (nodes after the failure depend on a world
// the failure left unfinished, directly or through authored order) and the
// Result names every node's status. No rollback is attempted: IaC state is
// the recovery mechanism, and the honest re-run story is running the same
// command again.
func Execute(plan *Plan, deployer NodeDeployer, events Events) (*Result, error) {
	if plan.Report.Refused() {
		return nil, errors.Errorf("refusing to execute: the preflight report carries %d refusal(s)", plan.Report.RefusalCount())
	}
	if plan.Order == nil {
		return nil, errors.New("refusing to execute: the plan has no deploy order")
	}

	result := &Result{
		Statuses: make([]NodeStatus, len(plan.Set.Nodes)),
		Outputs:  map[manifestgraph.Identity]*outputs.CaptureResult{},
	}
	for i := range result.Statuses {
		result.Statuses[i] = NodeStatusNeverStarted
	}

	lookup := func(id manifestgraph.Identity) (map[string]string, bool) {
		captured, ok := result.Outputs[id]
		if !ok || captured == nil {
			return nil, false
		}
		return captured.Flat, true
	}

	for position, idx := range plan.Order {
		node := plan.Nodes[idx]
		msg := plan.Set.Nodes[idx].Msg
		events.NodeStarting(position+1, len(plan.Order), node)

		// Warn BEFORE the value leaves our hands: a sensitive output resolved
		// into a manifest becomes plain config the engine may echo in its
		// diff — in CI, that is the log. The composition still deploys; the
		// user learns exactly which field to move to a provider-native secret
		// reference.
		warnSensitiveResolutions(msg, node, result.Outputs, events)

		_, findings := manifestgraph.ResolveRefs(msg, node.Identity.Env, lookup)
		if err := resolutionFailure(findings); err != nil {
			result.Statuses[idx] = NodeStatusFailed
			result.FailedErr = err
			events.NodeFailed(node, err)
			return result, nil
		}
		manifestPath, err := writeResolvedManifest(msg, node)
		if err != nil {
			result.Statuses[idx] = NodeStatusFailed
			result.FailedErr = err
			events.NodeFailed(node, err)
			return result, nil
		}

		captured, err := deployer.Deploy(node, manifestPath)
		os.Remove(manifestPath)
		if err != nil {
			result.Statuses[idx] = NodeStatusFailed
			result.FailedErr = err
			events.NodeFailed(node, err)
			return result, nil
		}

		result.Statuses[idx] = NodeStatusSucceeded
		if captured != nil {
			result.Outputs[node.Identity] = captured
		}
		events.NodeSucceeded(node, captured)
	}

	return result, nil
}

// resolutionFailure turns resolution findings into the node's failure. After
// a passing wall every reference target is in-set and topologically earlier,
// so a miss here means the producer deployed without exporting what the
// composition names — the message says exactly that.
func resolutionFailure(findings []manifestgraph.Finding) error {
	if len(findings) == 0 {
		return nil
	}
	lines := make([]string, 0, len(findings))
	for _, f := range findings {
		lines = append(lines, f.Message)
	}
	return errors.Errorf("reference resolution failed:\n  %s", strings.Join(lines, "\n  "))
}

// warnSensitiveResolutions inspects the node's references against the outputs
// captured so far and warns for each one about to inject a SENSITIVE output
// value into the manifest. Sensitivity lives on the CaptureResult, so the
// check reads the run's outputs map directly, with the same target and key
// derivation ResolveRefs applies.
func warnSensitiveResolutions(msg proto.Message, node NodePlan, captured map[manifestgraph.Identity]*outputs.CaptureResult, events Events) {
	for _, use := range manifestgraph.CollectRefUses(msg) {
		target, problems := manifestgraph.CheckRef(use)
		if len(problems) > 0 {
			continue
		}
		id := target.Identity(node.Identity.Env)
		producer, ok := captured[id]
		if !ok || producer == nil {
			continue
		}
		key := strings.TrimPrefix(target.FieldPath, "status.outputs.")
		if producer.IsSensitive(key) {
			events.NodeWarning(node, fmt.Sprintf(
				"%s resolves from %s's sensitive output %q — the value enters this manifest and may appear in IaC diff output; prefer a provider-native secret reference for secret material",
				use.FieldPath, id, key))
		}
	}
}

// writeResolvedManifest serializes the resolved message to a 0600 temp file
// for the engine handoff. Resolved manifests can carry secret material (that
// is what resolution does), so they are never world-readable and always
// removed after the handoff returns.
func writeResolvedManifest(msg proto.Message, node NodePlan) (string, error) {
	jsonBytes, err := protojson.Marshal(msg)
	if err != nil {
		return "", errors.Wrapf(err, "failed to serialize resolved manifest for %s", node.Identity)
	}
	yamlBytes, err := protobufyaml.JSONToYAML(jsonBytes)
	if err != nil {
		return "", errors.Wrapf(err, "failed to render resolved manifest for %s", node.Identity)
	}
	f, err := os.CreateTemp("", "planton-set-node-*.yaml")
	if err != nil {
		return "", errors.Wrap(err, "failed to create resolved manifest file")
	}
	defer f.Close()
	if err := f.Chmod(0o600); err != nil {
		return "", errors.Wrap(err, "failed to restrict resolved manifest permissions")
	}
	if _, err := f.Write(yamlBytes); err != nil {
		return "", errors.Wrapf(err, "failed to write resolved manifest for %s", node.Identity)
	}
	return f.Name(), nil
}
