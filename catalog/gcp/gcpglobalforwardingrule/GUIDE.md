# GcpGlobalForwardingRule Guide

The judgment this guide protects: the forwarding rule is the entry point
of the whole load-balancer chain — the piece that binds a reserved IP, a
port, and a target proxy into "traffic flows". It is also the piece whose
recreation is instantly user-visible.

## Bring your own IP, always

Let the rule reference a GcpGlobalAddress instead of auto-allocating.
An auto-allocated VIP dies with the rule; a referenced reservation
survives it, which turns "recreate the forwarding rule" from a DNS
incident into a blip. The rule's own recreation only breaks the binding
— with a stable IP the new rule picks up exactly where the old one
stopped.

## One rule per port contract

A production HTTPS frontend is typically TWO rules on one IP: 443 to the
HTTPS proxy, 80 to the HTTP proxy whose URL map redirects. The spec's
`portRange` accepts a single port or range, but resist wide ranges on
EXTERNAL schemes — every open port is attack surface, and GCP's
port-range semantics differ between classic and envoy-based schemes.

## Scheme is destiny

`loadBalancingScheme` decides which targets are legal, whether a network
is required, and which extras (metadata filters, Service Directory) even
apply. The spec's `NONE` sentinel maps to the API's empty scheme for PSC
frontends — the module translates it, so never write the empty string
yourself. The Service Directory registration is PSC-only; the spec
enforces the pairing pre-deploy.

## Source-filtered rules are regional

The provider carries `source_ip_ranges` on this resource, but documents
it as usable only on REGIONAL EXTERNAL rules — on the global rule it is
a dead lever, which is why this kind deliberately does not model it
(recorded in the parity manifest). Needing source filtering on a global
frontend means Cloud Armor, not the forwarding rule.

## Teardown discipline

Deleting the rule stops traffic to that IP immediately — the proxies and
backends behind it stay healthy but unreachable on that VIP. `PREVENT`
suits any rule fronting production; `ABANDON` keeps traffic flowing while
dropping management (the escape hatch for handing a frontend to another
stack). The reserved address it referenced is governed by its OWN kind's
deletion policy — destroying the rule never releases the IP.
