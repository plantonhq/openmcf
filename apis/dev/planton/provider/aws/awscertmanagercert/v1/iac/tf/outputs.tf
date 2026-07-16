output "cert_arn" {
  description = "The certificate ARN -- the join key every TLS consumer references (listeners, CloudFront, Cognito, OpenSearch, Client VPN)."
  value       = aws_acm_certificate.this.arn
}

output "status" {
  description = "The certificate status at the end of the deployment (PENDING_VALIDATION until domain ownership is proven, then ISSUED)."
  value       = aws_acm_certificate.this.status
}

output "domain_validation_records" {
  description = "The DNS records that prove domain ownership -- create these in external DNS when the module does not manage them; keep them in place so ACM can renew. Empty for imported and private certificates."
  value = [
    for dvo in aws_acm_certificate.this.domain_validation_options : {
      domain_name  = dvo.domain_name
      record_name  = dvo.resource_record_name
      record_type  = dvo.resource_record_type
      record_value = dvo.resource_record_value
    }
  ]
}

output "not_before" {
  description = "Start of the certificate's validity window (RFC3339); empty until issued."
  value       = aws_acm_certificate.this.not_before
}

output "not_after" {
  description = "End of the certificate's validity window (RFC3339) -- when an imported certificate must be re-imported by; empty until issued."
  value       = aws_acm_certificate.this.not_after
}

output "certificate_type" {
  description = "How the certificate came to be: AMAZON_ISSUED (requested), IMPORTED, or PRIVATE (ACM-PCA)."
  value       = aws_acm_certificate.this.type
}
