# KubernetesKeycloak Guide

The judgment this guide carries: two configuration choices decide whether
Keycloak even starts and whether logins work — both surface upstream only
as a CrashLoopBackOff or "login randomly breaks," and both are decisions a
proposal must make deliberately, not leave to defaults.

## TLS-or-HTTP is a start-or-crash decision

The server refuses to start with neither TLS configured nor HTTP
explicitly enabled — upstream shows this only as a crash loop. This spec
turns it into a validation rule (the field docs on
[reference.md](reference.md)), so make the call in the manifest: a TLS
secret for direct HTTPS, or `http.httpEnabled` when a proxy terminates
TLS in front. An architecture that leaves both unset is proposing a pod
that will never come up.

## Strict hostname and the CORS trap

With strict hostname resolution (the server default) a hostname is
mandatory. Disabling strict is only correct behind a reverse proxy that
rewrites Host headers — and then `proxyHeaders` must be set, or browsers
fail logins with computed-origin CORS 403s (the classic intermittent-login
misconfiguration). If the architecture fronts Keycloak with a gateway,
this is part of the proposal, not an afterthought.

## Operator, credentials, exposure

KubernetesKeycloakOperator is the registry prerequisite and, with its
default namespaced watch, lives in the SAME namespace as this server
([operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md)).
The one-time `<name>-initial-admin` Secret is create-once and never
rotated — consume it by reference. The operator's own Ingress is always
disabled by the module; compose external exposure from Gateway API kinds
over the exported service handles.

## On the diagram

Server and operator render as separate nodes with no edge (the pattern
above); exposure kinds draw edges into the exported Service. The
TLS/hostname decisions are spec-internal — they do not add nodes, which
is exactly why they are easy to forget and worth stating in the proposal.

## Pairs well with

- KubernetesKeycloakOperator — required, same namespace (default watch).
- Gateway API kinds (KubernetesHttpRoute + a Gateway) — external exposure
  and the hostname the strict-hostname setting expects.
- KubernetesNamespace — the namespace owner when others share it.
