package module

import (
	"sort"
	"strconv"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kuberneteskeycloakv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskeycloak/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// keycloakCR renders the single k8s.keycloak.org/v2beta1 Keycloak CR the
// Keycloak Operator (a KubernetesKeycloakOperator install, the
// prerequisite) reconciles into the StatefulSet, `<name>-service`,
// `<name>-discovery`, the NetworkPolicy and — unless the spec brings its
// own bootstrap-admin Secret — the create-once `<name>-initial-admin`
// credential Secret.
//
// UNTYPED CR: the body assembles in ONE place (keycloakSpecBody) through
// apiextensions.CustomResource with an untyped spec map, in byte lockstep
// with the Terraform twin's locals (locals.tf `keycloak_spec`). Every key
// renders ONLY when the spec declares it, so the operator's defaulting
// stays authoritative for everything the manifest leaves unsaid — except
// networkPolicy/serviceMonitor/ingress, deliberately always rendered (see
// the builder).
//
// NO WAIT on the CR, deliberately: server readiness depends on the
// operator (image pulls, database schema migrations, cluster formation)
// that is not part of applying the resource — the verifier owns
// readiness, the same never-block-on-a-controller posture as the sibling
// operator-CR modules.
func keycloakCR(ctx *pulumi.Context, locals *Locals,
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
				// Keycloak CR's cascade — its finalizer tears down the
				// StatefulSet, Services and NetworkPolicy. Foreground
				// propagation would block the delete on children the
				// operator keeps reconciling. Terraform twin:
				// delete_cascade = "Background" on kubectl_manifest.
				Annotations: pulumi.StringMap{
					"pulumi.com/deletionPropagationPolicy": pulumi.String("background"),
				},
			},
			OtherFields: map[string]interface{}{
				"spec": keycloakSpecBody(locals),
			},
		},
		append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return errors.Wrap(err, "failed to render the Keycloak custom resource")
	}
	return nil
}

