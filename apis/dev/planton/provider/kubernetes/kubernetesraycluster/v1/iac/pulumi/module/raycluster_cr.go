package module

import (
	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesrayclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesraycluster/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// rayClusterCR renders the single ray.io/v1 RayCluster CR the KubeRay
// operator (a KubernetesKubeRayOperator install, the prerequisite)
// reconciles into the head pod, the worker group pods, the
// `<name>-head-svc` Service every exported endpoint rides, and — in
// token auth mode without a bring-your-own Secret — the generated
// bearer-token Secret named exactly after this resource.
//
// UNTYPED CR: the body assembles in ONE place (rayClusterSpecBody)
// through apiextensions.CustomResource with an untyped spec map, in byte
// lockstep with the Terraform twin's locals (locals.tf
// `raycluster_spec`). Every key renders ONLY when the spec declares it,
// so the operator's defaulting stays authoritative for everything the
// manifest leaves unsaid — except authOptions and the module-owned
// headGroupSpec.rayStartParams entries, deliberately always rendered
// (see the builders).
//
// NO WAIT on the CR, deliberately: cluster readiness depends on the
// operator (image pulls — Ray images are multi-GB — autoscaler sidecar
// injection, GCS startup) that is not part of applying the resource —
// the verifier owns readiness, the same never-block-on-a-controller
// posture as the sibling operator-CR modules.
func rayClusterCR(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) error {
	_, err := apiextensions.NewCustomResource(ctx, locals.ResourceName,
		&apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String(vars.ApiVersion),
			Kind:       pulumi.String(vars.Kind),
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ResourceName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
				// BACKGROUND deletion, explicitly: the OPERATOR owns the
				// RayCluster's cascade — its finalizer tears down the head
				// and worker pods, the head Service, and the generated
				// token Secret. Foreground propagation would block the
				// delete on children the operator keeps reconciling.
				// Terraform twin: delete_cascade = "Background" on
				// kubectl_manifest.
				Annotations: pulumi.StringMap{
					"pulumi.com/deletionPropagationPolicy": pulumi.String("background"),
				},
			},
			OtherFields: map[string]interface{}{
				"spec": rayClusterSpecBody(locals),
			},
		},
		append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return errors.Wrap(err, "failed to render the RayCluster custom resource")
	}
	return nil
}

// rayClusterSpecBody renders the RayCluster CR spec from the typed
// fields. Field names are the CRD's own JSON keys (verified against the
// pinned ray.io/v1 schema and raycluster_types.go). The Terraform twin
// (locals.tf `raycluster_spec`) renders the identical body — keep them
// in byte lockstep. No upgradeStrategy, no managedBy, no
// headServiceAnnotations — unmodeled.
func rayClusterSpecBody(locals *Locals) map[string]interface{} {
	spec := locals.Spec

	body := map[string]interface{}{
		"rayVersion": spec.GetRayVersion(),
	}

	// suspend deletes head and worker PODS but keeps the declaration
	// (and, with GCS fault tolerance, the external state).
	if spec.GetSuspend() {
		body["suspend"] = true
	}

	if spec.GetAutoscaling().GetEnabled() {
		body["enableInTreeAutoscaling"] = true
		if options := autoscalerOptionsBody(spec.GetAutoscaling()); len(options) > 0 {
			body["autoscalerOptions"] = options
		}
	}

	if spec.GetGcsFaultTolerance().GetEnabled() {
		body["gcsFaultToleranceOptions"] = gcsFaultToleranceBody(spec.GetGcsFaultTolerance())
	}

	body["authOptions"] = authOptionsBody(locals)

	body["headGroupSpec"] = headGroupSpecBody(locals)

	if workerGroups := spec.GetWorkerGroups(); len(workerGroups) > 0 {
		specs := []interface{}{}
		for _, group := range workerGroups {
			specs = append(specs, workerGroupSpecBody(locals, group))
		}
		body["workerGroupSpecs"] = specs
	}

	return body
}

// authOptionsBody renders spec.auth → authOptions. This catalog's
// default is TOKEN auth (secure-by-default) while the operator's own
// nil-authOptions default is DISABLED — so authOptions renders ALWAYS,
// never left to the CR default: an absent block would silently deploy
// the legacy open cluster (anyone reaching the dashboard port runs
// arbitrary code). Empty auth or empty mode means token; only an
// explicit "disabled" opts out.
func authOptionsBody(locals *Locals) map[string]interface{} {
	if !locals.TokenAuthEnabled {
		return map[string]interface{}{"mode": "disabled"}
	}
	out := map[string]interface{}{"mode": "token"}
	// secretName tells the operator to skip generating a token Secret
	// and read the bring-your-own one (data key `auth_token`) instead.
	if existing := locals.Spec.GetAuth().GetExistingTokenSecretName(); existing != "" {
		out["secretName"] = existing
	}
	return out
}

