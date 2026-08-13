# GcpWorkloadIdentityPoolProvider Guide

The judgment this guide protects: the provider is the trust boundary of
keyless auth. Its two text fields — the attribute mapping and the
attribute condition — ARE the security model; everything else is
plumbing around them.

## Always set attributeCondition for multi-tenant issuers

An OIDC provider for GitHub Actions trusts an issuer that vouches for
every repository on the planet. Without a condition, ANY GitHub workflow
can exchange tokens against your pool. `attributeCondition` is what
scopes trust:

```yaml
attributeCondition: assertion.repository_owner == "my-org"
```

This is not optional hardening — for shared issuers it is the difference
between "my org can federate" and "the internet can federate."

## The mapping decides what IAM can see

`attributeMapping` translates issuer claims into the attributes IAM
bindings target (`principalSet://…/attribute.<name>/<value>`). Map only
what bindings will actually use: every mapped attribute is an interface
you must keep stable across issuer changes. OIDC providers MUST map
`google.subject` (validation enforces it — GCP rejects the provider
without it); AWS, SAML, and X.509 fall back to issuer-specific defaults.

## Leave allowedAudiences empty unless the issuer forces you

An empty audience list means tokens must be minted for the provider's
own canonical resource name — the safest contract, because a token
minted for anything else is rejected. Add explicit audiences only for
issuers that cannot set a custom audience.

## Disable, don't delete

`disabled: true` is the kill switch: new token exchanges stop
immediately (already-issued Google credentials ride out their ~1h
expiry), and flipping back restores service. Deletion starts a ~30-day
soft-delete clock during which the provider ID cannot be reused — and a
create against a soft-deleted ID fails outright. For rotation or
incident response, disable; delete only what will never come back under
that ID.

## The issuer type is fixed at design time

The `aws` / `oidc` / `saml` / `x509` choice cannot change on a live
provider — the API rejects cross-type updates. Switching issuers means a
new provider (new ID, new audience, updated IAM bindings). The arm's
CONTENTS update in place: rotating SAML metadata, adding audiences, or
updating an X.509 trust store is an ordinary edit.

## Teardown discipline

`deletionPolicy: PREVENT` for any provider live pipelines federate
through — deleting it breaks every workload authenticating via this
issuer AND locks the ID for ~30 days. `ABANDON` leaves the provider
exchanging tokens unmanaged, which for a trust boundary is the worst of
both worlds; prefer `disabled: true` when the goal is "stop trusting
this issuer."
