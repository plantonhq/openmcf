output "user_pool_id" {
  description = "The user pool identifier ({region}_{poolId})"
  value       = aws_cognito_user_pool.this.id
}

output "user_pool_arn" {
  description = "The ARN of the user pool"
  value       = aws_cognito_user_pool.this.arn
}

output "user_pool_endpoint" {
  description = "The pool endpoint as AWS reports it -- host and path WITHOUT a scheme (cognito-idp.{region}.amazonaws.com/{id})"
  value       = aws_cognito_user_pool.this.endpoint
}

output "issuer" {
  description = "The OIDC issuer URL JWT authorizers validate the token 'iss' claim against"
  value       = "https://${aws_cognito_user_pool.this.endpoint}"
}

output "user_pool_domain" {
  description = "The hosted-UI domain exactly as configured (prefix or custom domain, no scheme); empty when no domain is configured"
  value       = local.has_domain ? var.spec.domain.domain : ""
}

output "hosted_ui_url" {
  description = "The full https:// URL of the hosted sign-in UI; empty when no domain is configured"
  value = !local.has_domain ? "" : (
    local.is_custom_domain ?
    "https://${var.spec.domain.domain}" :
    "https://${var.spec.domain.domain}.auth.${var.spec.region}.amazoncognito.com"
  )
}

output "cloudfront_distribution" {
  description = "The CloudFront distribution domain name fronting a custom domain (the DNS alias target); empty for prefix domains"
  value       = local.is_custom_domain ? aws_cognito_user_pool_domain.this[0].cloudfront_distribution : ""
}

output "cloudfront_distribution_arn" {
  description = "The ARN of the CloudFront distribution fronting a custom domain; empty for prefix domains"
  value       = local.is_custom_domain ? aws_cognito_user_pool_domain.this[0].cloudfront_distribution_arn : ""
}

output "cloudfront_hosted_zone_id" {
  description = "The Route53 hosted-zone ID of the CloudFront distribution (alias target zone); empty for prefix domains"
  value       = local.is_custom_domain ? aws_cognito_user_pool_domain.this[0].cloudfront_distribution_zone_id : ""
}
