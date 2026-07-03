---
title: "Managed Identity Grant on a Resource Group"
description: "This preset grants a user-assigned managed identity a built-in role on a resource group -- the most common authorization pattern in composed Azure environments: the identity a workload runs as gets..."
type: "preset"
rank: "01"
presetSlug: "01-identity-resource-group-grant"
componentSlug: "role-assignment"
componentTitle: "Role Assignment"
provider: "azure"
icon: "package"
order: 1
---

# Managed Identity Grant on a Resource Group

This preset grants a user-assigned managed identity a built-in role on a
resource group -- the most common authorization pattern in composed Azure
environments: the identity a workload runs as gets exactly the access it needs
on the environment it lives in.

Both references resolve through their default kinds: `scope` targets an
`AzureResourceGroup`'s ARM ID and `principalId` targets an
`AzureUserAssignedIdentity`'s principal ID, so only the resource names are
needed. `skipServicePrincipalAadCheck` is pre-set for the composed case where
the identity is created in the same deployment (freshly created principals
replicate through Azure AD asynchronously; the flag avoids minutes of
PrincipalNotFound retries).

## When to Use

- Granting a workload identity access to the resource group its dependencies live in
- Granting a CI/CD deploy identity scoped rights on an environment
- Any identity + grant pair created together in one infra chart

## Key Configuration Choices

- **Role** (`roleDefinitionName`) -- prefer the narrowest built-in role that
  satisfies the need: "Reader" for observation, "Contributor" for resource
  management (never grants RBAC rights), or a service-specific data-plane role
- **Scope** -- a resource group inherits the grant to everything inside it;
  switch the reference to a single resource (explicit `kind` + `fieldPath`)
  for tighter scoping

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | Metadata name of the AzureResourceGroup being granted on | Your infra chart / resource list |
| `<built-in-role-name>` | Built-in role, e.g. `Reader`, `Contributor`, `AcrPull` | [Azure built-in roles](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles) |
| `<your-managed-identity>` | Metadata name of the AzureUserAssignedIdentity being granted | Your infra chart / resource list |
| `<why-this-grant-exists>` | Audit note shown in the portal's IAM blade | Your runbook / change ticket |
