# AwsPrivateCa — Terraform/OpenTofu module

Manages one private CA (`aws_acmpca_certificate_authority`) with composed activation (`aws_acmpca_certificate` root_ca/subordinate_ca + `aws_acmpca_certificate_authority_certificate`), issued certificates (`aws_acmpca_certificate.issued`, keyed by name), the ACM permission (`aws_acmpca_permission`), and the resource policy (`aws_acmpca_policy`).

Module facts worth knowing before editing:

- **The ca_certificate output comes from the ACTIVATION path** — the CA resource's own certificate attribute is read at create, while the CA is still PENDING_CERTIFICATE (empty).
- **aws_acmpca_certificate plays three roles** (root self-sign / subordinate activation / issued) — the import bridge scopes each logical name separately.
- **Issued certificates depend on the activation** — issuing needs the CA ACTIVE; their delete REVOKES them.
- **Template ARNs are partition-scoped** (`data.aws_partition`) — never hardcode `arn:aws`.
- **permanent_deletion_time_in_days is never read back** — declared config-only in the import catalog.

Outputs mirror the Pulumi module key-for-key: `certificate_authority_arn` (import ID), `certificate_authority_id`, `ca_certificate`, `ca_certificate_chain`, `ca_csr`, `activation_certificate_arn`, `issued_certificate_arns`.
