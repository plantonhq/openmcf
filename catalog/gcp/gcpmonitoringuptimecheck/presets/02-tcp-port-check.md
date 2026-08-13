# TCP Port Check

Asserts a non-HTTP service accepts connections — message brokers,
databases with public endpoints, custom TCP protocols. The probe passes
when the port completes a TCP handshake.

## What it configures

- `tcpCheck.port` — the port to connect to (no protocol-level assertion;
  a handshake is the whole test).
- `period: 60s` — the tightest cadence GCP offers, appropriate for
  connection-level monitoring of stateful backends.

## Adjust before deploying

- **host / port** — the endpoint to probe. The target must be reachable
  from Google's public checker fleet; for private endpoints use
  `checkerType: VPC_CHECKERS`.

## When to choose something else

If the service speaks HTTP, the **Public HTTPS Check** preset asserts
far more (status, TLS, body). A TCP handshake passing proves the
listener is alive, not that the service behind it is healthy.
