package resources

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	PostgreSQLPort = 5432

	// PostgreSQLSuperuser is the user every platform service connects as.
	// One cluster, one superuser, self-provisioned databases is the
	// platform-wide contract: the control-plane fat-jar creates its own
	// databases at boot, Temporal's schema job creates its two, and the
	// identity server's init container ensures its one. Splitting consumers
	// onto least-privilege roles is a deliberate future hardening, not an
	// accident of this shape.
	PostgreSQLSuperuser = "postgres"

	// postgresqlImage pins the exact PostgreSQL build CloudNativePG runs.
	// CloudNativePG treats a floating tag as an operational hazard (upgrades
	// must be deliberate), so the pin is explicit and bumped on purpose. The
	// major tracks the PostgreSQL line the platform is proven on.
	postgresqlImage = "ghcr.io/cloudnative-pg/postgresql:18.4"

	// postgresqlMaxConnections: one cluster serves the whole platform -- the
	// control-plane monolith (a connection pool per logical database),
	// Temporal's four services, and the identity server. PostgreSQL's
	// default max_connections=100 saturates under that (proven live: "sorry,
	// too many clients already" starved the identity server's pool and
	// crash-looped the control plane's boot).
	postgresqlMaxConnections = "300"

	// DBBase is the ONE database the control plane owns (DB_NAME=planton in
	// every deployment shape); every domain is separated at the schema level
	// inside it. This is the exact contract the desktop daemon boots against
	// (its local boot environment file); the operator
	// and the daemon are the two providers of one control-plane boot
	// contract and must stay in sync.
	//
	// The database is created at cluster birth (bootstrap initdb below) and
	// the control plane self-provisions every OTHER database it needs at
	// boot, connecting as the superuser.
	DBBase = "planton"

	// PostgreSQLOwnerRole owns DBBase. Unused by consumers today (everything
	// connects as the superuser), it exists so the database is born with a
	// non-superuser owner -- the seam a least-privilege split will widen.
	PostgreSQLOwnerRole = "planton"

	// DBOpenFGA backs the OpenFGA authorization server. Created at cluster
	// birth (postInitSQL) so enabling authorization at ANY later point finds
	// its database waiting -- OpenFGA's own migrate job cannot create it.
	DBOpenFGA = "openfga"

	// PostgreSQLClusterKind mirrors CloudNativePG's Cluster naming contract:
	// every object it creates derives from the Cluster name -- instance pods
	// ("{name}-1", ...), the traffic Services ("{name}-rw" primary
	// read-write, "{name}-ro", "{name}-r"), and the credential Secrets
	// ("{name}-app", "{name}-superuser").
	postgresqlAPIGroup   = "postgresql.cnpg.io"
	postgresqlAPIVersion = "v1"
	postgresqlKind       = "Cluster"
)

// PostgreSQLClusterGVK is the GroupVersionKind for CloudNativePG Cluster CRs.
var PostgreSQLClusterGVK = schema.GroupVersionKind{
	Group:   postgresqlAPIGroup,
	Version: postgresqlAPIVersion,
	Kind:    postgresqlKind,
}

// PostgreSQLClusterName returns the CloudNativePG Cluster name for a
// PlantonPlatform install: "{crName}-postgres".
func PostgreSQLClusterName(crName string) string {
	return fmt.Sprintf("%s-postgres", crName)
}

// PostgreSQLSuperuserSecretName returns the superuser credential Secret
// CloudNativePG generates for the platform cluster (basic-auth shape with
// "username" and "password" keys): "{cluster}-superuser". CloudNativePG owns
// this Secret's lifecycle -- it lives and dies with the Cluster, so
// credentials can never outlive (or predate) the volumes they unlock.
func PostgreSQLSuperuserSecretName(crName string) string {
	return fmt.Sprintf("%s-superuser", PostgreSQLClusterName(crName))
}