// keycloakSpecBody renders the Keycloak CR spec from the typed fields.
// Field names are the CRD's own JSON keys (verified against the pinned
// v2beta1 schema at operator 26.7.0). The Terraform twin (locals.tf
// `keycloak_spec`) renders the identical body — keep them in byte
// lockstep.
func keycloakSpecBody(locals *Locals) map[string]interface{} {
	spec := locals.Spec

	body := map[string]interface{}{}

	if spec.Instances != nil {
		body["instances"] = int(spec.GetInstances())
	}
	if spec.GetImage() != "" {
		body["image"] = spec.GetImage()
	}
	if spec.GetStartOptimized() {
		body["startOptimized"] = true
	}

	body["db"] = dbBody(spec.GetDb())

	if http := httpBody(locals, spec.GetHttp()); len(http) > 0 {
		body["http"] = http
	}
	if hostname := hostnameBody(spec.GetHostname()); len(hostname) > 0 {
		body["hostname"] = hostname
	}
	if spec.GetProxyHeaders() != "" {
		body["proxy"] = map[string]interface{}{"headers": spec.GetProxyHeaders()}
	}
	if features := featuresBody(spec.GetFeatures()); len(features) > 0 {
		body["features"] = features
	}
	if spec.GetTransactionXaEnabled() {
		body["transaction"] = map[string]interface{}{"xaEnabled": true}
	}
	if cacheConfig := spec.GetCacheConfig(); cacheConfig != nil {
		configMapFile := map[string]interface{}{
			"name": cacheConfig.GetConfigMapName(),
		}
		if cacheConfig.Key != nil {
			configMapFile["key"] = cacheConfig.GetKey()
		}
		body["cache"] = map[string]interface{}{"configMapFile": configMapFile}
	}
	// truststores is a CRD map keyed by an arbitrary handle — the module
	// keys each entry by its own Secret name.
	if truststoreSecretNames := spec.GetTruststoreSecretNames(); len(truststoreSecretNames) > 0 {
		truststores := map[string]interface{}{}
		for _, secretName := range truststoreSecretNames {
			truststores[secretName] = map[string]interface{}{
				"secret": map[string]interface{}{"name": secretName},
			}
		}
		body["truststores"] = truststores
	}
	if additionalOptions := additionalOptionsBody(spec.GetAdditionalOptions()); additionalOptions != nil {
		body["additionalOptions"] = additionalOptions
	}
	if spec.GetBootstrapAdminSecretName() != "" {
		body["bootstrapAdmin"] = map[string]interface{}{
			"user": map[string]interface{}{"secret": spec.GetBootstrapAdminSecretName()},
		}
	}
	if resources := resourcesBody(spec.GetResources()); len(resources) > 0 {
		body["resources"] = resources
	}
	if scheduling := schedulingBody(spec.GetScheduling()); len(scheduling) > 0 {
		body["scheduling"] = scheduling
	}

	if probes := spec.GetProbes(); probes != nil {
		if probe := probeBody(probes.LivenessFailureThreshold, probes.LivenessPeriodSeconds); len(probe) > 0 {
			body["livenessProbe"] = probe
		}
		if probe := probeBody(probes.ReadinessFailureThreshold, probes.ReadinessPeriodSeconds); len(probe) > 0 {
			body["readinessProbe"] = probe
		}
		if probe := probeBody(probes.StartupFailureThreshold, probes.StartupPeriodSeconds); len(probe) > 0 {
			body["startupProbe"] = probe
		}
	}

	if spec.HttpManagementPort != nil {
		body["httpManagement"] = map[string]interface{}{"port": int(spec.GetHttpManagementPort())}
	}

	// networkPolicy/serviceMonitor render ALWAYS, with the spec's
	// effective value: the operator defaults both to true, and an honest
	// declaration states the value it relies on instead of leaning on a
	// default that can change under it.
	networkPolicyEnabled := true
	if spec.NetworkPolicyEnabled != nil {
		networkPolicyEnabled = spec.GetNetworkPolicyEnabled()
	}
	body["networkPolicy"] = map[string]interface{}{"enabled": networkPolicyEnabled}

	serviceMonitorEnabled := true
	if spec.ServiceMonitorEnabled != nil {
		serviceMonitorEnabled = spec.GetServiceMonitorEnabled()
	}
	body["serviceMonitor"] = map[string]interface{}{"enabled": serviceMonitorEnabled}

	if update := updateBody(spec.GetUpdate()); len(update) > 0 {
		body["update"] = update
	}
	if tracing := tracingBody(spec.GetTracing()); len(tracing) > 0 {
		body["tracing"] = tracing
	}

	// ALWAYS disable the operator's default Ingress — an ABSENT ingress
	// block means ENABLED (verified in operator source: absence defaults
	// the Ingress on), so the block must render explicitly. Exposure
	// composes from Gateway API kinds referencing the exported service
	// handles, never from this component.
	body["ingress"] = map[string]interface{}{"enabled": false}

	return body
}

// dbBody renders spec.db. The dev-file/dev-mem SANDBOX vendors run
// Keycloak's embedded H2 on each pod's own ephemeral storage — no
// connection details apply, so ONLY the vendor renders; everything else
// would be dead configuration the CR should not carry.
func dbBody(db *kuberneteskeycloakv1.KubernetesKeycloakDb) map[string]interface{} {
	out := map[string]interface{}{
		"vendor": db.GetVendor(),
	}
	if db.GetVendor() == "dev-file" || db.GetVendor() == "dev-mem" {
		return out
	}

	if db.GetHost().GetValue() != "" {
		out["host"] = db.GetHost().GetValue()
	}
	if db.Port != nil {
		out["port"] = int(db.GetPort())
	}
	if db.GetDatabase() != "" {
		out["database"] = db.GetDatabase()
	}
	if db.GetSchema() != "" {
		out["schema"] = db.GetSchema()
	}
	// The CRD calls the JDBC URL override `url`; when set the server
	// ignores host/port/database.
	if db.GetJdbcUrl() != "" {
		out["url"] = db.GetJdbcUrl()
	}
	if db.PoolMinSize != nil {
		out["poolMinSize"] = int(db.GetPoolMinSize())
	}
	if db.PoolMaxSize != nil {
		out["poolMaxSize"] = int(db.GetPoolMaxSize())
	}
	if usernameSecret := secretSelectorBody(db.GetUsernameSecret()); usernameSecret != nil {
		out["usernameSecret"] = usernameSecret
	}
	if passwordSecret := secretSelectorBody(db.GetPasswordSecret()); passwordSecret != nil {
		out["passwordSecret"] = passwordSecret
	}
	return out
}

