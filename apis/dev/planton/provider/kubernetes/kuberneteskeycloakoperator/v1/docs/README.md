# Kubernetes Keycloak Operator — design notes

## Grain

One resource = the official Keycloak Operator, installed from the
keycloak-k8s-resources release manifests at the pinned tag. Keycloak
ships NO official Helm chart — the operator is the first-party
Kubernetes distribution, which is what rules out the alternatives:
Bitnami's images went commercial (post-August-2025 subscription) and
the community charts are installers, not lifecycle managers. The
bundle is 16 plain-YAML documents (ServiceAccount, RBAC, Service, the
operator Deployment) plus 4 CRDs — no admission webhooks, no
cert-manager dependency, no install hooks — applied as ordered groups
so creation and teardown order correctly by construction. Every
resource carries upstream's fixed names: exactly one install per
namespace.

## The two-kind grain

The operator (this kind) and the server declaration
(`KubernetesKeycloak`) are separate resources — the same split as the
Tekton pair: the manager installs once and reconciles many
declarations; installing the operator alone deploys no server. The
declaration kind carries the product surface (database, TLS,
hostname, update strategy); this kind carries only the manager's own
posture (watch scope, images, resources, scheduling).

## Watch scope

Upstream publishes namespaced and cluster-wide variants whose ONLY
differences are the RBAC binding kinds (RoleBindings vs
per-controller ClusterRoleBindings) and the JOSDK env values
(`JOSDK_ALL_NAMESPACES`) — so the spec models the choice as one
boolean, `cluster_wide`. Namespaced (the default) keeps the operator
and its declarations in one namespace and lets teams run isolated
stacks side by side; cluster-wide serves every namespace but runs at
most once per cluster and — an upstream constraint — refuses custom
ServiceAccounts on Keycloak pod templates.

## No version field

The ruling: `KubernetesKeycloak` renders the Keycloak CR against the
CRD schema THIS bundle installs. A selectable operator version would
let the installed schema drift away from what the declaration kind
renders — so the module pins the release, and upgrades arrive as
module updates that move both sides together. The image overrides
(`operator_image`, and `default_keycloak_image` → the bundle's
`RELATED_IMAGE_KEYCLOAK`, the default server image for every CR that
sets none) exist for mirrors, keeping the pinned tags for the same
reason.

## The namespace stamp

Upstream expects kustomize to set the namespace: the module stamps
the target namespace onto every namespaced document and every RBAC
binding subject — including the one ClusterRoleBinding subject
upstream bakes to the `keycloak` namespace.

## Destroy

The 4 CRDs delete with the resource, which cascade-deletes any
Keycloak declarations still in the cluster — destroy
`KubernetesKeycloak` resources first, the operator after.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