// autoscalerOptionsBody renders spec.autoscaling → autoscalerOptions.
// Keys render only when declared — the operator's defaults (60s idle
// timeout, Default upscaling) stay authoritative otherwise.
func autoscalerOptionsBody(autoscaling *kubernetesrayclusterv1.KubernetesRayClusterAutoscaling) map[string]interface{} {
	out := map[string]interface{}{}
	if autoscaling.IdleTimeoutSeconds != nil {
		out["idleTimeoutSeconds"] = int(autoscaling.GetIdleTimeoutSeconds())
	}
	if autoscaling.GetUpscalingMode() != "" {
		out["upscalingMode"] = autoscaling.GetUpscalingMode()
	}
	if resources := resourcesBody(autoscaling.GetResources()); len(resources) > 0 {
		out["resources"] = resources
	}
	return out
}

// gcsFaultToleranceBody renders spec.gcs_fault_tolerance →
// gcsFaultToleranceOptions. STATE TRUTH: without this block the head
// pod's GCS holds the cluster's control state in memory — losing the
// head loses every job, actor, and worker registration. With it, state
// lives in the external Redis-protocol store and a replaced head
// RECOVERS.
func gcsFaultToleranceBody(gcsFt *kubernetesrayclusterv1.KubernetesRayClusterGcsFaultTolerance) map[string]interface{} {
	// redisAddress arrives pre-resolved (literal or KubernetesValkey
	// reference) and is ALWAYS rendered inside the block — the spec's
	// CEL guarantees it is set whenever enabled is true.
	out := map[string]interface{}{
		"redisAddress": gcsFt.GetRedisAddress().GetValue(),
	}
	if passwordSecret := gcsFt.GetRedisPasswordSecret(); passwordSecret != nil {
		out["redisPassword"] = map[string]interface{}{
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"name": passwordSecret.GetName().GetValue(),
					"key":  passwordSecret.GetKey(),
				},
			},
		}
	}
	// Empty means the operator derives one from the cluster's UID (safe
	// default); set explicitly only when state must survive
	// delete-and-recreate of the CR itself.
	if gcsFt.GetExternalStorageNamespace() != "" {
		out["externalStorageNamespace"] = gcsFt.GetExternalStorageNamespace()
	}
	return out
}

// headGroupSpecBody renders spec.head → headGroupSpec.
//
// The module OWNS two rayStartParams and writes them LAST so user
// entries cannot override them:
//
//	dashboard-host "0.0.0.0" — the dashboard binds localhost otherwise
//	  and the head Service cannot answer (sample-verified: upstream
//	  ray-cluster.complete.yaml sets exactly this);
//	num-cpus "0" (unless schedule_tasks_on_head) — keeps the Ray
//	  scheduler from placing application work on the head (a
//	  task-loaded head starves the GCS).
func headGroupSpecBody(locals *Locals) map[string]interface{} {
	head := locals.Spec.GetHead()

	rayStartParams := map[string]interface{}{}
	for key, value := range head.GetRayStartParams() {
		rayStartParams[key] = value
	}
	rayStartParams["dashboard-host"] = "0.0.0.0"
	if !head.GetScheduleTasksOnHead() {
		rayStartParams["num-cpus"] = "0"
	}

	container := map[string]interface{}{
		"name":  "ray-head",
		"image": locals.Image,
	}
	if resources := resourcesBody(head.GetResources()); len(resources) > 0 {
		container["resources"] = resources
	}

	podSpec := map[string]interface{}{
		"containers": []interface{}{container},
	}
	schedulingEntries(podSpec, head.GetScheduling())

	return map[string]interface{}{
		"rayStartParams": rayStartParams,
		"template": map[string]interface{}{
			"spec": podSpec,
		},
	}
}

