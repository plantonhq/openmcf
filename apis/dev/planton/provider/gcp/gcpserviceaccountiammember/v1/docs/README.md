# Service Account IAM: Granting ON the Identity

## An Account Is Also a Resource

Every GCP service account lives a double life. As an *identity*, it is granted roles on projects, buckets, and keys — that side is covered by project- and resource-scoped grant kinds. As a *resource*, the account itself carries an IAM policy that answers a different question: who may USE this identity? That policy is where impersonation lives, and it is what this component writes.

The distinction matters because the highest-value grants in a modern GCP estate are exactly these account-scoped ones:

- **`roles/iam.workloadIdentityUser`** — lets a federated principal impersonate the account. This is the terminal hop of every keyless pattern: a GitHub Actions workflow, a GKE pod, or an external workload presents its own identity token, and this grant is what lets that identity become the service account.
- **`roles/iam.serviceAccountTokenCreator`** — lets the member mint short-lived access and ID tokens as the account. The building block of token-broker and privilege-elevation-on-demand designs.
- **`roles/iam.serviceAccountUser`** — the actAs permission. Cloud Run, GCE, Cloud Functions, and Dataflow all require the deployer to hold it on the runtime account being attached; it is the "may deploy things that run as this identity" right.

## Why Account-Scoped Beats Project-Scoped

All three roles can also be granted at project level, and that is precisely the mistake worth designing away from: a project-level `serviceAccountUser` grant allows acting as EVERY account in the project, and a project-level `workloadIdentityUser` grant lets the federated principal impersonate ANY account. Scoping the grant to one account is not a stylistic preference — it is the difference between "this CI pipeline can deploy as the deployer account" and "this CI pipeline can become anything in the project".

Modeled as its own node, the grant also makes the impersonation topology *visible* in the resource graph:

```
GcpWorkloadIdentityPoolProvider ──principal──▶ GcpServiceAccountIamMember
                                                      │ role: workloadIdentityUser
                                                      ▼
                                              GcpServiceAccount (the impersonated identity)
```

## Additive Is the Only Composable Mode

Like every IAM surface in GCP, the account policy can be written three ways: additive member, authoritative per-role binding, and authoritative whole-policy. Only the additive member is safe for composition — the authoritative modes clobber every grant they do not list, so two independent tools managing the same role silently remove each other's access. Planton deliberately models only the additive grant; the omission of binding/policy modes is a stopping line, not a coverage gap.

## Why Every Field Is Immutable

The IAM API has no "update grant" operation. What looks like an edit is always a remove-and-add on the policy, and any tooling that pretends otherwise invents an update semantic the platform beneath it does not have. This component mirrors reality: all fields are create-time, and any change replaces the grant atomically. A condition, when present, is part of the grant's identity — the same role granted with and without a condition are two independent policy entries.

## Choosing Between the Grant Kinds

- **This component** — the general primitive for any principal on any service account: federated principal sets, cross-account impersonation, users, groups, conditional grants.
- **GcpGkeWorkloadIdentityBinding** — the GKE convenience for the same underlying grant: it derives the `serviceAccount:<project>.svc.id.goog[<namespace>/<ksa>]` principal from cluster coordinates so the principal string is never hand-assembled. Prefer it when the impersonating workload is a GKE pod.
- **GcpProjectIamMember** — grants on the *project's* policy, including roles the account holds as an identity. Use it for what the account may do; use this component for who may use the account.
