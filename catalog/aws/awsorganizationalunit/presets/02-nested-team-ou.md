# Nested Team OU

This preset creates a second-level OU under a parent OU — the
per-team subdivision inside a Workloads-style container.

## When to Use

- Estates that subdivide first-level OUs by team or product line
- When a team needs its own guardrail scope (an SCP attached here
  governs only this team's accounts)

## What You Get

- One OU nested under the referenced parent OU (the reference resolves
  the parent's `ou_id` — no hand-copied IDs)

## Customize

- The parent is IMMUTABLE — re-parenting replaces the OU, so settle
  the tree shape before populating it (accounts move freely; OUs do
  not)
- AWS allows five OU levels under the root — prefer shallow trees;
  inheritance reasoning gets hard long before the quota does
- A literal `parentId: {value: "ou-..."}` carries pre-existing trees
  built outside Planton