// workerGroupSpecBody renders one spec.worker_groups entry →
// workerGroupSpecs item.
//
// replicas/minReplicas/maxReplicas render only when declared — the v1
// CRD defaults all three (replicas 0, minReplicas 0, maxReplicas
// maxInt32; verified in the pinned ray.io_rayclusters.yaml schema).
//
// rayStartParams renders ALWAYS, {} when empty: NOT required by the v1
// CRD (only the retired v1alpha1 schema listed it under `required`;
// verified in the pinned CRD), but the operator's own Go type
// serializes it unconditionally and every upstream sample renders
// `rayStartParams: {}` explicitly — matching that keeps SSA diffs
// quiet.
func workerGroupSpecBody(locals *Locals, group *kubernetesrayclusterv1.KubernetesRayClusterWorkerGroup) map[string]interface{} {
	out := map[string]interface{}{
		"groupName": group.GetName(),
	}
	if group.Replicas != nil {
		out["replicas"] = int(group.GetReplicas())
	}
	if group.MinReplicas != nil {
		out["minReplicas"] = int(group.GetMinReplicas())
	}
	if group.MaxReplicas != nil {
		out["maxReplicas"] = int(group.GetMaxReplicas())
	}

	rayStartParams := map[string]interface{}{}
	for key, value := range group.GetRayStartParams() {
		rayStartParams[key] = value
	}
	out["rayStartParams"] = rayStartParams

	// extra_resource_limits land in LIMITS ONLY: extended resources
	// (nvidia.com/gpu and friends) must not be requested without limits
	// — Kubernetes rejects requests-without-limits for extended
	// resources, and Ray discovers accelerators from the container
	// LIMITS.
	limits := cpuMemoryBody(group.GetResources().GetLimits())
	for key, value := range group.GetExtraResourceLimits() {
		limits[key] = value
	}
	requests := cpuMemoryBody(group.GetResources().GetRequests())

	container := map[string]interface{}{
		"name":  "ray-worker",
		"image": locals.Image,
	}
	if len(limits) > 0 || len(requests) > 0 {
		resources := map[string]interface{}{}
		if len(limits) > 0 {
			resources["limits"] = limits
		}
		if len(requests) > 0 {
			resources["requests"] = requests
		}
		container["resources"] = resources
	}

	podSpec := map[string]interface{}{
		"containers": []interface{}{container},
	}
	schedulingEntries(podSpec, group.GetScheduling())

	out["template"] = map[string]interface{}{
		"spec": podSpec,
	}
	return out
}

// schedulingEntries contributes spec scheduling onto a pod-template
// spec. Unlike CRs that model only affinity, the RayCluster embeds a
// full corev1 pod template — nodeSelector renders verbatim.
func schedulingEntries(podSpec map[string]interface{}, scheduling *kubernetesrayclusterv1.KubernetesRayClusterScheduling) {
	if scheduling == nil {
		return
	}

	if nodeSelector := scheduling.GetNodeSelector(); len(nodeSelector) > 0 {
		entries := map[string]interface{}{}
		for key, value := range nodeSelector {
			entries[key] = value
		}
		podSpec["nodeSelector"] = entries
	}

	if tolerations := scheduling.GetTolerations(); len(tolerations) > 0 {
		entries := []interface{}{}
		for _, toleration := range tolerations {
			entry := map[string]interface{}{}
			if toleration.GetKey() != "" {
				entry["key"] = toleration.GetKey()
			}
			if toleration.GetOperator() != "" {
				entry["operator"] = toleration.GetOperator()
			}
			if toleration.GetValue() != "" {
				entry["value"] = toleration.GetValue()
			}
			if toleration.GetEffect() != "" {
				entry["effect"] = toleration.GetEffect()
			}
			if toleration.TolerationSeconds != nil {
				entry["tolerationSeconds"] = int(toleration.GetTolerationSeconds())
			}
			entries = append(entries, entry)
		}
		podSpec["tolerations"] = entries
	}

	if scheduling.GetPriorityClassName() != "" {
		podSpec["priorityClassName"] = scheduling.GetPriorityClassName()
	}
}

func resourcesBody(resources *kubernetesprovider.ContainerResources) map[string]interface{} {
	if resources == nil {
		return nil
	}
	out := map[string]interface{}{}
	if limits := cpuMemoryBody(resources.GetLimits()); len(limits) > 0 {
		out["limits"] = limits
	}
	if requests := cpuMemoryBody(resources.GetRequests()); len(requests) > 0 {
		out["requests"] = requests
	}
	return out
}

func cpuMemoryBody(cpuMemory *kubernetesprovider.CpuMemory) map[string]interface{} {
	out := map[string]interface{}{}
	if cpuMemory == nil {
		return out
	}
	if cpuMemory.GetCpu() != "" {
		out["cpu"] = cpuMemory.GetCpu()
	}
	if cpuMemory.GetMemory() != "" {
		out["memory"] = cpuMemory.GetMemory()
	}
	return out
}
