package module

var vars = struct {
	// names & repo
	SystemNamespace  string
	GatewayNamespace string
	HelmRepo         string
	BaseChart        string
	IstiodChart      string
	GatewayChart     string

	// default chart versions
	DefaultStableVersion string
	DefaultLatestVersion string
}{
	// namespaces
	SystemNamespace:  "istio-system",
	GatewayNamespace: "istio-ingress",

	// upstream repo & charts
	HelmRepo:     "https://istio-release.storage.googleapis.com/charts",
	BaseChart:    "base",
	IstiodChart:  "istiod",
	GatewayChart: "gateway",

	// version pins – update when you move channels forward
	DefaultStableVersion: "1.22.3",
	DefaultLatestVersion: "1.23.0",
}
