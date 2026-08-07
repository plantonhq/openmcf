package module

var vars = struct {
	// ApiVersion / Kind of the rendered custom resource.
	ApiVersion string
	Kind       string

	// NameBudget is the ceiling on metadata.name: the operator derives
	// child names by suffixing (`-network-policy` is the longest at 15
	// characters) and StatefulSet pod hostnames must stay DNS-legal
	// (63-character labels) — 48 keeps every derived name inside the
	// budget. Both engines fail loudly past it.
	NameBudget int

	// Operator naming contract (verified at operator 26.7.0): every
	// child derives from the CR name — the StatefulSet is the name
	// itself, the main Service takes `-service`, the headless JGroups
	// discovery Service `-discovery`, and the generated one-time
	// bootstrap-admin Secret `-initial-admin`.
	ServiceSuffix            string
	DiscoveryServiceSuffix   string
	InitialAdminSecretSuffix string

	// Server listener defaults (the spec's own option defaults, equal
	// to Keycloak's): https 8443 / http 8080 / management 9000.
	DefaultHttpPort       int
	DefaultHttpsPort      int
	DefaultManagementPort int
}{
	ApiVersion:               "k8s.keycloak.org/v2beta1",
	Kind:                     "Keycloak",
	NameBudget:               48,
	ServiceSuffix:            "-service",
	DiscoveryServiceSuffix:   "-discovery",
	InitialAdminSecretSuffix: "-initial-admin",
	DefaultHttpPort:          8080,
	DefaultHttpsPort:         8443,
	DefaultManagementPort:    9000,
}
