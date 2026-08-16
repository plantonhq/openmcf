output "knowledge_base_id" {
  description = "The unique knowledge base identifier"
  value       = aws_bedrockagent_knowledge_base.this.id
}

output "knowledge_base_arn" {
  description = "The Amazon Resource Name of the knowledge base"
  value       = aws_bedrockagent_knowledge_base.this.arn
}

output "data_source_ids" {
  description = "Data source IDs keyed by each data_sources entry's name"
  value       = { for k, d in aws_bedrockagent_data_source.this : k => d.data_source_id }
}
