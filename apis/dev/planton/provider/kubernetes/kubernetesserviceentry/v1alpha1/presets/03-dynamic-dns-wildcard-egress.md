# Dynamic-DNS Wildcard Egress

Registers a whole wildcard domain (`*.example-saas.com`) in the mesh's service
registry with `DYNAMIC_DNS` resolution: the proxy resolves the ACTUAL hostname
each request asks for, at request time — a dynamic forward proxy.

Use this when workloads reach many per-tenant or per-region subdomains of an
external SaaS (e.g. `eu1.example-saas.com`, `tenant42.example-saas.com`) and
enumerating them is impossible. With the mesh's outbound traffic policy set to
`REGISTRY_ONLY` (egress lockdown), this entry is what allows that traffic at
all — everything else stays blocked.

The mode's contract (enforced at validation, mirroring istiod):

- every host must be wildcarded (`*.` prefix),
- no `addresses` and no `endpoints` — destinations derive purely from the
  requested hostname,
- ports must be HTTP-family or TLS (HTTP, HTTPS via TLS, GRPC, HTTP2).

Replace `<namespace>` with the namespace that should own the entry (add
`export_to` to narrow visibility; default is mesh-wide).
