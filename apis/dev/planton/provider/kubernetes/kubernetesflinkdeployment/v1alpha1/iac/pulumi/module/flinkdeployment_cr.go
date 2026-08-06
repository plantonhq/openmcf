package module

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesflinkdeploymentv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesflinkdeployment/v1alpha1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// flinkDeploymentCR renders the single flink.apache.org/v1beta1
// FlinkDeployment CR the Flink Kubernetes Operator (a
// KubernetesFlinkOperator install, the prerequisite) reconciles into the
// JobManager, its TaskManagers, the `<name>-rest` Service and (in
// application mode) the job they run.
//
// UNTYPED CR: the body assembles in ONE place (flinkDeploymentSpecBody)
// through apiextensions.CustomResource with an untyped spec map, in byte
// lockstep with the Terraform twin's locals (locals.tf
// `flinkdeployment_spec`). Every key renders ONLY when the spec declares
// it, so the operator's defaulting stays authoritative for everything
// the manifest leaves unsaid — except flinkVersion/image/serviceAccount
// (ALWAYS: the deployment's identity) and jobManager/taskManager
// (ALWAYS: the operator's validator requires resource on both tiers).
//
// NO WAIT on the CR, deliberately: cluster readiness depends on the
// operator (image pulls, job submission, TaskManager registration) that
// is not part of applying the resource — the verifier owns readiness,
// the same never-block-on-a-controller posture as the sibling
// operator-CR modules.
func flinkDeploymentCR(ctx *pulumi.Context, locals *Locals,
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
				// FlinkDeployment's cascade — its finalizer tears down
				// the JobManager, TaskManagers and Services. Foreground
				// propagation would block the delete on children the
				// operator keeps reconciling. Terraform twin:
				// delete_cascade = "Background" on kubectl_manifest.
				Annotations: pulumi.StringMap{
					"pulumi.com/deletionPropagationPolicy": pulumi.String("background"),
				},
			},
			OtherFields: map[string]interface{}{
				"spec": flinkDeploymentSpecBody(locals),
			},
		},
		append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return errors.Wrap(err, "failed to render the FlinkDeployment custom resource")
	}
	return nil
}

// flinkDeploymentSpecBody renders the FlinkDeployment CR spec from the
// typed fields. Field names are the CR's own JSON keys (verified against
// the pinned operator's API classes: FlinkDeploymentSpec / JobSpec /
// JobManagerSpec / TaskManagerSpec / Resource). The Terraform twin
// (locals.tf `flinkdeployment_spec`) renders the identical body — keep
// them in byte lockstep.
func flinkDeploymentSpecBody(locals *Locals) map[string]interface{} {
	spec := locals.Spec

	body := map[string]interface{}{
		"flinkVersion":   spec.GetFlinkVersion(),
		"image":          resolveImage(spec),
		"serviceAccount": resolveServiceAccount(spec),
	}

	// "native" is the CR default — mode renders only on divergence.
	if spec.GetMode() == "standalone" {
		body["mode"] = "standalone"
	}

	if configuration := flinkConfigurationBody(spec); len(configuration) > 0 {
		body["flinkConfiguration"] = configuration
	}

	if job := jobBody(spec.GetJob()); len(job) > 0 {
		body["job"] = job
	}

	if spec.GetRestartNonce() != 0 {
		body["restartNonce"] = spec.GetRestartNonce()
	}

	// Both tiers render ALWAYS — the operator's validator REQUIRES
	// resource memory on each ("JobManager resource memory must be
	// defined…" / "TaskManager resource memory must be defined…"),
	// defaulting to 1 CPU / 2Gi (the spec's default_container_resources
	// values) when the tier or its resources are unset.
	body["jobManager"] = jobManagerBody(spec.GetJobManager())
	body["taskManager"] = taskManagerBody(spec.GetTaskManager())

	if logConfiguration := spec.GetLogConfiguration(); len(logConfiguration) > 0 {
		body["logConfiguration"] = stringMapToInterface(logConfiguration)
	}

	if podTemplate := podTemplateBody(spec); podTemplate != nil {
		body["podTemplate"] = podTemplate
	}

	return body
}

