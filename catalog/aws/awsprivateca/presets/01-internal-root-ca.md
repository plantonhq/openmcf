# Internal Root CA

The one-deploy trust anchor: a self-activating 10-year root that comes up ACTIVE and issuing, with ACM pre-authorized to auto-renew the certificates it requests (the grant renewals silently depend on). Distribute the `ca_certificate` output to trust stores and issue internal TLS through AwsCertManagerCert against the `certificate_authority_arn`. GENERAL_PURPOSE mode: USD 400/month prorated — one root serves the whole organization.
