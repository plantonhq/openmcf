# KubernetesLocust

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesLocustSpec** deploys Locust — the open-source load
testing tool that simulates user traffic against your own
applications, with test behavior written as plain Python
(https://locust.io).

WHAT GETS INSTALLED: the locust Helm chart (deliveryhero org,
OCI-served; runs the OFFICIAL locustio/locust image) renders a
master Deployment (the web UI + REST API on port 8089 and the
coordination endpoint workers dial on port 5557) and a worker
Deployment (the load generators — each worker runs your locustfile
and reports stats to the master). Your test scripts reach the pods
as ConfigMap mounts; changing the scripts rolls the pods
automatically.

SECURED BY DEFAULT: upstream ships the web UI OPEN — anyone who
can reach the Service can start load tests against any host it
can see. This kind enables the web-UI login by default with a
module-generated credential (`<name>-auth` Secret): the module
delivers a small login backend alongside your locustfile (Locust's
own documented extension seam — the UI and every REST route
require the session). Disable only for headless runs or fenced
dev clusters.

TARGET WHAT YOU DECLARE: `load_test.target_host` accepts a literal
URL or a reference to another resource's output — point it at a
service you already declare on the platform and the wiring
composes.

EXPOSURE: the Service stays ClusterIP; expose it via first-class
kinds over the exported handle.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesLocust
metadata:
  name: hack-locust
spec:
  namespace:
    value: hack-locust
  create_namespace: true
  load_test:
    inline:
      locustfile_content: |
        from locust import HttpUser, task, between

        from lib import pages


        class WebsiteUser(HttpUser):
            wait_time = between(1, 2)

            @task
            def get_index(self):
                self.client.get("/")

            @task(3)
            def get_random_page(self):
                self.client.get(pages.choose_random_page())
      lib_files:
        "__init__.py": ""
        "pages.py": |
          import random


          def choose_random_page():
              return random.choice(["/about/", "/contact/", "/search/"])
    target_host:
      value: http://hack-target.hack-locust.svc.cluster.local:9898
  workers:
    replicas: 2
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
      limits:
        cpu: "1"
        memory: 512Mi
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.createNamespace` | `bool` |  |  |  |
| `spec.image` | `KubernetesLocustImage` |  |  |  |
| `spec.image.repository` | `string` |  | `locustio/locust` |  |
| `spec.image.tag` | `string` |  | `2.32.2` |  |
| `spec.imagePullSecrets` | `[]string` |  |  |  |
| `spec.loadTest` | `KubernetesLocustLoadTest` | yes |  |  |
| `spec.loadTest.name` | `string` |  |  |  |
| `spec.loadTest.inline` | `KubernetesLocustInlineScripts` |  |  |  |
| `spec.loadTest.inline.locustfileContent` | `string` | yes |  |  |
| `spec.loadTest.inline.libFiles` | `map<string, string>` |  |  |  |
| `spec.loadTest.existingConfigMaps` | `KubernetesLocustExistingScriptConfigMaps` |  |  |  |
| `spec.loadTest.existingConfigMaps.locustfileConfigMap` | `string` | yes |  |  |
| `spec.loadTest.existingConfigMaps.locustfileName` | `string` |  | `main.py` |  |
| `spec.loadTest.existingConfigMaps.libConfigMap` | `string` |  |  |  |
| `spec.loadTest.targetHost` | `string \| valueFrom` |  |  |  |
| `spec.loadTest.pipPackages` | `[]string` |  |  |  |
| `spec.loadTest.pipRequirementsConfigMap` | `string` |  |  |  |
| `spec.loadTest.environment` | `map<string, string>` |  |  |  |
| `spec.loadTest.envFromSecrets` | `[]string` |  |  |  |
| `spec.loadTest.envFromSecretKeys` | `[]KubernetesLocustSecretEnv` |  |  |  |
| `spec.loadTest.envFromSecretKeys[].secretName` | `string` | yes |  |  |
| `spec.loadTest.envFromSecretKeys[].keys` | `[]string` | yes |  |  |
| `spec.loadTest.tags` | `[]string` |  |  |  |
| `spec.loadTest.excludeTags` | `[]string` |  |  |  |
| `spec.loadTest.headless` | `bool` |  |  |  |
| `spec.master` | `KubernetesLocustMaster` |  |  |  |
| `spec.master.resources` | `ContainerResources` |  |  |  |
| `spec.master.resources.limits` | `CpuMemory` |  |  |  |
| `spec.master.resources.limits.cpu` | `string` |  |  |  |
| `spec.master.resources.limits.memory` | `string` |  |  |  |
| `spec.master.resources.requests` | `CpuMemory` |  |  |  |
| `spec.master.resources.requests.cpu` | `string` |  |  |  |
| `spec.master.resources.requests.memory` | `string` |  |  |  |
| `spec.master.logLevel` | `string` |  | `INFO` |  |
| `spec.master.pdbEnabled` | `bool` |  |  |  |
| `spec.master.scheduling` | `KubernetesLocustScheduling` |  |  |  |
| `spec.master.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.master.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.master.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.master.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.master.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.master.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.master.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.workers` | `KubernetesLocustWorkers` |  |  |  |
| `spec.workers.replicas` | `int32` |  | `1` |  |
| `spec.workers.resources` | `ContainerResources` |  |  |  |
| `spec.workers.resources.limits` | `CpuMemory` |  |  |  |
| `spec.workers.resources.limits.cpu` | `string` |  |  |  |
| `spec.workers.resources.limits.memory` | `string` |  |  |  |
| `spec.workers.resources.requests` | `CpuMemory` |  |  |  |
| `spec.workers.resources.requests.cpu` | `string` |  |  |  |
| `spec.workers.resources.requests.memory` | `string` |  |  |  |
| `spec.workers.logLevel` | `string` |  | `INFO` |  |
| `spec.workers.pdbEnabled` | `bool` |  |  |  |
| `spec.workers.scheduling` | `KubernetesLocustScheduling` |  |  |  |
| `spec.workers.scheduling.nodeSelector` | `map<string, string>` |  |  |  |
| `spec.workers.scheduling.tolerations` | `[]WorkloadToleration` |  |  |  |
| `spec.workers.scheduling.tolerations[].key` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].operator` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].value` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].effect` | `string` |  |  |  |
| `spec.workers.scheduling.tolerations[].tolerationSeconds` | `int64` |  |  |  |
| `spec.workers.hpa` | `KubernetesLocustWorkerHpa` |  |  |  |
| `spec.workers.hpa.minReplicas` | `int32` |  | `1` |  |
| `spec.workers.hpa.maxReplicas` | `int32` |  |  |  |
| `spec.workers.hpa.targetCpuUtilizationPercent` | `int32` |  | `40` |  |
| `spec.workers.keda` | `KubernetesLocustWorkerKeda` |  |  |  |
| `spec.workers.keda.minReplicas` | `int32` |  | `1` |  |
| `spec.workers.keda.maxReplicas` | `int32` |  |  |  |
| `spec.workers.keda.targetUsersPerWorker` | `int32` |  | `50` |  |
| `spec.workers.keda.pollingIntervalSeconds` | `int32` |  |  |  |
| `spec.workers.keda.cooldownPeriodSeconds` | `int32` |  |  |  |
| `spec.workers.keda.customTriggers` | `string` |  |  |  |
| `spec.webUiAuth` | `KubernetesLocustWebUiAuth` |  |  |  |
| `spec.webUiAuth.enabled` | `bool` |  | `true` |  |
| `spec.webUiAuth.username` | `string` |  | `locust` |  |
| `spec.service` | `KubernetesLocustService` |  |  |  |
| `spec.service.type` | `string` |  | `ClusterIP` |  |
| `spec.service.annotations` | `map<string, string>` |  |  |  |
| `spec.helmValues` | `string` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace to install into. Accepts a literal namespace name or a
reference to a KubernetesNamespace resource.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.createNamespace