// resolveImage applies the VERSION/IMAGE LOCKSTEP: the default image
// derives from flink_version (`flink:<major>.<minor>` — v2_1 →
// flink:2.1) by stripping the "v" and replacing "_" with ".". A custom
// image must carry exactly that Flink version — the operator shapes its
// submission protocol from the declared version and a mismatch fails at
// runtime, not at apply.
func resolveImage(spec *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentSpec) string {
	if spec.GetImage() != "" {
		return spec.GetImage()
	}
	return fmt.Sprintf("flink:%s",
		strings.ReplaceAll(strings.TrimPrefix(spec.GetFlinkVersion(), "v"), "_", "."))
}

// resolveServiceAccount: empty = "flink" — the account the
// KubernetesFlinkOperator's chart creates with reconcile RBAC.
func resolveServiceAccount(spec *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentSpec) string {
	if spec.GetServiceAccount() != "" {
		return spec.GetServiceAccount()
	}
	return vars.DefaultServiceAccount
}

// flinkConfigurationBody merges (a) the user's flink_configuration
// entries with (b) the module-owned state keys rendered from spec.state.
// The typed fields are the truth: they merge LAST, so a colliding user
// entry loses, deliberately.
//
// "high-availability.type: kubernetes" is the CURRENT key form at the
// pinned operator (its own test/e2e manifests use exactly this key —
// test-deployment-key-value-configuration.yaml, multi-sessionjob.yaml —
// and examples/basic-checkpoint-ha.yaml carries the same key in nested
// YAML form); the legacy `high-availability: <factory-class>` form is
// not used.
//
// s3.path.style.access renders with its effective value whenever s3 is
// set: the spec default (true) is correct for in-cluster object stores;
// explicit false is the AWS-S3-itself posture. Credentials NEVER go here
// — flinkConfiguration renders into a ConfigMap in clear text; they ride
// pod env from Secret refs (see the pod template).
func flinkConfigurationBody(spec *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentSpec) map[string]interface{} {
	merged := map[string]interface{}{}
	for key, value := range spec.GetFlinkConfiguration() {
		merged[key] = value
	}

	state := spec.GetState()
	if state == nil {
		return merged
	}

	if state.GetCheckpointsDir() != "" {
		merged["state.checkpoints.dir"] = state.GetCheckpointsDir()
	}
	if state.GetSavepointsDir() != "" {
		merged["state.savepoints.dir"] = state.GetSavepointsDir()
	}
	if state.GetHighAvailability().GetEnabled() {
		merged["high-availability.type"] = "kubernetes"
		merged["high-availability.storageDir"] = state.GetHighAvailability().GetStorageDir()
	}
	if s3 := state.GetS3(); s3 != nil {
		merged["s3.endpoint"] = s3.GetEndpoint().GetValue()
		pathStyleAccess := true
		if s3.PathStyleAccess != nil {
			pathStyleAccess = s3.GetPathStyleAccess()
		}
		merged["s3.path.style.access"] = strconv.FormatBool(pathStyleAccess)
	}
	return merged
}

// jobBody renders spec.job. Set = APPLICATION cluster (the cluster runs
// exactly this job); absent = SESSION cluster. state renders only on
// divergence from the CR default "running" (i.e. "suspended"),
// upgradeMode only on divergence from "stateless" (i.e.
// "last-state"/"savepoint").
func jobBody(job *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentJob) map[string]interface{} {
	if job == nil {
		return nil
	}
	out := map[string]interface{}{
		"jarURI": job.GetJarUri(),
	}
	if job.GetEntryClass() != "" {
		out["entryClass"] = job.GetEntryClass()
	}
	if args := job.GetArgs(); len(args) > 0 {
		out["args"] = stringSliceToInterface(args)
	}
	if job.Parallelism != nil {
		out["parallelism"] = int(job.GetParallelism())
	}
	if job.GetState() == "suspended" {
		out["state"] = "suspended"
	}
	if upgradeMode := job.GetUpgradeMode(); upgradeMode == "last-state" || upgradeMode == "savepoint" {
		out["upgradeMode"] = upgradeMode
	}
	if job.GetInitialSavepointPath() != "" {
		out["initialSavepointPath"] = job.GetInitialSavepointPath()
	}
	if job.GetAllowNonRestoredState() {
		out["allowNonRestoredState"] = true
	}
	if job.GetSavepointTriggerNonce() != 0 {
		out["savepointTriggerNonce"] = job.GetSavepointTriggerNonce()
	}
	return out
}

