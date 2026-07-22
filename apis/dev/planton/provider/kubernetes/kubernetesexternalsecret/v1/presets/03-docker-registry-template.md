# Docker Registry Pull Secret (Template)

This preset syncs registry credentials from the backend and TEMPLATES them into a `kubernetes.io/dockerconfigjson` Secret — the shape `imagePullSecrets` expects. The synced `username`/`password` keys feed the Go template; only the rendered `.dockerconfigjson` key lands in the Secret (the template's default Replace merge policy). The registry password stays in the backend as the single source of truth, and rotation flows to the cluster on the refresh interval.

## When to Use

- Pods pull images from a private registry and the registry credentials live in a secret backend
- You want pull-secret rotation driven from the backend, not hand-edited Secrets
- Any private registry: the auths JSON carries the registry server, so Docker Hub, GHCR, quay.io, or a self-hosted registry all work

## Key Configuration Choices

- **Typed Secret** (`template.type: kubernetes.io/dockerconfigjson`) -- kubelet only honors pull secrets of this type
- **Templated rendering** (`template.data`) -- a Go template over the synced keys builds the auths JSON; `{{ .username }}` and `{{ .password }}` refer to the `data` entries' `secretKey` names
- **Replace merge policy** (the default, deliberately) -- the raw `username`/`password` keys do NOT land in the Secret; only the rendered dockerconfigjson does
- **Explicit target name** (`target.name: registry-pull-secret`) -- the stable handle pods reference in `imagePullSecrets`; also exported as the `secret_name` output

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<app-namespace>` | Namespace the ExternalSecret and pull Secret live in | `kubectl get namespaces` or `KubernetesNamespace` outputs |
| `<cluster-secret-store-name>` | The store to sync from | `KubernetesClusterSecretStore` `store_name` output |
| `<registry-server>` | Registry host (e.g. `ghcr.io`, `https://index.docker.io/v1/` for Docker Hub) | Your registry's documentation |
| `<registry-credentials-backend-key>` | Backend entry holding `username` and `password` properties | Your secret backend's console/CLI |

## Related Presets

- **01-explicit-keys** -- The plain form, when no reshaping is needed
- **02-extract-json-document** -- Use to pull ALL fields of a structured entry at once
