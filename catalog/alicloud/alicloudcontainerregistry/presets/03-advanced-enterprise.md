# Advanced Enterprise Registry

This preset creates an Advanced-edition Container Registry instance on a 12-month subscription, with internal, shared, and deliberately public namespaces. It suits an organization that distributes some images publicly while keeping internal images locked down.

## When to Use

- Enterprise image hosting with high repository counts and advanced features
- Organizations that publish public images (SDKs, agents, samples) alongside private ones
- Platform teams providing a shared registry to many internal consumers

## Key Configuration Choices

- **Advanced edition** (`instanceType: Advanced`) — the enterprise tier with the highest quotas and feature set
- **12-month subscription** (`paymentType: Subscription`, `period: 12`) — commitment pricing for a long-lived platform service
- **Three-way namespace split** — `internal` (auto-created, private), `shared` (private, repositories created deliberately), and `public-images` (`defaultVisibility: PUBLIC`) for content that is meant to be pulled anonymously
- **Deliberate repository creation** on `shared` and `public-images` (`autoCreate: false`) — nothing lands in the shared or public space by accident

## Placeholders to Replace

- `metadata.name` and `instanceName` — your registry's name
- `region` — the AliCloud region of your infrastructure (e.g., `cn-beijing`)
- `namespaces` — adjust the split to your organization's publishing model
