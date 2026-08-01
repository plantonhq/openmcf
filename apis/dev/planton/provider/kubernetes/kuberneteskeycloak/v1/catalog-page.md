# Keycloak

Open-source identity and access management: single sign-on, OIDC and
SAML, identity brokering and user federation. This component declares
the Keycloak CR the official Keycloak Operator reconciles into a
clustered StatefulSet — production posture as typed fields, not chart
values (a `KubernetesKeycloakOperator` watching the namespace is the
prerequisite).

## Highlights

- **Crash-loops become validation errors** — Keycloak refuses to
  start with neither TLS nor HTTP, or without its public hostname;
  upstream tells you via CrashLoopBackOff, this spec tells you at
  apply time.
- **The database is a decision, not a default** — a required typed
  vendor instead of the silent embedded-H2 fallback; a
  `KubernetesPostgres` composes naturally (host from its read-write
  Service, credentials as Secret selectors), and the sandbox vendors
  are honestly labeled: ephemeral, single-instance, never production.
- **Credentials never ride the CR** — database credentials are Secret
  selectors by upstream design, and the one-time
  `<name>-initial-admin` Secret the operator generates is exported as
  the credential handle.
- **The update truth, told** — image changes take a scale-to-zero
  recreate by default (two versions cannot share one cache cluster);
  `Auto` and `Explicit` strategies when you want the trade-off
  differently.
- **Exposure composes** — the operator's own Ingress is always
  disabled; Gateway API kinds reference the exported service handles,
  and the operator's NetworkPolicy/ServiceMonitor defaults surface as
  explicit fields.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
