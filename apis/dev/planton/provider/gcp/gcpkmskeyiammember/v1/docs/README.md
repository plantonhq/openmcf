# Key-Scoped KMS IAM: CMEK Grants as Graph Edges

## The CMEK Permission Everyone Hits

Customer-managed encryption has one universal prerequisite: before a service can encrypt anything with your key, that service's agent — a Google-managed service account like `service-<project_number>@gs-project-accounts.iam.gserviceaccount.com` for Cloud Storage — must hold `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key. Every CMEK-enabled bucket, dataset, disk, topic, and instance in the catalog shares this requirement, which makes the grant itself the most-composed IAM object in an encrypted estate.

## Why the Grant Is Key-Scoped

IAM on KMS flows down the resource hierarchy — project covers every ring, ring covers every key — so the same permission can be granted three ways. Key-scoped is the right default for two reasons:

1. **Least privilege.** A ring quickly accumulates keys with different consumers: the state-backend key, the database key, the artifact key. A ring- or project-level `cryptoKeyEncrypterDecrypter` grant hands every consumer's agent access to every key. The key-scoped grant hands each agent exactly the key it encrypts with.
2. **Ordering.** A first CMEK deploy has a hidden race: the encrypted resource is created moments after the IAM grant, and project-level IAM propagation can lose that race, failing the deploy with a permission error that self-heals minutes later. Modeled as its own node referencing the key, the grant becomes a real dependency edge — the bucket (or dataset, or disk) depends on the grant, the grant depends on the key, and orchestration serializes them correctly.

```
GcpKmsKey ──key_id──▶ GcpKmsKeyIamMember ◀──member── (service agent / GcpServiceAccount)
                              ▲
                    depends_on │
                  (encrypted resource: bucket, dataset, disk, ...)
```

## Additive Is the Only Composable Mode

Like every IAM surface in GCP, the key policy can be written three ways: additive member, authoritative per-role binding, and authoritative whole-policy. Only the additive member is safe for composition — the authoritative modes clobber every grant they do not list, so two independent tools managing the same role silently remove each other's access. Planton deliberately models only the additive grant; the omission of binding/policy modes is a stopping line, not a coverage gap.

## Why Every Field Is Immutable

The IAM API has no "update grant" operation. What looks like an edit is always a remove-and-add on the policy, and any tooling that pretends otherwise invents an update semantic the platform beneath it does not have. This component mirrors reality: all fields are create-time, and any change replaces the grant atomically. A condition, when present, is part of the grant's identity — the same role granted with and without a condition are two independent policy entries.

## Service Agents Are Lazily Created

A practical CMEK sharp edge worth knowing: many service agents do not exist until the service is first used in the project (some services expose an explicit "provision service identity" API call). If a grant fails because the member does not exist, the fix is usually to touch the consuming service once — or create any resource of that service — so GCP materializes the agent, then apply the grant.

## Scope Notes

Ring-scoped grants (covering every key in a ring) are expressible in the underlying provider but are not modeled as a kind: key-scoped grants cover the real composition need, and a ring-wide grant is usually the least-privilege smell this component exists to avoid. The scope can be revisited if a concrete architecture pulls for it.
