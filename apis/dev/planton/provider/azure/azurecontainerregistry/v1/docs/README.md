# Azure Container Registry: The Registry as Production Infrastructure

## Introduction: More Than a Place to Put Images

Every containerized platform has a registry, and most teams treat it as a solved problem: create it once, push images, move on. Then production arrives. A regional outage takes image pulls down with it. An auditor asks why the registry is reachable from the entire internet. A compliance framework demands customer-controlled encryption keys. CI has pushed forty thousand untagged manifests and storage costs are climbing. A partner's egress firewall team asks for the exact hostnames image pulls will use, and the answer turns out to be "a shared storage endpoint we don't control."

None of these are exotic scenarios — they are the ordinary lifecycle of a registry that started as a dev/test convenience and became production infrastructure. Azure Container Registry (ACR) has first-class answers to all of them, but they are unevenly distributed across its pricing tiers and coupled to decisions that must be made at creation time. Understanding that structure *before* deploying is the difference between a registry that grows with the platform and one that must be replaced — and a registry replacement is uniquely painful, because **registry contents do not migrate**: every image must be re-pushed.

This document covers the ACR domain — SKU tiers, geo-replication, network isolation, supply-chain policies, encryption, and authentication — and explains how Planton's declarative API models it.

## The SKU Is the Feature Gate

ACR's three SKUs differ in quota and throughput, but the decision that matters is which *capabilities* each tier unlocks:

| Capability | Basic | Standard | Premium |
|------------|-------|----------|---------|
| Included storage | 10 GiB | 100 GiB | 500 GiB |
| Core push/pull, Entra auth, webhooks | Yes | Yes | Yes |
| Anonymous (unauthenticated) pull | No | Yes | Yes |
| Geo-replication | No | No | Yes |
| Zone redundancy | No | No | Yes |
| Network rule set (IP allowlist) | No | No | Yes |
| Disabling public network access / private endpoints | No | No | Yes |
| Dedicated data endpoints | No | No | Yes |
| Quarantine, retention, and content-trust policies | No | No | Yes |
| Customer-managed-key encryption | No | No | Yes |

**Basic** is a cost-optimized dev/test tier with the full API surface but low storage and throughput limits. **Standard** is the production baseline for single-region workloads that don't need network isolation. **Premium** is not a performance upgrade so much as a *category* change: everything related to availability, isolation, and compliance lives there.

Two consequences follow:

1. **Pick the tier by capability, not by cost.** Premium's base price is modest next to the engineering cost of discovering mid-incident that private endpoints or geo-replication require a tier change. The SKU does change in place — upgrading is trivial — but downgrading requires every Premium-only feature to be unset first, in the same order ARM enforces.
2. **Validation should enforce the gates before ARM does.** ARM rejects a Basic registry with geo-replications at deploy time, minutes into an apply. Planton's spec enforces the identical gates at validation time: a manifest that sets `zoneRedundancyEnabled`, `networkRuleSet`, `retentionPolicyInDays`, or any other Premium-only field on a non-Premium registry is rejected before any infrastructure operation starts.

## Geo-Replication: One Registry, Many Regions

A geo-replicated registry behaves as a single logical registry with physical replicas in multiple regions. Pushes go to any replica and propagate to all of them automatically; pulls are served by the nearest replica via Azure's traffic routing. The payoffs:

- **Pull latency**: AKS nodes in West Europe pull from West Europe, not across the Atlantic. Pod startup times improve materially for large images.
- **Egress cost**: cross-region image pulls incur bandwidth charges; local replicas eliminate them.
- **Availability**: if a replicated region's storage goes down, the remaining replicas keep serving pulls.

Each replica is its own tracked ARM resource with its own configuration: zone redundancy is declared per replica (the home replica's `zoneRedundancyEnabled` and each replication's are independent), each replica can expose its own regional endpoint that clients address directly, and each carries its own tags. The home region is implicit — the replication list must contain only *additional* regions, and validation rejects a list that includes the registry's own region.

One provider-level quirk is worth knowing because it shapes how tooling must handle replications: the `azurerm` provider expects the inline replication list in **alphabetical order by location**, and produces a perpetual diff otherwise. Planton's modules absorb this — both engines sort replications by location internally — so manifests can list regions in any order without triggering spurious changes. Replicas add and remove in place; changing a replica's zone redundancy replaces that replication (a re-sync), not the registry.

## Network Isolation: Three Distinct Levers

ACR offers three network controls that are frequently conflated. They solve different problems and compose differently.

### The Network Rule Set (Public, but Allowlisted)

The registry stays publicly addressable, but a default-deny rule set drops every connection not originating from an explicit IPv4 CIDR allowlist (office egress, CI runners, NAT gateway addresses). This is the middle ground between "open to the internet" and "fully private": no private DNS, no endpoint NICs, no VNet coupling — just a reviewable allowlist. Two details matter:

- The default action must be `DENY` for the rule set to do anything; Azure's default (`ALLOW`) makes it a no-op.
- ARM only supports *allow* rules, so entries carry no per-rule action — the rule set is exactly a default action plus an allowlist.

