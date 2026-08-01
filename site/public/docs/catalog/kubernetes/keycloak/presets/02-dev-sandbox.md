---
title: "Dev-sandbox preset"
description: "The smallest Keycloak that starts: the `dev-mem` embedded H2 database (in memory), the plain-HTTP listener, and strict hostname resolution off so the server answers on whatever host reaches it..."
type: "preset"
rank: "02"
presetSlug: "02-dev-sandbox"
componentSlug: "keycloak"
componentTitle: "Keycloak"
provider: "kubernetes"
icon: "package"
order: 2
---

# Dev-sandbox preset

The smallest Keycloak that starts: the `dev-mem` embedded H2 database
(in memory), the plain-HTTP listener, and strict hostname resolution
off so the server answers on whatever host reaches it (port-forwards,
cluster DNS). Port-forward `keycloak-dev-service` on 8080 and sign in
with the generated `keycloak-dev-initial-admin` Secret's credentials.

Everything about this shape is deliberately disposable: dev-mem data
DIES WITH THE POD (realms, users, clients — everything), the spec caps
the sandbox vendors at a single instance (each pod would hold its own
divergent world), and plain HTTP belongs behind a TLS boundary you
trust. The spec's validation rules force real installs to make real
choices — this preset exists so evaluating the component never
requires them.

PREREQUISITE, same as always: a `KubernetesKeycloakOperator` watching
this namespace.

Change first: nothing. When the evaluation sticks, move to the
standard preset — a real database, TLS, and a hostname are the
decisions that make it production.

See [02-dev-sandbox.yaml](./02-dev-sandbox.yaml) for the manifest.
