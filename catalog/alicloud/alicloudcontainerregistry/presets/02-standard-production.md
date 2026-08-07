# Standard Production Registry

This preset creates a Standard-edition Container Registry instance on a 12-month subscription, with separate private namespaces for platform, backend, and frontend images. It fits a production team that organizes images by domain.

## When to Use

- Production image hosting for a multi-team or multi-domain codebase
- Workloads that justify a subscription commitment over pay-as-you-go
- Teams that want namespace-level separation without enterprise features

## Key Configuration Choices

- **Standard edition** (`instanceType: Standard`) — higher repository and namespace quotas than Basic
- **12-month subscription** (`paymentType: Subscription`, `period: 12`) — predictable cost for steady production usage
- **Domain-scoped namespaces** (`platform`, `backend`, `frontend`) — images are grouped by owning area; `frontend` has `autoCreate: false` so its repositories are created deliberately
- **Private everywhere** (`defaultVisibility: PRIVATE`) — production images are never publicly pullable by default

## Placeholders to Replace

- `metadata.name` and `instanceName` — your registry's name
- `region` — the AliCloud region of your production clusters (e.g., `cn-shanghai`)
- `namespaces` — rename or extend to match your team's image organization
