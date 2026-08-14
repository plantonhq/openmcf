---
title: "Team Workspace"
description: "This preset creates one AI Foundry project with a system-assigned identity -- the standard per-team workspace on a shared `AIServices` account: the team's agents, evaluations and files live here,..."
type: "preset"
rank: "01"
presetSlug: "01-team-workspace"
componentSlug: "cognitive-account-project"
componentTitle: "Cognitive Account Project"
provider: "azure"
icon: "package"
order: 1
---

# Team Workspace

This preset creates one AI Foundry project with a system-assigned identity -- the standard per-team workspace on a shared `AIServices` account: the team's agents, evaluations and files live here, isolated from sibling projects.

## When to Use

- One workspace per team on a platform-owned account
- Environment separation (a dev project beside a prod project)
- The default starting shape for any AI Foundry adoption

## Key Configuration Choices

- **System-assigned identity** -- created and rotated with the project; grant storage/search access to its principal ID from the outputs
- **Display name and description set at creation** -- ARM cannot clear either in place later (clearing replaces the project)
- **Region matches the account's** by convention

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-cognitive-account-id>` | ARM ID of the parent account (kind AIServices, project management enabled) | `AzureCognitiveAccount` status outputs (`cognitive_account_id`), or reference it with valueFrom |
