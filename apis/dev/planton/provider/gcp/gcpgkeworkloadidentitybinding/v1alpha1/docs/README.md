# GcpGkeWorkloadIdentityBinding: Design Notes

## What This Component Models

One additive IAM grant on a Google Service Account:
`roles/iam.workloadIdentityUser` to the workload-identity principal
`serviceAccount:<pool-project>.svc.id.goog[<namespace>/<ksa>]`. This is the
GCP half of GKE Workload Identity — the mechanism that lets a Kubernetes
ServiceAccount mint short-lived tokens as a Google Service Account with no
exported key material.

The underlying resource is a service-account-scoped IAM member — the same
additive-grant semantics as the generic project-scoped
`GcpProjectIamMember`, attached to a service account's policy instead of a
project's. The kind exists as its own specialized component because the
principal string is a structured value worth constructing, not typing.

## Why the Principal Is Constructed, Not Typed

The principal format is precise and unforgiving:

```
serviceAccount:<project>.svc.id.goog[<namespace>/<ksa-name>]
```

A single typo produces a grant that is syntactically valid, deploys
cleanly, and silently never matches the workload — the hardest kind of
failure to debug. The spec therefore takes the parts (`projectId`,
`ksaNamespace`, `ksaName`) and both engines assemble the string
identically. Namespace and name are validated against Kubernetes naming
rules (RFC 1123 label and DNS-subdomain respectively) so the invalid-parts
class is rejected before deploy.

## The Kubernetes Half Is Deliberately Out of Scope

Completing the handshake requires the KSA to carry the
`iam.gke.io/gcp-service-account: <gsa-email>` annotation. That annotation
lives on the Kubernetes object, in whatever chart or manifest owns the KSA
— which is also the right owner for the annotation: the workload's
deployment declares which GSA it intends to run as, and this component
grants that intent on the GCP side. Reaching into the cluster to mutate a
Kubernetes object from a GCP component would create a hidden cross-system
write with its own failure modes (cluster credentials, ordering, drift) and
would duplicate ownership of the KSA. The `service_account_email` output is
exactly the value the annotation needs, so composition stays explicit.

## Cross-Project Semantics

Two projects can appear in one binding:

- **The pool project (`projectId`)** — the GKE cluster's project. Every
  project has one implicit workload-identity pool named
  `<project>.svc.id.goog`; the principal always names the cluster project.
- **The GSA's own project** — embedded in its email. The modules pass the
  provider the fully-qualified resource name using the IAM API's `-`
  project wildcard (`projects/-/serviceAccounts/<email>`), which infers the
  SA's project from the email — correct even when the GSA lives outside
  the cluster project.

## IAM Condition

The optional `condition` mirrors the provider's condition block (title,
CEL expression, description). A condition is part of the grant's identity:
the same grant with and without a condition are two independent bindings
that do not interfere — useful for time-boxed migration grants.

## Immutability

Every field is immutable, mirroring the underlying API: an IAM grant has no
update. Any spec change replaces the grant atomically (destroy the old
pair, create the new), which is exactly how the IAM policy itself behaves.
