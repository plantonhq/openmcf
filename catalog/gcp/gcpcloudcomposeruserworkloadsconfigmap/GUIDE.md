# GcpCloudComposerUserWorkloadsConfigMap Guide

The judgment this guide protects: the line between this kind and
`GcpCloudComposerUserWorkloadsSecret` is the line between
configuration and credentials. Feature flags, endpoints, batch sizes,
and tuning parameters belong here — anything you would not paste into
a code review does not.

## Plain values, deliberately

Unlike the Secret sibling, `data` values are plain strings — readable
in plans, diffs, and reviews. That visibility is the point: DAG
behavior changes (a flag flip, an endpoint switch) should be reviewable
as text. The moment a value needs base64 or redaction, it has outgrown
this kind — move it to the Secret.

## Updates land in place; consumers re-read on their schedule

`data` is the only mutable field. A ConfigMap edit applies immediately
in Kubernetes, but running tasks keep the values they started with —
KubernetesPodOperator pods read at pod start, DAG-level reads at parse
or task time. Roll a change by applying it and letting the next
scheduled runs pick it up; nothing restarts by itself.

## Name for the role, reference the output

Name, environment, region, and project are immutable — a rename is a
new ConfigMap plus a DAG reference update. Name for the configuration
ROLE (`etl-tuning`, `feature-flags`); downstream wiring should consume
the `config_map_name` output, never a re-typed literal.

## PREVENT fits load-bearing configuration

A destroyed ConfigMap fails every DAG that mounts it — often less
visibly than a missing Secret (defaults kick in, behavior silently
changes). `deletionPolicy: PREVENT` fits ConfigMaps whose absence
changes pipeline BEHAVIOR rather than breaking it loudly; `ABANDON`
leaves the object in the cluster for a management handover.

## What is deliberately absent

No `google_project_service` enablement — the Composer API is
necessarily already on in the environment's project (a ConfigMap
cannot exist without an environment). Credentials and connection URIs
are the Secret sibling's territory, with base64 and state-secret
handling this kind deliberately does not have.
