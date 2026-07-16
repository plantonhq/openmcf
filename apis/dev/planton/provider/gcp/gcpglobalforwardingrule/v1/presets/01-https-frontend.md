# HTTPS Frontend VIP

The serving half of a production frontend: a reserved static IP bound on port 443 to a target HTTPS proxy, on the envoy-based `EXTERNAL_MANAGED` global external Application Load Balancer (advanced traffic management, the modern default for new deployments).

## When to Use

- Any new internet-facing HTTPS application — this is the rule DNS points at
- With preset 02 sharing the same `ipAddress` to complete the http→https story

## Remix Notes

- Reference the `GcpGlobalAddress` (its `address` output) and `GcpTargetHttpsProxy` via `valueFrom` instead of literals — the VIP then survives any frontend rebuild.
- Choose `EXTERNAL` instead only when a feature you need (e.g. the backend-bucket migration canary source state) is still classic-only; the scheme must match what the backend services were created for.
- `target` is the one mutable field: repoint it at a new proxy for a zero-downtime frontend swap — DNS never changes.