// jobManagerBody renders spec.jobManager. Replicas render only past 1
// (standbys — the spec-level rule requires state.high_availability for
// them).
func jobManagerBody(jobManager *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentJobManager) map[string]interface{} {
	out := map[string]interface{}{
		"resource": resourceBody(jobManager.GetResources()),
	}
	if jobManager != nil && jobManager.Replicas != nil && jobManager.GetReplicas() > 1 {
		out["replicas"] = int(jobManager.GetReplicas())
	}
	return out
}

// taskManagerBody renders spec.taskManager. Replicas render whenever
// declared — meaningful in standalone mode only (native mode derives
// worker count from the job's parallelism).
func taskManagerBody(taskManager *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentTaskManager) map[string]interface{} {
	out := map[string]interface{}{
		"resource": resourceBody(taskManager.GetResources()),
	}
	if taskManager != nil && taskManager.Replicas != nil {
		out["replicas"] = int(taskManager.GetReplicas())
	}
	return out
}

// resourceBody converts the spec's ContainerResources into the CR's
// Resource shape — cpu as a NUMBER (the CR's Java Double), memory as a
// string. cpu prefers requests.cpu, else limits.cpu; memory prefers
// limits.memory (the ceiling Flink sizes its JVM from), else
// requests.memory. Defaults: 1 CPU / 2Gi (the spec's
// default_container_resources values).
func resourceBody(resources *kubernetesprovider.ContainerResources) map[string]interface{} {
	cpuQuantity := resources.GetRequests().GetCpu()
	if cpuQuantity == "" {
		cpuQuantity = resources.GetLimits().GetCpu()
	}
	cpu := vars.DefaultCpuCores
	if cpuQuantity != "" {
		cpu = quantityToCores(cpuQuantity)
	}

	memory := resources.GetLimits().GetMemory()
	if memory == "" {
		memory = resources.GetRequests().GetMemory()
	}
	if memory == "" {
		memory = vars.DefaultMemory
	}

	return map[string]interface{}{
		"cpu":    cpu,
		"memory": memory,
	}
}

// quantityToCores converts a Kubernetes CPU quantity to the CR's numeric
// core count: "1000m" → 1.0, "500m" → 0.5, "2" → 2.0. Millicores divide
// by 1000; anything else parses as a plain number. Identical semantics
// in the Terraform twin (locals.tf jm_cpu/tm_cpu). Malformed quantities
// were rejected up front by validateCpuQuantities (Resources runs it
// before building the CR — the fail-loud twin of the Terraform plan's
// tonumber() failure), so this parse cannot fail; the defensive return
// keeps the compiler honest, never a rendered value.
func quantityToCores(quantity string) float64 {
	if strings.HasSuffix(quantity, "m") {
		milliCores, err := strconv.ParseFloat(strings.TrimSuffix(quantity, "m"), 64)
		if err != nil {
			return vars.DefaultCpuCores
		}
		return milliCores / 1000
	}
	cores, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return vars.DefaultCpuCores
	}
	return cores
}

