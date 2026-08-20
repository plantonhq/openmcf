package module

var vars = struct {
	// ApiVersion / Kind identify the custom resource the Planton operator
	// reconciles. The CRD is installed by KubernetesPlantonOperator
	// (module-owned there); this module only declares instances of it.
	ApiVersion string
	Kind       string

	// GatewayServiceSuffix / SetupCodeSecretSuffix / SetupCodeSecretKey
	// mirror the operator's deterministic per-platform naming — every
	// object the operator creates is "{platform name}-<suffix>". The
	// outputs derive from these so consumers (people, the desktop app's
	// connect-existing flow) get working handles from the first apply.
	GatewayServiceSuffix  string
	SetupCodeSecretSuffix string
	SetupCodeSecretKey    string
	GatewayDefaultPort    int
	GatewayServicePort    int

	// DeleteTimeout bounds destroy. Platform teardown is Kubernetes
	// garbage collection (every operator-created object is
	// owner-referenced to the CR), so the CR's own deletion normally
	// returns quickly — the budget exists for API-server pressure and
	// any future operator finalizer, never as an expected wait.
	DeleteTimeout string
}{
	ApiVersion:            "planton.ai/v1",
	Kind:                  "PlantonPlatform",
	GatewayServiceSuffix:  "-gateway",
	SetupCodeSecretSuffix: "-identity-setup-code",
	SetupCodeSecretKey:    "setup-code",
	GatewayDefaultPort:    8080,
	GatewayServicePort:    80,
	DeleteTimeout:         "15m",
}
