output "load_balancer_arn" {
  description = "ARN of the Network Load Balancer (what listeners attach through)."
  value       = aws_lb.this.arn
}

output "load_balancer_name" {
  description = "Name of the Network Load Balancer (metadata.name, truncated to AWS's 32-character limit)."
  value       = aws_lb.this.name
}

output "load_balancer_dns_name" {
  description = "DNS name assigned by AWS to the Network Load Balancer."
  value       = aws_lb.this.dns_name
}

output "load_balancer_hosted_zone_id" {
  description = "Route53 hosted zone ID for the NLB's DNS name."
  value       = aws_lb.this.zone_id
}
