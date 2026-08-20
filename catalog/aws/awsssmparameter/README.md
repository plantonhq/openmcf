<p align="center">
  <img src="logo.svg" alt="AWS SSM Parameter" width="80"/>
</p>

# AWS SSM Parameter

Manage an [AWS Systems Manager Parameter Store entry](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html)
— one named configuration value (String, StringList, or SecureString)
applications read at runtime.

## What Gets Managed

- **The parameter** (`spec.parameterName` is the name — an explicit
  field because hierarchical paths like `/prod/db/url` carry slashes
  `metadata.name` cannot): its **type**, **value** (exactly one of the
  plain `value` arm — readable in plans, the common config case — or
  the secret `secureValue` arm, stored as a managed-secret reference
  and resolved just-in-time at deploy), **description**, a
  server-enforced **allowedPattern**, the **tier** (Standard free /
  Advanced 8KB+policies / Intelligent-Tiering resolved by AWS per
  write), the **KMS key** for SecureString encryption, the **dataType**
  (`aws:ec2:image` validates values as AMI IDs), and first-create
  **overwrite** behavior.

SecureString parameters MUST use `secureValue` — the spec rejects the
plain arm for them, so plaintext never reaches engine plans or state.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
