# AWS Aurora DSQL

PostgreSQL that never runs out and never idles at cost: Aurora DSQL is AWS's serverless distributed SQL database — no instances to size, scale-to-zero pricing, and synchronous active-active writes across regions when you pair clusters.

## What Gets Managed

- The cluster: deletion protection, force-destroy posture, and your-own-key encryption.
- The multi-region pairing: the witness region and the peer clusters that form one logical database.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with DSQL permissions.

### AWS Prerequisites

- None for a single-region cluster. Multi-region needs the peer clusters already created in their regions (each deployed as its own instance of this kind).

## After You Deploy

- Connect with standard PostgreSQL drivers at the `endpoint` output, authenticating with IAM auth tokens (`aws dsql generate-db-connect-admin-auth-token`) — DSQL has no native database passwords.
- Private connectivity rides PrivateLink: create an interface VPC endpoint against the `vpc_endpoint_service_name` output.

## Common Changes

- Toggle deletion protection or swap the KMS key (both in-place; a key change re-encrypts without replacement).
- Multi-region pairing is create-time only — pair fresh clusters, never retrofit a live one (the witness region replaces the cluster).
