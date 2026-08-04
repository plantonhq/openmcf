package module

import (
	"strconv"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildCrdHelmValues renders the karpenter-crd release's values. The CRD
// chart's whole values surface is ONE knob (additionalAnnotations, stamped
// onto every CRD it templates) — the keep_on_uninstall flag rides it as the
// standard Helm resource-policy annotation. Without it a plain uninstall
// cascade-deletes every NodePool/EC2NodeClass/NodeClaim in the cluster
// along with the CRDs.
//
// PARITY: the Terraform module renders the identical document into its CRD
// release's values list.
func buildCrdHelmValues(locals *Locals) map[string]interface{} {
	if !locals.CrdsKeep {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"additionalAnnotations": map[string]interface{}{
			"helm.sh/resource-policy": "keep",
		},
	}
}

// buildHelmValues renders the typed spec into the controller chart's values
// map, then merges the spec's helm_values escape hatch over it with Helm
// `-f` semantics (maps deep-merge with the later document winning, lists
// replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// Rendering posture: fields whose spec default mirrors the chart default
// (replicas, logLevel, batching, scheduling policies, featureGates,
// priorityClassName, and the AWS arm's reservedENIs /
// vmMemoryOverheadPercent) render with the default APPLIED — explicit and
// byte-identical across engines regardless of whether the platform's
// defaulting middleware ran. Purely optional fields (endpoint, CA bundle,
// IRSA, queue, booleans that are chart-default-false) render only when set.
//
// No fullnameOverride: the release name ("karpenter") contains the chart
// name, so the chart's fullname template resolves to the release name —
// there is nothing for an override to pin.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}
	settings := map[string]interface{}{}

	// ---- cluster identity ----------------------------------------------------
	// clusterName is the one value the chart REFUSES to render without
	// (deployment.yaml wraps it in `required`) — always rendered, never
	// conditional.
	cluster := spec.GetCluster()
	settings["clusterName"] = cluster.GetName()
	if cluster.GetEndpoint() != "" {
		settings["clusterEndpoint"] = cluster.GetEndpoint()
	}
	if cluster.GetEksControlPlane() {
		settings["eksControlPlane"] = true
	}
	if cluster.GetCaBundle() != "" {
		settings["clusterCABundle"] = cluster.GetCaBundle()
	}

	// ---- AWS provider arm ------------------------------------------------------
	if aws := spec.GetAws(); aws != nil {
		// IRSA rides the service-account annotation the EKS webhook watches
		// — not a settings entry. Empty means EKS Pod Identity (association
		// configured cloud-side, no annotation needed). A valueFrom
		// reference is resolved to the AwsIamRole's role_arn before the
		// module runs.
		if aws.GetIrsaRoleArn().GetValue() != "" {
			values["serviceAccount"] = map[string]interface{}{
				"annotations": map[string]interface{}{
					"eks.amazonaws.com/role-arn": aws.GetIrsaRoleArn().GetValue(),
				},
			}
		}
		if aws.GetInterruptionQueue() != "" {
			settings["interruptionQueue"] = aws.GetInterruptionQueue()
		}
		if aws.GetIsolatedVpc() {
			settings["isolatedVPC"] = true
		}
		if aws.GetEnableZonalShift() {
			settings["enableZonalShift"] = true
		}

		// TYPE FIDELITY: the chart's reservedENIs default is the STRING "0"
		// — render the int32 as a string so the value keeps the served
		// chart's type. The Terraform module applies tostring() for the
		// same reason.
		reservedEnis := vars.DefaultReservedEnis
		if aws.ReservedEnis != nil {
			reservedEnis = int(aws.GetReservedEnis())
		}
		settings["reservedENIs"] = strconv.Itoa(reservedEnis)

		// TYPE FIDELITY (inverse case): the chart's vmMemoryOverheadPercent
		// default is the NUMBER 0.075 — the spec carries it as a string
		// (proto3 has no optional-double ergonomics for "unset"), so parse
		// and render a number. The Terraform module applies tonumber() for
		// the same reason.
		overhead := vars.DefaultVmMemoryOverheadPercent
		if aws.VmMemoryOverheadPercent != nil && aws.GetVmMemoryOverheadPercent() != "" {
			overhead = aws.GetVmMemoryOverheadPercent()
		}
		overheadNumber, err := strconv.ParseFloat(overhead, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse vm_memory_overhead_percent %q as a number", overhead)
		}
		settings["vmMemoryOverheadPercent"] = overheadNumber
	}

	// ---- controller sizing ----------------------------------------------------
	replicas := vars.DefaultReplicas
	logLevel := vars.DefaultLogLevel
	if c := spec.GetController(); c != nil {
		if c.Replicas != nil {
			replicas = int(c.GetReplicas())
		}
		if c.LogLevel != nil && c.GetLogLevel() != "" {
			logLevel = c.GetLogLevel()
		}
		// Resources live under the controller CONTAINER block
		// (controller.resources), unlike replicas/logLevel which are
		// top-level — the chart's layout, not ours.
		if r := resourcesMap(c.GetResources()); r != nil {
			values["controller"] = map[string]interface{}{
				"resources": r,
			}
		}
	}
	values["replicas"] = replicas
	values["logLevel"] = logLevel

	// ---- batching windows -------------------------------------------------------
	batchMax := vars.DefaultBatchMaxDuration
	batchIdle := vars.DefaultBatchIdleDuration
	if b := spec.GetBatching(); b != nil {
		if b.MaxDuration != nil && b.GetMaxDuration() != "" {
			batchMax = b.GetMaxDuration()
		}
		if b.IdleDuration != nil && b.GetIdleDuration() != "" {
			batchIdle = b.GetIdleDuration()
		}
	}
	settings["batchMaxDuration"] = batchMax
	settings["batchIdleDuration"] = batchIdle

	// ---- scheduler-simulation posture ---------------------------------------------
	preferencePolicy := vars.DefaultPreferencePolicy
	minValuesPolicy := vars.DefaultMinValuesPolicy
	if s := spec.GetScheduling(); s != nil {
		if s.PreferencePolicy != nil && s.GetPreferencePolicy() != "" {
			preferencePolicy = s.GetPreferencePolicy()
		}
		if s.MinValuesPolicy != nil && s.GetMinValuesPolicy() != "" {
			minValuesPolicy = s.GetMinValuesPolicy()
		}
	}
	settings["preferencePolicy"] = preferencePolicy
	settings["minValuesPolicy"] = minValuesPolicy

	// ---- feature gates ---------------------------------------------------------------
	// The WHOLE map renders with defaults applied — the deployment template
	// composes the FEATURE_GATES env var from all six keys unconditionally,
	// so explicit is safer than sparse (and reservedCapacity's default is
	// TRUE, unlike the other five).
	gates := spec.GetFeatureGates()
	reservedCapacity := true
	if gates != nil && gates.ReservedCapacity != nil {
		reservedCapacity = gates.GetReservedCapacity()
	}
	settings["featureGates"] = map[string]interface{}{
		"nodeRepair":              gates.GetNodeRepair(),
		"nodeOverlay":             gates.GetNodeOverlay(),
		"reservedCapacity":        reservedCapacity,
		"spotToSpotConsolidation": gates.GetSpotToSpotConsolidation(),
		"staticCapacity":          gates.GetStaticCapacity(),
		"capacityBuffer":          gates.GetCapacityBuffer(),
	}

	values["settings"] = settings

	// ---- controller-pod scheduling ------------------------------------------------
	// Where Karpenter itself runs — NOT what it provisions. The chart
	// already pins controller pods away from Karpenter-provisioned nodes
	// (nodeAffinity on karpenter.sh/nodepool DoesNotExist).
	cs := spec.GetControllerScheduling()
	priorityClassName := vars.DefaultPriorityClassName
	if cs != nil && cs.PriorityClassName != nil && cs.GetPriorityClassName() != "" {
		priorityClassName = cs.GetPriorityClassName()
	}
	values["priorityClassName"] = priorityClassName
	// nodeSelector MERGES onto the chart's kubernetes.io/os=linux default
	// (Helm deep-merges maps) — entries here narrow, never replace.
	if len(cs.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(cs.GetNodeSelector())
	}
	// tolerations REPLACE the chart's default list (CriticalAddonsOnly) —
	// Helm replaces lists wholesale — so render only when the spec provides
	// them; rendering an empty list would silently DROP the default.
	if len(cs.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(cs.GetTolerations())
	}
	if cs.GetHostNetwork() {
		values["hostNetwork"] = true
	}

	// ---- own telemetry ------------------------------------------------------------------
	if spec.GetPrometheus().GetServiceMonitor() {
		values["serviceMonitor"] = map[string]interface{}{
			"enabled": true,
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ------------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}
