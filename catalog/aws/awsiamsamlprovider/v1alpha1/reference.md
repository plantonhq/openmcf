# AwsIamSamlProvider

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

AwsIamSamlProviderSpec defines one IAM SAML identity provider: the
account's trust anchor for SAML 2.0 federation. It is created from
the identity provider's metadata document (the XML the IdP - Okta,
Entra ID, Google Workspace, ... - publishes, carrying its issuer
and signing certificates), and IAM roles then trust it via
sts:AssumeRoleWithSAML in their trust policies.

The provider's name comes from metadata.name and is WRITE-ONCE at
AWS - renaming replaces the provider (and every role trust policy
naming its ARN must be updated). The metadata document updates in
place - certificate rotations are ordinary updates.

IAM is a GLOBAL service: the provider exists account-wide, and the
spec's region is only the provider endpoint region.

## Example

```yaml
# Canonical AwsIamSamlProvider example (hack/dev manifest and refgen
# Example source): a federation trust anchor created from an IdP
# metadata document. The document below is a structurally realistic
# stand-in over AWS's 1000-character floor -- a real IdP (Okta, Entra
# ID, Google Workspace) publishes the equivalent XML at a metadata
# URL.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsIamSamlProvider
metadata:
  name: corporate-idp
  id: corporate-idp
  org: test-org
  env: dev
spec:
  region: us-east-1
  samlMetadataDocument: |
    <?xml version="1.0" encoding="UTF-8"?>
    <md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://idp.example.com/saml/metadata" validUntil="2030-01-01T00:00:00Z">
      <md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:KeyDescriptor use="signing">
          <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
            <ds:X509Data>
              <ds:X509Certificate>MIIDdzCCAl+gAwIBAgIEbGRkZDANBgkqhkiG9w0BAQsFADBsMRAwDgYDVQQGEwdF
              eGFtcGxlMRAwDgYDVQQIEwdFeGFtcGxlMRAwDgYDVQQHEwdFeGFtcGxlMRAwDgYD
              VQQKEwdFeGFtcGxlMRAwDgYDVQQLEwdFeGFtcGxlMRAwDgYDVQQDEwdFeGFtcGxl
              MB4XDTI2MDEwMTAwMDAwMFoXDTMwMDEwMTAwMDAwMFowbDEQMA4GA1UEBhMHRXhh
              bXBsZTEQMA4GA1UECBMHRXhhbXBsZTEQMA4GA1UEBxMHRXhhbXBsZTEQMA4GA1UE
              ChMHRXhhbXBsZTEQMA4GA1UECxMHRXhhbXBsZTEQMA4GA1UEAxMHRXhhbXBsZTCC
              ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKzExampleOnlyNotARealKey
              </ds:X509Certificate>
            </ds:X509Data>
          </ds:KeyInfo>
        </md:KeyDescriptor>
        <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://idp.example.com/saml/sso"/>
        <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://idp.example.com/saml/sso"/>
      </md:IDPSSODescriptor>
    </md:EntityDescriptor>
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.samlMetadataDocument` | `string` | yes |  |  |

## Field Details

### spec.region

`string` · required

The AWS region used by the provider while managing the SAML
provider. IAM is a global service - every AWS API call is still
made against a regional endpoint. Example: "us-east-1".

- rule: {"string":{"minLen":"1"}}

### spec.samlMetadataDocument

`string` · required

The IdP's SAML 2.0 metadata document - the XML verbatim, as
published by the identity provider (issuer, SSO endpoints, and
public signing certificates). This is a PUBLIC trust document,
not a secret - IdPs serve it unauthenticated. AWS requires
1000-10000000 characters; update it in place when the IdP rotates
its signing certificates (see valid_until in the outputs).

- rule: {"string":{"minLen":"1000","maxLen":"10000000"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsIamSamlProvider, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.provider_arn` | `string` | The provider's ARN (also the provider's import ID, and the value role trust policies reference via the SAML principal). |
| `status.outputs.saml_provider_uuid` | `string` | The provider's AWS-assigned UUID. |
| `status.outputs.valid_until` | `string` | When the metadata document's certificates expire - rotate the document before this date or federation breaks. |

## See Also

- [Overview](../README.md)
