package module

var vars = struct {
	// ApiVersion / Kind of the rendered custom resource.
	ApiVersion string
	Kind       string

	// NameBudget is the ceiling on metadata.name: the operator derives
	// `<name>-head-svc` (9-character suffix) and per-group worker pod
	// names (`<name>-<group>-worker-<random>`), and Kubernetes names cap
	// at 63 characters — 40 keeps every derived name inside the budget
	// with short group names. Both engines fail loudly past it.
	NameBudget int

	// HeadServiceSuffix is the operator's naming contract for the head
	// Service: GenerateHeadServiceName (controllers/ray/utils/util.go)
	// renders "%s-%s-%s" of the cluster name, "head", and "svc".
	HeadServiceSuffix string

	// DefaultImageRepository derives the default image
	// (`rayproject/ray:<ray_version>`) — the official image,
	// version-locked to ray_version by construction.
	DefaultImageRepository string

	// AuthTokenSecretKey is the operator's RAY_AUTH_TOKEN_SECRET_KEY
	// constant — the data key inside the bearer-token Secret.
	AuthTokenSecretKey string

	// The head node's fixed API ports: CLIENT (ray.init("ray://…")),
	// DASHBOARD (web UI + Job Submission API), GCS (Ray's own
	// coordination — what `ray start --address` joins).
	ClientPort    int
	DashboardPort int
	GcsPort       int
}{
	ApiVersion:             "ray.io/v1",
	Kind:                   "RayCluster",
	NameBudget:             40,
	HeadServiceSuffix:      "-head-svc",
	DefaultImageRepository: "rayproject/ray",
	AuthTokenSecretKey:     "auth_token",
	ClientPort:             10001,
	DashboardPort:          8265,
	GcsPort:                6379,
}
