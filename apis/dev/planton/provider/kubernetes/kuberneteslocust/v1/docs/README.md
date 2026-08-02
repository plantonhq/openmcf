# KubernetesLocust — design notes

## Distribution

The deliveryhero Helm chart, chart 0.35.0 = Locust 2.32.2 (chart
Apache-2.0; the application MIT; the pod image is the OFFICIAL
`locustio/locust`). THE SERVING HOME IS OCI:
`oci://ghcr.io/deliveryhero/helm-charts/locust` — the classic index at
charts.deliveryhero.io stalls at 0.31.6 (2024) while ghcr.io serves
the live line; both engines pin the OCI repository. The module pins
`fullnameOverride` to the resource name, so child names are
deterministic: `<name>-master`, `<name>-worker`, the bare `<name>`
master Service, per-component ServiceAccounts.

## The web-UI login architecture

Locust's `--web-login` flag (2.21+) protects every web and REST route
behind a session but deliberately leaves the credential backend to
locustfile code — the documented extension seam. The chart's
`master.auth` username/password values feed ONLY the legacy pre-2.21
path that renders `--web-auth=user:password` — credentials as a
LITERAL POD ARGUMENT — so this kind never engages them:

- The module composes a small login backend (`planton_auth.py`) into
  the `<name>-web-auth` ConfigMap, mounted through the chart's
  `extraConfigMaps` seam, and appends it to the master's `-f`
  argument next to the locustfile (command-line `-f` overrides the
  `LOCUST_LOCUSTFILE` env; workers keep the env default and never
  load it).
- The credential and the Flask session-signing key live in the
  `<name>-auth` Secret, projected as FILES through the chart's
  `mount_external_secret` seam — never environment, rendered values
  or process arguments.
- The session key is generated once and stable — sessions survive pod
  restarts.
- Both engines fail loud when an image tag below 2.21 (or a
  non-numeric tag) is paired with the login.

## Script delivery truths (chart templates at the pin)

- The chart's bundled example ConfigMaps render ONLY while
  `loadtest.name` and the ConfigMap values keep their literal
  defaults (a `Files.Glob` keyed on the name) — a fragile coupling
  the module never engages: script ConfigMaps are always named
  explicitly, `""` disables the lib mount.
- The chart checksums only its OWN ConfigMaps into the pod templates
  — module-owned script content changes would roll NOTHING. The
  module stamps a content hash onto `master.annotations` /
  `worker.annotations` (pod-template surfaces the templates read),
  so script edits roll the pods deterministically, byte-identical
  across engines.
- The `/config` entrypoint pip-installs `pip_packages` and the
  requirements-file mount AT POD START — internet at pod start, a
  PyPI outage becomes a pod-start failure (taught on the fields).

## Env seams

`loadtest.environment_secret` renders a Secret FROM VALUES — the
forbidden literal-credential path: the module force-empties it after
the escape-hatch merge (Pulumi) / nulls it in the re-pin document
(Terraform, Helm null-deletion). The reference arms —
`environment_external_secret` (selected keys; env name = key name by
chart contract) and `environment_load_from_secrets` (whole Secrets)
— are the modeled surfaces.

## Autoscaling truths

- The chart's KEDA ScaledObject reuses `worker.hpa.minReplicas` /
  `maxReplicas` for its bounds while `hpa.enabled` stays false, and
  the worker Deployment still renders `replicas` on the KEDA arm —
  the module pins replicas to the KEDA floor so Helm upgrades reset
  scaling to the floor, not an unrelated count.
- The default KEDA trigger reads the live `user_count` from the
  master's own `/stats/requests` — an API the web-UI login locks out
  and headless mode never serves. The spec-level CEL forbids the
  default trigger unless the login is explicitly off on a
  non-headless run; custom triggers lift the constraint.

## Name budget

The longest derived child name is the module's own
`<name>-locustfile` ConfigMap (11-char suffix) → names ≤ 52 chars,
checked fail-loud in both engines.

## Pins

Chart 0.35.0 = Locust 2.32.2; image `locustio/locust:2.32.2`
(COMBINED repository+tag form, tag always pinned explicitly). Web/REST
port 8089; worker-connect 5557/5558.
