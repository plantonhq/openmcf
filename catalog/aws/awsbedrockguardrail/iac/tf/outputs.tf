output "guardrail_id" {
  description = "The unique guardrail identifier -- the join key model invocations and agents use together with a version."
  value       = aws_bedrock_guardrail.this.guardrail_id
}

output "guardrail_arn" {
  description = "The ARN of the guardrail -- the canonical key for IAM policies and cross-account references."
  value       = aws_bedrock_guardrail.this.guardrail_arn
}

output "draft_version" {
  description = "The guardrail's mutable working version -- always the literal DRAFT."
  value       = aws_bedrock_guardrail.this.version
}

output "version_numbers" {
  description = "AWS-assigned numbers of the versions published through spec.versions, keyed by each entry's name."
  value       = { for name, v in aws_bedrock_guardrail_version.published : name => v.version }
}
