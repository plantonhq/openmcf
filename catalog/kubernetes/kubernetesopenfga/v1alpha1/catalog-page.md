# OpenFGA

Zanzibar-style authorization for any team: define types and
relations, store relationship tuples, and ask one question at
millisecond speed — "is this user related to this object in the
required way?" The CNCF relationship-based access control engine,
deployed from the official chart on a datastore you pick.

## Highlights

- **Engine here, data as resources** — stores, authorization models
  and relationship tuples are API-managed: the `OpenFgaStore` /
  `OpenFgaAuthorizationModel` / `OpenFgaRelationshipTuple` kinds
  compose against the exported endpoint.
- **Datastore composed by reference** — PostgreSQL recommended (a
  `KubernetesPostgres` pairs naturally), MySQL supported, memory as
  the honest evaluation arm; credentials ride Secrets and the
  connection URI renders credential-free.
- **Migrations that cannot deadlock** — `openfga migrate` as an
  idempotent init container in every pod, instead of the chart's
  hook Job that stalls rollout-waiting engines and dials the
  database during uninstall.
- **Security defaults upstream agrees with** — the demo playground
  is always off (the v1.18 server default), pre-shared keys ride
  Secrets in both arms, and OIDC enforces issuer plus audience.
- **Scales like the stateless service it is** — replicas scale
  checks linearly, the HPA arm hands the count to an autoscaler, and
  typed tuning covers query limits, deadlines and check-result
  caching.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
