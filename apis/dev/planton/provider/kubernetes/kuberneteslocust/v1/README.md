# Kubernetes Locust

Deploys Locust — the open-source load testing tool — from the
deliveryhero Helm chart (the official `locustio/locust` image). A
master serves the web UI and coordinates; a worker fleet generates the
load; your test behavior is plain Python that ships as ConfigMaps —
no custom images to build.

## Secured by default

Upstream Locust ships the web UI OPEN — anyone who can reach the
Service can start load tests against any host the cluster can see.
That never ships from here: an empty `web_ui_auth` block means the
login is ON with a module-generated credential (`<name>-auth` Secret,
exported as the credential handle). The login backend is
platform-managed code delivered alongside your locustfile through
Locust's own extension seam — the chart's legacy path that renders
credentials as pod arguments is never engaged. Headless runs start no
web UI, so the login machinery is skipped there.

## The load test is the product

Write the locustfile inline (`load_test.inline`) — supporting modules
ride `lib_files` and mount at `lib/` — or reference ConfigMaps your
CI already ships (`load_test.existing_config_maps`). Point
`target_host` at a literal URL or another resource's exported
endpoint: load-testing the services you already declare is the
composed story. Script changes roll the pods automatically (the
module stamps a content hash onto the pod templates).

Test credentials arrive as environment variables your Python reads:
`env_from_secrets` loads whole Secrets, `env_from_secret_keys` picks
keys. KNOW THIS: `pip_packages` installs from PyPI at every pod start
— convenient for extras, but production-grade setups bake packages
into a custom image instead.

## Sizing and scaling

Each worker is one Python process — roughly one CPU core of load
generation; size `workers.replicas` and resources to the load you
need, and give workers CPU requests before enabling HPA. KEDA scales
workers on the LIVE USER COUNT the master reports (its default
trigger reads the master's stats API, which requires the login to be
off — the validation enforces the pairing); custom triggers replace
it wholesale.

## Exposure

The master Service stays ClusterIP; compose exposure kinds over the
exported `master_service` handle, or use the exported
`port_forward_command` from a workstation. Additional workers
(including from other namespaces) can dial the exported
`master_bind_endpoint`.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
