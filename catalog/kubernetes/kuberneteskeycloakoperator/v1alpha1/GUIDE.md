# KubernetesKeycloakOperator Guide

The judgment this guide carries: this operator is fixed-name and
namespaced — exactly one per namespace, and Keycloak servers live in the
SAME namespace as the operator that reconciles them. Installing it alone
deploys nothing.

## One per namespace; co-located with its servers

Every resource the bundle installs is upstream-fixed-name
(`keycloak-operator`, ...), so a second install cannot share a namespace,
and with the default namespaced watch the operator reconciles only
KubernetesKeycloak resources beside it. The composition shape is
therefore co-located: operator and servers in one namespace, per team or
per environment. The invisible-edge mechanism:
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md).

## No version field, by design

There is no version selector: the KubernetesKeycloak CR rendering is built
against the CRD schema THIS bundle installs, so a selectable version would
drift the schema from what the declaration kind renders (the reference
page states this). Upgrades arrive as module updates — do not go looking
for a version knob.

## No cert-manager, no webhooks

Unlike the RabbitMQ operator, this bundle ships no admission webhooks and
has no cert-manager dependency — a plain set of manifests. Nothing extra
to compose ahead of it beyond its namespace.

## Namespace ownership

Because servers are co-located, the namespace is shared — wire the
operator and its Keycloak resources to a dedicated KubernetesNamespace
rather than `createNamespace: true`
([namespace-ownership pattern](../../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesKeycloak — the servers this operator reconciles, in the same
  namespace (see its [guide](../../kuberneteskeycloak/v1alpha1/GUIDE.md)).
