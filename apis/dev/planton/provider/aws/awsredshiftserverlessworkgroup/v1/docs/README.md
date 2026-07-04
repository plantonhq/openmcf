# AWS Redshift Serverless Workgroup — Architecture and Design

## Overview

A Redshift Serverless workgroup is the compute plane of AWS's serverless
data warehouse. It owns Redshift Processing Unit (RPU) capacity, VPC
placement, the connection endpoint, and query-level configuration -- and
nothing else. The data it serves lives on a namespace
(`AwsRedshiftServerlessNamespace`), which it attaches to by name at create
time.

RPU-hours accrue only while queries execute: an idle workgroup bills
nothing, which is what makes serverless the entry point most new analytics
adopters start on. The workgroup is where every cost and performance knob
lives.

## The Capacity Contract

The workgroup has exactly two capacity postures, and AWS rejects mixing
them (a constraint this spec enforces at validate time rather than at
deploy):

1. **Fixed baseline** -- `base_capacity` sets the RPU floor each query
   starts from. Empty keeps the AWS default (128 RPU). AWS validates the
   exact accepted RPU increments at deploy; they have changed over time
   (the floor dropped and the ceiling grew), so this spec deliberately
   does not freeze them.
2. **Price-performance target** -- AWS owns the baseline and adjusts it
   against a dial: level 1 leans cheapest, 100 leans fastest, 50 (the
   default) balances. While enabled, `base_capacity` must stay unset.

`max_capacity` caps scaling in either posture -- the spend guardrail. 0
leaves scaling uncapped (the AWS default); the provider treats a removed
cap specially (it sends -1 to the API), which both modules inherit by
simply not forwarding an unset value.

AWS serializes capacity changes -- the underlying API accepts one
parameter change per update call -- so the provider decomposes an edit
that touches several capacity fields into ordered single-field updates.
This is provider machinery, not module behavior, but it explains why a
multi-knob capacity edit takes several minutes to settle.

## Networking

Redshift Serverless requires the workgroup's subnets to span **three**
distinct availability zones (in regions with three or more), and the
free-IP requirement per subnet scales with base capacity. The spec states
the three-subnet minimum at validate time; AZ distinctness is checked by
AWS (subnet-to-AZ mapping is not knowable at validation).

The port is constrained to the only ranges the serverless API accepts:
5431-5455 and 8191-8215 (default 5439). This is a hard serverless
constraint, unlike the provisioned cluster where the narrow ranges apply
only to relocation-enabled clusters.

Ingress rules belong on referenced `AwsSecurityGroup` nodes. The
workgroup attaches groups; it never creates one.

## Query-Level Configuration

Serverless has no parameter groups -- `config_parameters` apply directly
to the workgroup. The accepted parameter set is small and closed (the
API rejects anything else), so the spec mirrors it exactly: session
defaults (`datestyle`, `search_path`, `query_group`), security posture
(`require_ssl`, `use_fips_ssl`), behavior toggles (`auto_mv`,
`enable_case_sensitive_identifier`, `enable_user_activity_logging`), and
the query-monitoring limits (`max_query_execution_time`,
`max_scan_row_count`, and the other `max_*` guards) that implement
query monitoring rules.

## Design Decisions

- **Name basis** -- the workgroup name is `metadata.name` (create-time
  immutable in AWS), the shared engine-parity basis.
- **Namespace by output-backed reference** -- the `namespace_name` ref
  points at the namespace's `status.outputs.namespace_name`; references
  resolve against stack outputs, never metadata.
- **`{name, value}` parameter shape** -- config parameters use the same
  message shape as the provisioned Redshift cluster's parameters, not the
  provider's `parameter_key`/`parameter_value` verbosity, so the Redshift
  family speaks one parameter vocabulary.
- **Endpoint as flat outputs** -- the provider nests the endpoint
  (address, port, VPC endpoint, network interfaces); the spec exports the
  two values downstream consumers actually wire (`endpoint_address`,
  `port`) as singular semantic outputs.

## Deliberately Skipped Provider Surface

| Provider surface | Verdict | Reason |
| --- | --- | --- |
| `aws_redshiftserverless_endpoint_access` | Defer | Additional managed VPC endpoints into a workgroup -- a separate access-plane surface (the provisioned-Redshift endpoint-access precedent); joins via the exported workgroup name on demand |
| `aws_redshiftserverless_custom_domain_association` | Defer | Custom domain + ACM certificate binding; needs a domain the deployment owns -- joins via the exported workgroup name on demand |
| `aws_redshiftserverless_usage_limit` | Defer | Cost-governance surface keyed by the workgroup ARN; `max_capacity` covers the 90/10 spend-guardrail need |

## References

- RPU capacity and billing: https://docs.aws.amazon.com/redshift/latest/mgmt/serverless-capacity.html
- Serverless networking requirements and known issues: https://docs.aws.amazon.com/redshift/latest/mgmt/serverless-known-issues.html
- Workgroup config parameters: https://docs.aws.amazon.com/redshift-serverless/latest/APIReference/API_CreateWorkgroup.html
