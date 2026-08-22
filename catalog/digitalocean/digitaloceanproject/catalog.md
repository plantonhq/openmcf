# Project on DigitalOcean

Creates a DigitalOcean project -- the account-level container that organizes droplets, load balancers, domains, buckets, and most other resources into named groups with a purpose and an environment. Membership is declared on the project itself, wiring member resources by reference to their exported URNs. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for membership wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Project** -- the named container with description, purpose, and environment
- **Membership** -- configured only when `resources` is set; each listed resource is moved into the project (and back to the account's default project when removed)

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### DigitalOcean Account

- Nothing: projects are free organizational containers.

## After You Deploy

Reference the `project_id` output from resources that carry a project field. Remember that destroying the project never destroys what is inside it -- members are relocated to the account's default project.
