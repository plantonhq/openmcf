---
title: "Standard preset"
description: "The production Keycloak shape: two instances clustering through the operator's discovery Service, a PostgreSQL database referenced from a `KubernetesPostgres` resource (host = its read-write Service,..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "keycloak"
componentTitle: "Keycloak"
provider: "kubernetes"
icon: "package"
order: 1
---

# Standard preset

The production Keycloak shape: two instances clustering through the
operator's discovery Service, a PostgreSQL database referenced from a
`KubernetesPostgres` resource (host = its read-write Service,
credentials = Secret selectors against the operator-maintained
`<cluster>-app` Secret — nothing credential-bearing rides the CR), TLS
served from a `KubernetesCertificate` reference, and a declared public
hostname with strict resolution (the server default — what tokens,
redirects and the OIDC discovery document advertise).

PREREQUISITE: a `KubernetesKeycloakOperator` watching this namespace
(under its default namespaced watch, the operator and this resource
live in the SAME namespace). Declare the `keycloak` database at the
Postgres resource's bootstrap (`initdb`), and co-locate Keycloak with
it — the credential secretKeyRefs read only their own namespace.

First login: the operator generates the one-time
`keycloak-initial-admin` Secret (username `temp-admin`) — exported as
this resource's credential handle. It seeds the FIRST admin only;
create your real admin inside Keycloak and treat the bootstrap Secret
as spent. Exposure composes from Gateway API kinds referencing the
exported service handles — the operator's own Ingress is always
disabled by this component.

Know the update truth: with the default strategy, changing the IMAGE
takes a full scale-to-zero recreate (two Keycloak versions cannot
share one cache cluster/schema) — an outage window by design.

Change first: the `my-postgres` and `keycloak-cert` references and the
hostname; then sizing (Keycloak is a JVM — production typically runs
1–2Gi).

See [01-standard.yaml](./01-standard.yaml) for the manifest.
