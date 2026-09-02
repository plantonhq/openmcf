---
display_name: Postgres mTLS
---

# PostgreSQL Hyperdrive with mTLS

A Hyperdrive config that connects to a PostgreSQL origin requiring mutual TLS,
verifying both the CA and the server hostname.

## When to use

- The origin database requires client certificates (mutual TLS).
- You need full certificate verification (`verify-full`) for the strongest
  transport security.

## Key choices

- `mtls.caCertificateId` / `mtls.mtlsCertificateId`: IDs of certificates already
  uploaded to Cloudflare (the CA used to verify the server, and the client cert
  Hyperdrive presents).
- `mtls.sslmode`: `verify-full` verifies the CA and hostname; `verify-ca`
  verifies only the CA; `require` encrypts without verification.
- `originConnectionLimit`: raise from the default to handle higher concurrency
  on paid plans (5–100).

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
| `<database-name>` | Name of the origin database |
| `<database-user>` | Database user Hyperdrive authenticates as |
| `<database-host>` | Hostname/IP of the origin database |
| `<database-password>` | Password for the database user (managed secret) |
| `<ca-certificate-id>` | ID of the uploaded CA certificate |
| `<client-certificate-id>` | ID of the uploaded client certificate |
