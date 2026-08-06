# Organization Docker-builds preset

An organization-wide fleet that can build containers: dind mode runs a
privileged Docker daemon beside every runner, so `docker build`,
`docker run` and container-based actions work exactly as they do on
GitHub-hosted runners. The runner group fences which repositories may
use the fleet — govern access in GitHub, not by deploying more fleets.

Two idle runners stay warm (`min_runners: 2`) because organization
fleets serve many repositories and the per-job cold start compounds;
resources are sized for BUILDS (they inherit the runner's limits), not
for the runner agent.

Know the privilege trade: dind requires privileged pods. On clusters
where that is off the table, `container_mode.mode: kubernetes` runs
container jobs as separate unprivileged pods through the container
hook — at the cost of a per-runner work volume and slightly different
job semantics (see the chart's documentation for hook limitations).

Change first: a GitHub App credential if this manifest still points at
a PAT — an organization fleet is exactly where PAT scope and expiry
become liabilities.

See [02-org-docker-builds.yaml](./02-org-docker-builds.yaml) for the
manifest.
