# Destination Address

Register a destination mailbox that Email Routing rules and catch-alls can forward
to. A verification email is sent to the address on creation.

## When to use

- Adding the real inbox a domain's mail should be forwarded to.

## Key choices

- `email` — the destination mailbox. The owner must click the verification link
  before the address can receive forwarded mail.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
| `destination@replaceme.example.com` | The destination mailbox address |
