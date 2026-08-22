# VPC Scraper for Static Targets

No EKS required: collectors placed into your private subnets scrape node exporters (or any Prometheus endpoint) by address. The VPC arm requires your own `scrapeConfiguration` — AWS's published default exists only for EKS sources. The security group must allow egress to the target ports.
