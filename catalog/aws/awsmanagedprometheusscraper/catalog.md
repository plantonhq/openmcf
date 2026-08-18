# AWS Managed Prometheus Scraper

Prometheus scraping without an agent to run: AWS provisions the collectors, discovers your EKS pods (or scrapes static VPC targets), and remote-writes into Managed Prometheus — no DaemonSet, no agent upgrades, no collector capacity planning.

## What Gets Managed

- The scraper: its source (an EKS cluster with subnets/security groups, or a bare VPC placement), its destination (an AMP workspace or CloudWatch dataset), the Prometheus scrape configuration (or AWS's published default for EKS), cross-account role pairs, and component logging.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AMP scraper permissions (plus EKS describe/access for EKS sources).

### AWS Prerequisites

- EKS arm: the cluster (reference an AwsEksCluster) — AWS also needs its access entries to allow the scraper's role, which AWS manages for same-account setups.
- VPC arm: subnets and at least one security group whose rules allow the collectors to reach your scrape targets.
- The destination workspace (reference an AwsManagedPrometheus's `workspace_arn`).

## After You Deploy

- Expect the create to take up to ~30 minutes (collector provisioning) and deletes to drain up to ~20 — plan pipeline timeouts accordingly.
- `scraper_role_arn` is the writer identity; cross-account destinations grant it remote-write.
- Scrapers bill per collector-hour plus the AMP ingestion they produce.

## Common Changes

- Scrape configuration edits apply in place; source changes (cluster, subnets) replace the scraper by design.
