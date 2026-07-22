package verify

import (
	"context"
	"strings"
)

// ResourceVerifier knows how to verify a specific Kubernetes resource type.
type ResourceVerifier interface {
	VerifyExists(ctx context.Context, kubeconfig string) error
	VerifyAbsent(ctx context.Context, kubeconfig string) error
}

// operatorKinds lists manifest kind values (lowercased) for operator/controller
// components. Operators install CRD controllers that watch resources but typically
// do not expose a Kubernetes Service. Verification checks namespace + running
// pods only (no service requirement).
var operatorKinds = map[string]bool{
	// Tier 2/3 fixture operators (tested since sessions 3-9)
	"kubernetesperconamongooperator":    true,
	"kubernetesperconamysqloperator":    true,
	"kubernetesperconapostgresoperator": true,
	"kubernetessolroperator":            true,
	"kubernetesaltinityoperator":        true,
	"kuberneteselasticoperator":         true,
	"kubernetesstrimzikafkaoperator":    true,
	"kuberneteszalandopostgresoperator": true,
	// Tier 4 operators with configurable namespace (session 010)
	"kubernetesexternalsecrets":             true,
	"kubernetesgharunnerscalesetcontroller": true,
	"kubernetesrookcephoperator":            true,
}

// crdWorkloadKinds lists manifest kind values (lowercased) for Tier 3
// operator-dependent components. These create Custom Resources (e.g.,
// Zalando Postgresql, Strimzi Kafka, ECK Elasticsearch) that are
// reconciled by their prerequisite operator into pods and services.
// Verification checks namespace + running pods + at least one service.
var crdWorkloadKinds = map[string]bool{
	"kubernetespostgres":      true,
	"kuberneteskafka":         true,
	"kuberneteselasticsearch": true,
	"kubernetesmongodb":       true,
	"kubernetessolr":          true,
	"kubernetesclickhouse":    true,
}

// helmTier2Kinds lists manifest kind values (lowercased) for Helm-based
// Kubernetes components that deploy applications with Services.
// These must match the CloudResourceKind enum names from cloud_resource_kind.proto
// (case-insensitive via lowercasing).
var helmTier2Kinds = map[string]bool{
	// Tier 2 Helm applications
	"kubernetesredis":    true,
	"kubernetesnats":     true,
	"kubernetesgrafana":  true,
	"kubernetesneo4j":    true,
	"kubernetesopenbao":  true,
	"kubernetesopenfga":  true,
	"kubernetesjenkins":  true,
	"kubernetestemporal": true,
	"kubernetesargocd":   true,
	"kubernetesharbor":   true,
	"kuberneteslocust":   true,
	"kubernetessignoz":   true,
	"kuberneteskeycloak": true,
	// Tier 4 Helm applications with configurable namespace (session 010)
	"kubernetesingressnginx": true,
	"kubernetesistio":        true,
}

// crdInstallKinds maps manifest kind values (lowercased) to their expected CRD
// names for components that only install cluster-scoped CRDs without deploying
// any pods or services.
var crdInstallKinds = map[string][]string{
	"kubernetesgatewayapicrds": {
		"gatewayclasses.gateway.networking.k8s.io",
		"gateways.gateway.networking.k8s.io",
		"httproutes.gateway.networking.k8s.io",
		"referencegrants.gateway.networking.k8s.io",
	},
	// KubernetesIstioBaseCrds installs the istio/base CRD bundle (no istiod). Verify the
	// CRDs backing the seven typed Istio components are present.
	"kubernetesistiobasecrds": {
		"destinationrules.networking.istio.io",
		"serviceentries.networking.istio.io",
		"envoyfilters.networking.istio.io",
		"peerauthentications.security.istio.io",
		"requestauthentications.security.istio.io",
		"authorizationpolicies.security.istio.io",
		"telemetries.telemetry.istio.io",
	},
}

// gatewayApiCustomResource describes how to verify a Gateway API custom resource
// created by one of the Gateway API deployment components. These components do
// not run pods; verification confirms the CR itself exists after apply and is
// gone after destroy. The CRDs are installed by the KubernetesGatewayApiCrds
// registry prerequisite before the component applies.
type gatewayApiCustomResource struct {
	// resource is the fully-qualified kubectl resource (plural.group), which is
	// stable across the served apiVersion (e.g. tcproutes are served at v1alpha2).
	resource string
	// clusterScoped is true for cluster-scoped kinds (GatewayClass), which must
	// be queried without a namespace.
	clusterScoped bool
}

