package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	// The nats chart's served index lives under /helm/charts/ — the
	// repo-root index.yaml 404s (verified at the pin).
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. The
	// chart and the NATS server move in lockstep: chart 2.14.2 runs
	// nats-server 2.14.2.
	DefaultChartVersion string
	// The chart's fullname budget: the longest derived child name is
	// `<fullname>-box-contents` (13 extra chars) and Helm truncates
	// derived names at 63 — names past 63-13=50 characters would
	// truncate SILENTLY and break the naming contract the exported
	// outputs are built on.
	FullnameBudget int
	// Listener ports (the chart's defaults; websocket/mqtt/leafnodes
	// ports are spec-configurable).
	ClientPort int
	// The env-var prefix for the per-user password references: the
	// rendered nats.conf carries `$NATS_PW_<i>` (unquoted via the
	// chart's `<< >>` syntax) and the server resolves it from the
	// container environment — verified in the server's config parser at
	// the pin (conf/parse.go, os.LookupEnv). Index-based names keep the
	// contract deterministic and env-var-safe for any username.
	PasswordEnvPrefix string
	// Where the chart mounts the client-listener TLS Secret; the CA
	// bundle key inside a cert-manager Secret.
	ClientTlsDir string
	CaFileName   string
	// Helm wait budget: a StatefulSet rollout (up to 9 servers) plus
	// nats-box; NATS starts in seconds.
	HelmTimeoutSeconds int
}{
	HelmChartName:       "nats",
	HelmChartRepo:       "https://nats-io.github.io/k8s/helm/charts/",
	DefaultChartVersion: "2.14.2",
	FullnameBudget:      50,
	ClientPort:          4222,
	PasswordEnvPrefix:   "NATS_PW_",
	ClientTlsDir:        "/etc/nats-certs/nats",
	CaFileName:          "ca.crt",
	HelmTimeoutSeconds:  300,
}
