package module

var vars = struct {
	// BundleRelease is the pinned keycloak/keycloak-k8s-resources
	// release tag.
	//
	// MUST stay in sync with the Terraform module's bundle_release
	// local. There is NO user-facing version field BY DESIGN: the
	// KubernetesKeycloak declaration kind's CR rendering is built
	// against the CRD schema this bundle installs — a selectable
	// operator version would drift the schema away from what the
	// declaration kind renders. The module pins the release; upgrades
	// arrive as module updates. Always an exact release TAG, never a
	// branch — tag pinning keeps installs reproducible.
	BundleRelease string

	// DeploymentName / ServiceName / ServiceAccountName are the
	// bundle's fixed object names (upstream's own `keycloak-operator`,
	// not derived from this resource's name). Fixed names mean exactly
	// ONE operator install fits per namespace.
	DeploymentName     string
	ServiceName        string
	ServiceAccountName string

	// RelatedImageKeycloakEnvName is the operator Deployment env var
	// carrying the DEFAULT Keycloak server image the operator stamps
	// into Keycloak StatefulSets whose declaration sets no image — the
	// spec's default_keycloak_image override patches its value.
	RelatedImageKeycloakEnvName string

	// CrdFiles are the four k8s.keycloak.org CRD files published
	// beside the bundle (SHARED by both watch-scope variants; each
	// file carries exactly one CRD document).
	CrdFiles []string
}{
	BundleRelease:               "26.7.0",
	DeploymentName:              "keycloak-operator",
	ServiceName:                 "keycloak-operator",
	ServiceAccountName:          "keycloak-operator",
	RelatedImageKeycloakEnvName: "RELATED_IMAGE_KEYCLOAK",
	CrdFiles: []string{
		"keycloaks.k8s.keycloak.org-v1.yml",
		"keycloakrealmimports.k8s.keycloak.org-v1.yml",
		"keycloakoidcclients.k8s.keycloak.org-v1.yml",
		"keycloaksamlclients.k8s.keycloak.org-v1.yml",
	},
}

// bundleBaseURL is the raw-content root of the pinned release tag —
// keycloak-k8s-resources publishes NO single-file release asset; the
// tagged tree's kubernetes/ directory IS the official distribution
// (Keycloak ships no official Helm chart either; the operator is the
// first-party Kubernetes distribution).
func bundleBaseURL() string {
	return "https://raw.githubusercontent.com/keycloak/keycloak-k8s-resources/" +
		vars.BundleRelease + "/kubernetes"
}

// BundleURL is the 16-document operator bundle for the requested watch
// scope: kubernetes.yml (the operator watches ONLY its own namespace —
// JOSDK_WATCH_CURRENT) or cluster-wide/kubernetes.yml (the operator
// watches ALL namespaces — per-controller ClusterRoleBindings and
// JOSDK_ALL_NAMESPACES).
func BundleURL(clusterWide bool) string {
	if clusterWide {
		return bundleBaseURL() + "/cluster-wide/kubernetes.yml"
	}
	return bundleBaseURL() + "/kubernetes.yml"
}

// CrdURL is one of the four CRD files under the pinned tag's
// kubernetes/ root (both variants share them).
func CrdURL(file string) string {
	return bundleBaseURL() + "/" + file
}
