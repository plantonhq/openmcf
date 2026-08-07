package module

// Output constants exported from the aws_cloud_front Pulumi module.
// Names match the Terraform module's outputs.tf key-for-key so both engines
// present one contract to consumers.
const (
	// OpDistributionId is the distribution ID (e.g. E2ABCDEF123456) -- what
	// invalidation requests and monitoring subscriptions key on.
	OpDistributionId = "distribution_id"
	// OpDistributionArn is the distribution ARN -- what WAF associations and
	// resource policies reference.
	OpDistributionArn = "distribution_arn"
	// OpDomainName is the CloudFront domain name (e.g. d123abc.cloudfront.net)
	// -- the target for Route53 alias records and CNAMEs.
	OpDomainName = "domain_name"
	// OpHostedZoneId is the Route53 hosted zone ID for CloudFront alias
	// records (always Z2FDTNDATAQYW2, exported so alias records compose
	// without hardcoding it).
	OpHostedZoneId = "hosted_zone_id"
	// OpStatus is the distribution status at the end of the deployment:
	// Deployed (propagated everywhere) or InProgress (still propagating when
	// wait_for_deployment is false).
	OpStatus = "status"
)
