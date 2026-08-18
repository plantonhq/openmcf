# AWS Cloud Map Namespace

Service discovery for everything: a registry where services publish themselves (or are registered statically) and consumers find them — by private DNS inside a VPC, by public DNS, or by API call with zero DNS at all. The registry ECS service discovery writes into.

## What Gets Managed

- The namespace: HTTP (API-only), PRIVATE_DNS (a private hosted zone visible in one VPC), or PUBLIC_DNS (an internet-resolvable zone) — the name is the domain (e.g. `corp.internal`).
- Services keyed by name: the DNS records instances publish (A/AAAA/SRV/CNAME with TTL, multivalue or weighted routing), Route 53 health checks (public namespaces) or custom heartbeat health.
- Statically registered instances under each service: an IP and port, a CNAME to any endpoint, a Route 53 alias to a load balancer, or an EC2 instance — plus custom attributes for API-side discovery.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Cloud Map (servicediscovery) permissions.

### AWS Prerequisites

- For PRIVATE_DNS: the VPC the namespace's zone should be visible in (DNS resolution enabled).
- Nothing else — namespaces, services, and registrations are self-contained.

## After You Deploy

- DNS namespaces resolve immediately inside their scope: service `api` in namespace `corp.internal` answers at `api.corp.internal`.
- HTTP namespaces answer DiscoverInstances calls via the `http_name` output.
- ECS services reference the `service_arns` output to auto-register their tasks.

## Common Changes

- Add services or static instances (in-place list edits); re-registering the same instance id updates it (AWS upserts).
- Namespace name and type are fixed for life; the HTTP namespace's description change REPLACES it (the DNS namespaces update in place).
- Never set a service's force_destroy where ECS registers tasks — it deregisters everything, including registrations this manifest never made.