func secretSelectorBody(selector *kuberneteskeycloakv1.KubernetesKeycloakSecretSelector) map[string]interface{} {
	if selector == nil {
		return nil
	}
	return map[string]interface{}{
		"name": selector.GetName().GetValue(),
		"key":  selector.GetKey(),
	}
}

// httpBody renders spec.http — the TLS-or-HTTP posture the server
// REFUSES TO START without (upstream surfaces the misconfiguration only
// as a CrashLoopBackOff; the spec's validation rule catches it first).
func httpBody(locals *Locals, http *kuberneteskeycloakv1.KubernetesKeycloakHttp) map[string]interface{} {
	if http == nil {
		return nil
	}
	out := map[string]interface{}{}
	if locals.TlsSecretName != "" {
		out["tlsSecret"] = locals.TlsSecretName
	}
	if http.GetHttpEnabled() {
		out["httpEnabled"] = true
	}
	if http.HttpPort != nil {
		out["httpPort"] = int(http.GetHttpPort())
	}
	if http.HttpsPort != nil {
		out["httpsPort"] = int(http.GetHttpsPort())
	}
	return out
}

func hostnameBody(hostname *kuberneteskeycloakv1.KubernetesKeycloakHostname) map[string]interface{} {
	if hostname == nil {
		return nil
	}
	out := map[string]interface{}{}
	if hostname.GetHostname() != "" {
		out["hostname"] = hostname.GetHostname()
	}
	if hostname.GetAdmin() != "" {
		out["admin"] = hostname.GetAdmin()
	}
	// strict is tri-state: unset leaves the server default (true);
	// declared renders either way — `strict: false` is the meaningful
	// behind-a-proxy posture.
	if hostname.Strict != nil {
		out["strict"] = hostname.GetStrict()
	}
	if hostname.GetBackchannelDynamic() {
		out["backchannelDynamic"] = true
	}
	return out
}

func featuresBody(features *kuberneteskeycloakv1.KubernetesKeycloakFeatures) map[string]interface{} {
	if features == nil {
		return nil
	}
	out := map[string]interface{}{}
	if len(features.GetEnabled()) > 0 {
		out["enabled"] = stringSliceToInterface(features.GetEnabled())
	}
	if len(features.GetDisabled()) > 0 {
		out["disabled"] = stringSliceToInterface(features.GetDisabled())
	}
	return out
}

func additionalOptionsBody(options []*kuberneteskeycloakv1.KubernetesKeycloakAdditionalOption) []interface{} {
	if len(options) == 0 {
		return nil
	}
	out := []interface{}{}
	for _, option := range options {
		entry := map[string]interface{}{
			"name": option.GetName(),
		}
		// The spec validates value XOR secret — exactly one arm renders.
		if option.GetValue() != "" {
			entry["value"] = option.GetValue()
		}
		if secret := secretSelectorBody(option.GetSecret()); secret != nil {
			entry["secret"] = secret
		}
		out = append(out, entry)
	}
	return out
}