Trusted Azure services (ACR Tasks, Microsoft Defender) bypass the restrictions by default (`networkRuleBypassOption` unspecified = AzureServices); setting `NONE` closes even that door, at the cost of breaking first-party integrations.

### Private Endpoints (Fully Private)

Setting `publicNetworkAccessEnabled: false` removes the registry from the public internet entirely — it becomes reachable only through private endpoints, which project the registry into a VNet with a private IP and private DNS. This is the regulated-industry posture: pulls and pushes traverse Azure's backbone, never the internet. The endpoints themselves are separate resources composed alongside the registry, not properties of it. A fully private registry also unlocks disabling the export policy — a data-exfiltration control that blocks transferring artifacts out of the registry, which ARM (and Planton validation) only permits when public access is explicitly off.

### Dedicated Data Endpoints (Exact Allowlisting)

Registry traffic has two halves: the REST endpoint (`{name}.azurecr.io`) and the blob data, which by default is served from *shared* regional storage endpoints with unpredictable hostnames. That default is what breaks egress-firewall allowlisting: allowing "the registry" means allowing a broad storage wildcard. Enabling dedicated data endpoints moves blob traffic to `{name}.{region}.data.azurecr.io` — one deterministic hostname per region (home plus each replica), surfaced in the `data_endpoint_host_names` output — so the firewall allowlist becomes exact. This lever is about traffic *predictability*, not access control, and it pairs naturally with either of the other two.

## Supply-Chain Policies: Quarantine, Retention, Trust

Premium registries carry three built-in policies aimed at image hygiene and provenance:

**Quarantine** holds every newly pushed image in a quarantined state until a scanner (or other automation) marks it passed; clients without quarantine access cannot pull it. The workflow itself — scanning, passing, failing — is driven through the registry's data-plane API by the scanning tooling; the registry just enforces the gate. This turns "we scan images" from a convention into a guarantee: an unscanned image is unpullable.

**Retention** automatically purges *untagged* manifests after a configurable window (0-365 days; 0 means immediately). This is the answer to CI churn: every rebuild that re-tags `latest` orphans the previous manifest, and without a retention policy those orphans accumulate forever. A 30-day window keeps storage bounded without risking anything still referenced. Note the scope: only untagged manifests are purged — cleaning up old *tagged* images is a different problem (deliberately out of scope here; see the backlog note below).

**Content trust** enables Docker Content Trust: clients with content trust enabled can push signed images and verify signatures at pull, providing publisher-side provenance. It requires client-side participation (`DOCKER_CONTENT_TRUST=1`) to mean anything — the registry stores and serves signatures; enforcement happens at the client.

## Customer-Managed-Key Encryption

ACR encrypts at rest by default with Microsoft-managed keys. Compliance regimes that require customer-controlled keys (rotate on your schedule, revoke to render data unreadable, audit key usage in your own Key Vault logs) need CMK encryption, and it has the strictest prerequisites of any registry feature:

1. **Premium only**, and **fixed at creation** — a registry cannot gain or lose CMK later; that decision is part of the registry's identity.
2. **A user-assigned managed identity is mandatory.** The registry must unwrap the encryption key *before* it exists, so a system-assigned identity (which is created with the registry) cannot do the job. The identity must exist first, be attached via the registry's identity configuration, and be named (by client ID) in the encryption block.
3. **The identity needs key permissions on the vault**: `get`, `wrapKey`, and `unwrapKey` on the Key Vault holding the key (via access policy or the equivalent RBAC role). A missing permission fails the deploy with an opaque ARM error, so grant first, create second.
4. **The key is referenced by its full Key Vault ID** (`https://{vault}.vault.azure.net/keys/{name}`); pinning a version suffix freezes rotation to that version.

The ordering constraint in (2) and (3) is exactly why composition-by-reference matters: the identity and its vault grant must be independent resources that exist before the registry does.

## Authentication: Three Surfaces, One Production Answer

**Microsoft Entra ID is the production path.** Azure RBAC roles — `AcrPull`, `AcrPush`, `AcrDelete` — assigned to managed identities (AKS kubelets, Container Apps), service principals (CI pipelines), or federated OIDC credentials (GitHub Actions without stored secrets) give short-lived, audited, least-privilege access. Every grant is a visible role assignment that can be reviewed and revoked.

**The admin account is the escape hatch.** One username (the registry name), two rotatable passwords, full registry access, no identity attribution in audit logs. Azure disables it by default and so does Planton. Its legitimate use is narrow: consumers that can accept nothing but a static username/password (some app-hosting image-pull integrations). When enabled, the credentials surface as stack outputs — and only then.

**Repository-scoped tokens** sit between the two: static credentials scoped to specific repositories and actions, for external parties that cannot use Entra but shouldn't get registry-wide access. They are part of the ACR data plane rather than the registry's ARM configuration.

