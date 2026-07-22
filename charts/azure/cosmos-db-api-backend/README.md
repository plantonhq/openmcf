# Azure Cosmos DB API Backend

A planet-scale API data layer in the posture Azure's own security baseline describes and few teams assemble: key-based data-plane access disabled entirely, a least-privilege custom role granted to a workload identity, a private endpoint with no public data plane, continuous backup, and the throttling alert that catches an undersized RU ceiling before customers do. There is no connection string anywhere in this architecture — nothing to leak, rotate, or vault.

## Who this is for

A product team building an API on Cosmos DB who wants the keyless RBAC-only posture without reverse-engineering it. The pieces are individually documented and collectively subtle: disabling local auth without first wiring a data-plane role locks everyone out, the custom-role verbs are allow-only (no deny carve-outs), and the private endpoint needs the right subresource name and DNS zone to keep SDK clients working unchanged. Deploy this chart and the posture is simply already correct.

## Architecture

```
                      VNet 10.40.0.0/16
   ┌───────────────────────────────────────────────┐
   │  private-endpoints subnet 10.40.1.0/24        │      privatelink.documents.azure.com
   │                                               │◀──── private DNS zone + VNet link
   │   private endpoint ── subresource "Sql"       │
   └─────────┬─────────────────────────────────────┘
             │ the only data-plane path
   ┌─────────▼─────────────────────────────────────┐
   │  Cosmos DB account (SQL API)                  │      app identity (UAI)
   │   · localAuthenticationEnabled: false         │        │ Entra token auth
   │   · publicNetworkAccessEnabled: false         │◀───────┘
   │   · continuous backup (7-day PITR)            │   custom role "api-backend"
   │   └─ database ─ container (autoscale RU/s)    │   (read + create/replace/upsert)
   └─────────┬─────────────────────────────────────┘
     diagnostics (CDBDataPlaneRequests) ──▶ Log Analytics
     alert: 429s > 100 in 15m ──▶ ops email
```

Design decisions worth knowing:

- **Keyless is the architecture, not a flag.** `localAuthenticationEnabled: false` makes Entra RBAC the only door, and the chart supplies everything that door needs: the workload identity, a custom role with exactly the verbs an API uses, and the grant binding them. Applications construct their Cosmos client with the account endpoint and a token credential — the identity's client id is the only configuration.
- **The custom role spells out its verbs.** Cosmos RBAC is allow-only, so least privilege means listing what the workload does: the full read surface plus create/replace/upsert. `items/*` would silently include delete — the omission here is a decision the template comments make visible.
- **Autoscale with a visible ceiling.** The container scales between 10% and `autoscale_max_throughput` on demand; the throttle alert fires when demand presses past the ceiling, making "raise the ceiling or fix the hot partition" an informed operator decision instead of a customer-reported incident.
- **Multi-region is one toggle.** `secondary_region_enabled` adds a read region with automatic failover — reads keep serving through a regional outage and writes resume after failover. It doubles RU and storage cost, which is why it is a toggle and not the default.

## Resources

| Kind | Name | Purpose |
| --- | --- | --- |
| AzureResourceGroup | `{env}-cosmos-api` | One container for the estate |
| AzureLogAnalyticsWorkspace | `{env}-cosmos-logs` | Data-plane request logs, auth failures |
| AzureUserAssignedIdentity | `{env}-cosmos-app-identity` | The API workload's identity — the only data-plane principal |
| AzureCosmosdbAccount | `{env}-cosmos` | SQL API, keyless, private, continuous backup |
| AzureCosmosdbSqlDatabase | `{env}-cosmos-db` | The API's database (throughput lives on the container) |
| AzureCosmosdbSqlContainer | `{env}-cosmos-container` | Partitioned, autoscale RU/s, explicit indexing posture |
| AzureCosmosdbSqlRoleDefinition | `{env}-cosmos-app-role` | Least-privilege data-plane role (no delete) |
| AzureCosmosdbSqlRoleAssignment | `{env}-cosmos-app-grant` | Grants the role to the identity, account-wide |
| AzureVirtualNetwork | `{env}-cosmos-vnet` | The data-plane network |
| AzureSubnet | `{env}-cosmos-pe-subnet` | Private-endpoint subnet |
| AzurePrivateDnsZone + link | `{env}-cosmos-dns` | Private resolution of the account FQDN |
| AzurePrivateEndpoint | `{env}-cosmos-pe` | The SQL data plane's only path |
| AzureMonitorActionGroup | `{env}-cosmos-ops` | Alert routing |
| AzureMonitorDiagnosticSetting | `{env}-cosmos-diag` | Request logs + metrics into the workspace |
| AzureMonitorMetricAlert | `{env}-cosmos-throttle-alert` | Sustained 429s — the RU ceiling alert |

## Parameters

| Parameter | Description | Default | Must change |
| --- | --- | --- | --- |
| `region` | Azure region | `centralus` | |
| `cosmos_account_name` | Globally unique account name (becomes the DNS name) | `my-api-cosmos` | yes |
| `database_name` | The SQL database | `appdb` | |
| `container_name` | The primary container | `items` | |
| `partition_key_path` | Partition key path — fixed for the container's life | `/tenantId` | eventually |
| `autoscale_max_throughput` | Autoscale ceiling in RU/s (multiples of 1000) | `1000` | |
| `consistency_level` | Account default consistency | `SESSION` | |
| `zone_redundant_enabled` | Zone-spread the primary region's replicas (+25% RU) | `false` | |
| `secondary_region_enabled` | Read region + automatic failover (doubles cost) | `false` | |
| `secondary_region` | The failover region (paired region recommended) | `eastus2` | |
| `vnet_cidr` / `pe_subnet_cidr` | Network layout | `10.40.0.0/16` / `10.40.1.0/24` | |
| `ops_email` | Alert recipient | `ops@example.com` | yes |
| `log_retention_days` | Workspace retention | `30` | |

## After deploying

1. **Bind the identity to the workload** — assign `{env}-cosmos-app-identity` to the compute running the API (AKS workload identity, App Service identity assignment, VM/VMSS identity). The identity's `client_id` output is what the application configures.
2. **Construct the client keylessly** — `CosmosClient(endpoint, tokenCredential)` in every SDK; the account `endpoint` output carries the URI. There is no key parameter, and key-based calls are rejected by the account.
3. **Bring applications into the network** — the data plane only answers through the private endpoint. Run workloads in `{env}-cosmos-vnet`, or peer their VNet in **and** link it to the private DNS zone (resolution does not cross peering by itself).
4. **Watch the first weeks of RU consumption** — the `Requests` metrics and `CDBDataPlaneRequests` logs in the workspace show per-operation RU charges; tune `autoscale_max_throughput` from data.

## Day 2

- **More containers** — add `AzureCosmosdbSqlContainer` resources per workload, each with its own partition key and autoscale ceiling. The account-wide role definition already covers them; add narrower grants if a workload should see only its container.
- **Narrow the grant** — the assignment's `scope` accepts a literal database (`{account-id}/dbs/{db}`) or container (`.../colls/{container}`) path when one identity should not see everything.
- **Tune the indexing posture** — the container indexes everything except `_etag` (Azure's default, stated explicitly). Write-heavy payload subtrees belong in `excludedPaths`; composite indexes serve multi-property ORDER BY.
- **Point-in-time restore** — continuous backup restores to any second in the last 7 days (`CONTINUOUS_30_DAYS` is a one-field upgrade). Restores create a NEW account — plan the DNS/endpoint swap as part of the drill.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
