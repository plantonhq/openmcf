# Producer Hierarchy Allowlist

A security-hardened policy that restricts which resource-hierarchy
locations producer instances may live in — only producers under the
listed organizations, folders, or projects can connect into the network.

## When to use

Regulated environments where the network team must guarantee that PSC
endpoints only ever route to producer instances inside the company's own
organization (or an explicitly reviewed set of projects).

## What to customize

- `allowedGoogleProducersResourceHierarchyLevels` — your organization
  number, or a narrower folders/projects list. Entries use the
  `projects/{id}`, `folders/{number}`, `organizations/{number}` forms.
- The allowlist only takes effect with
  `producerInstanceLocation: CUSTOM_RESOURCE_HIERARCHY_LEVELS` — the spec
  enforces that the two move together.

## Composes with

`GcpVpcNetwork` and `GcpSubnetwork` upstream; any PSC-first managed
service instance downstream.
