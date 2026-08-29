# Keyless CI — Workload Identity Bindings and CI-Step Registration

External CI (GitHub Actions, gitlab.com CI) can call Planton without any stored secret: the CI job presents its provider's own short-lived OIDC token, and Planton exchanges it for a short-lived scoped credential. Read this when someone asks how CI authenticates to Planton, how to register a service from a CI step, why a federation exchange was refused, or how to stop trusting a repository.

## The trust binding

A **workload identity binding** is the org-registered rule that makes the exchange possible: "tokens from THIS provider, matching THESE conditions, may act as THIS service account." It is declarative YAML (`WorkloadIdentityBinding`, applied like any resource), and reads/writes require **manage-access on the organization** — trusting an external identity is an access-management act, so ordinary org write access is deliberately not enough. Hold these facts when explaining one:

- The binding names an existing **service account** — the identity CI acts as, the name audit trails show. Create the service account first; a binding naming a missing one is refused at apply.
- Conditions are **structured**, matched against the provider's own token claims: the repository (GitHub `owner/name`, GitLab `group/project`, always required, matched exactly — one binding trusts one repository), and optionally a ref or environment. There are no patterns or wildcards to author.
- The **audience is required** and should be a value your CI explicitly requests the token with (GitHub's `audience` parameter). It is what makes Planton the token's only valid consumer.
- Each provider's **issuer is pinned in Planton's code** — a binding never carries an issuer URL. GitHub Actions and gitlab.com are supported; self-hosted GitLab instances are not yet (their issuers are customer URLs, which need a hardened key-fetch design — say so plainly rather than improvising).
- A binding **never widens authority**: exchanged credentials carry the external-CI class's fixed rule set (registering services, today), and binding a powerful service account grants CI nothing extra — work-credential authorization never consults the account's standing grants.
- `disabled: true` suspends the trust without deleting the record; deleting the binding is the off-switch (already-exchanged credentials still expire on their own short clock, ~15 minutes).

## The exchange, in a CI step

```bash
export PLANTON_API_KEY=$(planton iam federate --org acme --token "$OIDC_TOKEN")
planton service register -f service.yaml
```

`planton iam federate` prints the raw credential to stdout and nothing else, so capture works exactly as above; the token comes from `--token`, stdin (`--token -`), or `$PLANTON_OIDC_TOKEN`. In GitHub Actions the job needs `permissions: id-token: write` and requests its token with the binding's audience. The exchanged credential is a standard bearer — `PLANTON_API_KEY` accepts it.

Two reading rules for refusals:

- **Every refusal is identical by design** ("the presented token was not accepted by any workload identity binding in this organization") — the endpoint never reveals whether the org, a binding, or a repository exists. When someone asks why, walk the checklist instead of the message: does a binding exist in that org for that provider? Is it disabled? Does the token's repository match exactly? Was the token requested with the binding's audience? Is the token fresh (they expire in minutes)? Does the bound service account still exist?
- The CI token is **exchange-only**: it can never be used directly as a Planton credential. Only the exchanged `pwk_` credential calls APIs.

## Proven repository identity

A service registered with a federated credential may only declare the repository the credential's own verified token names — a mismatch is refused, naming both repositories. This is the point of CI-step registration: a catalog entry whose repository was typed can drift or lie; one registered from the repository's own CI is **proven**. A service that declares no repository (a pure catalog entry) registers fine from any lane.
