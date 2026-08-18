# AwsManagedPrometheusScraper — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Budget the lifecycle, not just the resource

Creates run up to 30 minutes and deletes drain up to 20 — a scraper in a tight CI window will look hung when it is merely provisioning. The modules pin the provider's long timeouts; give pipelines the same headroom.

## The default configuration is EKS-only

Leaving `scrape_configuration` unset resolves AWS's published default at deploy — a sensible kubelet/cAdvisor/service-discovery baseline for EKS sources. VPC sources have no default (nothing to discover); the spec requires your own scrape_configs there.

## The VPC arm is the cheap smoke test

A VPC-sourced scraper against a static target proves the whole lifecycle without an EKS cluster — collectors provision, an empty scrape is a healthy scraper. Useful for platform validation before pointing at production clusters.

## Source replaces; plan for the gap

Changing anything under the source (cluster, subnets, security groups) re-provisions the collector — a scrape gap of many minutes. Roll a second scraper first when continuity matters, then retire the old one.

## Cross-account scraping is a role PAIR

`role_configuration` takes both halves (source account's role + destination account's role) or neither — and the AWS-managed `scraper_role_arn` still needs remote-write granted on the destination workspace's resource policy.
