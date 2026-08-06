output "identity_arn" {
  description = "The ARN of the email identity -- the resource identity policies and IAM statements scope to."
  value       = aws_sesv2_email_identity.this.arn
}

output "email_identity" {
  description = "The identity string (domain or email address) -- the join key for composing DNS record names."
  value       = aws_sesv2_email_identity.this.email_identity
}

output "identity_type" {
  description = "The identity type AWS classified this as: DOMAIN or EMAIL_ADDRESS."
  value       = aws_sesv2_email_identity.this.identity_type
}

output "verification_status" {
  description = "The identity's verification status at deploy time (PENDING until DNS/mailbox verification completes, then SUCCESS)."
  value       = aws_sesv2_email_identity.this.verification_status
}

output "dkim_tokens" {
  description = "Easy DKIM's CNAME tokens: publish each as <token>._domainkey.<domain> CNAME <token>.dkim.amazonses.com. Empty for BYODKIM and email-address identities."
  value       = try(aws_sesv2_email_identity.this.dkim_signing_attributes[0].tokens, [])
}
