output "allocation_id" {
  description = "The allocation ID of the Elastic IP (e.g., eipalloc-xxx)."
  value       = aws_eip.this.allocation_id
}

output "public_ip" {
  description = "The public IPv4 address assigned to this Elastic IP."
  value       = aws_eip.this.public_ip
}

output "arn" {
  description = "The ARN of the Elastic IP."
  value       = aws_eip.this.arn
}

output "public_dns" {
  description = "The public DNS hostname associated with the Elastic IP."
  value       = aws_eip.this.public_dns
}

output "association_id" {
  description = "The association ID (eipassoc-xxx) when the spec attaches this EIP to an instance or network interface; empty when unattached."
  value       = aws_eip.this.association_id
}

output "ptr_record" {
  description = "The reverse DNS (PTR) record AWS granted when reverse_dns_domain_name is set; empty otherwise."
  value       = local.reverse_dns_domain_name != null ? aws_eip_domain_name.this[0].ptr_record : ""
}
