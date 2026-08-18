# AWS Private Certificate Authority

Your own certificate authority, minus the HSM room: AWS Private CA issues the internal TLS certificates behind mTLS services, MSK client auth, service meshes, and device fleets — with the root-of-trust ceremony reduced to one deploy.

## What Gets Managed

- The CA: root or subordinate, key and signing algorithms (post-quantum ML-DSA included), X.500 subject, usage mode, HSM standard.
- Activation: a root self-signs at apply; a subordinate is signed by a parent AwsPrivateCa by reference.
- Revocation: CRLs published to S3 and/or the managed OCSP responder.
- Certificates issued from your CSRs, the ACM auto-renewal permission, and the cross-account resource policy.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with ACM PCA permissions.

### AWS Prerequisites

- None for a plain root CA. CRL publishing needs an S3 bucket whose policy grants the ACM PCA service getBucketAcl/putObject BEFORE the CA is created (AWS validates at create).

## After You Deploy

- **Billing starts at create and stops at delete**: USD 400/month (GENERAL_PURPOSE) or USD 50/month (SHORT_LIVED_CERTIFICATE), prorated hourly, plus per-certificate issuance fees.
- Relying parties trust the `ca_certificate` output; consumers reference `certificate_authority_arn`.

## Common Changes

- Issue or revoke certificates (add/remove `issued_certificates` entries), toggle revocation arms, grant/revoke the ACM permission — all in-place.
- Algorithms, subject, and type are fixed for life; a deleted CA stays restorable for `permanent_deletion_time_in_days` (7-30).