`bool`

When true, the namespace is created (with the standard Planton
governance labels) before installing and deleted with the
resource. When false, the namespace must already exist. KNOW
THIS: Secrets referenced by `env_from_secrets` /
`env_from_secret_keys` and existing script ConfigMaps are read
by the pods at runtime and must live in the SAME namespace —
co-locate Locust with the credentials its tests use.

### spec.image

`KubernetesLocustImage`

Container image for master and workers. Empty = the official
`locustio/locust` image at the Locust release this kind is
built against.

### spec.image.repository

`string` · optional (explicit presence)

Image repository (including registry host for mirrors, e.g.
"my-registry.example.com/locustio/locust"). Empty =
"locustio/locust" (Docker Hub).

- default: `locustio/locust`

### spec.image.tag

`string` · optional (explicit presence)

Image tag. Empty = "2.32.2" (the Locust release this kind is
built against). KNOW THIS: the web-UI login mechanism requires
Locust >= 2.21.0 — older tags fall onto a chart code path that
renders credentials as pod arguments, which this kind refuses.

- default: `2.32.2`

### spec.imagePullSecrets

`[]string`

Names of image-pull Secrets in the same namespace, for pulling
from private mirrors.

### spec.loadTest

`KubernetesLocustLoadTest` · required

The load test — scripts, target, packages, environment. This is
the reason the deployment exists.

