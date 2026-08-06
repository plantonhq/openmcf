# Kubernetes JupyterHub

Deploys JupyterHub — the multi-user notebook platform — from the
official Zero to JupyterHub Helm chart. Every user who signs in gets
their OWN JupyterLab server pod and (by default) their own persistent
home volume; the hub spawns them, the proxy routes to them, the idle
culler stops the abandoned ones.

## When NOT to Use This

A single person who just wants a notebook does not need a hub — run a
notebook image as a KubernetesDeployment. This kind is for TEAMS:
classes, research groups, data-science orgs — anywhere "everyone gets
their own isolated notebook with their own storage" is the product.

## Secured by default

The chart's own default authenticator accepts ANY username with NO
password. That never ships from here: an empty `authentication` block
means shared-password sign-in with a module-generated password
(`<name>-auth` Secret, exported as the credential handle). Real
identity comes from the typed arms — GitHub, Google, any OIDC provider
(a KubernetesKeycloak realm slots straight in), or self-service
accounts with admin approval. OAuth client secrets reach the hub as
environment variables from your Secrets — never through rendered chart
values, which this chart embeds readably inside its own hub Secret.

## Hub state

The hub database records users, running servers and tokens. The default
is sqlite on a small PVC — right for most installs, because the hub is
single-replica by design (upstream has no HA story). The postgres arm
composes a KubernetesPostgres: the connection URL renders
credential-free and the password rides a mounted Secret.

## User pods and volumes are RUNTIME artifacts

KubeSpawner creates per-user pods and `claim-<username>` PVCs when
users sign in — they are the hub's children, not this resource's.
Destroying this resource removes the hub, proxy and scheduling
machinery; running user pods die with the proxy, but USER HOME VOLUMES
SURVIVE until you delete them (or the namespace) explicitly. User homes
are data — plan their backup and their cleanup deliberately.

## Capacity machinery

The user-scheduler PACKS user pods onto the busiest fitting node so
cluster autoscalers can actually reclaim the rest; placeholder pods
keep warm capacity a real user can evict instantly; the pre-puller
pushes the notebook image onto every node before the install completes
(first spawns are fast because the install already paid for the pull).
The spawn menu (`single_user.profiles`) lets users pick their machine
size at login.

## Exposure

The proxy-public Service is the ONLY front door — all user traffic,
login included, enters there. This kind keeps it ClusterIP (the chart's
own default is LoadBalancer — deliberately overridden): compose
exposure from first-class kinds (KubernetesService, Gateway API routes)
over the exported `proxy_public_service` handle. Names are chart-fixed
(hub, proxy-public…), so one JupyterHub per namespace.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
