# Multi-tenant shared Loki

One Loki serving several teams with isolation: every push and query
carries an `X-Scope-OrgID` tenant header, and the gateway enforces HTTP
basic auth per tenant. Passwords are supplied as bcrypt htpasswd hashes
(one-way — safe to commit; the plaintext never appears), so no Secret
fixture is needed.

**When to use:** a platform team running logging as a shared service with
per-team separation and access control.

**Wiring:** each team points its Grafana `loki` datasource and its
KubernetesOtelCollector pipeline at the gateway with its own tenant name
and password.
