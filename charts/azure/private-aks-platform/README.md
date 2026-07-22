# Private AKS Platform

Production Kubernetes on Azure, private and keyless. One deploy produces
an AKS cluster whose API server is unreachable from the internet, whose
nodes egress through one stable NAT address, whose images pull from a
private registry with no password anywhere, and whose pods read Key Vault
secrets through workload identity federation -- the full secret-less
chain, wired end to end. It is the AKS posture every platform team wants
and few assemble correctly, because the pieces only work when they
reference each other exactly right.

## The architecture

Everything is keyless, and everything is referenced:

- **The cluster** deploys with a private API server (AKS manages the
  private DNS zone), Entra-integrated authentication with Azure RBAC
  authorization, and workload identity enabled on the cluster's OIDC
  issuer. It carries only its mandatory system pool -- tainted
  `CriticalAddonsOnly` so application pods cannot starve CoreDNS.
- **The workload pool** is a first-class `AzureAksNodePool`: it scales,
  upgrades, and rotates without touching the control plane, and it can
  flip to spot capacity (`workload_pool_spot_enabled`) for fault-tolerant
  workloads.
- **The network is yours, not AKS-managed.** Nodes live in a dedicated
  subnet with implicit outbound access off; every outbound connection
  leaves through an explicit NAT gateway and one static public IP --
  the address partners allowlist.
- **Image pulls are keyless.** The registry's admin account stays off;
  the cluster's kubelet identity holds exactly `AcrPull`, granted against
  the registry's own ID output. A compromised node can pull images --
  never push them.
- **Pod-to-Azure access is keyless.** A dedicated workload identity holds
  `Key Vault Secrets User` on the platform vault (RBAC mode), and a
  federated identity credential trusts exactly one namespace/service-
  account pair against the cluster's OIDC issuer URL -- referenced from
  the cluster's own output, so the trust can never disagree with the
  cluster.
- **Telemetry**: control-plane logs (API server, admin audit, schedulers,
  autoscaler, Entra authorization) and Container Insights stream into the
  platform's Log Analytics workspace; applications trace into the
  workspace-based Application Insights beside it.

## What is on by default

- **Private API server** (`private_api_enabled`): the endpoint gets only
  a private IP -- reach it via peering, VPN, or ExpressRoute. Disabling
  it opens a public endpoint; gate that with `api_authorized_ip_ranges`.
- **Workload identity + OIDC issuer**: on, and the federated credential
  below depends on them.
- **Azure RBAC for Kubernetes**: cluster authorization is Azure role
  assignments, not hand-kept RoleBindings. The local admin account stays
  available until `local_account_disabled` is raised -- do that once
  `entra_admin_group_object_ids` (or equivalent grants) are in place.
- **No implicit outbound**: the node subnet sets
  `defaultOutboundAccessEnabled: false`; the NAT gateway is the only
  exit.
- **NetworkPolicy enforcement** (Azure engine): policies teams write are
  enforced, not silently inert.
- **STANDARD control plane** (`cluster_sku_tier`): the 99.95% SLA tier;
  zones on (`zones_enabled`) is what that SLA assumes.

## Parameters worth understanding

- **`registry_name` and `key_vault_name`** are GLOBALLY unique across
  Azure (they become DNS names). Change both before deploying -- the
  defaults will collide with any other deployment of this chart.
- **`kubernetes_version`**: empty deploys the latest AKS-recommended GA
  version. Pin one for production so upgrades happen when you choose;
  upgrades move one minor at a time.
- **`workload_namespace` / `workload_service_account`**: the one
  Kubernetes workload trusted to act as the platform's workload identity.
  Additional workloads should get their own
  `AzureUserAssignedIdentity` + `AzureFederatedIdentityCredential` +
  grant triple -- one auditable trust per workload, composed from the
  same first-class resources.
- **`key_vault_purge_protection`**: the production posture, but
  irreversible -- once on, a deleted vault's name stays locked for the
  soft-delete retention window. Raise it when the vault carries anything
  real.
- **`api_authorized_ip_ranges`** only applies to a PUBLIC API server;
  a private cluster ignores it by design.

## After deployment

The cluster takes 10-15 minutes to provision; the full platform
typically lands in under 25.

- **Wire the workload**: annotate the trusted service account and label
  its pods --

  ```yaml
  apiVersion: v1
  kind: ServiceAccount
  metadata:
    name: <workload_service_account>
    namespace: <workload_namespace>
    annotations:
      azure.workload.identity/client-id: <the identity's client_id output>
  ```

  with `azure.workload.identity/use: "true"` on the pod template. The
  Azure SDK inside those pods then reaches Key Vault with no secret.
- **Push an image** to the registry (`az acr login` with your own Entra
  identity) and deploy it -- nodes pull it through the kubelet identity's
  `AcrPull` grant.
- **Reach a private API server** from inside the network (a peered hub,
  VPN, or a jump host in the VNet); `az aks command invoke` works from
  anywhere as the audited break-glass path.
- **Natural next steps**: peer the platform VNet into a hub-spoke
  foundation (its NAT egress then gives way to firewall egress by
  attaching a route table to the node subnet and switching the cluster's
  outbound type to USER_DEFINED_ROUTING); add per-workload identities as
  applications multiply; and take the registry and vault private with
  `AzurePrivateEndpoint`s once the network has private DNS zones for
  them.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
