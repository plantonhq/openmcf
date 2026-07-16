# Public Index Endpoint

The simplest serving surface: a public Vector Search index endpoint in
the ambient project.

## What this preset creates

An index endpoint named `Vector Serving` in `us-central1` with the public
querying arm enabled. Once a `GcpVertexAiDeployedIndex` places an index
onto it, nearest-neighbor queries go to the endpoint's
`public_endpoint_domain_name` output over the public internet
(authenticated by IAM).

## When to use

- Standard vector-search serving without private-networking requirements
- Development and staging environments
- Applications querying from outside GCP (or across many VPCs) where
  peering everything is impractical

## Remix ideas

- Swap the public arm for `network` (a GcpVpcNetwork reference) to keep
  queries inside a peered VPC — see the vpc-peered preset.
- Swap it for `privateServiceConnectConfig` for consumer-project
  isolation without peering — see the psc preset.
- Add cost-attribution labels; they merge beneath the platform's
  attribution labels on both engines.
