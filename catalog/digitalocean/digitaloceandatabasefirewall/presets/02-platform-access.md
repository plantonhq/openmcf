# Platform Workloads Only

This preset trusts a Kubernetes cluster and an App Platform app -- both by reference -- and nothing else. The database becomes reachable exclusively from platform-managed compute, with zero IP management.

## When to Use

- Databases consumed only by DOKS workloads and App Platform services
- Charts wiring database + compute + firewall as one deployable unit
- Environments where IP allowlists rot (NAT changes, autoscaling)

## Key Configuration Choices

- **References, not IDs** -- `valueFrom` resolves each platform resource's id at deploy time, so the rule set always points at the live resource, never a stale hand-copied UUID.
- **No ipRules at all** -- valid by design (the spec requires at least one rule across ALL lists); adding an operator bastion later is one `ipRules` entry.

## What You Get

A cluster whose trusted sources are platform identities. DigitalOcean tracks the members (cluster nodes, app instances) automatically as they scale.
