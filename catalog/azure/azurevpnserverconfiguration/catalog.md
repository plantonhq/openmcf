# Azure VPN Server Configuration

Deploys a VPN Server Configuration -- the reusable "who may connect and how" policy for point-to-site VPN: the authentication methods remote users sign in with (Entra ID, certificate, RADIUS), the trusted and revoked certificates, the tunnel protocols offered, and optional policy groups for user segmentation. The configuration is free and deploys in seconds; a Point-to-Site VPN Gateway is born pointing at one, and many gateways can share it. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VPN Server Configuration** -- the authentication policy object (Entra ID / certificate / RADIUS parameters, IPsec proposal, tunnel protocols)
- **Policy Groups** (optional) -- one ARM child per `policyGroups` entry, keyed by name (the `policy_group_ids` output republishes each group's ARM ID)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **The authentication material for each enabled type**: the Entra ID tenant/audience/issuer values for "AAD", a root certificate's public data for "Certificate", or reachable RADIUS servers and their shared secrets for "Radius".

## Deploy

### Console

Open the deployment store, find **Azure VPN Server Configuration**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Entra ID Remote Workforce** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVpnServerConfiguration
metadata:
  name: remote-workforce
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: network-rg
      fieldPath: status.outputs.resource_group_name
  name: remote-workforce
  vpnAuthenticationTypes:
    - AAD
  aadAuthentication:
    audience: 41b23e61-6c1e-4545-b367-cd054e0ed4b4
    issuer: https://sts.windows.net/<your-tenant-id>/
    tenant: https://login.microsoftonline.com/<your-tenant-id>
  vpnProtocols:
    - OpenVPN
```

```shell
planton apply -f azure-vpn-server-configuration.yaml
```

The configuration provisions in seconds and is free.

### InfraChart

In a remote-access chart the order is: WAN → hub → **server configuration** → point-to-site gateway, the gateway wiring to the configuration by reference.

## Key Configuration

These are the most important decisions when configuring the policy. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Authentication types** -- the central choice. "AAD" (Entra ID) gives managed sign-in with conditional access; "Certificate" works offline against your own root; "Radius" reuses existing enterprise auth. Each enabled type requires its block -- the spec rejects a type without its parameters upfront.

**Tunnel protocols** -- leave `vpnProtocols` empty for Azure's default. Set `["OpenVPN"]` when you use policy groups or plan multiple gateway address pools -- both require OpenVPN.

**Policy groups** -- named member-matching rules (Entra ID group ID, certificate common name, RADIUS group ID). The gateway maps groups to address pools; each group's ARM ID surfaces in `policy_group_ids` keyed by the group's name.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpn_server_configuration_id` | ARM ID of the configuration | A point-to-site gateway's `vpnServerConfigurationId` |
| `vpn_server_configuration_name` | Name of the configuration | Operational tooling |
| `policy_group_ids` | ARM ID of each policy group, keyed by group name | Group-aware automation (`status.outputs.policy_group_ids.<group-name>`) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Entra ID remote workforce** -- managed sign-in, no certificate distribution. Start from the **Entra ID Remote Workforce** preset.

**Certificate auth with policy groups** -- offline-capable auth with user segmentation. Start from the **Certificate Auth with Policy Groups** preset.

## Works With

- [**Azure Point-to-Site VPN Gateway**](/cloud-catalog/azure-point-to-site-vpn-gateway) -- the hub gateway that attaches this policy
- [**Azure Virtual Hub**](/cloud-catalog/azure-virtual-hub) -- where the gateway lives
- [**Azure Virtual WAN**](/cloud-catalog/azure-virtual-wan) -- the managed network umbrella
