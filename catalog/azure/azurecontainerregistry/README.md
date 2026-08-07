# Overview

The **Azure Container Registry API Resource** provides a consistent and standardized interface for deploying and managing Azure Container Registry (ACR) instances within our infrastructure. ACR is the managed, private OCI registry that stores the container images and artifacts a platform's workloads run — this resource makes provisioning it declarative, validated, and composable.

## Purpose

We developed this API resource to streamline the deployment of production container registries on Azure. The spec mirrors Azure's own SKU tiering rather than hiding it, enabling:

- **Fail-Fast Configuration**: Spec-level validation enforces the same SKU gates ARM does, so a misconfigured manifest (e.g. geo-replication on a Standard registry) fails at validation, not at apply
- **Full Premium Surface**: Geo-replication, zone redundancy, network isolation, content-trust/quarantine/retention policies, and customer-managed-key encryption — all declared in one spec
- **Secure Defaults**: Admin account off by default; Microsoft Entra ID (managed identities, service principals) is the production authentication path
- **Composition by Reference**: The identities and grants the registry works with are first-class catalog resources, never bundled

## Key Features

- **Consistent Interface**: Aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **SKU as the Feature Gate**: `BASIC` for dev/test, `STANDARD` as the production baseline (applied when the SKU is left unspecified), `PREMIUM` for the enterprise surface — with every Premium-only field validated against the chosen tier
- **Geo-Replication**: Per-replica configuration (region, zone redundancy, regional endpoint, tags); the home region is implicit and excluded from the list
- **Network Isolation Options**: IP allowlisting via `networkRuleSet`, fully private via `publicNetworkAccessEnabled: false`, and dedicated data endpoints for exact firewall allowlisting
- **CMK Encryption**: Customer-managed-key encryption through a referenced first-class `AzureUserAssignedIdentity` that unwraps the Key Vault key
- **AKS Integration**: AKS clusters reference the registry via its `container_registry_id` output; `AcrPull`/`AcrPush` grants are composed with the standalone `AzureRoleAssignment` resource

## Use Cases

- **Private Container Images**: Store and distribute the images AKS clusters, App Service, and Container Apps run
- **CI/CD Pipelines**: Push targets for Azure DevOps, GitHub Actions, and GitLab CI builds, authenticated via Entra service principals or OIDC
- **Global Deployments**: Premium geo-replication serves pulls locally in every deployment region — lower latency, no cross-region egress fees, and pull availability through a regional outage
- **Locked-Down Environments**: IP-allowlisted or fully private registries with dedicated data endpoints, quarantine-gated pushes, signed images, and customer-managed encryption keys
- **Registry Hygiene**: Retention policies that automatically purge untagged manifests so CI push churn cannot grow storage without bound

## Future Enhancements

Future updates will include:

- **Repository-Scoped Tokens and Scope Maps**: Fine-grained, repo-level credentials for systems that cannot use Entra ID
- **Webhooks**: Push/delete event notifications for downstream automation
- **Cache Rules**: Pull-through caching of upstream registries (Docker Hub, MCR)
- **ACR Tasks**: Cloud-native image builds and scheduled maintenance jobs

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
