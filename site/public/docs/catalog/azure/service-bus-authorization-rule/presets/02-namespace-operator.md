---
title: "Namespace Operator Credential"
description: "This preset mints a namespace-wide manage credential -- for platform tooling that creates and deletes entities at runtime (dynamic queue provisioning, tenant onboarding automation). A deliberate,..."
type: "preset"
rank: "02"
presetSlug: "02-namespace-operator"
componentSlug: "service-bus-authorization-rule"
componentTitle: "Service Bus Authorization Rule"
provider: "azure"
icon: "package"
order: 2
---

# Namespace Operator Credential

This preset mints a namespace-wide manage credential -- for platform
tooling that creates and deletes entities at runtime (dynamic queue
provisioning, tenant onboarding automation). A deliberate, named
alternative to using the root key.

## When to Use

- Tooling that provisions Service Bus entities dynamically
- Replacing ad-hoc use of `RootManageSharedAccessKey` with a named,
  revocable credential

## Key Configuration Choices

- **`manage: true` with `listen` + `send`** -- Azure requires the full
  trio; manage is a superset, never standalone
- **A dedicated name** -- revoking or rotating this credential never
  touches the root rule or application credentials

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-bus` | The AzureServiceBusNamespace's Planton resource name | Your messaging composition |