// PostgreSQLHost returns the in-cluster DNS hostname of the cluster's
// primary read-write Service: "{cluster}-rw.{namespace}.svc.cluster.local".
// Applications always connect through the -rw Service, never a pod: after a
// failover it re-points to the new primary automatically.
func PostgreSQLHost(crName, namespace string) string {
	return fmt.Sprintf("%s-rw.%s.svc.cluster.local", PostgreSQLClusterName(crName), namespace)
}

// PostgreSQLClusterOptions configures the platform's CloudNativePG Cluster.
type PostgreSQLClusterOptions struct {
	// CRName is the PlantonPlatform CR name the cluster derives its name from.
	CRName string

	// Namespace is the Kubernetes namespace for the cluster.
	Namespace string

	// Instances is the number of PostgreSQL instances (a primary plus hot
	// standbys with automated failover -- HA is one knob on the one cluster).
	Instances int32

	// StorageSize is the per-instance PVC size (e.g., "10Gi"). Sizes can only
	// grow -- CloudNativePG rejects shrinks.
	StorageSize string

	// StorageClassName pins the cluster's volumes to a StorageClass.
	// Empty means the key is omitted and the cluster default provisions.
	StorageClassName string

	// OwnerRef ties the Cluster to the PlantonPlatform CR for garbage
	// collection. CloudNativePG in turn owns the Cluster's Secrets and PVCs,
	// so deleting the platform removes the database WITH its credentials --
	// one coherent lifecycle.
	OwnerRef *metav1.OwnerReference
}

// NewPostgreSQLCluster builds the platform's postgresql.cnpg.io/v1 Cluster as
// an unstructured object.
//
// The resource floor exists because the control-plane's first boot
// self-provisions every database and runs Flyway migrations across all of
// them at once -- enough concurrent connections and shared buffers to
// OOM-kill a tiny default and cascade the application boot into failure.
// Memory limit only (no CPU limit) so migrations are never CPU-throttled.
//
// Absent quantities are OMITTED, never rendered empty: CloudNativePG's
// mutating webhook rejects "" with a quantity-format error.
func NewPostgreSQLCluster(opts PostgreSQLClusterOptions) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(PostgreSQLClusterGVK)
	obj.SetName(PostgreSQLClusterName(opts.CRName))
	obj.SetNamespace(opts.Namespace)

	if opts.OwnerRef != nil {
		obj.SetOwnerReferences([]metav1.OwnerReference{*opts.OwnerRef})
	}

	storage := map[string]any{
		"size": opts.StorageSize,
	}
	if opts.StorageClassName != "" {
		storage["storageClass"] = opts.StorageClassName
	}

	spec := map[string]any{
		"instances": int64(opts.Instances),
		"imageName": postgresqlImage,
		// Superuser access over the network is off by default in
		// CloudNativePG; the platform contract runs on it (see
		// PostgreSQLSuperuser), so it is enabled and the generated
		// {cluster}-superuser Secret is the one credential every consumer
		// references.
		"enableSuperuserAccess": true,
		"storage":               storage,
		"postgresql": map[string]any{
			"parameters": map[string]any{
				"max_connections": postgresqlMaxConnections,
			},
		},
		"resources": map[string]any{
			"requests": map[string]any{
				"cpu":    "250m",
				"memory": "512Mi",
			},
			"limits": map[string]any{
				"memory": "2Gi",
			},
		},
		"bootstrap": map[string]any{
			"initdb": map[string]any{
				"database": DBBase,
				"owner":    PostgreSQLOwnerRole,
				// Databases whose consumers cannot create them are born with
				// the cluster: OpenFGA's migrate job expects its database to
				// exist. Runs once, as the superuser, against the postgres
				// maintenance database.
				"postInitSQL": []any{
					fmt.Sprintf("CREATE DATABASE %s", DBOpenFGA),
				},
			},
		},
	}

	obj.Object["spec"] = spec

	return obj
}