// validateCpuQuantities rejects malformed CPU quantity strings on the
// JobManager and TaskManager tiers before any resource is built —
// keeping this engine's failure semantics identical to the Terraform
// twin, whose tonumber() fails the plan on the same inputs.
func validateCpuQuantities(spec *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentSpec) error {
	check := func(tier, quantity string) error {
		if quantity == "" {
			return nil
		}
		numeric := quantity
		if strings.HasSuffix(quantity, "m") {
			numeric = strings.TrimSuffix(quantity, "m")
		}
		if _, err := strconv.ParseFloat(numeric, 64); err != nil {
			return errors.Errorf(
				"%s cpu quantity %q is not a valid Kubernetes CPU quantity (use forms like \"500m\" or \"2\")",
				tier, quantity)
		}
		return nil
	}
	for tier, resources := range map[string]*kubernetesprovider.ContainerResources{
		"job_manager":  spec.GetJobManager().GetResources(),
		"task_manager": spec.GetTaskManager().GetResources(),
	} {
		if err := check(tier+".resources.requests", resources.GetRequests().GetCpu()); err != nil {
			return err
		}
		if err := check(tier+".resources.limits", resources.GetLimits().GetCpu()); err != nil {
			return err
		}
	}
	return nil
}

// podTemplateBody renders spec.podTemplate — only when it would carry
// something: scheduling, image-pull Secrets, or the S3 credential env.
// "flink-main-container" is the operator's merge contract for the main
// container (its own examples/pod-template.yaml: "Do not change the
// main container name").
func podTemplateBody(spec *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentSpec) map[string]interface{} {
	scheduling := spec.GetScheduling()
	imagePullSecrets := spec.GetImagePullSecrets()
	s3 := spec.GetState().GetS3()

	if scheduling == nil && len(imagePullSecrets) == 0 && s3 == nil {
		return nil
	}

	podSpec := map[string]interface{}{}

	if len(imagePullSecrets) > 0 {
		entries := []interface{}{}
		for _, secretName := range imagePullSecrets {
			entries = append(entries, map[string]interface{}{"name": secretName})
		}
		podSpec["imagePullSecrets"] = entries
	}

	if scheduling != nil {
		if nodeSelector := scheduling.GetNodeSelector(); len(nodeSelector) > 0 {
			podSpec["nodeSelector"] = stringMapToInterface(nodeSelector)
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

	if env := s3EnvBody(s3); len(env) > 0 {
		podSpec["containers"] = []interface{}{
			map[string]interface{}{
				"name": vars.MainContainerName,
				"env":  env,
			},
		}
	}

	return map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata":   map[string]interface{}{"name": vars.PodTemplateName},
		"spec":       podSpec,
	}
}

// s3EnvBody renders the S3 credential seam: credentials ride pod env
// from Secret refs — NEVER into flinkConfiguration, which renders into a
// ConfigMap in clear text. ENABLE_BUILT_IN_PLUGINS activates the named
// S3 filesystem plugin from the image's bundled (disabled-by-default)
// plugin set.
func s3EnvBody(s3 *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentS3) []interface{} {
	if s3 == nil {
		return nil
	}
	env := []interface{}{
		secretEnvVar("AWS_ACCESS_KEY_ID", s3.GetAccessKeySecret()),
		secretEnvVar("AWS_SECRET_ACCESS_KEY", s3.GetSecretKeySecret()),
	}
	if s3.GetBuiltinPluginJar() != "" {
		env = append(env, map[string]interface{}{
			"name":  "ENABLE_BUILT_IN_PLUGINS",
			"value": s3.GetBuiltinPluginJar(),
		})
	}
	return env
}

func secretEnvVar(name string, selector *kubernetesflinkdeploymentv1alpha1.KubernetesFlinkDeploymentSecretSelector) map[string]interface{} {
	return map[string]interface{}{
		"name": name,
		"valueFrom": map[string]interface{}{
			"secretKeyRef": map[string]interface{}{
				"name": selector.GetName().GetValue(),
				"key":  selector.GetKey(),
			},
		},
	}
}

func stringMapToInterface(in map[string]string) map[string]interface{} {
	out := map[string]interface{}{}
	for key, value := range in {
		out[key] = value
	}
	return out
}

func stringSliceToInterface(in []string) []interface{} {
	out := []interface{}{}
	for _, value := range in {
		out = append(out, value)
	}
	return out
}
