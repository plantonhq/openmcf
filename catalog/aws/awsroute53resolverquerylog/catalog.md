# AWS Route 53 Resolver Query Logging

Every DNS query your VPCs make, logged: who asked, what they asked, and what the resolver answered — including DNS Firewall verdicts. The visibility layer for DNS-level security and debugging.

## What Gets Managed

- The logging configuration: its destination (CloudWatch Logs, S3, or Kinesis Data Firehose by ARN).
- The VPC associations that turn logging on per VPC.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Route 53 Resolver permissions.

### AWS Prerequisites

- The destination, writable by the resolver: same-account CloudWatch log groups work out of the box; S3 buckets and Firehose streams need the documented resource policies.

## After You Deploy

- Queries from associated VPCs start appearing in the destination within minutes.
- If an association shows FAILED, the destination was unwritable — fix permissions and re-associate (the deploy-time verifier catches this class).

## Common Changes

- Associate more VPCs (in-place list edit).
- Change the destination: this replaces the configuration (name and destination are both fixed-for-life); prior logs stay where they were written.
