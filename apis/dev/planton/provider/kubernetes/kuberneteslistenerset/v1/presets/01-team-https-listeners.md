# Team HTTPS Listeners

Attach a team-owned HTTPS listener to a centrally managed Gateway without
editing the Gateway itself. The platform team owns the Gateway (address, load
balancer, GatewayClass); each application team owns a ListenerSet in its own
namespace carrying its hostname and TLS certificate. This is the per-team
delegation pattern for shared gateways.

## When to Use

- A platform team runs one shared Gateway and application teams need to add
  their own hostnames and certificates.
- You want per-team TLS material to stay in the team's namespace (a ListenerSet
  can reference Secrets in its own namespace without a ReferenceGrant).
- Team routes should attach to the team's listener, not the Gateway at large.

## Key Configuration Choices

- **`parentRef.name`** -- the Gateway these listeners merge into. It is a foreign key: this preset uses `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) so the ListenerSet deploys after its Planton-managed Gateway; write `value: <literal>` instead for an externally created Gateway.
- **`tls.certificateRefs[0].name`** -- also a foreign key, wired with `valueFrom:` to a `KubernetesCertificate`'s exported Secret (`status.outputs.secret_name`), so the listener terminates with the cert-manager-issued certificate. Use `value:` for a Secret provisioned outside Planton.
- **`allowedRoutes.namespaces.from` (`Same`)** -- only routes in the ListenerSet's OWN namespace may attach (for ListenerSets, `Same` means the ListenerSet's namespace, not the Gateway's).
- **`protocol` (`HTTPS`) + `tls.mode` (`Terminate`)** -- the listener decrypts at the edge; HTTPRoutes attach for host/path routing.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`, standard
  channel as of v1.6 -- ListenerSet joined the standard channel in v1.5).
- The parent `KubernetesGateway` exists and opts in to ListenerSet attachment
  (`allowedListeners.namespaces.from` set to `All`, `Selector`, or `Same` --
  Gateways allow none by default).
- The team namespace exists (`KubernetesNamespace`).
- The `KubernetesCertificate` referenced by `certificateRefs` exists in the
  team namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<shared-gateway-resource-name>` | Name of the Planton `KubernetesGateway` resource this ListenerSet attaches to. |
| `<certificate-resource-name>` | Name of the Planton `KubernetesCertificate` whose issued Secret terminates the listener. |
| `team-a.example.com` | The team's hostname (a literal example value -- replace with your real host). |

Set `spec.namespace.value` and `parentRef.namespace` to your team and Gateway
namespaces, or replace the namespace with a `valueFrom` reference to a
`KubernetesNamespace`.
