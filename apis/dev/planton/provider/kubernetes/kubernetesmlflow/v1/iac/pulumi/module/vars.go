package module

// Deployment identity — must stay byte-identical with the Terraform
// module's locals: this kind renders MODULE-OWNED manifests (MLflow
// publishes no Helm chart; the official image is the distribution), so
// these constants ARE the deployment contract on both engines.
var vars = struct {
	// The OFFICIAL MLflow image (ghcr.io/mlflow/mlflow) and the release
	// this kind is built against. Fallbacks when spec.server.image is
	// unset — mirror of the proto defaults and the Terraform coalesces.
	DefaultImageRepository string
	DefaultImageTag        string
	// The tracking server's HTTP port and default uvicorn worker count
	// (upstream's own deployment reference values).
	ServerPort     int
	DefaultWorkers int
	// Filesystem paths inside the container: the PVC mount for the
	// sqlite backend (tracking db + auth db) and the PVC mount for
	// locally-stored artifacts (the server proxies all artifact
	// traffic).
	DataMountPath      string
	ArtifactsMountPath string
	// Where the auth configuration Secret mounts and the ini file name
	// the server reads through MLFLOW_AUTH_CONFIG_PATH.
	AuthConfigMountPath string
	AuthConfigFileName  string
	// Where the GCS service-account key Secret mounts (the declared-key
	// arm; keyless rides ambient Workload Identity).
	GcsCredentialsMountPath string
	// Module-owned Secret names/keys. The backend URI Secret carries
	// the composed SQLAlchemy URI under `uri`; the admin Secret carries
	// the bootstrap admin password under `password`; the auth-config
	// Secret carries the basic-auth ini.
	BackendUriSecretSuffix string
	BackendUriKey          string
	AdminAuthSecretSuffix  string
	AdminPasswordKey       string
	AuthConfigSecretSuffix string
	// Fixed child-resource name suffixes (Deployment/Service render
	// bare `<name>`; satellites append these).
	GcCronJobSuffix       string
	ArtifactsPvcSuffix    string
	DataPvcSuffix         string
	ServiceMonitorSuffix  string
	MetricsExportPathFlag string
}{
	DefaultImageRepository:  "ghcr.io/mlflow/mlflow",
	DefaultImageTag:         "v3.15.0",
	ServerPort:              5000,
	DefaultWorkers:          4,
	DataMountPath:           "/mlflow/data",
	ArtifactsMountPath:      "/mlflow/artifacts",
	AuthConfigMountPath:     "/etc/mlflow/auth",
	AuthConfigFileName:      "basic_auth.ini",
	GcsCredentialsMountPath: "/etc/mlflow/gcs",
	BackendUriSecretSuffix:  "-backend-uri",
	BackendUriKey:           "uri",
	AdminAuthSecretSuffix:   "-admin-auth",
	AdminPasswordKey:        "password",
	AuthConfigSecretSuffix:  "-auth-config",
	GcCronJobSuffix:         "-gc",
	ArtifactsPvcSuffix:      "-artifacts",
	DataPvcSuffix:           "-data",
	ServiceMonitorSuffix:    "-metrics",
	MetricsExportPathFlag:   "/tmp/metrics",
}
