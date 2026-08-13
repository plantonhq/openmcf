# Shared Backend Link

This preset is the production shape: one VPC link per NLB, shared by
every REST API in the environment that needs that private backend.

## When to Use

- Several REST APIs fronting the same internal service
- Environments that should not create a link per API

## What You Get

- A named, described v1 VPC link on the shared NLB
- One `vpc_link_id` every API's private integrations reference

## Customize

- Name the link after the backend (`orders-nlb-link`), not after an
  API
- Changing the NLB replaces the link — stand up the new one, repoint
  integrations, then destroy the old
