---
title: "Application Configuration"
description: "This preset creates a ConfigMap carrying a typical application's configuration: a couple of scalar settings (consumed as environment variables via `envFrom`/`configMapKeyRef`) and a properties file..."
type: "preset"
rank: "01"
presetSlug: "01-app-config"
componentSlug: "configmap"
componentTitle: "ConfigMap"
provider: "kubernetes"
icon: "package"
order: 1
---

# Application Configuration

This preset creates a ConfigMap carrying a typical application's configuration: a couple of scalar settings (consumed as environment variables via `envFrom`/`configMapKeyRef`) and a properties file (consumed as a mounted file via a `configMap` volume).

## When to Use

- Externalizing application settings from container images (log levels, timeouts, feature flags)
- Shipping a file-shaped config (`application.properties`, `nginx.conf`, `config.yaml`) that the app reads from disk
- Any non-confidential configuration — for passwords, tokens, or keys use KubernetesSecret instead

## Key Configuration Choices

- **Mixed key shapes** — env-var-style keys (`LOG_LEVEL`, `CACHE_TTL_SECONDS`) for environment consumption, a file-style key (`application.properties`) for volume mounting; shape each key after how it will be consumed
- **Mutable by default** — edits propagate to volume mounts eventually (kubelet sync period) but never to environment variables without a pod restart; switch to the immutable-versioned preset for production rollout discipline
- **All values are strings** — YAML numbers must be quoted (`"300"`); the Kubernetes `data` map is `string → string`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace for the ConfigMap (must match the consuming workload's namespace) | Your namespace management |
| `<your-log-level>` | Application log level, e.g. `info`, `debug`, `warn` | Your application's configuration reference |
| `<your-cache-ttl>` | Cache TTL in seconds, e.g. `"300"` (quoted — all values are strings) | Your application's configuration reference |

The `application.properties` content is a working example — replace it with your application's actual file content.

## Related Presets

- **02-immutable-versioned** — locked, versioned config for production roll-forward
- **03-binary-payload** — base64-encoded binary entries alongside text
