# Locust

Declares one Locust load-testing cluster -- the open-source tool that simulates thousands of concurrent users against your own applications, with test behavior written as plain Python. The deliveryhero Helm chart (0.35.0, OCI-served from ghcr.io -- the classic index stalled in 2024 -- running the OFFICIAL locustio/locust image at 2.32.2) renders a master Deployment (the web UI and REST API on port 8089, worker coordination on 5557) and a worker fleet, each worker roughly ONE CPU CORE of load generation. Test scripts reach the pods as ConfigMaps: written inline in the spec (the module renders them and script changes ROLL THE PODS through a content-hash annotation) or referenced from ConfigMaps your CI already ships. SECURED BY DEFAULT: upstream ships the web UI OPEN -- anyone who can reach the Service can fire load at any host the cluster sees -- and this kind never deploys that: the login is ON from the first apply with a module-generated credential in the `<name>-auth` Secret, delivered through Locust's own extension seam, never as pod arguments. Uses a Kubernetes Provider Connection for cluster access.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes Namespace** -- created only when `create_namespace` is `true`; otherwise deploys into an existing namespace
- **Script ConfigMaps** (inline scripts) -- `<name>-locustfile` holding your locustfile as `main.py`, and `<name>-lib` holding supporting modules mounted at `lib/`; a content-hash annotation on the pod template rolls the pods when scripts change
- **The Helm release** -- the deliveryhero `locust` chart, rendering:
  - The master Deployment -- the web UI and REST API (port 8089) and the coordination endpoint workers dial (5557); exactly one master by design
  - The worker Deployment -- the load generators, at your fixed count (0 is a legal PAUSED fleet) or autoscaled
  - The master Service (ClusterIP by default) carrying the web and worker-connect port families
- **The web-login backend** (login on, non-headless) -- a module-owned ConfigMap delivering Locust's documented auth extension, plus the `<name>-auth` Secret holding the generated credential (keys `username`/`password`) -- ONE mounted source of truth; rotating the credential is a single-Secret edit
- **Autoscaling objects** (arm-dependent) -- an HPA scaling workers on CPU utilization, or a KEDA ScaledObject scaling on the live user count the master reports
- **PodDisruptionBudgets** (opt-in per role) -- `maxUnavailable 0`, so node drains never kill the master or workers mid-test

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with credentials for the target cluster. Map it as the default for your environment, or pick it explicitly in the wizard.

### Kubernetes Cluster

- **Internet or mirror access from the pods** (only when `pip_packages` / `pip_requirements_config_map` are set) -- pip installs at EVERY pod start; a PyPI outage becomes a pod-start failure. Bake a custom image for production gates.
- **KEDA operator** (only for the KEDA autoscaling arm) -- a KubernetesKeda composes naturally. NOTE the pairing the spec enforces: KEDA's DEFAULT trigger polls the master's `/stats/requests` API, which the web-UI login locks out and headless mode never serves -- provide custom triggers, or explicitly disable the login on a non-headless run.
- **Metrics server** (only for the HPA arm) -- plus worker CPU requests, or the HPA has no denominator.
- **Same-namespace credentials** -- Secrets named by `env_from_secrets` / `env_from_secret_keys` and existing script ConfigMaps are read by the pods at RUNTIME and must live in the install namespace.

## Deploy

### Console

Open the deployment store, find **Locust**, and click **Deploy**. The wizard walks the full run declaration: placement, the scripts (inline or existing ConfigMaps), the target, the test environment, dependencies, the master and worker fleet, autoscaling, sign-in, images, and the front door. Start from the **web-load-test** preset in the Presets tab for an interactive swarm, or **headless-ci** for a CI performance gate.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesLocust
metadata:
  name: load-test
  org: acme-corp
  env: dev
spec:
  namespace:
    value: load-test
  create_namespace: true
  load_test:
    inline:
      locustfile_content: |
        from locust import HttpUser, task, between

        class WebsiteUser(HttpUser):
            wait_time = between(1, 2)

            @task
            def get_index(self):
                self.client.get("/")
    target_host:
      value: http://my-app.my-namespace.svc.cluster.local:8080
  workers:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 256Mi
```

```bash
planton apply -f load-test.yaml
```

Everything else is the component's defaults -- and the defaults are deliberately SECURED: the web-UI login is ON with a module-generated credential (Secret `load-test-auth`, key `password`). Read it with:

```bash
kubectl get secret load-test-auth -n load-test \
  -o jsonpath='{.data.password}' | base64 -d
