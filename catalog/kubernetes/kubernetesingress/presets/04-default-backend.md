# Default Backend Only

This preset declares no host or path rules at all — just a `default_backend`. Every request the controller routes to this Ingress, regardless of hostname or path, goes to one Service. It is the catch-all shape.

## When to Use

- Exposing a single Service on ALL traffic reaching the controller, with no host matching
- A cluster-wide fallback — a branded 404/maintenance page for requests no other Ingress claims
- Wildcard-DNS setups where one router Service handles every hostname itself

## Key Configuration Choices

- **`default_backend` instead of `rules`** — the spec requires at least one of the two (an Ingress with neither routes nothing and is rejected at validation); this preset takes the rules-free branch
- **Interaction with other Ingresses** — controllers merge all Ingresses they serve; host-matched rules elsewhere take precedence, and this backend receives what nothing else claims. When rules ARE present on an Ingress, most controllers also fall back to their own global default backend if `default_backend` is unset
- **`first_host` output stays empty** — there is no rule host to export; the load-balancer address outputs still populate once a controller reconciles the object
- **Same-namespace backend** — the Service must live in the Ingress's namespace (`web`)

## Values to Replace

| Value | Description |
|---|---|
| `web` | Namespace of the Ingress and the backend Service |
| `web-svc` | The Service receiving all unmatched traffic |
| `8080` | The Service port receiving traffic |
| `nginx` | Your cluster's IngressClass (`kubectl get ingressclass`) |

## Related Presets

- **01-single-host** — host-matched routing to one backend
- **03-fanout-paths** — path-based routing to several backends
