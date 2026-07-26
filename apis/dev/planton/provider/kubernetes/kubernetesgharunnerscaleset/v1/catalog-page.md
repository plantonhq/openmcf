# GHA Runner Scale Set

Self-hosted GitHub Actions runners that scale with your queue: declare
a fleet against a repository or organization, put its name in
`runs-on:`, and every queued job gets a fresh ephemeral runner pod on
your own cluster — your hardware, your network, your images.

## Highlights

- **Ephemeral by architecture** — one runner, one job, replaced; no
  state bleeds between builds.
- **Credentials done right** — PAT or GitHub App, always in a Secret;
  inline declarations are materialized by the module and never appear
  in rendered values.
- **Docker builds, two ways** — privileged dind for full Docker
  compatibility, or the unprivileged Kubernetes container hook for
  security-restricted clusters.
- **Scale bounds as configuration** — warm minimum against cold-start
  latency, hard maximum against runaway queues; excess jobs wait in
  GitHub, not in your cluster.
- **Enterprise plumbing included** — outbound proxies with
  authenticated credentials, private-CA trust for GitHub Enterprise
  Server, runner groups for repository access control.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
