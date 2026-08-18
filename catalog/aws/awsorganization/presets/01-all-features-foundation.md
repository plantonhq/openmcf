# All-Features Foundation

This preset creates the standard multi-account foundation: an
all-features organization with trusted access for the governance
services and the two most common policy types enabled.

## When to Use

- The first deploy of any multi-account AWS estate — everything else
  in the Organizations family (OUs, member accounts, policies) hangs
  off this
- Estates planning SCP guardrails or org-wide CloudTrail/Config from
  day one

## What You Get

- An ALL-features organization (the level every advanced capability
  requires)
- Trusted access for CloudTrail, Config, and Account Management (the
  last one is what lets member-account contacts and regions be managed
  centrally)
- SCP and tag-policy types enabled on the root — attachments work
  immediately

## Customize

- Add `delegatedAdministrators` entries once member accounts exist
  (e.g. delegate GuardDuty to a security account)
- Add a `resourcePolicy` document to delegate organization reads to a
  tooling account
- Trim `awsServiceAccessPrincipals` to exactly the services you run —
  trusted access lets a service create roles in member accounts
