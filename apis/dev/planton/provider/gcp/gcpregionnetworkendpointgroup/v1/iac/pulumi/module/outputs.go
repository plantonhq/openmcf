package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpSelfLink                 = "self_link"
	OpNetworkEndpointGroupName = "network_endpoint_group_name"
	OpNetworkEndpointType      = "network_endpoint_type"
	OpRegion                   = "region"
)
