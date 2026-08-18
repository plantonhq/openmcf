# AwsPrivateCa

An AWS Private Certificate Authority — the managed CA behind internal TLS (mTLS services, MSK/EKS client auth, service meshes) — with its activation, issued certificates, ACM renewal permission, and resource policy managed in-line.

## Highlights

- **Activation is composed, never user-choreographed**: a ROOT CA self-signs at apply (the CSR→issue→install dance the raw provider makes you wire by hand), so one deploy produces an ACTIVE, issuing CA; a SUBORDINATE activates from a parent AwsPrivateCa by reference, or waits in PENDING_CERTIFICATE for an external parent (the `ca_csr` output is what that parent signs).
- **Cost is modeled honestly**: GENERAL_PURPOSE bills USD 400/month prorated hourly from creation to deletion; SHORT_LIVED_CERTIFICATE mode is USD 50/month with 7-day certificates and no revocation — the mode choice is a billing decision the spec teaches.
- **Certificates without key custody**: `issued_certificates` take YOUR CSRs — private keys never touch AWS or the manifest; delete revokes (the provider's documented semantic).

## Both Engines

Both modules compose the CA, activation, certificates, permission, and policy identically and export the same outputs: `certificate_authority_arn` (import ID and the universal join key), `certificate_authority_id`, `ca_certificate`/`ca_certificate_chain` (from the ACTIVATION path — the CA attribute reads empty at create), `ca_csr`, `activation_certificate_arn`, and the `issued_certificate_arns` map keyed by name.

## Chart Wiring

`certificate_authority_arn` → AwsMskCluster's `tls_certificate_authority_arns`, AwsCertManagerCert's issuing CA, or a subordinate AwsPrivateCa's `subordinate_activation.parent_ca_arn`; `revocation.crl.s3_bucket_name` → AwsS3Bucket `bucket_id`. Set `acm_renewal_permission` whenever ACM certificates are issued from this CA — without it their renewals fail silently at expiry.
