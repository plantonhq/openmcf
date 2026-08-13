# Team Foundry Project

The standard team workspace inside a shared AI Foundry hub: one
project per team, with its own system-assigned identity so grants can
be scoped per team while everything else (vault, storage, network,
encryption) is inherited from the hub.

## When to Use

- Giving one team its own working space inside the shared hub
- Per-team access scoping (grants bind to the project's own identity)
- The default shape for every new AI team

## Key Configuration Choices

- `aiServicesHubId` references the hub by name -- swap `team-hub` for
  your AzureAiFoundry resource's name, or use `value:` with a literal
  ARM ID. Fixed at creation: a project cannot move between hubs.
- There is NO resource group here by design -- the project deploys
  into the hub's group (the provider derives it).
- `identity.type: SYSTEM_ASSIGNED` -- Azure manages the project's
  identity; grant it data access per the team's needs.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `team-hub` (valueFrom name) | Your AzureAiFoundry resource's name | The hub manifest's `metadata.name` |
| `name` / `friendlyName` / `description` | The team's own naming | Your team |
