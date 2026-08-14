---
title: "Customer-Managed Key"
description: "This preset encrypts the region's AgentCore token vault with your own KMS key — you control rotation, policy, and revocation for every credential agents store."
type: "preset"
rank: "01"
presetSlug: "01-customer-managed-key"
componentSlug: "bedrock-agentcore-token-vault"
componentTitle: "Bedrock AgentCore Token Vault"
provider: "aws"
icon: "package"
order: 1
---

# Customer-Managed Key

This preset encrypts the region's AgentCore token vault with your own
KMS key — you control rotation, policy, and revocation for every
credential agents store.

## When to Use

- Compliance postures requiring customer-owned encryption for
  credentials
- Organizations centralizing key custody in KMS

## What You Get

- The default vault encrypted under the named AwsKmsKey
- Revocation power: disabling the key locks AgentCore out of every
  stored credential (an outage-grade lever — see the GUIDE)

## Customize

- Point `kmsKeyArn` at a symmetric, same-region key
- To return to AWS-managed encryption later, apply preset 02 — destroy
  alone does NOT revert the setting
