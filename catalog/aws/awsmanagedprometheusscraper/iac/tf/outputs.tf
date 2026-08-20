output "scraper_id" {
  description = "The scraper's AWS-generated ID (the provider's import ID)"
  value       = aws_prometheus_scraper.this.id
}

output "scraper_arn" {
  description = "The scraper's ARN"
  value       = aws_prometheus_scraper.this.arn
}

output "scraper_role_arn" {
  description = "The AWS-managed role the scraper writes to its destination with - grant it remote-write on cross-account destinations"
  value       = aws_prometheus_scraper.this.role_arn
}
