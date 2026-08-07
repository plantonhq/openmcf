package module

var vars = struct {
	// HelmRepo is the official Istio chart repository. Chart identity
	// (repo/name/version) must stay byte-identical with the Terraform module's
	// locals (cross-engine parity contract).
	HelmRepo string

	// The four control-plane charts. Istio versions its charts in lockstep with
	// the product, so one version pin drives all releases this module installs.
	BaseChart    string
	IstiodChart  string
	CniChart     string
	ZtunnelChart string

	// Fixed release names — one control plane per cluster (the CRDs and the
	// validation-webhook plumbing are cluster singletons). The istiod release
	// gains a "-<revision>" suffix when spec.revision is set, matching the
	// chart's own resource naming.
	BaseReleaseName    string
	IstiodReleaseName  string
	CniReleaseName     string
	ZtunnelReleaseName string

	// DefaultVersion is the fallback Istio version used only when spec.version
	// is unset (spec.version is normally defaulted at manifest-load time). Keep
	// in sync with the crd2pulumi SDK pin in pkg/kubernetes/kubernetestypes/
	// Makefile and the KubernetesIstioBaseCrds module.
	DefaultVersion string

	// GatewayClassName is the GatewayClass istiod serves — the composition
	// handle exported for KubernetesGateway resources.
	GatewayClassName string

	// CrdBundleURLTemplate is the istio/base CRDs-only bundle at a release
	// tag (%s = the version). The module applies the CRDs ITSELF via
	// server-side apply and installs the base chart with every CRD excluded:
	// Helm-owned CRDs cannot adopt CRDs that already exist without Helm
	// ownership metadata, so a cluster running the CRDs-only
	// KubernetesIstioBaseCrds kind could never upgrade to the full mesh —
	// module-owned SSA CRDs are co-ownable by both kinds, making that
	// migration a plain redeploy.
	CrdBundleURLTemplate string

	// CrdNames is the exclusion list handed to the base chart
	// (base.excludedCRDs) so the release templates NO CRDs. Pinned knowledge
	// of the DefaultVersion bundle — reconcile together with the version pin
	// (and the crd2pulumi pin in pkg/kubernetes/kubernetestypes/Makefile).
	CrdNames []string

	// DefaultTrustDomain mirrors upstream MeshConfig's default.
	DefaultTrustDomain string
}{
	HelmRepo: "https://istio-release.storage.googleapis.com/charts",

	BaseChart:    "base",
	IstiodChart:  "istiod",
	CniChart:     "cni",
	ZtunnelChart: "ztunnel",

	BaseReleaseName:    "istio-base",
	IstiodReleaseName:  "istiod",
	CniReleaseName:     "istio-cni",
	ZtunnelReleaseName: "ztunnel",

	// version pin — bump together with the SDK/CRD pin when the release target moves
	DefaultVersion: "1.30.3",

	GatewayClassName: "istio",

	CrdBundleURLTemplate: "https://raw.githubusercontent.com/istio/istio/%s/manifests/charts/base/files/crd-all.gen.yaml",

	CrdNames: []string{
		"authorizationpolicies.security.istio.io",
		"destinationrules.networking.istio.io",
		"envoyfilters.networking.istio.io",
		"gateways.networking.istio.io",
		"peerauthentications.security.istio.io",
		"proxyconfigs.networking.istio.io",
		"requestauthentications.security.istio.io",
		"serviceentries.networking.istio.io",
		"sidecars.networking.istio.io",
		"telemetries.telemetry.istio.io",
		"trafficextensions.extensions.istio.io",
		"virtualservices.networking.istio.io",
		"wasmplugins.extensions.istio.io",
		"workloadentries.networking.istio.io",
		"workloadgroups.networking.istio.io",
	},

	DefaultTrustDomain: "cluster.local",
}
