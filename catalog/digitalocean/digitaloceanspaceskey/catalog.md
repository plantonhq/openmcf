# Spaces Access Keys on DigitalOcean

Creates the S3-style credential workloads use against Spaces, DigitalOcean's object storage -- scoped to exactly the buckets each workload needs through per-bucket read/readwrite grants, or account-wide with a single full-access grant. Integrates with Planton's Provider Connections for DigitalOcean API token management; buckets are wired by reference.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Spaces Access Key** -- the access-key/secret-key pair, carrying the grant rows you declare

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Buckets (optional)** -- DigitalOceanBucket resources for per-bucket grants; a full-access grant needs none.

### DigitalOcean Account

- Nothing extra: keys are free and managed through the same API token.

## After You Deploy

Read the key pair from the outputs -- and treat `secret_key` with care: DigitalOcean shows it exactly once, at creation, and it can never be fetched again. Point workloads at the bucket's regional Spaces endpoint with these credentials as standard S3 access/secret keys.
