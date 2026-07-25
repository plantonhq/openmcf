# Argo CD

Declarative GitOps for Kubernetes. Argo CD watches your Git
repositories and keeps clusters converged on what they declare — drift
is detected, shown, and (with automated sync) corrected. The de-facto
standard delivery engine for the Kubernetes ecosystem.

## Highlights

- **The full control plane, typed** — controller, API/UI server, repo
  server, ApplicationSet controller, notifications, dex and the
  manifest-hydration commit server, each with its own knobs.
- **Secrets stay secrets** — SSO client secrets and repo credentials
  ride Argo CD's own runtime Secret indirection; nothing sensitive
  ever lands in rendered values.
- **SSO + RBAC first-class** — direct OIDC (PKCE or secret-backed) or
  dex connectors, with Argo CD's CSV policy language typed in.
- **Cache by need** — bundled Redis for most installs, the Sentinel HA
  subchart for platforms, or an external endpoint (a `KubernetesValkey`
  pairs by reference).
- **Safe by default** — CRDs are kept on uninstall (deleting them would
  cascade to every Application), and the generated admin credential is
  exported as a Secret handle, never printed.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
