package module

// Keys exported by the gcpfilestoreinstance Pulumi module — byte-identical
// to the Terraform module's output names.
const (
	OpInstanceId      = "instance_id"       // Fully qualified resource ID
	OpInstanceName    = "instance_name"     // Short name of the instance
	OpIpAddresses     = "ip_addresses"      // Addresses on the VPC network
	OpFileShareName   = "file_share_name"   // File share (NFS mount path)
	OpCreateTime      = "create_time"       // RFC3339 creation timestamp
	OpReservedIpRange = "reserved_ip_range" // GCP-resolved /29 block
	OpEtag            = "etag"              // Concurrency guard ETag
)
