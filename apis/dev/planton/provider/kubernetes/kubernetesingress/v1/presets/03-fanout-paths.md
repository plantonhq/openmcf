# Path Fan-Out

This preset routes one hostname to multiple Service backends by URL path: `/` to the frontend, `/api` to the API. The classic single-domain, multi-service layout — one certificate, one DNS record, several workloads behind it.

## When to Use

- A frontend and its API served under one domain (`app.example.com` and `app.example.com/api`)
- Consolidating several services behind one hostname instead of one subdomain each
- Incrementally splitting a monolith — carve paths out to new backends one at a time

## Key Configuration Choices

- **Longest match wins** — a request to `/api/users` matches both paths, and the more specific `/api` backend receives it; order in the list does not matter
- **Prefix matching is per path element** — `/api` matches `/api` and `/api/users` but NOT `/apiary`
- **Mixed port forms** — the frontend routes by number (`3000`), the API by name (`http`). Prefer names when the Service defines them; the reference survives port-number changes
- **No rewriting** — the backend receives the path as-is (`/api/users`, not `/users`); if the API expects stripped paths, add a controller rewrite annotation (e.g. `nginx.ingress.kubernetes.io/rewrite-target`)
- **All backends in one namespace** — every Service must live in the Ingress's namespace (`app`); a Kubernetes API constraint

## Values to Replace

| Value | Description |
|---|---|
| `app.example.com` | The shared public hostname |
| `app` | Namespace of the Ingress and all backend Services |
| `frontend-svc` / `3000` | Frontend Service and port number |
| `api-svc` / `http` | API Service and its named port |

## Related Presets

- **01-single-host** — one host, one backend
- **02-tls-cert-manager** — add HTTPS to this layout with a `tls` block and issuer annotation
