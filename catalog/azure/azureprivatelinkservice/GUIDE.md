# Azure Private Link Service -- Operational Guide

Judgment that does not fit in field references.

## The subnet flag is the deploy-time trap

Every NAT configuration's subnet must have
`privateLinkServiceNetworkPoliciesEnabled: false` -- a flag on the
SUBNET, not on this service. Nothing validates it offline: the manifest
is clean, the plan is clean, and ARM rejects the create. When adding
Private Link Service to an existing network, flip the flag on the
subnet first (it is an in-place subnet update).

## Aliases are the sharing currency

The `alias` output is a globally unique string consumers use to request
a connection with ZERO access to your subscription. Treat it as the
service's public name: share it freely with intended consumers (it is
not a secret -- visibility and approval are the actual gates), and put
it in your service's onboarding docs.

## Visibility vs approval are two different gates

- **Visibility** decides who can DISCOVER the service and request a
  connection. Empty = your own tenant's RBAC only; UUIDs = those
  subscriptions; `"*"` = anyone holding the alias.
- **Auto-approval** decides whose requests connect WITHOUT you clicking
  approve. Keep it a subset of visibility (ARM requires listed
  subscriptions to be visible), and keep it short -- manual approval is
  the last human checkpoint before someone's VNet reaches your service.

## NAT sizing and immutability

One NAT configuration funds ~64k concurrent flows per consumer
endpoint -- add more only for genuinely high fan-in services. Two
ARM-enforced rules surface only at update time: a static NAT address,
once assigned, can never be cleared back to dynamic, and the PRIMARY
configuration's subnet can never change. Both need destroy/recreate --
so choose the primary's subnet deliberately, and prefer dynamic
addresses unless something downstream pins the address.

## PROXY protocol is a contract with the backend

`proxyProtocolEnabled` prepends a PROXY v2 header to every connection.
A backend that does not parse it sees garbage bytes at the start of
every stream -- connections break in ways that look like application
bugs. Enable it only after the backend is configured for it, and roll
it out backend-first.

## Load balancer prerequisites

The frontend must belong to a STANDARD SKU load balancer (Basic is
rejected), and the frontend set is fixed at creation -- moving the
service to a different LB is a replace. The LB needs at least one rule
on the frontend for traffic to actually flow to your backends;
a frontend with no rules approves connections that then time out.