func resourcesBody(resources *kubernetesprovider.ContainerResources) map[string]interface{} {
	if resources == nil {
		return nil
	}
	out := map[string]interface{}{}
	if requests := cpuMemoryBody(resources.GetRequests()); len(requests) > 0 {
		out["requests"] = requests
	}
	if limits := cpuMemoryBody(resources.GetLimits()); len(limits) > 0 {
		out["limits"] = limits
	}
	return out
}

func cpuMemoryBody(cpuMemory *kubernetesprovider.CpuMemory) map[string]interface{} {
	if cpuMemory == nil {
		return nil
	}
	out := map[string]interface{}{}
	if cpuMemory.GetCpu() != "" {
		out["cpu"] = cpuMemory.GetCpu()
	}
	if cpuMemory.GetMemory() != "" {
		out["memory"] = cpuMemory.GetMemory()
	}
	return out
}

// schedulingBody renders spec.scheduling. The Keycloak CR models affinity
// rather than nodeSelector, so node_selector translates to REQUIRED node
// affinity — one In-expression per label, the catalog's established
// translation. Keys sort so both engines render the same expression order.
func schedulingBody(scheduling *kuberneteskeycloakv1.KubernetesKeycloakScheduling) map[string]interface{} {
	if scheduling == nil {
		return nil
	}
	out := map[string]interface{}{}

	if nodeSelector := scheduling.GetNodeSelector(); len(nodeSelector) > 0 {
		keys := make([]string, 0, len(nodeSelector))
		for key := range nodeSelector {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		matchExpressions := []interface{}{}
		for _, key := range keys {
			matchExpressions = append(matchExpressions, map[string]interface{}{
				"key":      key,
				"operator": "In",
				"values":   []interface{}{nodeSelector[key]},
			})
		}
		out["affinity"] = map[string]interface{}{
			"nodeAffinity": map[string]interface{}{
				"requiredDuringSchedulingIgnoredDuringExecution": map[string]interface{}{
					"nodeSelectorTerms": []interface{}{
						map[string]interface{}{"matchExpressions": matchExpressions},
					},
				},
			},
		}
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
		out["tolerations"] = entries
	}

	if scheduling.GetPriorityClassName() != "" {
		out["priorityClassName"] = scheduling.GetPriorityClassName()
	}
	return out
}

func probeBody(failureThreshold *int32, periodSeconds *int32) map[string]interface{} {
	out := map[string]interface{}{}
	if failureThreshold != nil {
		out["failureThreshold"] = int(*failureThreshold)
	}
	if periodSeconds != nil {
		out["periodSeconds"] = int(*periodSeconds)
	}
	return out
}

// updateBody renders spec.update. KNOW THIS: with the default
// RecreateOnImageChange strategy, changing the image takes a full
// scale-to-zero recreate — an outage window by design (two Keycloak
// versions cannot share one cache cluster/schema).
func updateBody(update *kuberneteskeycloakv1.KubernetesKeycloakUpdate) map[string]interface{} {
	if update == nil {
		return nil
	}
	out := map[string]interface{}{}
	if update.Strategy != nil {
		out["strategy"] = update.GetStrategy()
	}
	if update.GetRevision() != "" {
		out["revision"] = update.GetRevision()
	}
	return out
}

func tracingBody(tracing *kuberneteskeycloakv1.KubernetesKeycloakTracing) map[string]interface{} {
	if tracing == nil {
		return nil
	}
	out := map[string]interface{}{}
	if tracing.GetEnabled() {
		out["enabled"] = true
	}
	if tracing.GetEndpoint() != "" {
		out["endpoint"] = tracing.GetEndpoint()
	}
	if tracing.Protocol != nil {
		out["protocol"] = tracing.GetProtocol()
	}
	// The CRD types samplerRatio as a NUMBER; the spec carries it as a
	// pattern-validated string (proto3 has no optional double with
	// presence) — parse here so the CR renders what the schema wants.
	if tracing.SamplerRatio != nil {
		if ratio, err := strconv.ParseFloat(tracing.GetSamplerRatio(), 64); err == nil {
			out["samplerRatio"] = ratio
		}
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
