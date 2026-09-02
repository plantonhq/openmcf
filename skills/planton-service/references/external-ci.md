---
title: Keyless CI — Workload Identity Bindings, CI-Step Registration, and the Deploy Step
description: External CI (GitHub Actions, gitlab.com) calls Planton with zero stored secrets — the org registers one workload identity binding, the CI job exchanges its provider's OIDC token via planton iam federate, and the credential can register and deploy exactly the service its token proves. Read when someone asks how CI authenticates, how to deploy from a CI step, how to use the Planton GitHub Action, or why a federation exchange was refused (every refusal is deliberately identical — walk the checklist, not the message).
---

# Keyless CI — Workload Identity Bindings, CI-Step Registration, and the Deploy Step

External CI (GitHub Actions, gitlab.com CI) can call Planton without any stored secret: the CI job presents its provider's own short-lived OIDC token, and Planton exchanges it for a short-lived scoped credential that can register and deploy the one service the token proves. Read this when someone asks how CI authenticates to Planton, how to register or deploy a service from a CI step, how to use the Planton GitHub Action, why a federation exchange was refused, or how to stop trusting a repository.

## The trust binding

A **workload identity binding** is the org-registered rule that makes the exchange possible: "tokens from THIS provider, matching THESE conditions, may act as THIS service account." It is declarative YAML (`WorkloadIdentityBinding`, applied like any resource), and reads/writes require **manage-access on the organization** — trusting an external identity is an access-management act, so ordinary org write access is deliberately not enough. Hold these facts when explaining one:

- The binding names an existing **service account** — the identity CI acts as, the name audit trails show. Create the service account first; a binding naming a missing one is refused at apply.
- Conditions are **structured**, matched against the provider's own token claims: the repository (GitHub `owner/name`, GitLab `group/project`, always required, matched exactly — one binding trusts one repository), and optionally a ref or environment. There are no patterns or wildcards to author.
- The **audience is required** and should be a value your CI explicitly requests the token with (GitHub's `audience` parameter). It is what makes Planton the token's only valid consumer.
- Each provider's **issuer is pinned in Planton's code** — a binding never carries an issuer URL. GitHub Actions and gitlab.com are supported; self-hosted GitLab instances are not yet (their issuers are customer URLs, which need a hardened key-fetch design — say so plainly rather than improvising).
- A binding **never widens authority**: exchanged credentials carry the external-CI class's fixed rule set (registering the service, deploying it, and reading the runs and receipts those verbs produce — never approving anything), and binding a powerful service account grants CI nothing extra — work-credential authorization never consults the account's standing grants.
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

## The deploy step

A CI job that built its own image deploys it with the same credential:

```bash
export PLANTON_API_KEY=$(planton iam federate --org acme --token "$OIDC_TOKEN")
planton service deploy checkout-api --org acme --env prod \
  --image ghcr.io/acme/checkout-api@sha256:... \
  --commit "$GITHUB_SHA" --branch "$GITHUB_REF_NAME" --follow
```

The deploy renders manifests from the service's CURRENT inline deploy declaration with the image injected, and rides the same engine as every deploy: protection gates, separation of duties (the deployer can never approve), URLs, and rollout verification. Facts to hold:

- **`--follow` tells the truth in CI**: the command exits non-zero when the run fails AND when rollout verification reports `failed` (the workload demonstrably never came online) — the failed checks print with the provider's own words. An honestly `unverifiable` rollout passes with a note. This is why a green deploy step means something.
- **Kustomize-source services refuse** the deploy verb (their rendering runs on the build lane) — the working paths are the git-push lane or promoting an existing deployment.
- **The environment must be declared** in the service's deploy configuration; the refusal names the declared set.
- **A protected environment parks the run at its approval gate** — a followed CI job waits for the approval; set the job's own timeout accordingly.

## The GitHub Action

`plantonhq/planton/actions/deploy` wraps the whole story — install the CLI (checksum-verified), mint the job's OIDC token, federate, optionally register, deploy, and wait honestly. The job needs `permissions: id-token: write`, and the `audience` input must be exactly the binding's audience. `register: true` applies the service manifest first (its repository proven by the same token). See the action's own README for the full input table. The same action also runs fully backendless: with `org` and `audience` absent it deploys the repository's own kustomize declaration through the open-source engine — offline mode, covered by the action's README.

## Proven repository identity

A service registered with a federated credential may only declare the repository the credential's own verified token names — a mismatch is refused, naming both repositories. This is the point of CI-step registration: a catalog entry whose repository was typed can drift or lie; one registered from the repository's own CI is **proven**. A service that declares no repository (a pure catalog entry) registers fine from any lane.

The same law narrows the delivery verbs, and MORE strictly: a federated credential may deploy, promote, or roll back **only the service whose declared repository matches its token's** — and here a repo-less service refuses (with no declared repository there is no proof relation at all, so the answer is to set the service's repository or deliver from a lane with standing access). The deploy rules themselves are org-wide by grammar; this per-service narrowing is what keeps one repository's workflow from deploying every service in the organization.
