---
title: Serving Domains — a Stable Hostname per Service, per Environment
description: An environment declares its serving domain once, and every service deployed there serves at {service}.{env-domain} — computed by the platform, injected into declared carrier resources, surfaced on deployment records, and probed as the domain_serving rollout check. Read when someone wants their service on their own domain, asks why a hostname is or isn't on a deployment, wants a different URL name than the service's, or a domain_serving check needs walking.
---

# Serving Domains — a Stable Hostname per Service, per Environment

Native endpoints (the run.app URL, the load balancer's DNS name) are infrastructure artifacts: generated, target-leaking, replaced when infrastructure is recreated. What a team publishes is a name on their own domain. That name is an ENVIRONMENT's fact, declared once — never per service, never per deploy.

## The declaration, and the two rules that never bend

An environment declares `spec.servingDomain.domain` — either a literal domain (`staging.acmecorp.com`; DNS can be managed anywhere, including outside Planton) or a `valueFrom` reference to a DNS zone resource managed through Planton. Two rules to hold in every answer:

1. **The declared value is the suffix, verbatim.** The platform never inserts the environment's name into a hostname and there is no prefixing convention to configure. Staging serves env-prefixed names because its declared domain SAYS `staging.…`; production serves clean names because its domain says so. When someone asks "how do I make prod URLs not say prod" — the answer is what they type on prod's declaration.
2. **The label defaults to the service's slug.** `checkout-api` in an environment declaring `acme.dev` serves at `checkout-api.acme.dev`, automatically. `spec.deploy.hostname` on the service overrides the label (`api` → `api.acme.dev`) — one label for ALL environments, deliberately. Slugs are unique per organization, so defaults never collide; two services resolving to one hostname in a shared environment refuse at apply, naming both.

What is deliberately NOT a service field: per-environment naming divergence, apex serving (`acmecorp.com` with no label), multiple hostnames, and CDN fronting. Those are edge infrastructure composed from the cloud catalog — real resources the user owns, not switches on the service.

## How the hostname reaches the infrastructure

The platform computes the hostname at deploy collection and INJECTS it into declared carrier manifests — field mutation of resources the author declared, exactly like the built image filling a blank container-image field. It never synthesizes a resource and never writes DNS records itself. The authoring contract: a blank host field is the injection slot. A Cloudflare Worker needs nothing at all (a custom-domain entry is appended; Cloudflare provisions DNS and TLS inside the zone it infers from the hostname). A Kubernetes ingress's blank `rules[].host` and empty `tls[].hosts` are filled — and an ingress whose hosts are all authored is untouched (the platform never appends a rule, because a rule needs a backend only its author knows).

Promote and rollback stay exact and correct: each environment's captured manifests were rendered with that environment's OWN hostname, so what moves between environments never carries another environment's name.

## Reading the record, and walking a domain_serving check

A deployment in a serving-domain environment carries the computed hostname on its URLs under the `serving_domain` tier (both schemes — whichever answers is the one that serves), beside the `native_endpoint` URLs. Rollout verification probes it from the organization's own runner as the `domain_serving` staged check. Dispatch on its status:

- **verified** — the hostname answers. Report it as THE address; the native endpoint is the debugging address.
- **failed** — the hostname did not answer, and the deployment itself still applied. The usual causes, distinguished by the evidence rows' own words and the checks beside them: DNS not created or not propagated yet (a healthy workload + an answering native endpoint + "no such host" on the domain check); or a PROMOTED deployment whose captured manifests predate a domain change (redeploying from a build adopts the current domain). Never relay it as "the deploy failed".
- **unverifiable, naming a missing carrier** — the environment declares a domain, but none of the service's resources can carry a hostname. The wiring is the author's to add: an ingress beside a Kubernetes workload, a domain-mapping resource beside a serverless service, or composed edge infrastructure. The computed hostname still shows on the URLs so the person knows what name awaits the wiring.

An environment with no declaration produces no serving URLs and no domain check at all — absent is the honest answer, not a gap to explain away.