- rule: {"required":true}

### spec.loadTest.name

`string` · optional (explicit presence)

Load-test name — labels every resource (`load_test: <name>`)
and groups stats. Empty = the resource name. KNOW THIS: the
name lands in the Deployments' immutable selector labels — it
cannot change after the first deploy (delete and recreate to
rename).

- rule: {"string":{"pattern":"^[a-z0-9]([a-z0-9-]*[a-z0-9])?$"}}

### spec.loadTest.inline

`KubernetesLocustInlineScripts`

Scripts written inline in this manifest. The module renders
them into ConfigMaps (`<name>-locustfile`, `<name>-lib`) and
mounts them where Locust looks; script changes roll the pods
automatically.

### spec.loadTest.inline.locustfileContent

`string` · required

The locustfile — the Python that defines your simulated users
and their tasks. Rendered into the `<name>-locustfile`
ConfigMap as `main.py`.

- rule: {"required":true}

### spec.loadTest.inline.libFiles

`map<string, string>`

Supporting Python modules: filename → content. Rendered into
the `<name>-lib` ConfigMap and mounted at `lib/` next to the
locustfile — import them as `from lib import <module>`.

- rule: {"map":{"keys":{"string":{"pattern":"^[A-Za-z0-9_.-]+$"}}}}

### spec.loadTest.existingConfigMaps

`KubernetesLocustExistingScriptConfigMaps`

Scripts you already ship as ConfigMaps (e.g. synced from your
repo by CI). Same-namespace constraint applies.

### spec.loadTest.existingConfigMaps.locustfileConfigMap

`string` · required

ConfigMap holding the locustfile. Same-namespace constraint
applies.

- rule: {"required":true}

### spec.loadTest.existingConfigMaps.locustfileName

`string` · optional (explicit presence)

The locustfile's key/filename inside that ConfigMap. Empty =
"main.py".

- default: `main.py`
- rule: {"string":{"pattern":"^[A-Za-z0-9_.-]+\\.py$"}}

### spec.loadTest.existingConfigMaps.libConfigMap

`string`

ConfigMap holding supporting modules, mounted at `lib/` next to
the locustfile. Empty = no lib mount.

### spec.loadTest.targetHost

`string | valueFrom`

The base URL the test targets (Locust's `--host` — what
`self.client.get("/")` resolves against). Accepts a literal URL
or a reference to another resource's output — point it at a
service you declare on the platform. Empty = the locustfile
must set `host` itself (Locust refuses to start a test without
one).

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.loadTest.pipPackages

`[]string`

Extra Python packages pip-installed AT POD START (both master
and workers). KNOW THIS: this downloads from PyPI every time a
pod starts — pods need internet (or a mirror via pip env
configuration in `environment`), and a PyPI outage becomes a
pod-start failure. For air-gapped or production-grade setups,
bake packages into a custom image instead.