// gatewayApiKinds maps manifest kind values (lowercased) to their Gateway API
// custom-resource verification descriptor.
var gatewayApiKinds = map[string]gatewayApiCustomResource{
	"kubernetesgatewayclass":   {resource: "gatewayclasses.gateway.networking.k8s.io", clusterScoped: true},
	"kubernetesgateway":        {resource: "gateways.gateway.networking.k8s.io"},
	"kuberneteshttproute":      {resource: "httproutes.gateway.networking.k8s.io"},
	"kubernetesgrpcroute":      {resource: "grpcroutes.gateway.networking.k8s.io"},
	"kubernetestcproute":       {resource: "tcproutes.gateway.networking.k8s.io"},
	"kubernetestlsroute":       {resource: "tlsroutes.gateway.networking.k8s.io"},
	"kubernetesreferencegrant": {resource: "referencegrants.gateway.networking.k8s.io"},
}

// istioApiKinds maps manifest kind values (lowercased) to the fully-qualified
// kubectl resource (plural.group) for the typed Istio API deployment components
// (861-867). Like the Gateway API kinds, these components do not run pods;
// verification confirms the CR itself exists after apply and is gone after
// destroy. The Istio CRDs are installed by the KubernetesIstioBaseCrds registry
// prerequisite before the component applies. All seven Istio kinds are
// namespaced, so no clusterScoped flag is needed.
var istioApiKinds = map[string]string{
	"kubernetespeerauthentication":    "peerauthentications.security.istio.io",
	"kubernetesrequestauthentication": "requestauthentications.security.istio.io",
	"kubernetesauthorizationpolicy":   "authorizationpolicies.security.istio.io",
	"kubernetesserviceentry":          "serviceentries.networking.istio.io",
	"kubernetesdestinationrule":       "destinationrules.networking.istio.io",
	"kubernetesenvoyfilter":           "envoyfilters.networking.istio.io",
	"kubernetestelemetry":             "telemetries.telemetry.istio.io",
}

