# Lab preset

The smallest useful Ray cluster: a head that also RUNS tasks
(`scheduleTasksOnHead: true` — the deliberate lab inversion of the
production default, where the head advertises zero CPUs to keep
application work off the control plane), no worker groups, token auth
on by the catalog default.

PREREQUISITE: a `KubernetesKubeRayOperator` whose watch scope covers
this namespace (cluster-wide with its defaults).

Connect to it: notebooks and applications dial the CLIENT endpoint
(`ray://<name>-head-svc.<namespace>:10001`); jobs submit through the
DASHBOARD endpoint (port 8265, the Job Submission API). Both require
the bearer token from the operator-provisioned Secret named exactly
after this resource (key `auth_token`) — the exported credential
handle. STATE TRUTH: without GCS
fault tolerance, losing the head pod loses every job and actor
registration; a lab accepts that, production composes a
`KubernetesValkey` (see the production preset).

Change first: `rayVersion` (and keep any custom `image` version-locked
to it — the operator shapes its commands from the declared version and
a mismatch fails at runtime, not at apply).

See [01-lab.yaml](./01-lab.yaml) for the manifest.
