# Azure Firewall Policy Rule Collection Group

Deploys an Azure Firewall Policy Rule Collection Group — an ordered document of rule collections (application, network, and DNAT) nested inside a firewall policy, enforced by every firewall that attaches the policy. A policy carries many groups — typically one per team or per application — each deployed and updated independently, which is exactly why the group is its own resource rather than a fold inside the policy: the security team's baseline and each application's rules move on their own schedules.

**Processing order**: groups evaluate by group priority, collections by collection priority WITHIN a type — and across types Azure always processes DNAT rules first, then network rules, then application rules, regardless of priorities. Lower numbers run first.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Rule Collection Group** -- the ordered rule document with its application, network, and DNAT collections, nested under the parent policy

The group is an ARM CHILD of the policy: it has no region, resource group, or tags of its own (the policy owns placement), and individual rules have no ARM identity — the group is the unit of deployment, so rules travel with it as one document. An **empty group is legal**: a declare-then-fill anchor whose rules land later.

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Firewall Policy** for the group to nest under — reference an AzureFirewallPolicy Cloud Resource via ValueFromRef.
- **A priority plan**: leave gaps between groups (100, 200, 300…) so future documents slot in without renumbering; the security baseline takes the low numbers so nothing outranks it.
- **For FQDN network rules and FQDN translation targets**: the parent policy's DNS proxy must be enabled.

## Deploy

### Console

Open the deployment store, find **Azure Firewall Policy Rule Collection Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Egress Baseline Rules** preset in the [Presets](#presets) tab for the classic allow-list document.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFirewallPolicyRuleCollectionGroup
metadata:
  name: payments-app
  org: acme-corp
  env: prod
spec:
  firewallPolicyId:
    valueFrom:
      name: egress-baseline
  name: payments-app
  priority: 1000
  applicationRuleCollections:
    - name: allow-saas
      priority: 1000
      action: ALLOW
      rules:
        - name: allow-github
          protocols:
            - type: HTTPS
              port: 443
          sourceAddresses:
            - "10.0.0.0/16"
          destinationFqdns:
            - "*.github.com"
```

```shell
planton apply -f rule-collection-group.yaml
```

Every firewall attached to `egress-baseline` now enforces the document — and the payments team redeploys it without touching the policy or anyone else's rules.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the group to its dependencies:

```yaml
spec:
  firewallPolicyId:
    valueFrom:
      kind: AzureFirewallPolicy
      name: egress-baseline
      fieldPath: status.outputs.firewall_policy_id
  networkRuleCollections:
    - name: core-egress
      priority: 500
      action: ALLOW
      rules:
        - name: trusted-branches
          protocols: [TCP]
          sourceIpGroups:
            - valueFrom:
                kind: AzureIpGroup
                name: branch-offices
                fieldPath: status.outputs.ip_group_id
          destinationAddresses: ["*"]
          destinationPorts: ["443"]
```

The InfraPipeline resolves the dependency graph, deploys the policy (and any referenced IP Groups) first, then provisions the group.

## Key Configuration

These are the most important decisions when configuring a Rule Collection Group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Priorities** -- 100-65000 at BOTH levels: the group's priority orders it among the policy's groups; each collection's priority orders it among the group's collections of the same type. Across types, DNAT → network → application always. Priorities update in place.

**Application rules** -- L7 destinations: FQDNs (SNI-matched; the durable egress grammar), Azure FQDN tags ("AzureKubernetesService" — Microsoft tracks the endpoints), URLs and web categories (Premium; HTTPS paths need `terminateTls` and the policy's TLS certificate). Every rule needs at least one source and one destination (ARM's apply-time contract).

**Network rules** -- protocol/address/port filtering with Azure vocabulary: service tags, IP Groups on BOTH sides (the only rule type with destination groups), and FQDN destinations (requires the policy's DNS proxy).

**DNAT rules** -- publish internal services: one public arrival (firewall public IP + a single port entry — ARM's current cap) translates to exactly one internal target (address XOR FQDN) and port. TCP/UDP only, and a match implicitly ALLOWS the flow — the DNAT rule is the whole permission.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFirewallPolicy** | `firewallPolicyId` | `status.outputs.firewall_policy_id` |
| **AzureIpGroup** | every rule type's `sourceIpGroups[]` (network rules also `destinationIpGroups[]`) | `status.outputs.ip_group_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `rule_collection_group_id` | Azure Resource Manager ID of the group | Automation scripts, inventory |
| `rule_collection_group_name` | Name of the group | Automation scripts, inventory |

The group is a composition LEAF: nothing references a rule document — enforcement flows entirely through the parent policy's firewall attachments.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Egress baseline rules** -- the security team's allow-list document: network rules for core protocols, application rules for sanctioned SaaS. Start from the **Egress Baseline Rules** preset.

**Published service (DNAT)** -- one public-to-internal translation with a tightly-scoped source. Start from the **DNAT Publish Service** preset.

## Works With

- [**Azure Firewall Policy**](/cloud-catalog/azure-firewall-policy) -- the parent this rule document nests under
- [**Azure Firewall**](/cloud-catalog/azure-firewall) -- the enforcement instance that attaches the parent policy
- [**Azure IP Group**](/cloud-catalog/azure-ip-group) -- named address sets the rules reference (source and destination)