// GetVerifierFromManifest creates the appropriate verifier by parsing the manifest.
func GetVerifierFromManifest(manifestPath string) (ResourceVerifier, error) {
	info, err := ParseManifestInfo(manifestPath)
	if err != nil {
		return nil, err
	}

	component := strings.ToLower(info.Kind)

	switch component {
	case "kubernetesnamespace":
		return &NamespaceVerifier{Name: info.Name}, nil

	case "kubernetesdeployment":
		return &WorkloadVerifier{
			Namespace: info.Namespace,
			Kind:      "deployment",
			Name:      info.Name,
		}, nil

	case "kubernetesstatefulset":
		return &WorkloadVerifier{
			Namespace: info.Namespace,
			Kind:      "statefulset",
			Name:      info.Name,
		}, nil

	case "kubernetessecret":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "secret",
			Name:      info.Name,
		}, nil

	case "kubernetesconfigmap":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "configmap",
			Name:      info.Name,
		}, nil

	case "kubernetesserviceaccount":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "serviceaccount",
			Name:      info.Name,
		}, nil

	// The RBAC grant deploys role/binding objects whose names derive from the
	// grant's shape, so its verifier re-reads the manifest itself.
	case "kubernetesrbac":
		return &RbacVerifier{ManifestPath: manifestPath}, nil

	case "kubernetesservice":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "service",
			Name:      info.Name,
		}, nil

	// An Ingress is a valid API object with or without a controller running in
	// the cluster, so existence (not load-balancer status) is the correct
	// verification on every lane.
	case "kubernetesingress":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "ingress",
			Name:      info.Name,
		}, nil

	// A NetworkPolicy has no runtime status of its own (enforcement lives in
	// the CNI); the object's existence is the verifiable contract.
	case "kubernetesnetworkpolicy":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "networkpolicy",
			Name:      info.Name,
		}, nil

	// A claim under a WaitForFirstConsumer StorageClass is correctly Pending
	// until a pod consumes it, so existence (not Bound) is the verifiable
	// contract; the composed consumer scenario proves real binding.
	case "kubernetespersistentvolumeclaim":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "pvc",
			Name:      info.Name,
		}, nil

	// Cluster-scoped: kubectl ignores the empty namespace argument.
	case "kubernetesstorageclass":
		return &ResourceExistenceVerifier{
			Kind: "storageclass",
			Name: info.Name,
		}, nil

	// The quota object is the verifiable contract; the companion LimitRange
	// shares its name and is asserted by the with-limit-defaults scenario's
	// import round-trip (both resources ride the same state).
	case "kubernetesresourcequota":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "resourcequota",
			Name:      info.Name,
		}, nil

	// Cluster-scoped; preemption BEHAVIOR needs scheduling pressure and
	// rides the real-cluster lanes — the object is the verifiable contract.
	case "kubernetespriorityclass":
		return &ResourceExistenceVerifier{
			Kind: "priorityclass",
			Name: info.Name,
		}, nil

	case "kubernetespoddisruptionbudget":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "pdb",
			Name:      info.Name,
		}, nil

	// Scaling BEHAVIOR needs a metrics source (kind has no metrics-server);
	// the object is the verifiable contract on the kind lane.
	case "kuberneteshorizontalpodautoscaler":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "hpa",
			Name:      info.Name,
		}, nil

	case "kubernetescronjob":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "cronjob",
			Name:      info.Name,
		}, nil

	case "kubernetesjob":
		return &JobVerifier{
			Namespace: info.Namespace,
			Name:      info.Name,
		}, nil

	case "kubernetesdaemonset":
		return &WorkloadVerifier{
			Namespace: info.Namespace,
			Kind:      "daemonset",
			Name:      info.Name,
		}, nil

	// cert-manager installation: the three component Deployments must be
	// Available and the core CRDs Established — the preconditions every
	// Issuer/Certificate apply depends on.
	case "kubernetescertmanager":
		return &CertManagerInstallVerifier{
			Namespace:     info.Namespace,
			ComponentName: info.Name,
		}, nil

	// Issuers with in-cluster backends (self-signed, CA) reach Ready with
	// no external dependency, so Ready is the verifiable contract on the
	// kind cluster. ClusterIssuer is cluster-scoped (no namespace).
	case "kubernetesclusterissuer":
		return &IssuerVerifier{
			Kind: "clusterissuer",
			Name: info.Name,
		}, nil

	case "kubernetesissuer":
		return &IssuerVerifier{
			Kind:      "issuer",
			Name:      info.Name,
			Namespace: info.Namespace,
		}, nil

	// Real issuance proof: Ready condition + signed material present in the
	// target TLS Secret.
	case "kubernetescertificate":
		secretName, err := manifestSpecString(manifestPath, "secretName")
		if err != nil {
			return nil, err
		}
		return &CertificateVerifier{
			Name:       info.Name,
			Namespace:  info.Namespace,
			SecretName: secretName,
		}, nil

	case "kubernetesmanifest":
		return &ConfigGroupVerifier{
			Namespace:    info.Namespace,
			ManifestPath: manifestPath,
		}, nil

	// A real Helm release of an arbitrary chart: the Helm verifier's
	// namespace + running-pods + service assertions hold for any chart that
	// deploys a workload (the e2e chart, podinfo, deploys one pod and one
	// service).
	case "kuberneteshelmrelease":
		return &HelmComponentVerifier{
			Namespace:     info.Namespace,
			ComponentName: info.Name,
		}, nil

	// Fixed-namespace components: the proto spec has no namespace field because
	// the upstream tooling (Tekton, etc.) uses hardcoded namespaces. The manifest
	// YAML cannot carry a namespace hint because protojson.Unmarshal rejects
	// unknown fields. Namespace is therefore embedded here.
	case "kubernetestekton":
		return &HelmComponentVerifier{
			Namespace:     "tekton-pipelines",
			ComponentName: info.Name,
		}, nil

	case "kubernetestektonoperator":
		return &OperatorComponentVerifier{
			Namespace:     "tekton-operator",
			ComponentName: info.Name,
		}, nil

	default:
		if crdNames, ok := crdInstallKinds[component]; ok {
			return &CRDInstallVerifier{
				ComponentName: info.Name,
				CRDNames:      crdNames,
			}, nil
		}
		if gw, ok := gatewayApiKinds[component]; ok {
			namespace := info.Namespace
			if gw.clusterScoped {
				namespace = ""
			}
			return &ResourceExistenceVerifier{
				Namespace: namespace,
				Kind:      gw.resource,
				Name:      info.Name,
			}, nil
		}
		if resource, ok := istioApiKinds[component]; ok {
			return &ResourceExistenceVerifier{
				Namespace: info.Namespace,
				Kind:      resource,
				Name:      info.Name,
			}, nil
		}
		if operatorKinds[component] {
			return &OperatorComponentVerifier{
				Namespace:     info.Namespace,
				ComponentName: info.Name,
			}, nil
		}
		if crdWorkloadKinds[component] {
			return &CRDWorkloadVerifier{
				Namespace:     info.Namespace,
				ComponentName: info.Name,
			}, nil
		}
		if helmTier2Kinds[component] {
			return &HelmComponentVerifier{
				Namespace:     info.Namespace,
				ComponentName: info.Name,
			}, nil
		}
		return &GenericVerifier{Component: component}, nil
	}
}
