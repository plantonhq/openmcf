package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpDomain          = "domain"
	OpRegion          = "region"
	OpResourceRecords = "resource_records"
	OpMappedRouteName = "mapped_route_name"
)
