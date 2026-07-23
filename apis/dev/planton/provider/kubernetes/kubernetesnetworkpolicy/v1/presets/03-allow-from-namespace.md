# Allow From Namespace

This preset allows inbound traffic to one workload's pods from ANY pod in a specific other namespace, on one port. Selecting the source namespace uses the automatic `kubernetes.io/metadata.name: <name>` label that Kubernetes stamps on every namespace — a guaranteed by-name handle that needs no labelling convention of your own.

## When to Use

- Granting a known cross-namespace caller (an ingress controller, a monitoring stack, another team's service) access to one workload
- As the targeted exception layered on top of **01-default-deny-all** or **02-allow-same-namespace**
- Any time "namespace X may call service Y on port Z" is the requirement

## Key Configuration Choices

- **`pod_selector.match_labels.app`** — targets one Planton workload's pods; every workload kind stamps `app: <workload-metadata-name>` on its pods as immutable selection identity
- **`namespace_selector` alone in the peer** — allows ALL pods in the matched namespaces. To restrict to specific pods within that namespace, add a `pod_selector` in the SAME peer entry: the two selectors then AND (matching pods in matching namespaces). Listing them as two separate peers would OR instead — any pod in the namespace, plus same-namespace pods matching the selector — which is the classic NetworkPolicy authoring mistake
- **`ports` restricts the allow to TCP 8080** — the peer AND the port must both match; remove the list to allow all ports
- **Additive** — this policy only widens; combine freely with the other presets governing the same pods

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace of the workload being protected | Your namespace management |
| `<your-workload-name>` | The target workload's `metadata.name` (its pods carry `app: <name>`) | The workload's manifest |
| `<source-namespace>` | The namespace whose pods are allowed in | Your namespace management |

The port `"8080"` is a working example — replace it with the target workload's actual container port.

## Related Presets

- **01-default-deny-all** — the baseline this exception is designed to punch through
- **02-allow-same-namespace** — namespace-internal traffic, the complementary allow
