package module

// Output constants exported from the aws_cert_manager_cert Pulumi module.
// Names match the Terraform module's outputs.tf key-for-key so both engines
// present one contract to consumers.
const (
	// OpCertArn is the certificate ARN -- the join key every TLS consumer
	// references.
	OpCertArn = "cert_arn"
	// OpStatus is the certificate status at the end of the deployment
	// (PENDING_VALIDATION until domain ownership is proven, then ISSUED).
	OpStatus = "status"
	// OpDomainValidationRecords is the list of DNS records that prove domain
	// ownership -- created externally when the module does not manage them.
	OpDomainValidationRecords = "domain_validation_records"
	// OpNotBefore is the start of the certificate's validity window (RFC3339).
	OpNotBefore = "not_before"
	// OpNotAfter is the end of the certificate's validity window (RFC3339).
	OpNotAfter = "not_after"
	// OpCertificateType is how the certificate came to be: AMAZON_ISSUED,
	// IMPORTED, or PRIVATE.
	OpCertificateType = "certificate_type"
)
