# KubernetesRabbitMqOperator Guide

The judgment this guide carries: this operator breaks three assumptions
the other operator kinds set. It watches ALL namespaces by default (not
its own), it needs cert-manager already running, and — the dangerous one
— destroying it deletes every RabbitMQ cluster on the cluster.

## Destroying the operator destroys the data — this one is different

Unlike the keep-on-uninstall operators (Altinity, Percona, CloudNativePG),
the RabbitmqCluster CRD is part of this operator's applied manifest and is
REMOVED with the resource — which cascade-deletes every RabbitmqCluster,
and their data, cluster-wide (the reference page states the lifecycle).
Never destroy this operator while any KubernetesRabbitMq exists. Fold that
warning into any teardown or migration proposal.

## cert-manager first, or every cluster admission fails

cert-manager is a hard registry prerequisite: the operator ships
`failurePolicy: Fail` admission webhooks whose serving certificate is a
cert-manager Certificate. Without cert-manager running, the webhook cert
is never issued and every RabbitmqCluster admission fails. Declare
KubernetesCertManager first — the
[cert-manager guide](../../kubernetescertmanager/v1alpha1/GUIDE.md) is its
composition story.

## Exactly one install; watches everything by default

The operator installs into a fixed `rabbitmq-system` namespace and its
cluster-scoped admission webhooks are fixed-name singletons — a second
install cannot coexist. Its default watch scope is ALL namespaces (the
opposite of the Altinity and external-secrets operators); set
`watchNamespaces` only to fence it. The invisible-edge mechanism:
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md).

## Pairs well with

- KubernetesRabbitMq — the clusters this operator reconciles (see its
  [guide](../../kubernetesrabbitmq/v1alpha1/GUIDE.md)).
- KubernetesCertManager — hard prerequisite for the admission webhooks.
