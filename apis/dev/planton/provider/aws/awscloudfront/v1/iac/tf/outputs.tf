output "distribution_id" {
  description = "The distribution ID (e.g. E2ABCDEF123456) -- what invalidation requests and monitoring subscriptions key on."
  value       = aws_cloudfront_distribution.this.id
}

output "distribution_arn" {
  description = "The distribution ARN -- what WAF associations and resource policies reference."
  value       = aws_cloudfront_distribution.this.arn
}

output "domain_name" {
  description = "The CloudFront domain name (e.g. d123abc.cloudfront.net) -- the target for Route53 alias records and CNAMEs."
  value       = aws_cloudfront_distribution.this.domain_name
}

output "hosted_zone_id" {
  description = "The Route53 hosted zone ID for CloudFront alias records (always Z2FDTNDATAQYW2, exported so alias records compose without hardcoding it)."
  value       = aws_cloudfront_distribution.this.hosted_zone_id
}

output "status" {
  description = "The distribution status at the end of the deployment: Deployed (propagated everywhere) or InProgress (still propagating when wait_for_deployment is false)."
  value       = aws_cloudfront_distribution.this.status
}
