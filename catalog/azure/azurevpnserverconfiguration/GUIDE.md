# Azure VPN Server Configuration -- Operational Guide

Judgment that does not fit in field references.

## Pick the authentication family by who manages the endpoints

Entra ID ("AAD") is the right default for a managed workforce: sign-in
rides conditional access and MFA, and revoking a person is an identity
operation, not a certificate hunt. Certificates fit unmanaged or
offline endpoints (labs, contractors, appliances) -- but YOU own the
root, the distribution, and the revocation list. RADIUS exists to
reuse auth you already operate; it adds a network dependency (the
gateway must reach your RADIUS server through the hub) that the other
two families do not have.

## The configuration is shared state -- edit it like shared state

Many gateways can point at one configuration, and edits apply in place
to ALL of them: swap a root certificate here and every attached
gateway starts rejecting the old chain as clients reconnect. Fleet-wide
policy is exactly why the object exists -- but a "quick fix" for one
office is a new configuration and a gateway repoint (which REPLACES
that gateway), not an edit to the shared one.

## Revoke certificates by thumbprint before rotating roots

A lost laptop needs `clientRevokedCertificates` (one thumbprint entry,
in-place update, effective on reconnect) -- not a root swap, which
disconnects everyone anchored to it. Rotate the root on schedule by
ADDING the new root alongside the old, re-issuing clients, then
removing the old -- both lists update in place.

## Policy groups only matter when the gateway maps them

Declaring policy groups changes nothing by itself: they are matching
rules, and the mapping to address pools lives on the point-to-site
gateway's connection configurations. Policy groups (and multiple
gateway address pools) also require OpenVPN -- set `vpnProtocols:
["OpenVPN"]` on the configuration before designing a segmented rollout,
and note IKEv2-only native clients cannot join OpenVPN-only setups.

## The RADIUS secret never comes back

ARM does not return `radius.servers[].secret` on reads -- Azure
compares your value on writes but never echoes it. Keep the secret in
a secret manager and reference it; treat "what secret is deployed?"
as unanswerable from Azure's side (drift detection tolerates it by
design).

## AAD values are tenant-plumbing, not preferences

The Microsoft-managed Azure VPN Client's audience is a fixed
application ID (`41b23e61-6c1e-4545-b367-cd054e0ed4b4`), and the
tenant/issuer URLs must embed YOUR directory ID exactly (the issuer
keeps its trailing slash). A typo here deploys fine and then fails
every sign-in -- validate against a test user before rolling out.
