output "memory_id" {
  description = "The unique memory identifier"
  value       = aws_bedrockagentcore_memory.this.id
}

output "memory_arn" {
  description = "The Amazon Resource Name of the memory"
  value       = aws_bedrockagentcore_memory.this.arn
}

output "strategy_ids" {
  description = "Strategy IDs keyed by each strategies entry's name"
  value       = { for k, s in aws_bedrockagentcore_memory_strategy.this : k => s.memory_strategy_id }
}
