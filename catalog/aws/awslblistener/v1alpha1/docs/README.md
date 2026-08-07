# AWS LB Listener: The Entry Point of the Routing Graph

## What a Listener Is

An ELBv2 listener is the process that accepts connections on one
port/protocol pair of a load balancer and decides what happens to them. The
load balancer itself is just placement and capacity -- subnets, IPs, security
groups; every routing decision starts at a listener. A typical ALB carries at
least two: port 80, whose only job is redirecting to HTTPS, and port 443,
which holds the TLS material and the real routing.

`AwsLbListener` models the listener as its own component because it is the
anchor two different things hang off. Certificates attach to a listener (not
to the load balancer), and listener rules attach to a listener as services
deploy. Folding listeners into the load balancer definition would force every
certificate rotation and every new service route through an edit of the most
shared resource in the network stack.

## One Kind for Both Families

The same component serves Application and Network Load Balancers, matching
AWS's own model -- a listener is a listener; the protocol decides the family
and the allowed feature set:

- **ALB listeners** (`HTTP`, `HTTPS`) take the full default-action set --
  forward, redirect, fixed-response, authenticate-cognito, authenticate-oidc,
  jwt-validation -- plus mutual TLS and HTTP header controls.
- **NLB listeners** (`TCP`, `UDP`, `TCP_UDP`, `TLS`) accept exactly one
  `forward` action. AWS rejects every other action type at Layer 4: there is
  no HTTP to redirect, no response to fabricate, no session to authenticate.

The `loadBalancerArn` field defaults to referencing an `AwsAlb`; attaching to
an `AwsNlb` (or an externally managed load balancer) takes an explicit
`valueFrom` or a literal ARN. The default follows frequency: rule-based ALB
routing is the reason listeners are managed standalone most often.

## Why SNI Certificates Fold Into the Listener

AWS models each additional certificate as a separate listener-certificate
attachment resource. This component folds them into
`additionalCertificateArns` instead, because an attachment is pure glue: no
configuration of its own, no referenceable identity, no lifecycle beyond
"this certificate is on this listener". The list is the honest shape --
SNI selection is by hostname match at handshake time, `certificateArn` serves
clients that match nothing, and each entry is typically an
`AwsCertManagerCert` reference.

## Action Chains

`defaultActions` is a list because authentication actions are *middleware*:
an `authenticate-oidc` or `authenticate-cognito` step runs first, then the
terminal action (exactly one of `forward`, `redirect`, or `fixed-response`)
handles the authenticated request. List position sets the order; the
optional `order` field exists for explicit control but is rarely needed.

Each action is a discriminated union -- `type` names the action, and exactly
one matching configuration message must be set, the same shape the AWS API
uses. Validation enforces the pairing at manifest level, so a `type:
forward` with a `redirect` block fails before any API call.

Three authentication styles cover different clients:

- **authenticate-cognito** / **authenticate-oidc**: browser-facing. The ALB
  redirects unauthenticated users to a hosted login, then manages a session
  cookie. Cognito wires by `AwsCognitoUserPool` reference; OIDC takes the
  issuer and endpoint set of any provider (Okta, Auth0, Entra ID, Google).
- **jwt-validation**: API-facing. Each request must carry a valid JWT bearer
  token verified against a JWKS endpoint; no redirects, no cookies, rejects
  happen at the load balancer before targets spend any capacity.

A practical default for a rules-based listener: make the default action a
`fixed-response` 404, and let every real route live in `AwsLbListenerRule`
resources -- unmatched hostnames get an explicit answer instead of hitting an
arbitrary target group.

## Mutability

The load balancer is the listener's only create-only field: moving a listener
to a different load balancer replaces it (and rules re-attach to the
successor). Everything else -- port, protocol, certificates, SSL policy,
actions, header configuration -- updates in place. In practice this means
certificate rotation, TLS policy upgrades, and default-action changes are
zero-downtime edits, while consolidating listeners onto a new load balancer
is a re-create operation to plan deliberately.

## Scoping Notes

- **Gateway Load Balancer listeners are not modeled**: there is no
  gateway-load-balancer kind to compose them with, so the GENEVE surface
  would be dead weight.
- **ELBv2 trust stores are not a kind**: mutual TLS `verify` mode takes a
  literal trust store ARN (or an explicit reference). Trust stores change
  rarely and carry only a CA bundle; a dedicated kind can be added if that
  changes.
- **`tcpIdleTimeoutSeconds` is TCP-only**: raise it (60-6000, default 350)
  for protocols with long quiet periods -- database wire protocols, MQTT --
  to stop the NLB reaping live-but-idle connections.

## Dual-Engine Implementation

`AwsLbListener` ships both a Terraform/OpenTofu module and a Pulumi (Go)
module at behavioral parity. Both enforce the certificate requirement for
HTTPS/TLS protocols at deploy time (message-level CEL cannot inspect
reference-typed fields, so the modules carry the check), attach SNI
certificates identically, and export the same `listener_arn` output.
Whichever engine a team standardizes on, the listener behaves identically.
