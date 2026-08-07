# AzureFrontDoorCustomDomain

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFrontDoorCustomDomainSpec** defines the configuration for
creating a custom domain inside an Azure Front Door (Standard/Premium)
profile -- your own hostname (e.g. www.example.com) served by Front
Door instead of the generated *.azurefd.net endpoint hostname.

A custom domain has a two-step lifecycle:

 1. **Create** -- the domain deploys immediately in a
    pending-validation state and exports a validation_token.
 2. **Validate** -- publish the token as a DNS TXT record at
    `_dnsauth.<host_name>` (Azure checks it and flips the domain to
    approved), then CNAME the hostname at your DNS provider to the
    endpoint's host_name output so traffic actually arrives.

When dns_zone_id references an AzureDnsZone, Azure watches that zone
for the records; the records themselves are DNS resources you manage
in the zone. After validation, attach the domain to routes via
AzureFrontDoorRoute's custom_domain_ids -- the route side owns the
attachment.

TLS is always on for custom domains: either Azure manages the
certificate end to end (MANAGED_CERTIFICATE, the default -- free,
auto-rotated), or the domain uses a bring-your-own certificate wrapped
by an AzureFrontDoorSecret (CUSTOMER_CERTIFICATE -- for EV/OV
certificates, pinned chains, or org-mandated CAs).

**ForceNew fields**: `profile_id`, `domain_name`, `host_name` -- they
fix the domain's ARM identity and the hostname it proves ownership
of. The TLS block and dns_zone_id update in place.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFrontDoorCustomDomain
metadata:
  name: test-front-door-custom-domain
spec:
  profileId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor
  domainName: www-example-com
  hostName: www.example.com
  # Exercises the Azure DNS zone seam (Azure watches the zone for the
  # validation TXT record).
  dnsZoneId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/dnsZones/example.com
  # Exercises the customer-certificate pairing and the CUSTOMIZED cipher
  # policy with both TLS 1.3 suites pinned.
  tls:
    certificateType: CUSTOMER_CERTIFICATE
    secretId:
      value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor/secrets/wildcard-example-com
    cipherSuite:
      type: CUSTOMIZED
      customCiphers:
        tls12: [ECDHE_RSA_AES128_GCM_SHA256, ECDHE_RSA_AES256_GCM_SHA384]
        tls13: [TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384]
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.profileId` | `string \| valueFrom` | yes |  | AzureFrontDoorProfile (`status.outputs.profile_id`) |
| `spec.domainName` | `string` | yes |  |  |
| `spec.hostName` | `string` | yes |  |  |
| `spec.dnsZoneId` | `string \| valueFrom` |  |  | AzureDnsZone (`status.outputs.zone_id`) |
| `spec.tls` | `AzureFrontDoorCustomDomainTls` | yes |  |  |
| `spec.tls.certificateType` | `enum` |  |  |  |
| `spec.tls.secretId` | `string \| valueFrom` |  |  | AzureFrontDoorSecret (`status.outputs.secret_id`) |
| `spec.tls.cipherSuite` | `AzureFrontDoorCustomDomainCipherSuite` |  |  |  |
| `spec.tls.cipherSuite.type` | `enum` |  |  |  |
| `spec.tls.cipherSuite.customCiphers` | `AzureFrontDoorCustomDomainCustomCiphers` |  |  |  |
| `spec.tls.cipherSuite.customCiphers.tls12` | `[]string` | yes |  |  |
| `spec.tls.cipherSuite.customCiphers.tls13` | `[]string` |  |  |  |

## Field Details

### spec.profileId

`string | valueFrom` · required

The Front Door profile the custom domain lives in, by ARM ID.
References an AzureFrontDoorProfile's profile_id output so the
profile and its domains compose in one manifest set. Fixed at
creation.

- references: AzureFrontDoorProfile (`status.outputs.profile_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorProfile, name: <that resource's name>, fieldPath: status.outputs.profile_id}} -- a bare string does not parse

### spec.domainName

`string` · required

The domain resource's name -- unique within the profile. This is
the ARM resource name, NOT the hostname (that is host_name);
convention is the hostname with dots replaced by hyphens, e.g.
"www-example-com".

