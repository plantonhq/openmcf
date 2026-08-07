---
title: "IAM OIDC Provider"
description: "IAM OIDC Provider deployment documentation"
icon: "package"
order: 100
componentName: "awsiamoidcprovider"
---

# AWS IAM OIDC Provider

Registers an OpenID Connect (OIDC) identity provider in AWS IAM -- the trust anchor for keyless, web-identity federation. It lets an external issuer's short-lived tokens be exchanged for AWS credentials through STS `AssumeRoleWithWebIdentity`, so workloads and pipelines never hold long-lived AWS access keys. You define it from an issuer URL, a list of allowed client IDs (audiences), and optional CA thumbprints; it exports the provider ARN that IAM roles trust as a `Federated` principal.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **IAM OIDC Provider** -- an `aws_iam_openid_connect_provider` (Pulumi: `iam.OpenIdConnectProvider`) registered under the issuer `url`, scoped to the supplied `clientIdList`, and optionally pinned to `thumbprintList`
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the provider

That single resource is the trust anchor. Access itself is granted by a separate **AWS IAM Role** whose trust policy references this provider's ARN.

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **An OIDC issuer** -- for EKS IRSA, an `AwsEksCluster` you can reference (its `status.outputs.oidc_issuer_url` flows in automatically); for CI/CD, the platform issuer URL (e.g. `https://token.actions.githubusercontent.com`).

### AWS Account

- **IAM permissions** -- the credentials used by the Provider Connection must allow `iam:CreateOpenIDConnectProvider`, `iam:TagOpenIDConnectProvider`, and related IAM permissions.
- **One provider per URL** -- AWS allows at most one OIDC provider per unique issuer URL per account.

## Deploy

### Console

Open the deployment store, find **AWS IAM OIDC Provider**, and click **Deploy**. The two-step creation wizard walks you through the **Issuer** (issuer URL + region) and **Trust** (audiences + optional thumbprints) decisions, with a live trust-policy preview. Start from the **EKS IRSA** preset in the [Presets](#presets) tab to pre-populate a working configuration that references your cluster.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsIamOidcProvider
metadata:
  name: github-actions-oidc
  org: acme-corp
  env: prod
spec:
  region: us-east-1
  url:
    value: https://token.actions.githubusercontent.com
  clientIdList:
    - sts.amazonaws.com
```

```shell
planton apply -f oidc-provider.yaml
```

This registers GitHub Actions as a trusted OIDC issuer. Next, create an **AWS IAM Role** whose trust policy references the exported `provider_arn`. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an OIDC provider. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Issuer URL (`url`)** -- The OIDC issuer (`iss` claim) AWS will trust. Reference an `AwsEksCluster` to wire IRSA without copying the issuer URL by hand (it resolves from the cluster's `status.outputs.oidc_issuer_url` at deploy time), or paste a literal HTTPS URL for CI issuers. The URL is fixed after creation -- changing it replaces the provider.

**Client IDs / audiences (`clientIdList`)** -- The values tokens carry in their `aud` claim. For both EKS IRSA and GitHub Actions this is `sts.amazonaws.com` (preselected by default). At least one is required; each must be 1-255 characters and unique.

**CA thumbprints (`thumbprintList`)** -- Optional SHA-1 fingerprints (40 hex chars) of the issuer's root CA. Leave empty for well-known issuers (EKS, GitHub Actions, Google) -- AWS secures the TLS connection with its trusted CA store and derives the thumbprint for you. Supply thumbprints only for an issuer whose CA is not publicly trusted.

## Outputs and Dependencies

### What This Component Consumes

| Reference | Target | Purpose |
|-----------|--------|---------|
| `url` (optional ref) | `AwsEksCluster.status.outputs.oidc_issuer_url` | Resolves the cluster's OIDC issuer so IRSA setup is composable rather than copy-paste |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `provider_arn` | ARN of the IAM OIDC provider | Referenced as a `Federated` principal in an AWS IAM Role trust policy |
| `provider_url` | Issuer URL with the `https://` scheme stripped | Builds the `<provider_url>:sub` / `<provider_url>:aud` trust-policy condition keys |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**EKS IRSA** -- Register an EKS cluster's OIDC issuer (by reference) with the `sts.amazonaws.com` audience, so Kubernetes ServiceAccounts can assume IAM roles. The required first step before creating IRSA roles. Start from the **EKS IRSA** preset.

**GitHub Actions federation** -- Register `https://token.actions.githubusercontent.com` so pipelines assume a deploy role without storing AWS keys. Start from the **GitHub Actions** preset.

**Generic issuer with a pinned thumbprint** -- For an issuer whose CA is not publicly trusted, supply the root CA's SHA-1 thumbprint. Start from the **Generic Issuer** preset.

## Works With

- **AWS IAM Role** -- the role whose trust policy references this provider's `provider_arn` to grant web-identity (IRSA / CI federation) access.
- **AWS EKS Cluster** -- exports the `oidc_issuer_url` this provider consumes for IRSA; reference it from the `url` field.
