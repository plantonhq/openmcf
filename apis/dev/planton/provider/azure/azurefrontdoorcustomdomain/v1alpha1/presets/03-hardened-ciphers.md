# Hardened Cipher Policy

This preset creates a managed-certificate domain with a hand-picked
cipher policy: GCM-only TLS 1.2 suites (dropping the CBC pair) plus the
mandatory TLS 1.3 suites -- for compliance baselines that name allowed
ciphers explicitly.

## When to Use

- PCI-DSS, FedRAMP, or internal crypto baselines that forbid
  CBC-mode suites
- When an auditor asks "which ciphers does this endpoint negotiate?"
  and the answer must be a fixed list

## Key Configuration Choices

- **`type: CUSTOMIZED` requires `customCiphers`** (and forbids it on the
  predefined sets) -- the spec enforces the pairing
- **Prefer `TLS12_2023` over CUSTOMIZED** when Azure's newest
  predefined set satisfies the baseline -- predefined sets get
  Microsoft's updates automatically
- **`tls13` lists BOTH suites** -- TLS 1.3 mandates them when pinned;
  leave the list out entirely to serve Azure's TLS 1.3 defaults
- **The TLS floor is always 1.2** -- Azure retired TLS 1.0/1.1, so there
  is no version knob to harden

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `hostName` (example value) | Your real hostname | Your DNS zone |
