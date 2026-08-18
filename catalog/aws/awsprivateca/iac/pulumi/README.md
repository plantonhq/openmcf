# AwsPrivateCa — Pulumi module

Manages one private CA (`acmpca.CertificateAuthority`) with composed activation (`acmpca.Certificate` + `acmpca.CertificateAuthorityCertificate`), issued certificates (one `acmpca.Certificate` per spec entry), the ACM permission (`acmpca.Permission`), and the resource policy (`acmpca.Policy`).

Module facts worth knowing before editing:

- **The ca_certificate output comes from the ACTIVATION path** — the CA resource's own certificate attribute is read at create, while the CA is still PENDING_CERTIFICATE (empty).
- **acmpca.Certificate plays three roles** (root-ca-certificate / subordinate-ca-certificate / certificate-{name}) — resource names match the Terraform module's logical names role-for-role.
- **Issued certificates DependOn the activation** — issuing needs the CA ACTIVE; their delete REVOKES them.
- **Template ARNs are partition-scoped** (`aws.GetPartition`) — never hardcode `arn:aws`.
- **permanent_deletion_time_in_days is never read back** — declared config-only in the import catalog.

Outputs mirror the Terraform module key-for-key: `certificate_authority_arn` (import ID), `certificate_authority_id`, `ca_certificate`, `ca_certificate_chain`, `ca_csr`, `activation_certificate_arn`, `issued_certificate_arns`.
