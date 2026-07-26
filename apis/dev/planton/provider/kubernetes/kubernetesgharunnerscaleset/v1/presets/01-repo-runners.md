# Repository runners preset

Self-hosted runners for one repository, scaled to zero when idle: a
queued job creates an ephemeral runner pod, the pod runs exactly that
job and is replaced. Workflows target the fleet by NAME — this
manifest registers as `build-runners`, so jobs say
`runs-on: build-runners` (names, not labels, are how scale sets route).

Requires a KubernetesGhaRunnerScaleSetController on the cluster first.
The credential here is a pre-created Secret reference — the
recommended posture; the PAT never appears in any manifest or rendered
value. A GitHub App (`auth.github_app` or App keys in the same Secret)
is the stronger production choice: fine-grained permissions and
expiring installation tokens.

Know what the default runner CANNOT do: `docker build` and container
jobs need a container mode (see the dind preset) — the plain runner
runs shell/tool steps only.

Change first: `min_runners: 1` if the per-job pod cold start bothers
developers; `runner.resources` sized for the JOBS, because builds
inherit those limits.

See [01-repo-runners.yaml](./01-repo-runners.yaml) for the manifest.
