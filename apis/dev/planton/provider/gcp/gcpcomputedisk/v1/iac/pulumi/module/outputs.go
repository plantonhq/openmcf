package module

// Keys exported by the gcpcomputedisk Pulumi module — byte-identical to
// the Terraform module's output names.
const (
	OpName     = "name"      // Name of the disk in GCP
	OpDiskId   = "disk_id"   // Server-assigned unique numeric identifier
	OpSelfLink = "self_link" // Self-link URL — the attachment composition key
	OpZone     = "zone"      // Zone (plain zone name)
	OpSizeGb   = "size_gb"   // Provisioned size in GB
	OpType     = "type"      // Disk type (plain type name)
)
