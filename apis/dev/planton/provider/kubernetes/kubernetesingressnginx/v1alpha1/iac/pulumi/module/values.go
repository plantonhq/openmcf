package module

import (
	"strings"

	"github.com/pkg/errors"
	kubernetesingressnginxv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesingressnginx/v1alpha1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	// ---- controller block ------------------------------------------------
	controller := map[string]interface{}{}

	// The IngressClass this instance owns. `enabled` is the chart default
	// (true) — the class is the instance's identity, always created.
	ingressClassResource := map[string]interface{}{
		"name":            locals.IngressClassName,
		"controllerValue": locals.IngressClassControllerValue,
	}
	if spec.GetIngressClass().GetIsDefaultClass() {
		ingressClassResource["default"] = true
	}
	controller["ingressClassResource"] = ingressClassResource
	// The legacy ingress.class-annotation vocabulary must track the class
	// name too — the chart defaults it to "nginx" independently of the
	// IngressClass resource name, which would mis-route annotation-based
	// Ingresses on a non-default-named instance.
	controller["ingressClass"] = locals.IngressClassName
	if spec.GetIngressClass().GetWatchIngressWithoutClass() {
		controller["watchIngressWithoutClass"] = true
	}

	// Replicas vs autoscaling: when the HPA owns the count, do not also
	// pin replicaCount — the chart's Deployment template omits replicas
	// under autoscaling, avoiding a rollout tug-of-war.
	if a := spec.GetAutoscaling(); a.GetEnabled() {
		autoscaling := map[string]interface{}{"enabled": true}
		if a.MinReplicas != nil {
			autoscaling["minReplicas"] = a.GetMinReplicas()
		}
		if a.MaxReplicas != nil {
			autoscaling["maxReplicas"] = a.GetMaxReplicas()
		}
		if a.TargetCpuUtilizationPercent != nil {
			autoscaling["targetCPUUtilizationPercentage"] = a.GetTargetCpuUtilizationPercent()
		}
		if a.TargetMemoryUtilizationPercent != nil {
			autoscaling["targetMemoryUtilizationPercentage"] = a.GetTargetMemoryUtilizationPercent()
		}
		controller["autoscaling"] = autoscaling
	} else if spec.Replicas != nil {
		controller["replicaCount"] = spec.GetReplicas()
	}

	if r := resourcesMap(spec.GetResources()); r != nil {
		controller["resources"] = r
	}

	// ---- controller service ---------------------------------------------
	service := map[string]interface{}{}
	switch spec.GetService().GetType() {
	case kubernetesingressnginxv1alpha1.KubernetesIngressNginxServiceType_node_port:
		service["type"] = "NodePort"
	case kubernetesingressnginxv1alpha1.KubernetesIngressNginxServiceType_cluster_ip:
		service["type"] = "ClusterIP"
	default:
		// LoadBalancer is the chart default; set it explicitly so the
		// rendered values are self-describing.
		service["type"] = "LoadBalancer"
	}
	if annotations := spec.GetService().GetAnnotations(); len(annotations) > 0 {
		service["annotations"] = stringMapToInterface(annotations)
	}
	if spec.GetService().GetExternalTrafficPolicy() == kubernetesingressnginxv1alpha1.KubernetesIngressNginxExternalTrafficPolicy_local {
		service["externalTrafficPolicy"] = "Local"
	}
	if ranges := spec.GetService().GetLoadBalancerSourceRanges(); len(ranges) > 0 {
		service["loadBalancerSourceRanges"] = toInterfaceSlice(ranges)
	}
	if spec.GetService().GetLoadBalancerClass() != "" {
		service["loadBalancerClass"] = spec.GetService().GetLoadBalancerClass()
	}
	if spec.GetService().EnableHttp != nil && !spec.GetService().GetEnableHttp() {
		service["enableHttp"] = false
	}
	if spec.GetService().EnableHttps != nil && !spec.GetService().GetEnableHttps() {
		service["enableHttps"] = false
	}
	nodePorts := map[string]interface{}{}
	if spec.GetService().HttpNodePort != nil {
		nodePorts["http"] = spec.GetService().GetHttpNodePort()
	}
	if spec.GetService().HttpsNodePort != nil {
		nodePorts["https"] = spec.GetService().GetHttpsNodePort()
	}
	if len(nodePorts) > 0 {
		service["nodePorts"] = nodePorts
	}
	if internal := spec.GetService().GetInternal(); internal.GetEnabled() {
		service["internal"] = map[string]interface{}{
			"enabled":     true,
			"annotations": stringMapToInterface(internal.GetAnnotations()),
		}
	}
	controller["service"] = service

	// ---- workload shape ---------------------------------------------------
	if spec.GetControllerKind() == kubernetesingressnginxv1alpha1.KubernetesIngressNginxControllerKind_daemon_set {
		controller["kind"] = "DaemonSet"
	}
	if spec.GetHostNetwork() {
		controller["hostNetwork"] = true
		// Keep in-cluster name resolution working from the host network —
		// without this the controller resolves through the node's DNS and
		// cannot see cluster Services.
		controller["dnsPolicy"] = "ClusterFirstWithHostNet"
	}
	if spec.GetHostPorts() {
		controller["hostPort"] = map[string]interface{}{"enabled": true}
	}

	// ---- NGINX behavior ----------------------------------------------------
	if len(spec.GetNginxConfig()) > 0 {
		controller["config"] = stringMapToInterface(spec.GetNginxConfig())
	}
	if spec.GetAllowSnippetAnnotations() {
		controller["allowSnippetAnnotations"] = true
	}
	if tls := spec.GetDefaultTlsCertificate(); tls.GetSecretName().GetValue() != "" {
		ns := tls.GetNamespace()
		if ns == "" {
			ns = locals.Namespace
		}
		// The chart exposes no first-class key for the default certificate;
		// upstream's own documented mechanism is the controller flag.
		controller["extraArgs"] = map[string]interface{}{
			"default-ssl-certificate": ns + "/" + tls.GetSecretName().GetValue(),
		}
	}

	// ---- admission webhook -------------------------------------------------
	if aw := spec.GetAdmissionWebhooks(); aw != nil {
		admission := map[string]interface{}{}
		if aw.Enabled != nil && !aw.GetEnabled() {
			admission["enabled"] = false
		}
		if aw.GetFailurePolicy() == kubernetesingressnginxv1alpha1.KubernetesIngressNginxWebhookFailurePolicy_ignore {
			admission["failurePolicy"] = "Ignore"
		}
		if aw.TimeoutSeconds != nil {
			admission["timeoutSeconds"] = aw.GetTimeoutSeconds()
		}
		if len(admission) > 0 {
			controller["admissionWebhooks"] = admission
		}
	}

	// ---- observability -----------------------------------------------------
	if m := spec.GetMetrics(); m.GetEnabled() {
		metrics := map[string]interface{}{"enabled": true}
		if m.GetServiceMonitor() {
			sm := map[string]interface{}{"enabled": true}
			if m.GetServiceMonitorInterval() != "" {
				sm["scrapeInterval"] = m.GetServiceMonitorInterval()
			}
			if len(m.GetServiceMonitorLabels()) > 0 {
				sm["additionalLabels"] = stringMapToInterface(m.GetServiceMonitorLabels())
			}
			metrics["serviceMonitor"] = sm
		}
		controller["metrics"] = metrics
	}

	// ---- scheduling ----------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		controller["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		controller["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}
	if spec.GetPriorityClassName() != "" {
		controller["priorityClassName"] = spec.GetPriorityClassName()
	}

	// ---- assemble root -------------------------------------------------------
	values := map[string]interface{}{
		// Pin the chart's fullname to the release name (= metadata.name):
		// every chart object then carries a deterministic, manifest-derived
		// name — what verification, imports, and multi-instance coexistence
		// (including per-instance leader election, "<fullname>-leader") all
		// key off.
		"fullnameOverride": locals.ReleaseName,
		"controller":       controller,
	}

	if spec.GetImageRegistry() != "" {
		values["global"] = map[string]interface{}{
			"image": map[string]interface{}{"registry": spec.GetImageRegistry()},
		}
	}

	if db := spec.GetDefaultBackend(); db.GetEnabled() {
		defaultBackend := map[string]interface{}{"enabled": true}
		if db.Replicas != nil {
			defaultBackend["replicaCount"] = db.GetReplicas()
		}
		if db.GetImage() != "" {
			// Spec carries "repository:tag"; the chart wants them split.
			image := map[string]interface{}{}
			if idx := strings.LastIndex(db.GetImage(), ":"); idx > 0 {
				image["repository"] = db.GetImage()[:idx]
				image["tag"] = db.GetImage()[idx+1:]
			} else {
				image["repository"] = db.GetImage()
			}
			defaultBackend["image"] = image
		}
		if r := resourcesMap(db.GetResources()); r != nil {
			defaultBackend["resources"] = r
		}
		values["defaultBackend"] = defaultBackend
	}

	if len(spec.GetTcpServices()) > 0 {
		values["tcp"] = stringMapToInterface(spec.GetTcpServices())
	}
	if len(spec.GetUdpServices()) > 0 {
		values["udp"] = stringMapToInterface(spec.GetUdpServices())
	}

	// ---- escape hatch (merged LAST, helm -f semantics) -----------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}
