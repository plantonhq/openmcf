package module

import (
	"github.com/pkg/errors"
	kubernetesargoworkflowsv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesargoworkflows/v1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// SECRET DISCIPLINE (load-bearing): nothing rendered here carries
// credential material. Artifact-store keys and archive-database
// credentials ride the chart's own secret SELECTORS ({name, key} pairs the
// workloads resolve at runtime); keyless arms lean on ambient pod identity
// (useSDKCreds) on the RUNNER service account.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// fullnameOverride pins every child name (`<name>-server`,
	// `<name>-workflow-controller`, ...) to the resource name; the
	// exported outputs are built from that contract.
	values["fullnameOverride"] = locals.ReleaseName

	// ---- CRD lifecycle -----------------------------------------------------
	// Full-schema CRDs arrive via the chart's hook Job, which DOWNLOADS
	// them from the chart's GitHub release at install time
	// (internet-at-install); crds.full=false falls back to templated
	// minified CRDs for air-gapped clusters.
	if crds := spec.GetCrds(); crds != nil {
		install := true
		if crds.Install != nil {
			install = crds.GetInstall()
		}
		keep := true
		if crds.Keep != nil {
			keep = crds.GetKeep()
		}
		full := true
		if crds.FullSchema != nil {
			full = crds.GetFullSchema()
		}
		crdsBlock := map[string]interface{}{
			"install": install,
			"keep":    keep,
			"full":    full,
		}
		if crds.GetBaseUrl() != "" {
			crdsBlock["upgradeJob"] = map[string]interface{}{
				"crdBaseURL": crds.GetBaseUrl(),
			}
		}
		values["crds"] = crdsBlock
	}

	// ---- image override ------------------------------------------------------
	// The registry/tag override maps onto the chart's SPLIT image values:
	// images.tag is the shared tag; each component's image.registry moves
	// to the mirror while the upstream repository paths stay (the
	// split-image discipline — a combined mapping would break every
	// mirror override identically in both engines).
	img := spec.GetImage()
	if img != nil && (img.GetTag() != "" || img.GetPullSecretName() != "") {
		images := map[string]interface{}{}
		if img.GetTag() != "" {
			images["tag"] = img.GetTag()
		}
		if img.GetPullSecretName() != "" {
			images["pullSecrets"] = []interface{}{
				map[string]interface{}{"name": img.GetPullSecretName()},
			}
		}
		values["images"] = images
	}
	registryOverride := map[string]interface{}{}
	if img.GetRegistry() != "" {
		registryOverride["image"] = map[string]interface{}{"registry": img.GetRegistry()}
	}

	// ---- scheduling (rendered per component — the chart has no global block) --
	scheduling := map[string]interface{}{}
	if sched := spec.GetScheduling(); sched != nil {
		if len(sched.GetNodeSelector()) > 0 {
			scheduling["nodeSelector"] = stringMapToInterface(sched.GetNodeSelector())
		}
		if len(sched.GetTolerations()) > 0 {
			scheduling["tolerations"] = tolerationsSlice(sched.GetTolerations())
		}
		if sched.GetPriorityClassName() != "" {
			scheduling["priorityClassName"] = sched.GetPriorityClassName()
		}
	}

	// ---- controller ------------------------------------------------------------
	// workflowNamespaces is ALWAYS rendered: the chart's own default is
	// ["default"], and its workflow-role template creates the runner
	// SA/Role/RoleBinding in every listed namespace PLUS the release
	// namespace — leaving the default in place makes every install leak
	// runner RBAC into the cluster's `default` namespace. An explicit
	// empty list keeps the runner identity release-namespace-only (lists
	// REPLACE under Helm -f merge semantics).
	controller := map[string]interface{}{}
	workflowNamespaces := []interface{}{}
	if c := spec.GetController(); c != nil {
		if c.Replicas != nil {
			controller["replicas"] = int(c.GetReplicas())
		}
		if resources := resourcesBlock(c.GetResources()); resources != nil {
			controller["resources"] = resources
		}
		for _, ns := range c.GetWorkflowNamespaces() {
			workflowNamespaces = append(workflowNamespaces, ns)
		}
		if c.GetInstanceId() != "" {
			// instanceID is a STRUCTURED chart block ({enabled,
			// useReleaseName, explicitID} — templates read .enabled
			// directly); the spec's plain string maps to the
			// enabled+explicitID shape, never a bare string.
			controller["instanceID"] = map[string]interface{}{
				"enabled":    true,
				"explicitID": c.GetInstanceId(),
			}
		}
		if c.Parallelism != nil {
			controller["parallelism"] = int(c.GetParallelism())
		}
		if c.NamespaceParallelism != nil {
			controller["namespaceParallelism"] = int(c.GetNamespaceParallelism())
		}
	}
	controller["workflowNamespaces"] = workflowNamespaces

	// The workflow archive (controller.persistence): the engine section
	// is keyed by name (postgresql | mysql) with secret SELECTORS for the
	// credentials — resolved by the controller at runtime, never rendered
	// as values.
	if archive := spec.GetArchive(); archive != nil {
		usernameKey := archive.GetCredentialsSecret().GetUsernameKey()
		if usernameKey == "" {
			usernameKey = vars.ArchiveUsernameKey
		}
		passwordKey := archive.GetCredentialsSecret().GetPasswordKey()
		if passwordKey == "" {
			passwordKey = vars.ArchivePasswordKey
		}

		engineKey := "postgresql"
		defaultPort := 5432
		if archive.GetEngine() == kubernetesargoworkflowsv1.KubernetesArgoWorkflowsArchiveEngine_mysql {
			engineKey = "mysql"
			defaultPort = 3306
		}
		port := defaultPort
		if archive.Port != nil {
			port = int(archive.GetPort())
		}

		engineBlock := map[string]interface{}{
			"host":      archive.GetHost().GetValue(),
			"port":      port,
			"database":  archive.GetDatabase(),
			"tableName": vars.ArchiveTableName,
			"userNameSecret": map[string]interface{}{
				"name": archive.GetCredentialsSecret().GetName(),
				"key":  usernameKey,
			},
			"passwordSecret": map[string]interface{}{
				"name": archive.GetCredentialsSecret().GetName(),
				"key":  passwordKey,
			},
		}
		if archive.GetSslMode() != "" {
			engineBlock["sslMode"] = archive.GetSslMode()
			if archive.GetSslMode() != "disable" {
				engineBlock["ssl"] = true
			}
		}

		controller["persistence"] = map[string]interface{}{
			"archive": true,
			engineKey: engineBlock,
		}
	}

	if rp := spec.GetRetentionPolicy(); rp != nil {
		retention := map[string]interface{}{}
		if rp.Completed != nil {
			retention["completed"] = int(rp.GetCompleted())
		}
		if rp.Failed != nil {
			retention["failed"] = int(rp.GetFailed())
		}
		if rp.Errored != nil {
			retention["errored"] = int(rp.GetErrored())
		}
		if len(retention) > 0 {
			controller["retentionPolicy"] = retention
		}
	}

	if spec.GetServiceMonitorEnabled() {
		controller["serviceMonitor"] = map[string]interface{}{"enabled": true}
	}

	mergeInto(controller, scheduling)
	mergeInto(controller, registryOverride)
	if len(controller) > 0 {
		values["controller"] = controller
	}

	// ---- server ------------------------------------------------------------------
	server := map[string]interface{}{}
	if !locals.ServerEnabled {
		server["enabled"] = false
	} else {
		if s := spec.GetServer(); s != nil {
			if s.Replicas != nil {
				server["replicas"] = int(s.GetReplicas())
			}
			if resources := resourcesBlock(s.GetResources()); resources != nil {
				server["resources"] = resources
			}
			if len(s.GetAuthModes()) > 0 {
				modes := make([]interface{}, 0, len(s.GetAuthModes()))
				for _, m := range s.GetAuthModes() {
					modes = append(modes, m)
				}
				server["authModes"] = modes
			}
			if s.GetSecure() {
				server["secure"] = true
			}
			if s.GetBaseHref() != "" {
				server["baseHref"] = s.GetBaseHref()
			}
		}
		mergeInto(server, scheduling)
		mergeInto(server, registryOverride)
	}
	if len(server) > 0 {
		values["server"] = server
	}

	// ---- executor (registry override only; tuning via helm_values) -----------------
	if len(registryOverride) > 0 {
		values["executor"] = map[string]interface{}{
			"image": map[string]interface{}{"registry": img.GetRegistry()},
		}
	}

	// ---- the workflow runner ServiceAccount -----------------------------------------
	// The chart DEFAULTS workflow.serviceAccount.create to FALSE — an
	// engine whose runner ServiceAccount does not exist rejects every
	// submission. The module always creates it (with the configured
	// name); the chart's workflow.rbac.create default (true) then binds
	// the runner Role to it in every watched namespace.
	values["workflow"] = map[string]interface{}{
		"serviceAccount": map[string]interface{}{
			"create": true,
			"name":   locals.WorkflowServiceAccount,
		},
	}

	// ---- artifact repository -----------------------------------------------------------
	// Exactly one backend renders (proto oneof). Secret material rides
	// the chart's secret selectors; keyless arms lean on ambient pod
	// identity (useSDKCreds) on the runner service account.
	if ar := spec.GetArtifactRepository(); ar != nil {
		repository := map[string]interface{}{}
		if ar.GetArchiveLogs() {
			repository["archiveLogs"] = true
		}
		switch backend := ar.GetBackend().(type) {
		case *kubernetesargoworkflowsv1.KubernetesArgoWorkflowsArtifactRepository_S3:
			s3 := map[string]interface{}{"bucket": backend.S3.GetBucket()}
			if backend.S3.GetEndpoint().GetValue() != "" {
				s3["endpoint"] = backend.S3.GetEndpoint().GetValue()
			}
			if backend.S3.GetRegion() != "" {
				s3["region"] = backend.S3.GetRegion()
			}
			if backend.S3.GetInsecure() {
				s3["insecure"] = true
			}
			if backend.S3.GetUseAmbientCredentials() {
				// Keyless: sign with the pod's ambient identity (IRSA /
				// workload identity on the runner service account).
				s3["useSDKCreds"] = true
			} else {
				// The Secret name resolves a KubernetesSeaweedFs
				// credentials Secret by reference (or a literal); the key
				// names default to that kind's generated `-s3-secret`
				// convention so it composes with zero key configuration.
				creds := backend.S3.GetCredentialsSecret()
				accessKeyIdKey := creds.GetAccessKeyIdKey()
				if accessKeyIdKey == "" {
					accessKeyIdKey = vars.S3AccessKeyIdKeyDefault
				}
				secretAccessKeyKey := creds.GetSecretAccessKeyKey()
				if secretAccessKeyKey == "" {
					secretAccessKeyKey = vars.S3SecretAccessKeyKeyDefault
				}
				s3["accessKeySecret"] = map[string]interface{}{
					"name": creds.GetSecretName().GetValue(),
					"key":  accessKeyIdKey,
				}
				s3["secretKeySecret"] = map[string]interface{}{
					"name": creds.GetSecretName().GetValue(),
					"key":  secretAccessKeyKey,
				}
			}
			repository["s3"] = s3
		case *kubernetesargoworkflowsv1.KubernetesArgoWorkflowsArtifactRepository_Gcs:
			gcs := map[string]interface{}{"bucket": backend.Gcs.GetBucket()}
			if backend.Gcs.GetCredentialsSecretName() != "" {
				gcs["serviceAccountKeySecret"] = map[string]interface{}{
					"name": backend.Gcs.GetCredentialsSecretName(),
					"key":  vars.GcsServiceAccountKey,
				}
			}
			repository["gcs"] = gcs
		case *kubernetesargoworkflowsv1.KubernetesArgoWorkflowsArtifactRepository_Azure:
			azure := map[string]interface{}{
				"endpoint":  backend.Azure.GetEndpoint(),
				"container": backend.Azure.GetContainer(),
			}
			if backend.Azure.GetCredentialsSecretName() != "" {
				azure["accountKeySecret"] = map[string]interface{}{
					"name": backend.Azure.GetCredentialsSecretName(),
					"key":  vars.AzureAccountKeyKey,
				}
			} else {
				// Keyless: managed identity / workload identity when no
				// account key is declared.
				azure["useSDKCreds"] = true
			}
			repository["azure"] = azure
		}
		if len(repository) > 0 {
			values["artifactRepository"] = repository
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ----------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// fullnameOverride re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). Every child name — and
	// the exported outputs built from them — derive from the fullname;
	// letting an override move it would break every output.
	values["fullnameOverride"] = locals.ReleaseName

	return values, nil
}
