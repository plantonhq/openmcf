---
title: "External DNS"
description: "External DNS deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesexternaldns"
---

# External DNS on Kubernetes

Publishes DNS records for your cluster's Services, Ingresses and Gateway API
routes into a real DNS provider — so a hostname exists the moment the workload
does, and disappears with it.

## What it does

Installs the ExternalDNS controller from the official chart. The controller
watches the Kubernetes objects you nominate, reads the hostnames they carry, and
reconciles matching records in your DNS zone.

## The dual-provider model

Kubernetes here is deliberately two-sided: the cluster runs in one environment
while the zone often lives in another. An EKS cluster publishing into Cloudflare
is an ordinary case, not an exception. The spec keeps the halves separate:

- **Where records go** — the DNS provider arm (Route 53, Cloud DNS, Azure DNS,
  Cloudflare, a webhook provider, or the in-memory sandbox). Exactly one.
- **How the controller authenticates** — a workload identity for the cloud
  providers, a token for Cloudflare, or per-provider static credentials.

One installation manages one provider. Clusters publishing to several deploy
several instances of this component; give each its own owner ID so their
registries never fight.

## Choosing a provider arm

| Arm | Authentication | Notes |
|---|---|---|
| Route 53 | EKS IRSA (preferred), node role, or static keys | Cross-account zones through an assumed role |
| Cloud DNS | GKE Workload Identity (preferred), node SA, or a key | Project is required |
| Azure DNS | AKS Workload Identity, managed identity, or a service principal | Public or private zones |
| Cloudflare | API token, always | Works from any cluster — the cross-cloud arm |
| Webhook | The provider's own image, as a sidecar | Upstream's extension path for everything not built in |
| In-memory | None | Sandbox: records live in the pod and vanish with it |

## The two ways this goes wrong

**Writing to the wrong zone.** With no domain filter the controller may act on
every zone its credentials can see. In a shared cloud account that is how one
cluster starts rewriting another team's records. Set zone filters on the
provider arm, domain filters, or both.

**Fighting another instance.** The registry decides which records belong to this
installation. Under the `sync` policy the controller DELETES records it believes
it owns — so two instances sharing an owner ID will each clean up the other's
work. Every instance sharing a zone needs a distinct `txt_owner_id`.

## Policy at a glance

| Policy | Creates | Updates | Deletes |
|---|---|---|---|
| `upsert-only` (default) | yes | yes | no |
| `sync` | yes | yes | yes — records this instance owns |
| `create-only` | yes | no | no |

## Works with

| Kind | Relationship |
|---|---|
| KubernetesNamespace | Where the controller runs |
| KubernetesService, KubernetesIngress | Common sources — their hostnames become records |
| Gateway API routes | Additional sources when hostnames live on routes |
| AwsRoute53Zone, GcpDnsZone, AzureDnsZone, CloudflareDnsZone | Zone filters, referenced by output |
| AwsIamRole, GcpServiceAccount, AzureUserAssignedIdentity | The keyless identity the ServiceAccount federates with |
| GcpProject | The project owning Cloud DNS zones |

## Outputs

| Output | Description |
|---|---|
| `namespace` | The namespace the controller runs in |
| `release_name` | The Helm release name |
| `service_account_name` | The controller's ServiceAccount |

The ServiceAccount name is not decoration: together with the namespace it is
exactly what an IAM role trust policy, a Google service-account binding or an
Azure federated credential references when you wire keyless authentication.
