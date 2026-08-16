output "prompt_id" {
  description = "The unique prompt identifier"
  value       = aws_bedrockagent_prompt.this.id
}

output "prompt_arn" {
  description = "The Amazon Resource Name of the prompt"
  value       = aws_bedrockagent_prompt.this.arn
}

output "draft_version" {
  description = "The prompt's mutable working version (always DRAFT)"
  value       = "DRAFT"
}
