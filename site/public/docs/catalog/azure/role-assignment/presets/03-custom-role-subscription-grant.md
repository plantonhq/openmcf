---
title: "Custom Role at Subscription Scope"
description: "This preset assigns a custom role -- referenced by its role definition ID, the exact form custom roles require -- to a group at subscription scope. This is the governance-team pattern: define a..."
type: "preset"
rank: "03"
presetSlug: "03-custom-role-subscription-grant"
componentSlug: "role-assignment"
componentTitle: "Role Assignment"
provider: "azure"
icon: "package"
order: 3
---

# Custom Role at Subscription Scope

This preset assigns a custom role -- referenced by its role definition ID, the
exact form custom roles require -- to a group at subscription scope. This is
the governance-team pattern: define a tailored role once (e.g. a cost auditor
that can read billing and usage but touch nothing), then grant it to a security
group so membership changes in Azure AD control access without further
deployments.

## When to Use

- Assigning organization-specific custom roles (built-in names cannot address them)
- Subscription- or management-group-wide governance grants
- Group-based access management (grant once to the group, manage membership in
  Azure AD -- also the mitigation for the 4,000-assignments-per-subscription limit)

## Key Configuration Choices

- **Role by ID** -- custom roles live at a scope and are addressed by their
  fully-scoped definition ID; the built-in `roleDefinitionName` lookup does not
  apply to them
- **Explicit principalType** -- declared here because group grants under
  ABAC-constrained delegated administrators require the type on the request;
  it is also self-documenting
- **Scope breadth** -- subscription scope inherits to every resource group and
  resource beneath it; prefer a management group only when the grant genuinely
  spans subscriptions

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<subscription-id>` | Target subscription GUID | `az account show --query id` |
| `<role-definition-guid>` | The custom role definition's GUID | `az role definition list --custom-role-only true` |
| `<principal-object-id>` | Azure AD OBJECT ID of the group | Entra ID portal → Groups → Object ID |
| `<why-this-grant-exists>` | Audit note shown in the portal's IAM blade | Your runbook / change ticket |
