# AwsPrivateCa — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The mode decision is a 8x billing decision

GENERAL_PURPOSE (USD 400/month prorated) exists for certificates that live long enough to need revocation. SHORT_LIVED_CERTIFICATE (USD 50/month) caps end-entity validity at 7 days and forbids CRL/OCSP — revocation by expiry. High-churn workload identity (mesh mTLS, batch jobs) belongs in short-lived mode; anything a human installs belongs in general-purpose. The mode is fixed for life — and the provider silently ignores an update attempt, so get it right at create.

## Deleting is cheap; the slot lingers

Billing stops the moment the CA deletes, but the CA parks restorable for `permanent_deletion_time_in_days` (default 30). Ephemeral environments should set 7 and use run-scoped names — a recreate under the same subject works fine, but restore-window records clutter the console and hold ARNs.

The window is fixed at create (live-proven 2026-08-26): AWS stores nothing — only the delete call consumes it — and the provider cannot apply a window-only edit (it composes an empty update AWS rejects), so both engines ignore post-create changes and the delete uses the create-time value. Decide the window when you create the CA. An imported CA always deletes with the provider-default 30-day window regardless of the manifest value; to delete it faster, use `aws acm-pca delete-certificate-authority --permanent-deletion-time-in-days 7` directly.

## ACM renewal needs the permission, and failures are silent

AwsCertManagerCert certificates issued from this CA renew automatically ONLY while acm.amazonaws.com holds all three grant actions. The `acm_renewal_permission` flag manages the full grant; without it, renewals fail with no deploy-time symptom — the certificate just expires months later. Set it on any CA that issues ACM certificates, always.

## CRL buckets are pre-work, enforced at create

AWS validates the CRL bucket's policy when the CA is CREATED (the service needs getBucketAcl/putObject) — a missing policy fails the deploy, not the first CRL. Public-read objects are the default posture; a private bucket needs `s3_object_acl: BUCKET_OWNER_FULL_CONTROL` plus a `custom_cname` fronting it (CloudFront), because relying parties must reach the CRL URL baked into every issued certificate.

## Same-family hierarchies, by design

The composed subordinate activation signs with THIS spec's signing_algorithm at the parent — correct when parent and child share a key family (the norm). A cross-family hierarchy (RSA root signing EC subordinates) needs out-of-band activation via the `ca_csr` output.

## Templates decide what a certificate can DO

The `template_arn` on issued certificates picks the X.509 profile: EndEntityCertificate (plain TLS), EndEntityClientAuthCertificate (mTLS clients), CodeSigningCertificate, or the CA-path-length templates that mint sub-CAs. The CSR asks; the template decides — a CSR with CA:TRUE extensions issued under an end-entity template comes out a plain leaf.

## Certificates are born, not edited (live-proven 2026-08-26)

Every issue request's parameters — CSR, signing algorithm, validity, template, API passthrough — are the certificate's birth certificate: AWS reads them once and never reports them back, so both engines deliberately ignore later edits to them. Editing an issued certificate's entry does nothing; to get different parameters, add a NEW entry (its name keys a new certificate) and retire the old one when consumers have moved. The same contract protects adopters at import: bringing an existing CA under management plans zero changes — without it, the first plan after import would re-issue the CA's own trust anchor, a certificate your trust stores pin by fingerprint. Reissuing or renewing the CA's certificate itself is an operational act through ACM PCA, never a manifest edit.