- rule: {"repeated":{"items":{"string":{"pattern":"^[^\\s]+$"}}}}

### spec.loadTest.pipRequirementsConfigMap

`string`

An existing ConfigMap holding a `requirements.txt` — the
file-based alternative to `pip_packages`, same
pod-start-install caveat. Same-namespace constraint applies.

### spec.loadTest.environment

`map<string, string>`

Environment variables for the test (master AND workers) — feed
configuration to your locustfile (`os.environ[...]`). Plain
values only; for credentials use `env_from_secrets` /
`env_from_secret_keys`.

### spec.loadTest.envFromSecrets

`[]string`

Existing Secrets loaded WHOLE into the test environment (every
key becomes an environment variable named after it) — the
simple way to hand test credentials to your locustfile.
Same-namespace constraint applies.

### spec.loadTest.envFromSecretKeys

`[]KubernetesLocustSecretEnv`

Selected keys from existing Secrets injected into the test
environment — each key becomes an environment variable named
exactly after the key (the chart's contract). Use when a Secret
carries more keys than the test should see.

### spec.loadTest.envFromSecretKeys[].secretName

`string` · required

Name of the Secret. Same-namespace constraint applies.

- rule: {"required":true}

### spec.loadTest.envFromSecretKeys[].keys

`[]string` · required

Key NAMES within the Secret (references, not secret material) to
inject — each becomes an environment variable with the SAME name
(name your Secret keys like environment variables, e.g.
"API_TOKEN").

- rule: {"repeated":{"minItems":"1"}}

### spec.loadTest.tags

`[]string`

Run only tasks tagged with any of these tags (Locust
`--tags`). Empty = all tasks.

- rule: {"repeated":{"items":{"string":{"pattern":"^[^\\s]+$"}}}}

### spec.loadTest.excludeTags

`[]string`

Skip tasks tagged with any of these tags (Locust
`--exclude-tags`).

- rule: {"repeated":{"items":{"string":{"pattern":"^[^\\s]+$"}}}}

### spec.loadTest.headless

`bool`

Run headless: no web UI — the test starts immediately and runs
until stopped. Drive run shape (users, spawn rate, duration)
through `environment` (LOCUST_USERS, LOCUST_SPAWN_RATE,
LOCUST_RUN_TIME). With headless on, `web_ui_auth` has nothing
to protect and is ignored.

### spec.master

`KubernetesLocustMaster`

The master: web UI, REST API, worker coordination. Empty = chart
defaults (no resource requests — size for real load tests).

### spec.master.resources

`ContainerResources`

CPU/memory for the master container. Empty = no requests (fine
for trying it out; size for real tests — the master aggregates
stats from every worker).

### spec.master.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.master.resources.limits.cpu

`string`

### spec.master.resources.limits.memory

`string`

### spec.master.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.master.resources.requests.cpu

`string`

### spec.master.resources.requests.memory

`string`

### spec.master.logLevel

`string` · optional (explicit presence)

Log level. Empty = "INFO".

- default: `INFO`
- rule: {"string":{"in":["DEBUG","INFO","WARNING","ERROR","CRITICAL"]}}

### spec.master.pdbEnabled

`bool`

Create a PodDisruptionBudget (maxUnavailable 0) so voluntary
disruptions (node drains) never kill the master mid-test.

### spec.master.scheduling

`KubernetesLocustScheduling`

Pod scheduling for the master.

### spec.master.scheduling.nodeSelector

`map<string, string>`

Node selector.

### spec.master.scheduling.tolerations

`[]WorkloadToleration`

Tolerations.

### spec.master.scheduling.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.master.scheduling.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.master.scheduling.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.master.scheduling.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.master.scheduling.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.workers

`KubernetesLocustWorkers`

The worker fleet: the actual load generators. Empty = 1 worker
with chart defaults. Worker COUNT bounds the load you can
generate — each worker is one Python process (roughly one CPU
core of load generation).

### spec.workers.replicas

`int32` · optional (explicit presence)

Worker count. Empty = 1. Ignored when an autoscaling arm is
set.

- default: `1`
- rule: {"int32":{"gte":0}}

### spec.workers.resources

`ContainerResources`

CPU/memory per worker container. Empty = no requests. KNOW
THIS: a worker saturates around one CPU core — more load per
worker needs more CPU, more concurrent users need more memory;
HPA additionally requires CPU REQUESTS to compute utilization
against.

### spec.workers.resources.limits

`CpuMemory`

The resource limits for the container.
Specify the maximum amount of CPU and memory that the container can use.

### spec.workers.resources.limits.cpu

`string`

### spec.workers.resources.limits.memory

`string`

### spec.workers.resources.requests

`CpuMemory`

The resource requests for the container.
Specify the minimum amount of CPU and memory that the container is guaranteed.

### spec.workers.resources.requests.cpu

`string`

### spec.workers.resources.requests.memory

`string`

### spec.workers.logLevel

`string` · optional (explicit presence)

Log level. Empty = "INFO".

- default: `INFO`
- rule: {"string":{"in":["DEBUG","INFO","WARNING","ERROR","CRITICAL"]}}

### spec.workers.pdbEnabled

`bool`

Create a PodDisruptionBudget (maxUnavailable 0) for the
workers.

### spec.workers.scheduling

`KubernetesLocustScheduling`

Pod scheduling for workers.

### spec.workers.scheduling.nodeSelector

`map<string, string>`

Node selector.

### spec.workers.scheduling.tolerations

`[]WorkloadToleration`

Tolerations.

### spec.workers.scheduling.tolerations[].key

`string`

Taint key to tolerate. Empty key with operator "Exists" tolerates every taint.

### spec.workers.scheduling.tolerations[].operator

`string`

How key/value match: "Equal" (default — value must match too) or "Exists"
(key presence alone matches).

- rule: Toleration operator must be either "Equal" or "Exists"

### spec.workers.scheduling.tolerations[].value

`string`

Taint value to match when operator is "Equal".

### spec.workers.scheduling.tolerations[].effect

`string`

Which taint effect is tolerated: "NoSchedule", "PreferNoSchedule", or
"NoExecute". Empty tolerates all effects for the key.

- rule: Toleration effect must be one of "NoSchedule", "PreferNoSchedule", or "NoExecute"

### spec.workers.scheduling.tolerations[].tolerationSeconds

`int64` · optional (explicit presence)

For "NoExecute" taints only: how many seconds already-running pods stay bound
after the taint appears. Unset means tolerate forever.

### spec.workers.hpa

`KubernetesLocustWorkerHpa`

Horizontal Pod Autoscaler on CPU utilization (needs a metrics
server on the cluster and worker CPU REQUESTS).

### spec.workers.hpa.minReplicas

`int32` · optional (explicit presence)

Lower replica bound. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.workers.hpa.maxReplicas

`int32`

Upper replica bound. Required.

- rule: {"int32":{"gte":1}}

### spec.workers.hpa.targetCpuUtilizationPercent

`int32` · optional (explicit presence)

Target average CPU utilization percentage (of worker CPU
requests). Empty = 40 (the chart default — workers should scale
BEFORE saturating, or they distort the load they generate).

- default: `40`
- rule: {"int32":{"lte":100,"gte":1}}

### spec.workers.keda

`KubernetesLocustWorkerKeda`

KEDA event-driven autoscaling — scales workers on the LIVE
USER COUNT the master reports (more simulated users = more
workers). Requires the KEDA operator on the cluster — a
KubernetesKeda composes naturally.

### spec.workers.keda.minReplicas

`int32` · optional (explicit presence)

Lower replica bound. Empty = 1.

- default: `1`
- rule: {"int32":{"gte":0}}

### spec.workers.keda.maxReplicas

`int32`

Upper replica bound. Required.

- rule: {"int32":{"gte":1}}

### spec.workers.keda.targetUsersPerWorker

`int32` · optional (explicit presence)

Simulated users one worker should carry before another is added
— the scaling target on the master's live user count. Empty =
50 (the chart default).

- default: `50`
- rule: {"int32":{"gte":1}}

### spec.workers.keda.pollingIntervalSeconds

`int32` · optional (explicit presence)

Seconds between trigger evaluations. Empty = 15 (the chart
default).

- rule: {"int32":{"gte":1}}

### spec.workers.keda.cooldownPeriodSeconds

`int32` · optional (explicit presence)

Seconds to wait after the last active trigger before scaling
down. Empty = 30 (the chart default).

- rule: {"int32":{"gte":0}}

### spec.workers.keda.customTriggers

`string`

Replace the default user-count trigger with your own KEDA
trigger list as raw YAML (the `triggers:` array content). Empty
= the user-count trigger described above.

### spec.webUiAuth

`KubernetesLocustWebUiAuth`

Web-UI login. Empty = ENABLED with a module-generated credential
— the open, anyone-can-start-load-tests UI never ships. Has no
effect when `load_test.headless` is true (headless runs start no
web UI at all).

### spec.webUiAuth.enabled

`bool` · optional (explicit presence)

Require login on the web UI and REST API. Empty = true — the
open UI never ships by default. Disabling means anyone who can
reach the Service can start load tests against any host the
cluster can see.

- default: `true`

### spec.webUiAuth.username

`string` · optional (explicit presence)

The login username. Empty = "locust". The module generates this
user's password into the `<name>-auth` Secret (key `password`)
— exported as the credential handle for operators and
verifiers.

- default: `locust`
- rule: {"string":{"pattern":"^[a-z0-9][a-z0-9._-]*$"}}

### spec.service

`KubernetesLocustService`

The master Service (web UI 8089 + worker-connect 5557/5558) —
this kind keeps it ClusterIP (compose exposure kinds over the
exported handle); the annotations surface carries cloud LB
configuration when you flip the type.

### spec.service.type

`string` · optional (explicit presence)

Service type. Empty = "ClusterIP" — compose exposure kinds over
the exported handle instead of exposing directly.

- default: `ClusterIP`
- rule: {"string":{"in":["ClusterIP","LoadBalancer","NodePort"]}}

### spec.service.annotations

`map<string, string>`

Service annotations — the cloud load-balancer configuration
surface when `type` is LoadBalancer.

### spec.helmValues

`string`

Raw Helm values (YAML) merged LAST over everything this spec
renders — the escape hatch for chart surfaces the typed fields
do not model (probes, init containers, extra volumes, host
aliases). The module re-pins security-critical values after the
merge (the login wiring, the script delivery) — those cannot be
silently disabled from here.

## Validation Rules

- `workers.keda.default_trigger_needs_open_stats`: The default KEDA trigger scales on the master's /stats/requests API, which the web-UI login locks out and headless mode never serves — set web_ui_auth.enabled to false with a non-headless run, or provide keda.custom_triggers.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesLocust, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace Locust runs in. |
| `status.outputs.master_service` | `string` | Name of the master Service (`<name>` at default naming) — the handle exposure kinds route to. Serves the web UI (8089) and the worker-connect ports (5557/5558). |
| `status.outputs.web_endpoint` | `string` | In-cluster web endpoint, `http://<master_service>.<namespace>.svc.cluster.local:8089` — the web UI and REST API (with login on, pair it with the credential below). |
| `status.outputs.master_bind_endpoint` | `string` | In-cluster worker-connect endpoint, `<master_service>.<namespace>.svc.cluster.local:5557` — where additional workers (e.g. from other namespaces) register with the master. |
| `status.outputs.web_ui_username` | `string` | The web-UI login username (`web_ui_auth.username`, default "locust"). Empty when the login is disabled or the run is headless. |
| `status.outputs.web_ui_password_secret` | `KubernetesSecretKey` | The web-UI credential: the Secret and key holding the login password (module-generated `<name>-auth`, key `password`). Empty when the login is disabled or the run is headless. |
| `status.outputs.web_ui_password_secret.name` | `string` | The name of the Kubernetes Secret. |
| `status.outputs.web_ui_password_secret.key` | `string` | The key within the Kubernetes Secret. |
| `status.outputs.port_forward_command` | `string` | Port-forward command for reaching the web UI from a workstation when no exposure is composed. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
