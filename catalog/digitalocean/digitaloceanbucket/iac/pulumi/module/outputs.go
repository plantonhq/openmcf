package module

const (
	// OpBucketId exports the provider's resource id for the bucket, which
	// IS the bucket name.
	OpBucketId = "bucket_id"

	// OpEndpoint exports the region-level Spaces endpoint host.
	OpEndpoint = "endpoint"

	// OpRegion exports the region slug the bucket lives in.
	OpRegion = "region"

	// OpBucketDomainName exports the bucket's virtual-host-style FQDN.
	OpBucketDomainName = "bucket_domain_name"

	// OpUrn exports the uniform resource name (do:space:<name>).
	OpUrn = "urn"
)
