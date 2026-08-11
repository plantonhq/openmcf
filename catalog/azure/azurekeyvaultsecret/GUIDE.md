# Azure Key Vault Secret -- Operational Guide

Judgment calls that matter when you run Key Vault secrets in production.

## Reference the versionless ID unless compliance says otherwise

Every value change creates a NEW version; the versioned ID (`secret_id`) freezes consumers on the old value while the versionless ID follows the latest. Default every consumer to `versionless_id` -- rotation then propagates with zero intervention. Pin a version only when a compliance regime demands a frozen value, and own the re-pin as part of rotation.

## The value never belongs in a manifest

The `value` field is sensitive and reference-resolved: point it at a managed secret or another resource's output. A literal value in a manifest defeats the vault's purpose -- the manifest becomes the secret store, unencrypted and version-controlled. Treat any literal that slips through as leaked and rotate it.

## Expiry is an attribute, not an enforcement

`expirationDate` and `notBeforeDate` are advisory: Key Vault returns them, Azure Policy can audit them, and well-behaved consumers honor them -- but an RBAC-permitted reader can still fetch an expired value. Use expiry to make rotation auditable, and pair it with monitoring on near-expiry secrets rather than trusting it as a lock.

## Deleted names linger -- know your vault's posture

Soft delete reserves a deleted secret's name for the vault's retention window. The module purges on destroy so the name frees immediately -- EXCEPT when the vault has purge protection, where the reservation holds for the full window by design. In purge-protected vaults, treat secret names as append-only: version values, don't recycle names.

## Multi-line values need encoding

Key Vault strips raw newlines. Certificates-as-secrets, SSH keys, and JSON blobs should be base64-encoded (record it in `contentType`, e.g. "application/x-pem-file;base64") so consumers know to decode.
