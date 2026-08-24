# Preset: Standard Elastic IP

**Use case:** Allocate a static public IPv4 address from Amazon's default pool.

This is the most common pattern — a zero-configuration EIP that provides a stable public IP for use with NLBs, NAT Gateways, or EC2 instances.

## What You Get

- A VPC Elastic IP allocated from Amazon's public IP pool
- Outputs: `allocation_id`, `public_ip`, `arn`, `public_dns`

## When to Use

- You need a static IP for an NLB subnet mapping
- You need a static outbound IP via NAT Gateway
- You need a persistent public IP for an EC2 instance
- You need a whitelistable IP for external service integrations

## Cost

- **Free** when associated with a running resource (EC2, NLB, NAT Gateway)
- **Billed hourly** when idle (not associated) — leaving the address unattached is the only thing that costs money here

The verified figure for this preset lives in the component's generated estimate at `catalog/_pricing/estimates/awselasticip.yaml` — computed from the pinned price book, never hand-typed here.
