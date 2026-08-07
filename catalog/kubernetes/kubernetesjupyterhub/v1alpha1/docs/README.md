# KubernetesJupyterHub — research notes

Design substance for this kind comes from the official Zero to
JupyterHub chart (`jupyterhub` from https://hub.jupyter.org/helm-chart/,
chart 4.4.0 = JupyterHub 5.5.0; chart and application BSD-3-Clause).
The served chart is the truth for image tags — the repository's
vendored Chart.yaml and values carry chartpress build placeholders.

## Architecture at the pin

Four fixed tiers, all with CHART-FIXED bare resource names
(`fullnameOverride: ""`, module-re-pinned): the hub Deployment (
JupyterHub + KubeSpawner + the configured authenticator), the
configurable-http-proxy Deployment fronted by the `proxy-public`
Service (every user byte, login included, passes through CHP), the
user-scheduler (a kube-scheduler configured to PACK user pods via
MostAllocated scoring), and the pre-puller DaemonSets (hook + a
continuous one for autoscaled nodes). Bare names make the deployment a
per-namespace singleton — multi-instance separation is the namespace.

## Secret truths (why the module is shaped the way it is)

- The chart-owned `hub` Secret embeds the ENTIRE rendered values
  document (`values.yaml` key) — anything placed in Helm values is
  readable cluster state. Therefore no credential ever rides typed
  values: sign-in secrets reach the hub as env vars (extraEnv
  valueFrom) consumed by module-owned extraConfig python snippets.
- The chart's three internal auth materials (CHP auth token, cookie
  secret, CryptKeeper keys) generate LOOKUP-STABLE (reuse the existing
  Secret's values on upgrade) — deliberately left chart-managed;
  pushing module values through hub.config would land them readable in
  values, strictly worse.
- The external database password rides `hub.existingSecret`: the hub
  mounts it at /usr/local/etc/jupyterhub/existing-secret/ and exports
  PGPASSWORD/MYSQL_PWD from the `hub.db.password` key at startup
  (jupyterhub_config.py truth) — so `hub.db.url` renders
  CREDENTIAL-FREE.

## Authenticator surface

The hub image ships oauthenticator 17.4.0, jupyterhub-nativeauthenticator
1.3.0 and jupyterhub-ldapauthenticator 2.0.2 — the typed arms (dummy
with a required password, native, GitHub, Google, generic OIDC) are all
in-image; LDAP rides `helm_values`. JupyterHub 5 denies sign-in unless
an allow rule matches: the module renders `allowed_users` when a roster
is declared and `Authenticator.allow_all: true` explicitly otherwise.
Authenticator classes render as FULL class paths (entry-point
shortnames are registration details).

## Lifecycle truths

- KubeSpawner creates per-user pods and `claim-<username>` PVCs at
  RUNTIME — outside both engines' state. Uninstall removes the release
  resources; user home PVCs survive by Kubernetes semantics (nothing
  owns them but the namespace). The E2E destroy assertions treat
  surviving `claim-*` PVCs as DESIGNED and sweep them explicitly.
- The hub is single-replica by upstream design (Recreate strategy
  pinned in the chart; sqlite-pvc requires it, and JupyterHub itself
  has no HA story).
- The pre-puller hook runs INSIDE the install: a hook DaemonSet pulls
  the notebook image on every node and an awaiter Job gates the
  release — install duration is dominated by the image pull.

## Deliberate exclusions

- `proxy.https` / autohttps (Traefik + ACME) and the chart's
  ingress/httpRoute blocks: exposure and TLS compose from first-class
  kinds over the exported service handle.
- `hub.services`/`loadRoles`, per-component securityContext overrides,
  network-policy rule shaping, LDAP auth, `singleuser.extraFiles`:
  the `helm_values` escape hatch.
- `sqlite-memory` database type: hub state vanishing on every pod
  restart is not a shippable posture.

## Version pairing

Chart 4.4.0 = JupyterHub 5.5.0; hub/singleuser-sample/network-tools
images tag 4.4.0 (quay.io/jupyterhub), CHP 5.2.0, user-scheduler
registry.k8s.io/kube-scheduler v1.30.14, pause 3.10.1. kubeVersion
floor >= 1.28.
