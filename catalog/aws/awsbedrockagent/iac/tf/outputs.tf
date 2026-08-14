output "agent_id" {
  description = "The unique agent identifier"
  value       = aws_bedrockagent_agent.this.agent_id
}

output "agent_arn" {
  description = "The Amazon Resource Name of the agent"
  value       = aws_bedrockagent_agent.this.agent_arn
}

output "draft_version" {
  description = "The agent's mutable working version (always DRAFT)"
  value       = "DRAFT"
}

output "alias_ids" {
  description = "Alias IDs keyed by each aliases entry's name"
  value       = { for k, a in aws_bedrockagent_agent_alias.this : k => a.agent_alias_id }
}

output "alias_arns" {
  description = "Alias ARNs keyed by each aliases entry's name"
  value       = { for k, a in aws_bedrockagent_agent_alias.this : k => a.agent_alias_arn }
}

output "action_group_ids" {
  description = "Action group IDs keyed by each action_groups entry's name"
  value       = { for k, g in aws_bedrockagent_agent_action_group.this : k => g.action_group_id }
}

output "collaborator_ids" {
  description = "Collaborator IDs keyed by each collaborators entry's name"
  value       = { for k, c in aws_bedrockagent_agent_collaborator.this : k => c.collaborator_id }
}

output "associated_knowledge_base_ids" {
  description = "Associated knowledge base IDs keyed by each knowledge_base_associations entry's name"
  value       = { for k, a in aws_bedrockagent_agent_knowledge_base_association.this : k => a.knowledge_base_id }
}
