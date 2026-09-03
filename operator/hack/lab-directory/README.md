# The Lab Directory

A containerized corporate directory the platform rehearses enterprise identity against: a Samba Active Directory domain (`lab.example.internal`) with deterministically seeded users and groups, plus an "entra-sim" realm shape for the brokered arm. Every directory journey — connect, verify, map, first sign-in, drift, offboarding — is proven here before any real company's directory is ever touched.

## Why it exists

No other test infrastructure anywhere in the platform can answer "does sign-in against a corporate directory actually work?". The identity server's LDAP federation has an Active Directory vendor mode with AD-specific behavior (`sAMAccountName`, `member`/`memberOf` DNs, `objectGUID` as the immutable id, `userAccountControl` disabled bits, nested-group resolution) — so the lab is a real Samba AD directory, not a generic LDAP server with an AD-ish schema. Rehearsing against generic LDAP would validate the wrong code path.

## The one seed

`seed.yaml` is the single definition of the lab population. Every consumer seeds from it:

- **The operator's federation convergence suite** (Docker-gated, `make test-realm-convergence` in the operator) boots the container and seeds it via `samba-tool` — its loader lives with that suite.
- **Sync-engine and mapping suites** consume the same manifest with their own loaders.
- **The e2e lab cluster** deploys the same container and seed beside the full self-hosted stack.

Adding an edge case is a data change in `seed.yaml`, never new plumbing. A journey proven on one federation arm is comparable on the other because both rehearse the SAME population.

## The container

`smblds/smblds` — Samba AD Lightweight Directory Services, purpose-built for CI: a real Samba AD directory serving LDAP (389) and LDAPS (636) without full domain-controller overhead (no Kerberos KDC, no DNS), which is exactly the surface LDAP federation touches. It provisions itself from environment (`REALM`, `DOMAIN`, `ADMINPASS`) and mints its own CA at `/var/lib/samba/private/tls/ca.pem` — that CA **is** the lab's private-CA fixture: consumers extract it and wire it through the platform's CA-bundle path, so the classic enterprise blocker (private-CA LDAPS) is rehearsed on every run.

The admin bind DN is `CN=Administrator,CN=Users,DC=lab,DC=example,DC=internal` with the seed's `adminPassword`.

## The fixtures, each present deliberately

| Fixture | Edge it proves |
|---|---|
| `platform-eng` (14 members) | The primary mapping group for happy-path journeys |
| `all-engineering` (220 members, 206 generated) | Blast-radius preview realism |
| `platform-eng` ⊂ `engineering` ⊂ `everyone` | Transitive (nested) membership resolution |
| `sam.nomail` | The email-less posture — email is load-bearing platform-wide |
| `alex.kim` / `akim2` (both "Alex Kim") | Display-name collision rendering |
| `developers` | Group name colliding with a pre-existing Planton team |
| `dara.disabled` | Disabled account still holding memberships (`userAccountControl`) |
| `lee.lonely` | Signs in but lands with nothing (and the seat-cap interaction) |
| `it.admin` | Email collision with the seeded platform admin |
| `svc.jenkins` | A bot account without a person (or an email) behind it |

Generated users (`eng-001`…) carry no email deliberately — large-directory realism.

## The brokered arm (entra-sim)

The upstream identity provider for OIDC-brokering rehearsals is a second identity-server realm mimicking Entra ID's observable shapes — most importantly, group memberships emitted as object-id GUIDs, never names (the reason the platform maps groups by stable id). Coding-phase suites import that realm onto their existing test identity server; the e2e lab runs it as its own container, closer to production shape. Entra's proprietary quirks beyond OIDC stay with the one real-tenant dress rehearsal.

## What the lab is NOT

Not a load-test rig (scale questions are answered by the seat model's design — accounts materialize only at first sign-in), not a multi-forest AD, and not a SCIM endpoint (that slice brings its own fixtures when it comes).
