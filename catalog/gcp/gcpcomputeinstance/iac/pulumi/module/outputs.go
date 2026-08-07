package module

// Keys exported by the gcpcomputeinstance Pulumi module — byte-identical
// to the Terraform module's output names.
const (
	OpInstanceName = "instance_name" // Name of the Compute Engine instance
	OpInstanceId   = "instance_id"   // Server-assigned unique numeric identifier
	OpSelfLink     = "self_link"     // GCP resource self link
	OpInternalIp   = "internal_ip"   // Primary internal IP (first interface)
	OpExternalIp   = "external_ip"   // External IP of the first interface, "" when private
	OpStatus       = "status"        // Current status (RUNNING, TERMINATED, ...)
	OpZone         = "zone"          // Zone where the instance is located
	OpMachineType  = "machine_type"  // Machine type of the instance
	OpCpuPlatform  = "cpu_platform"  // CPU platform the instance landed on
)
