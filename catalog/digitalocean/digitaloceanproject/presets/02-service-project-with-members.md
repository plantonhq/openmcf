# Service Project with Managed Membership

This preset creates a per-service project that OWNS its membership: a droplet joins by reference to its `urn` output, and an existing bucket joins by literal URN. Removing a member from the list moves it back to the account's default project -- nothing is ever destroyed by membership changes.

## When to Use

- Grouping one service's resources so the console and billing views match the architecture
- Charts where the project and its members deploy together and membership should be code
- Mixing chart-managed members (by reference) with pre-existing ones (by literal URN)

## Key Configuration Choices

- **Membership by reference** (`valueFrom` with an explicit `kind`) -- the list is polymorphic across kinds, so each reference names its own kind and reads the producer's `urn` output.
- **A resource belongs to exactly one project** -- listing it here moves it from wherever it was; never list the same resource in two projects.

## What You Get

A staging-labeled project containing the droplet and the bucket, with `project_id` exported for resources that carry a project field.
