# AWS IAM SAML Provider

The account's SAML federation trust anchor: created from your
identity provider's metadata XML, it lets IAM roles trust corporate
sign-ins (Okta, Entra ID, Google Workspace) via
sts:AssumeRoleWithSAML.

## What Gets Managed

- The provider, created from the IdP's metadata document (issuer, SSO
  endpoints, public signing certificates — a public document, not a
  secret).
- In-place metadata updates for certificate rotations, with the
  expiry date surfaced as an output.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with IAM permissions.

### AWS Account

- The metadata XML from your IdP's metadata URL (most IdPs publish
  it unauthenticated).
- Roles that trust the provider come next
  ([AWS IAM Role](/cloud-catalog/aws-iam-role)).

## Deploy

### Console

Create the resource from the AWS catalog, paste the metadata
document, and deploy.

### CLI

```bash
planton apply -f saml-provider.yaml
```

## After Deploy

- Reference the provider's ARN output in role trust policies
  (`Federated` principal + sts:AssumeRoleWithSAML) and configure the
  IdP's AWS app with the role/provider ARN pairs.
- Rotate the metadata before `valid_until` or federation breaks.
