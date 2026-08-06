package module

import (
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kuberneteskafkav1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskafka/v1alpha1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createKafkaCluster renders the kafka.strimzi.io/v1 Kafka resource as an
// UNTYPED CustomResource. The typed crd2pulumi tree is structurally unable
// to carry this CR: the CRD types spec.kafka.config (and the listener and
// Cruise Control configuration blocks) with
// x-kubernetes-preserve-unknown-fields, which the generated types flatten
// into shapes that cannot hold the free-typed bodies — so no generated
// package is shipped for the Kafka family at all.
//
// The spec body built here is the exact twin of the Terraform module's
// local.kafka_manifest (locals.tf) — same keys rendered and omitted,
// numbers as ints, booleans as booleans. An unset optional is simply never
// inserted into the map (the Go twin of TF's null-prune), so the apiserver
// applies the CRD's own defaults. Shape errors still fail loudly without
// compile-time typing: the operator validates the applied spec against its
// schema, and the kind-cluster E2E lanes exercise the rendered arms live.
//
// No await machinery, deliberately: cluster readiness depends on the
// operator (image pulls, KRaft quorum formation, listener provisioning)
// that is not part of applying the resource — the same
// never-block-on-a-controller posture as the sibling database kinds.
func createKafkaCluster(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	spec := locals.Spec

	kafkaBody := map[string]interface{}{
		"listeners": listenersBody(spec.GetListeners()),
	}
	if spec.GetKafkaVersion() != "" {
		kafkaBody["version"] = spec.GetKafkaVersion()
	}
	if spec.GetMetadataVersion() != "" {
		kafkaBody["metadataVersion"] = spec.GetMetadataVersion()
	}
	if len(spec.GetConfig()) > 0 {
		kafkaBody["config"] = stringMapToInterface(spec.GetConfig())
	}
	if authorization := authorizationBody(spec.GetAuthorization()); authorization != nil {
		kafkaBody["authorization"] = authorization
	}
	if spec.GetRack().GetTopologyKey() != "" {
		kafkaBody["rack"] = map[string]interface{}{
			"topologyKey": spec.GetRack().GetTopologyKey(),
		}
	}
	if jvm := jvmBody(spec.GetJvm()); jvm != nil {
		kafkaBody["jvmOptions"] = jvm
	}
	if spec.GetMetrics().GetEnabled() {
		// The module owns the rules ConfigMap (metrics.go); the CR only
		// points at it.
		kafkaBody["metricsConfig"] = map[string]interface{}{
			"type": "jmxPrometheusExporter",
			"valueFrom": map[string]interface{}{
				"configMapKeyRef": map[string]interface{}{
					"name": locals.MetricsConfigMapName,
					"key":  vars.MetricsConfigKey,
				},
			},
		}
	}

	specBody := map[string]interface{}{
		"kafka": kafkaBody,
	}
	if entityOperator := entityOperatorBody(spec.GetEntityOperator()); entityOperator != nil {
		specBody["entityOperator"] = entityOperator
	}
	if cruiseControl := cruiseControlBody(spec.GetCruiseControl()); cruiseControl != nil {
		specBody["cruiseControl"] = cruiseControl
	}
	if exporter := kafkaExporterBody(spec.GetKafkaExporter()); exporter != nil {
		specBody["kafkaExporter"] = exporter
	}
	if ca := caBody(spec.GetClusterCa()); ca != nil {
		specBody["clusterCa"] = ca
	}
	if ca := caBody(spec.GetClientsCa()); ca != nil {
		specBody["clientsCa"] = ca
	}
	if len(spec.GetMaintenanceTimeWindows()) > 0 {
		specBody["maintenanceTimeWindows"] = stringSliceToInterface(spec.GetMaintenanceTimeWindows())
	}

	return apiextensions.NewCustomResource(ctx, locals.ClusterName,
		&apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String(vars.ApiVersion),
			Kind:       pulumi.String("Kafka"),
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ClusterName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": specBody,
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}

// listenersBody is the twin of TF's listeners rendering: every listener
// carries the CRD-required quartet (name/port/type/tls); authentication
// and configuration render only when declared.
func listenersBody(listeners []*kuberneteskafkav1alpha1.KubernetesKafkaListener) []interface{} {
	out := make([]interface{}, 0, len(listeners))
	for _, listener := range listeners {
		body := map[string]interface{}{
			"name": listener.GetName(),
			"port": int(listener.GetPort()),
			"type": listenerType(listener),
			"tls":  listener.GetTls(),
		}
		if auth := listenerAuthenticationBody(listener.GetAuthentication()); auth != nil {
			body["authentication"] = auth
		}
		if configuration := listenerConfigurationBody(listener.GetConfiguration()); configuration != nil {
			body["configuration"] = configuration
		}
		out = append(out, body)
	}
	return out
}

// listenerType mirrors TF's coalesce(type, "internal"): the CRD requires
// an explicit type on every listener, and the spec defaults to internal.
func listenerType(listener *kuberneteskafkav1alpha1.KubernetesKafkaListener) string {
	if listener.GetType() != "" {
		return listener.GetType()
	}
	return "internal"
}

func listenerAuthenticationBody(auth *kuberneteskafkav1alpha1.KubernetesKafkaListenerAuthentication) map[string]interface{} {
	if auth == nil {
		return nil
	}
	body := map[string]interface{}{
		"type": auth.GetType(),
	}
	if auth.GetType() == "custom" {
		body["sasl"] = auth.GetSasl()
		if len(auth.GetListenerConfig()) > 0 {
			body["listenerConfig"] = stringMapToInterface(auth.GetListenerConfig())
		}
	}
	return body
}

func listenerConfigurationBody(configuration *kuberneteskafkav1alpha1.KubernetesKafkaListenerConfiguration) map[string]interface{} {
	if configuration == nil {
		return nil
	}
	body := map[string]interface{}{}
	if configuration.GetClass() != "" {
		body["class"] = configuration.GetClass()
	}
	if configuration.ExternalTrafficPolicy != nil && configuration.GetExternalTrafficPolicy() != "" {
		body["externalTrafficPolicy"] = configuration.GetExternalTrafficPolicy()
	}
	if len(configuration.GetLoadBalancerSourceRanges()) > 0 {
		body["loadBalancerSourceRanges"] = stringSliceToInterface(configuration.GetLoadBalancerSourceRanges())
	}
	if configuration.AllocateLoadBalancerNodePorts != nil {
		body["allocateLoadBalancerNodePorts"] = configuration.GetAllocateLoadBalancerNodePorts()
	}
	if configuration.CreateBootstrapService != nil {
		body["createBootstrapService"] = configuration.GetCreateBootstrapService()
	}
	if configuration.GetUseServiceDnsDomain() {
		body["useServiceDnsDomain"] = true
	}
	if configuration.MaxConnections != nil {
		body["maxConnections"] = int(configuration.GetMaxConnections())
	}
	if configuration.MaxConnectionCreationRate != nil {
		body["maxConnectionCreationRate"] = int(configuration.GetMaxConnectionCreationRate())
	}
	if configuration.PreferredNodePortAddressType != nil && configuration.GetPreferredNodePortAddressType() != "" {
		body["preferredNodePortAddressType"] = configuration.GetPreferredNodePortAddressType()
	}
	if configuration.GetPublishNotReadyAddresses() {
		body["publishNotReadyAddresses"] = true
	}
	if cert := configuration.GetBrokerCertChainAndKey(); cert != nil {
		// cert-manager writes tls.crt/tls.key; the spec defaults mirror
		// that, so unset keys resolve through the getters' defaults.
		certificate := cert.GetCertificate()
		if certificate == "" {
			certificate = "tls.crt"
		}
		key := cert.GetKey()
		if key == "" {
			key = "tls.key"
		}
		body["brokerCertChainAndKey"] = map[string]interface{}{
			"secretName":  cert.GetSecretName().GetValue(),
			"certificate": certificate,
			"key":         key,
		}
	}
	if bootstrap := listenerBootstrapBody(configuration.GetBootstrap()); bootstrap != nil {
		body["bootstrap"] = bootstrap
	}
	if brokers := listenerBrokersBody(configuration.GetBrokers()); len(brokers) > 0 {
		body["brokers"] = brokers
	}
	if len(body) == 0 {
		return nil
	}
	return body
}

func listenerBootstrapBody(bootstrap *kuberneteskafkav1alpha1.KubernetesKafkaListenerBootstrap) map[string]interface{} {
	if bootstrap == nil {
		return nil
	}
	body := map[string]interface{}{}
	if bootstrap.GetHost() != "" {
		body["host"] = bootstrap.GetHost()
	}
	if len(bootstrap.GetAnnotations()) > 0 {
		body["annotations"] = stringMapToInterface(bootstrap.GetAnnotations())
	}
	if len(bootstrap.GetLabels()) > 0 {
		body["labels"] = stringMapToInterface(bootstrap.GetLabels())
	}
	if bootstrap.GetLoadBalancerIp() != "" {
		body["loadBalancerIP"] = bootstrap.GetLoadBalancerIp()
	}
	if bootstrap.NodePort != nil {
		body["nodePort"] = int(bootstrap.GetNodePort())
	}
	if len(bootstrap.GetAlternativeNames()) > 0 {
		body["alternativeNames"] = stringSliceToInterface(bootstrap.GetAlternativeNames())
	}
	if len(body) == 0 {
		return nil
	}
	return body
}

func listenerBrokersBody(brokers []*kuberneteskafkav1alpha1.KubernetesKafkaListenerBroker) []interface{} {
	out := make([]interface{}, 0, len(brokers))
	for _, broker := range brokers {
		body := map[string]interface{}{
			"broker": int(broker.GetBroker()),
		}
		if broker.GetHost() != "" {
			body["host"] = broker.GetHost()
		}
		if broker.GetAdvertisedHost() != "" {
			body["advertisedHost"] = broker.GetAdvertisedHost()
		}
		if broker.AdvertisedPort != nil {
			body["advertisedPort"] = int(broker.GetAdvertisedPort())
		}
		if len(broker.GetAnnotations()) > 0 {
			body["annotations"] = stringMapToInterface(broker.GetAnnotations())
		}
		if len(broker.GetLabels()) > 0 {
			body["labels"] = stringMapToInterface(broker.GetLabels())
		}
		if broker.GetLoadBalancerIp() != "" {
			body["loadBalancerIP"] = broker.GetLoadBalancerIp()
		}
		if broker.NodePort != nil {
			body["nodePort"] = int(broker.GetNodePort())
		}
		out = append(out, body)
	}
	return out
}

func authorizationBody(authorization *kuberneteskafkav1alpha1.KubernetesKafkaAuthorization) map[string]interface{} {
	if authorization == nil {
		return nil
	}
	body := map[string]interface{}{
		"type": authorization.GetType(),
	}
	if len(authorization.GetSuperUsers()) > 0 {
		body["superUsers"] = stringSliceToInterface(authorization.GetSuperUsers())
	}
	if authorization.GetType() == "custom" {
		body["authorizerClass"] = authorization.GetAuthorizerClass()
		if authorization.GetSupportsAdminApi() {
			body["supportsAdminApi"] = true
		}
	}
	return body
}

// entityOperatorBody renders the entity operator with each sub-operator
// present when enabled (the spec defaults both true). When BOTH are
// disabled the block is omitted entirely — Strimzi deploys no entity
// operator pod, and KafkaTopic/KafkaUser declarations for this cluster
// become inert (the spec comments warn about this).
func entityOperatorBody(entityOperator *kuberneteskafkav1alpha1.KubernetesKafkaEntityOperator) map[string]interface{} {
	topicEnabled := true
	userEnabled := true
	if entityOperator != nil {
		if entityOperator.TopicOperatorEnabled != nil {
			topicEnabled = entityOperator.GetTopicOperatorEnabled()
		}
		if entityOperator.UserOperatorEnabled != nil {
			userEnabled = entityOperator.GetUserOperatorEnabled()
		}
	}
	if !topicEnabled && !userEnabled {
		return nil
	}
	body := map[string]interface{}{}
	if topicEnabled {
		body["topicOperator"] = map[string]interface{}{}
	}
	if userEnabled {
		body["userOperator"] = map[string]interface{}{}
	}
	return body
}

func cruiseControlBody(cruiseControl *kuberneteskafkav1alpha1.KubernetesKafkaCruiseControl) map[string]interface{} {
	if cruiseControl == nil || !cruiseControl.GetEnabled() {
		return nil
	}
	body := map[string]interface{}{}
	if len(cruiseControl.GetConfig()) > 0 {
		body["config"] = stringMapToInterface(cruiseControl.GetConfig())
	}
	if resources := resourcesMap(cruiseControl.GetResources()); resources != nil {
		body["resources"] = resources
	}
	if len(cruiseControl.GetAutoRebalanceModes()) > 0 {
		modes := make([]interface{}, 0, len(cruiseControl.GetAutoRebalanceModes()))
		for _, mode := range cruiseControl.GetAutoRebalanceModes() {
			modes = append(modes, map[string]interface{}{"mode": mode})
		}
		body["autoRebalance"] = modes
	}
	return body
}

func kafkaExporterBody(exporter *kuberneteskafkav1alpha1.KubernetesKafkaExporter) map[string]interface{} {
	if exporter == nil || !exporter.GetEnabled() {
		return nil
	}
	body := map[string]interface{}{}
	if exporter.GetGroupRegex() != "" {
		body["groupRegex"] = exporter.GetGroupRegex()
	}
	if exporter.GetTopicRegex() != "" {
		body["topicRegex"] = exporter.GetTopicRegex()
	}
	if resources := resourcesMap(exporter.GetResources()); resources != nil {
		body["resources"] = resources
	}
	return body
}

func caBody(ca *kuberneteskafkav1alpha1.KubernetesKafkaCa) map[string]interface{} {
	if ca == nil {
		return nil
	}
	body := map[string]interface{}{}
	if ca.ValidityDays != nil {
		body["validityDays"] = int(ca.GetValidityDays())
	}
	if ca.RenewalDays != nil {
		body["renewalDays"] = int(ca.GetRenewalDays())
	}
	if len(body) == 0 {
		return nil
	}
	return body
}

func jvmBody(jvm *kuberneteskafkav1alpha1.KubernetesKafkaJvm) map[string]interface{} {
	if jvm == nil {
		return nil
	}
	body := map[string]interface{}{}
	if jvm.GetXms() != "" {
		body["-Xms"] = jvm.GetXms()
	}
	if jvm.GetXmx() != "" {
		body["-Xmx"] = jvm.GetXmx()
	}
	if len(body) == 0 {
		return nil
	}
	return body
}

// resourcesMap renders the shared ContainerResources message into the
// standard Kubernetes limits/requests shape. Returns nil when nothing is
// set.
func resourcesMap(r *kubernetesprovider.ContainerResources) map[string]interface{} {
	if r == nil {
		return nil
	}
	out := map[string]interface{}{}
	if l := r.GetLimits(); l != nil && (l.GetCpu() != "" || l.GetMemory() != "") {
		limits := map[string]interface{}{}
		if l.GetCpu() != "" {
			limits["cpu"] = l.GetCpu()
		}
		if l.GetMemory() != "" {
			limits["memory"] = l.GetMemory()
		}
		out["limits"] = limits
	}
	if q := r.GetRequests(); q != nil && (q.GetCpu() != "" || q.GetMemory() != "") {
		requests := map[string]interface{}{}
		if q.GetCpu() != "" {
			requests["cpu"] = q.GetCpu()
		}
		if q.GetMemory() != "" {
			requests["memory"] = q.GetMemory()
		}
		out["requests"] = requests
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func stringMapToInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func stringSliceToInterface(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}