**Anonymous pull** is the fourth, deliberately extreme setting: it makes every repository in the registry publicly readable, which is the right shape for exactly one use case — operating a public artifact-distribution registry — and a data leak for every other.

## What Is Deliberately Not Modeled

ACR's surface extends beyond the registry resource itself: **repository-scoped tokens and scope maps**, **webhooks** (push/delete event notifications), **cache rules** (pull-through caching of upstream registries), and **ACR Tasks** (cloud-native builds and scheduled jobs). None of these is modeled as spec surface, and the omission is a decision, not a gap: none of them is referenced by a foreign key anywhere in the catalog, and no composed-infrastructure scenario has demanded them. Modeling them speculatively would grow the spec's maintenance surface without a consumer. They live on the adoption backlog and become spec surface when a real composition needs them.

The same reasoning applies to private endpoints in the opposite direction: they *are* demanded by real scenarios, and they are modeled — as the separate resource they actually are in ARM, composed against the registry, rather than inlined into it.

## The Planton Approach

Planton provides a declarative, protobuf-based API for ACR whose design philosophy is that **the spec should mirror Azure's real feature structure, validated up front, with everything the registry composes with referenced rather than bundled**.

### The SKU Gates Are Spec-Level Validation

Every Premium-only field is guarded by a validation rule that mirrors ARM's own enforcement: geo-replications, zone redundancy, the network rule set, disabling public access, data endpoints, the quarantine/retention/trust policies, and CMK encryption all require `PREMIUM`; anonymous pull requires `STANDARD` or `PREMIUM`; disabling the export policy requires `PREMIUM` *and* public access explicitly off. A manifest violating any gate fails validation immediately — the feedback arrives in seconds, not minutes into an apply. An unspecified SKU deploys the `STANDARD` baseline.

### Composition by Reference

- **AKS clusters** reference the registry through its `container_registry_id` output — the registry doesn't know or list its consumers.
- **Pull and push grants** are standalone `AzureRoleAssignment` resources scoping `AcrPull`/`AcrPush` to the registry's ARM ID. Grants are never bundled into the registry, so access stays visible in the resource graph, auditable, and independently lifecycled.
- **The CMK identity** is a referenced first-class `AzureUserAssignedIdentity`: `identity.identityIds` carries its ARM ID, and the encryption block's `identityClientId` defaults to referencing the same identity's `client_id` output. Because the identity is independent, its Key Vault grant can be composed *before* the registry exists — which is exactly the ordering CMK requires. The Key Vault key ID is a plain value today and becomes referenceable when a Key Vault key kind exists in the catalog.

### Lifecycle Honesty

The spec's field documentation states what each change costs. Name and region are the registry's identity — changing either replaces the registry, and contents do not migrate. Zone redundancy (home replica) and CMK encryption are likewise fixed at creation. Nearly everything else — SKU, admin account, network posture, policies, replicas, tags — updates in place.

### Example: The Multi-Region Production Shape

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerRegistry
metadata:
  name: platform-registry
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: platform-rg
  registryName: platformregistry01
  sku: PREMIUM
  zoneRedundancyEnabled: true
  retentionPolicyInDays: 30
  georeplications:
    - location: westeurope
      zoneRedundancyEnabled: true
```

This configuration:
- Survives a zone outage in the home region and keeps serving pulls through a regional outage of either region
- Serves European pulls locally, with no cross-region egress fees
- Purges untagged CI churn after 30 days, keeping storage bounded

### Example: The Locked-Down Compliance Shape

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerRegistry
metadata:
  name: secure-registry
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: secure-rg
  registryName: secureregistry01
  sku: PREMIUM
  publicNetworkAccessEnabled: false
  exportPolicyEnabled: false
  quarantinePolicyEnabled: true
  identity:
    type: USER_ASSIGNED
    identityIds:
      - valueFrom:
          name: acr-cmk-identity
  encryption:
    identityClientId:
      valueFrom:
        name: acr-cmk-identity
    keyVaultKeyId:
      value: https://secure-vault.vault.azure.net/keys/acr-cmk
```

This configuration:
- Is reachable only through private endpoints, with artifact export blocked as an exfiltration control
- Gates every pushed image behind quarantine until scanning passes it
- Encrypts registry data with a customer-owned Key Vault key, unwrapped by a referenced user-assigned identity that was granted `get`/`wrapKey`/`unwrapKey` on the vault before the registry was created

## Conclusion: Decide at Creation What Cannot Change Later

Most of a registry's configuration is forgiving — SKUs upgrade in place, replicas come and go, policies toggle. Four things are not: the name, the region, home-replica zone redundancy, and CMK encryption. Those four define the registry's identity, and getting them wrong means a replacement whose contents do not follow.

Everything else is a matter of matching capability to tier: Standard for the single-region baseline, Premium the moment availability, isolation, or compliance enters the conversation. Planton's API keeps that decision explicit — the SKU gates are validated up front, the composition points (identities, grants, consumers) are first-class references, and the lifecycle cost of every field is documented where you set it.

Treat the registry as production infrastructure from the first manifest, and it never has to be replaced to become it.