```

## Key Configuration

- **`load_test.scripts`** (required oneof) -- `inline` (locustfile + lib modules rendered into ConfigMaps; changes roll the pods) XOR `existing_config_maps` (your CI ships the scripts; same-namespace rule applies).
- **`load_test.target_host`** -- a literal URL or a reference to another resource's exported endpoint ("load-test what you declare"). Empty = the locustfile must set `host` itself; Locust refuses to start a test without one from either source.
- **`load_test.headless`** -- no web UI; the test starts as soon as the pods are up. Run shape rides `LOCUST_USERS` / `LOCUST_SPAWN_RATE` / `LOCUST_RUN_TIME` in `load_test.environment`. The login is moot and the web-UI outputs honestly read empty.
- **`load_test.name`** -- labels every resource and lands in the Deployments' IMMUTABLE selector labels: renaming is delete-and-recreate (cheap -- the cluster is stateless). Empty = the resource name.
- **`workers.replicas`** -- 0 or more; 0 is a PAUSED fleet. Ignored while an autoscaling arm owns the count. Size from the target: requests-per-second ÷ what one CPU core of your test generates.
- **`workers.autoscaling`** -- fixed count (no arm) XOR `hpa` (CPU-based; the 40% default target is deliberate -- workers must scale BEFORE saturating or they distort the load they generate) XOR `keda` (the live user count; `min_replicas: 0` is legal scale-to-zero; `custom_triggers` replaces the default trigger).
- **`web_ui_auth.enabled`** -- unset = ON with the generated credential (the open UI never ships); the explicit `false` is a recorded decision for fenced dev clusters and the KEDA default-trigger pairing. The login mechanism requires image tags >= 2.21.0 -- below the floor the chart would render credentials as literal pod arguments, which the module refuses.
- **`service.type`** -- ClusterIP by default: compose real exposure (Gateway API, KubernetesIngress) over the exported service handle; the annotations map carries the cloud-LB vocabulary when you flip to LoadBalancer.
- **`helm_values`** -- merged LAST over everything the spec renders; the module RE-PINS the login wiring and script delivery after the merge, so the security posture cannot be silently disabled from the hatch.

## Outputs and Dependencies

### What This Component Consumes

- **Namespace** -- a literal name or a KubernetesNamespace reference.
- **Target endpoint** (optional) -- any resource's exported endpoint via `load_test.target_host`.
- **Same-namespace Secrets/ConfigMaps** (optional) -- test credentials (`env_from_secrets`, `env_from_secret_keys`) and CI-shipped scripts, read at runtime.

### What This Component Provides

After deployment (or import), the stack outputs carry:

- **`namespace`** -- where the cluster runs
- **`master_service`** -- the Service name exposure kinds route to
- **`web_endpoint`** -- the in-cluster web UI / REST API URL (port 8089); pair it with the credential when the login is on
- **`master_bind_endpoint`** -- where additional workers (other namespaces, other clusters) register with the master (port 5557)
- **`web_ui_username`** / **`web_ui_password_secret`** -- the login identity and the Secret+key holding the generated password; both honestly EMPTY when the login is disabled or the run is headless
- **`port_forward_command`** -- one line to reach the web UI from a workstation when no exposure is composed

## Common Patterns

- **Interactive tuning swarm** -- the untouched defaults: inline scripts, the secured web UI, a small fixed fleet; drive runs from the browser and download reports before teardown.
- **CI performance gate** -- `headless: true` with run shape in `environment`; the pipeline applies the manifest, watches the master logs, and tears the namespace down.
- **Scale-to-zero standing swarm** -- the KEDA arm with `min_replicas: 0` and custom cron triggers: the fleet exists only during scheduled runs.
- **Isolated generator pool** -- worker `scheduling` onto a dedicated node pool so the load generator and the system under test never fight for CPU.

## Works With

- **KubernetesNamespace** -- the placement unit this kind installs into
- **KubernetesKeda** -- the operator the KEDA autoscaling arm requires
- **KubernetesNetworkPolicy** -- fence who can reach the master Service (defense in depth beside the login)
- **Gateway API kinds / KubernetesIngress** -- compose real exposure over `master_service` when engineers need the UI beyond the cluster
- **Any endpoint-exporting kind** -- the natural `target_host` reference: load-test the services you already declare