2-260 characters; letters, digits, and hyphens; must start and end
with a letter or digit.

**ForceNew**: changing the name replaces the domain (re-validation
required).

- rule: domain_name must be 2-260 characters, start and end with a letter or digit, and contain only letters, digits, and hyphens (dots are not allowed -- this is the resource name, not the hostname)
- rule: {"required":true}

### spec.hostName

`string` · required

The hostname this domain serves, e.g. "www.example.com" or the
wildcard "*.example.com" (wildcards need a CUSTOMER_CERTIFICATE --
Azure's managed certificates do not cover them). Must be a fully
qualified domain name you control: validation publishes a TXT
record under it, and serving traffic requires a CNAME from it to
the endpoint's hostname.

**ForceNew**: the hostname IS the domain's identity.

- rule: host_name must be a fully qualified domain name of dot-separated labels (each 1-63 characters, letters/digits/hyphens, not starting or ending with a hyphen); only the FIRST label may be the wildcard '*'
- rule: {"required":true,"string":{"maxLen":"253"}}

### spec.dnsZoneId

`string | valueFrom`

The Azure DNS zone that hosts this hostname's records, by ARM ID.
References an AzureDnsZone's zone_id output. Optional -- set it
when the domain's DNS lives in Azure DNS so Front Door watches the
zone for the validation TXT record; leave unset when DNS is hosted
elsewhere (validation then relies on you publishing the token at
your external provider).

- references: AzureDnsZone (`status.outputs.zone_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.tls

`AzureFrontDoorCustomDomainTls` · required

The domain's TLS configuration. Required -- custom domains always
terminate TLS. Leave the block at its defaults for an Azure-managed
certificate.

- rule: {"required":true}
- rule: secret_id is required with CUSTOMER_CERTIFICATE (it names the certificate to serve) and must stay unset with MANAGED_CERTIFICATE (Azure provides the certificate)

### spec.tls.certificateType

`enum`

Who provides the certificate. MANAGED_CERTIFICATE (default when
unspecified): Azure issues, hosts, and auto-rotates a free DV
certificate for the exact hostname. CUSTOMER_CERTIFICATE: the
domain serves the certificate wrapped by the AzureFrontDoorSecret
referenced in secret_id -- required for wildcard hostnames and for
EV/OV or org-mandated certificates.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_front_door_custom_domain_certificate_type_unspecified` -- Not specified -- deploys MANAGED_CERTIFICATE, Azure's default.
- `MANAGED_CERTIFICATE` -- Azure issues, hosts, and auto-rotates a free DV certificate for the exact hostname (no wildcards; hostname up to 64 characters).
- `CUSTOMER_CERTIFICATE` -- The domain serves a bring-your-own certificate wrapped by the AzureFrontDoorSecret referenced in secret_id.

### spec.tls.secretId

`string | valueFrom`

The Front Door secret carrying the bring-your-own certificate, by
ARM ID. References an AzureFrontDoorSecret's secret_id output.
Required with CUSTOMER_CERTIFICATE; must stay unset with
MANAGED_CERTIFICATE (Azure provides the certificate -- there is
nothing to reference).

- references: AzureFrontDoorSecret (`status.outputs.secret_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorSecret, name: <that resource's name>, fieldPath: status.outputs.secret_id}} -- a bare string does not parse

### spec.tls.cipherSuite

`AzureFrontDoorCustomDomainCipherSuite`

A hardened cipher-suite policy for the domain's TLS handshakes.
Unset serves Azure's default suite set. Set it to pin one of
Azure's predefined suite sets or to hand-pick suites (CUSTOMIZED)
when a compliance baseline demands specific ciphers.

- rule: custom_ciphers is required when type is CUSTOMIZED (it lists the suites to serve) and must stay unset with the predefined TLS12_2022 / TLS12_2023 sets

### spec.tls.cipherSuite.type

`enum`

The suite-set policy. TLS12_2022 and TLS12_2023 are Azure's
predefined hardened sets (year-versioned; the custom_ciphers lists
must stay unset with them). CUSTOMIZED serves exactly the suites
listed in custom_ciphers.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_custom_domain_cipher_suite_set_type_unspecified` -- Not specified -- invalid; pick a predefined set or CUSTOMIZED.
- `TLS12_2022` -- Azure's 2022 hardened TLS 1.2 suite set.
- `TLS12_2023` -- Azure's 2023 hardened TLS 1.2 suite set (the most current predefined policy).
- `CUSTOMIZED` -- Serve exactly the suites listed in custom_ciphers.

### spec.tls.cipherSuite.customCiphers

`AzureFrontDoorCustomDomainCustomCiphers`

The hand-picked suites for the CUSTOMIZED type. Required with
CUSTOMIZED; forbidden with the predefined sets.

### spec.tls.cipherSuite.customCiphers.tls12

`[]string` · required

The TLS 1.2 suites to serve, from Azure's allowed ECDHE-RSA set. At
least one.

- rule: {"repeated":{"minItems":"1","unique":true,"items":{"cel":[{"id":"front_door_custom_cipher_tls12_value","message":"each TLS 1.2 suite must be one of Azure's allowed set: ECDHE_RSA_AES128_GCM_SHA256, ECDHE_RSA_AES256_GCM_SHA384, ECDHE_RSA_AES128_SHA256, ECDHE_RSA_AES256_SHA384","expression":"this in ['ECDHE_RSA_AES128_GCM_SHA256', 'ECDHE_RSA_AES256_GCM_SHA384', 'ECDHE_RSA_AES128_SHA256', 'ECDHE_RSA_AES256_SHA384']"}]}}}

### spec.tls.cipherSuite.customCiphers.tls13

`[]string`

The TLS 1.3 suites to serve. Leave empty to let Azure serve its
TLS 1.3 defaults; when set explicitly, Azure requires BOTH suites
to be listed (TLS 1.3 mandates them).

- rule: when tls13 is set it must list BOTH TLS_AES_128_GCM_SHA256 and TLS_AES_256_GCM_SHA384 -- TLS 1.3 mandates both suites; leave the list empty to serve Azure's defaults
- rule: {"repeated":{"unique":true,"items":{"cel":[{"id":"front_door_custom_cipher_tls13_value","message":"each TLS 1.3 suite must be one of: TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384","expression":"this in ['TLS_AES_128_GCM_SHA256', 'TLS_AES_256_GCM_SHA384']"}]}}}

## Validation Rules

- `front_door_custom_domain_managed_cert_host_name`: with an Azure-managed certificate the host_name cannot exceed 64 characters and cannot be a wildcard -- use a CUSTOMER_CERTIFICATE (an AzureFrontDoorSecret wrapping your own certificate) for wildcard or very long hostnames

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFrontDoorCustomDomain, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.custom_domain_id` | `string` | The Azure Resource Manager ID of the custom domain -- what AzureFrontDoorRoute's custom_domain_ids references to serve this hostname, and what a Front Door security policy scopes a WAF to. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cdn/profiles/{profile}/customDomains/{name} |
| `status.outputs.host_name` | `string` | The hostname the domain serves (echoed from the spec) -- the name to CNAME to the endpoint's host_name at your DNS provider once validation passes. |
| `status.outputs.validation_token` | `string` | The DNS validation challenge: publish this value as a TXT record at `_dnsauth.<host_name>` and Azure flips the domain from pending to approved. The token is public by design -- it proves control of the DNS name, not possession of a secret. |
| `status.outputs.expiration_date` | `string` | When the current validation token expires (RFC-3339 timestamp). A domain left unvalidated past this instant needs a fresh token (Azure regenerates it on re-validation). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.profileId` | AzureFrontDoorProfile | `status.outputs.profile_id` |
| `spec.dnsZoneId` | AzureDnsZone | `status.outputs.zone_id` |
| `spec.tls.secretId` | AzureFrontDoorSecret | `status.outputs.secret_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFrontDoorRoute | `spec.customDomainIds` | `status.outputs.custom_domain_id` |

## See Also

- [Overview](../README.md)
