# AWS Transit Gateway Family: Pure Hub, First-Class VPC Attachments, and Route Tables as Routing Domains

**Date**: July 8, 2026
**Type**: Feature, Breaking Change
**Components**: API Definitions, Provider Framework, Terraform Modules, Pulumi Modules, E2E Harness

## Summary

The Transit Gateway surface grows from one bundled kind to a composable three-kind family. **`AwsTransitGateway`** is rebuilt as the pure hub at the full provider surface — its embedded VPC attachments split out into the new **`AwsTransitGatewayVpcAttachment`**, and the previously unmodeled routing surface arrives as **`AwsTransitGatewayRouteTable`**, the isolated-routing-domain node that segmented enterprise networks (prod/non-prod isolation, inspection hair-pinning, central egress) are built from. Every kind ships generator-owned Terraform contracts, full Pulumi entrypoint anatomy, richly-commented modules in both engines, and first-ever family E2E artifacts.

## What Was Built

### `AwsTransitGateway` — the pure hub (breaking)

- The embedded `vpc_attachments` block (`min_items: 1`) and its `vpc_attachment_ids` map output are REMOVED. Attachments are many-per-gateway, independently lifecycled, and are what route tables associate, propagate, and route against — the first-class-node test is met three ways. A zero-attachment hub (pre-staged before spokes migrate in) is now representable.
- New surface: `encryption_support` (in-transit encryption posture; presence-honest tri-state — unset lets AWS compute the effective value). The family's Terraform floor is `>= 6.26.0`: the argument landed in 6.25.0 with a crash regression that 6.26.0 fixed.
- **Honesty fix**: the four default-enable dials (`dns_support`, `vpn_ecmp_support`, `default_route_table_association`, `default_route_table_propagation`) are now `optional bool` — the old plain-bool shape could not express "disable DNS support" distinguishably from "unset". Both engines map unset to a real omission so the AWS default applies.
- AWS's asymmetric replacement rule on the default-table dials (disable -> enable REPLACES the gateway; enable -> disable updates in place) is documented on the fields, and the CIDR-block constraints are now CEL-enforced (max 5; IPv4 /24 or larger, IPv6 /64 or larger, never 169.254.0.0/16).
- **Cross-engine parity fix**: the Terraform module tagged resources with `Name` while the Pulumi module did not — two engines produced visibly different AWS consoles. Both engines now emit the identical `Name` + `planton.ai/*` identity tag set.

### `AwsTransitGatewayVpcAttachment` — the spoke (new kind)

- Gateway + VPC references (create-time immutable, stated honestly), subnet references (in-place updatable — AZs can be added or removed without replacement), DNS/IPv6/appliance-mode dials, and `security_group_referencing_support`.
- The default-route-table membership pair and SG referencing stay **tri-state**: unset inherits the GATEWAY's posture exactly as the provider's Optional+Computed attributes model it; an explicit `false` is the segmented-topology posture where a custom route table owns the association.
- Outputs: `attachment_id` (the routing join key), `attachment_arn`, `vpc_owner_id`.

### `AwsTransitGatewayRouteTable` — the routing domain (new kind)

- One isolated routing domain per resource: folded per-name `associations` (which attachments USE the table — at most one association per attachment gateway-wide, AWS-enforced), `propagations` (which attachments ADVERTISE into it), static `routes` (destination CIDR + attachment XOR `blackhole`, as CEL; destinations unique), and `prefix_list_references` (route a managed prefix list's whole CIDR set; list IDs unique).
- Membership entries are `StringValueOrRef` lists that default to `AwsTransitGatewayVpcAttachment` references AND accept literal attachment IDs — so VPN, Direct Connect, and peering attachments created outside the Planton graph participate in routing domains.
- Every folded member materializes as its own provider resource keyed by a stable identifier (attachment ID, destination CIDR, prefix list ID), so a membership edit touches exactly one underlying resource in both engines.

## E2E

- Three state-aware verifiers on the EC2 SDK (deleting/deleted lifecycle states count as absent; typed NotFound codes mirror the provider's finders).
- Gateway and attachment prerequisite fixtures ship in the SEGMENTED posture (default-table dials off) so consumer lanes fully control route-table membership.
- The route-table scenario proves the isolated-domain story: explicit association + propagation + attachment-forwarded static route + blackhole.
- Live lanes are recorded as deferred in the three profiles (no usable AWS credentials on the ambient chain at session time); the artifacts are live-ready and re-runnable.

## Breaking Changes

- `AwsTransitGateway.spec.vpc_attachments` removed — create `AwsTransitGatewayVpcAttachment` resources referencing `status.outputs.transit_gateway_id` instead.
- `AwsTransitGateway.status.outputs.vpc_attachment_ids` removed — each attachment exports its own `attachment_id`.
- The four default-enable gateway dials changed from `bool` to `optional bool` (manifests that set them explicitly are unaffected).

## Validation

- Full offline gate: buf lint + stub regen, spec/CEL test suites for all three kinds, targeted Go builds + the Bazel repo build, foreign-key and secret-coverage guards, drift-guard + outputs-conformance enrollment, `tofu init`/`validate` for all three modules, 11 manifests CLI-validated, site catalog regenerated.
