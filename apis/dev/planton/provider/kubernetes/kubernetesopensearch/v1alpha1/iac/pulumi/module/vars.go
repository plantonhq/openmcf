package module

var vars = struct {
	// Vendor is the value the module always renders into general.vendor.
	// The operator accepts a small enum (Opensearch/Op/OP/os/opensearch)
	// and derives the default image repository from it — pinned here, not
	// user-configurable: this component IS OpenSearch.
	Vendor string

	// DashboardsPort is the fixed port the operator's Dashboards Service
	// listens on (hardcoded upstream in pkg/builders/dashboards.go —
	// `var port int32 = 5601`), not configurable through the CRD.
	DashboardsPort int
}{
	Vendor:         "opensearch",
	DashboardsPort: 5601,
}
