# AwsClientVpn — Design Notes

## Provider mapping

One spec folds four provider resources, all endpoint-scoped:

| Spec surface | Provider resource |
|---|---|
| The endpoint fields | `aws_ec2_client_vpn_endpoint` |
| `subnet_ids` | `aws_ec2_client_vpn_network_association` (one per subnet) |
| `authorization_rules` | `aws_ec2_client_vpn_authorization_rule` (one per rule) |
| `routes` | `aws_ec2_client_vpn_route` (one per route) |

The satellites fold because none has identity outside its endpoint and
nothing FK-references them; each remains its own provider resource keyed by
its AWS identity (subnet; CIDR+grantee; destination+subnet), so membership
edits apply in place.

## Design decisions

- **Typed authentication options (max 2), all create-time immutable.** The
  provider models authentication as a set of up to two blocks; the spec
  mirrors that with per-type CEL requiredness (certificate → the client CA
  chain ref; directory → the directory ID; federated → the SAML provider
  ARN). The client CA chain is deliberately its own reference — reusing
  the server certificate silently would conflate two different trust
  roles, even though a self-signed setup may point both at one imported
  certificate.
- **Connection logging is presence-is-the-switch.** The provider requires
  the `connection_log_options` block with an `enabled` boolean; the spec
  models only the block, and absence maps to `enabled = false` — there is
  no separate boolean for a manifest to contradict.
- **`client_cidr_block` is conditional.** Required in general, but AWS
  derives client addressing for pure-IPv6 tunnel traffic
  (`traffic_ip_address_type: ipv6`), where the field must be empty — the
  provider made the argument optional for exactly this case; the spec
  enforces the coupling as CEL.
- **The transit-gateway arm excludes the VPC surface.** The provider marks
  `transit_gateway_configuration` ConflictsWith `vpc_id` and
  `security_group_ids`; folded subnet associations belong to the VPC
  surface too, so all three exclusions are CEL rules (presence checks
  only — reference sub-fields are never dereferenced in CEL).
- **Routes depend on their subnet's association in BOTH engines.** AWS
  rejects a route whose target subnet is still associating
  (`InvalidClientVpnActiveAssociationNotFound`); both modules make the
  edge explicit instead of leaning on provider-side retries.
- **Port values are validated, port↔protocol pairing is not.** AWS accepts
  443 and 1194 with either transport protocol; only the value set is a
  real constraint.
- **Directory Service and IAM SAML providers are literal IDs/ARNs.**
  Neither has a Planton kind; the fields document that and take literals.

## Deliberately skipped / deferred (with reasons)

- **Directory Service and IAM SAML provider kinds** — separate identity
  surfaces; the auth arms compose by literal ID/ARN with zero rework if
  kinds ever exist.
- **A `status` output** — the provider exposes no endpoint status
  attribute (status codes exist only in its internal waiters).

## Update semantics

ForceNew (replaces the endpoint): `authentication_options`,
`client_cidr_block`, `transport_protocol`, `endpoint_ip_address_type`,
`traffic_ip_address_type`, `transit_gateway_configuration`. Everything
else updates in place. Individual satellites replace themselves (never the
endpoint) when their identity fields change.

## Operational notes

- Network associations are the slow path: AWS takes several minutes per
  attach/detach (the provider allows up to 30 minutes each, with a
  4-minute initial wait).
- Associations bill hourly whether or not clients are connected; a
  zero-association endpoint is nearly free and can be pre-staged.
- The exported client configuration
  (`aws ec2 export-client-vpn-client-configuration`) already prepends the
  required random subdomain label to `endpoint_dns_name`.
