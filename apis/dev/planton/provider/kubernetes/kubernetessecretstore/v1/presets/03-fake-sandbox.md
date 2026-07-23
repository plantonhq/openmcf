# Fake Backend Sandbox (Test-Only)

This preset creates a namespaced store backed by ESO's built-in fake backend: the store serves the literal key/value entries declared in the spec — no external account, no network, fully deterministic. **Test-only**: the values sit in plain text in the store spec, so never put real secrets here. It exists to exercise the full sync machinery (store → ExternalSecret → materialized Secret) in pipelines and sandboxes.

## When to Use

- CI pipelines and kind-cluster tests that need the ExternalSecret sync loop without any cloud account
- Evaluating the External Secrets family before wiring a real backend
- Local development namespaces where deterministic placeholder values are enough

## Key Configuration Choices

- **Fake backend** (`fake.data`) -- each entry's `key` is what an ExternalSecret's `remoteRef.key` addresses; the `value` is served verbatim
- **No authentication** -- there is nothing to authenticate to; the operator serves the entries itself
- **Namespaced grain** -- keeps the sandbox contained; the fake store and its consumers live in one namespace

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<sandbox-namespace>` | The namespace for the sandbox (store and ExternalSecrets live here) | `kubectl get namespaces` or `KubernetesNamespace` outputs |

The `data` entries are examples — replace keys and values with whatever your test consumers expect.

## Related Presets

- **01-team-gcp-secret-manager** -- The real-backend equivalent once you move past the sandbox
- **02-vault-approle** -- Use when the team's secrets live in Vault/OpenBao
