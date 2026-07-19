# Platform-Gap Issue Form and the Connection-Name Chart Convention

**Date**: July 19, 2026
**Type**: Enhancement
**Components**: GitHub Issue Templates, Chart Authoring Conventions

## Summary

Two small durable pieces land together. A new `platform-gap.yml` issue form
gives users (and the AI agents acting on their behalf) a structured way to
report "Planton fell short of what I needed" — distinct from a bug: the
platform worked as built, but the built thing wasn't enough. And the chart
authoring rule's annotation conventions now document the explicit
connection-name pattern for one-run cluster compositions: the cluster
manifest names the Kubernetes connection it publishes
(`planton.dev/connection-name`), and every in-cluster workload references
the identical values expression, so the producer and consumer ends of the
contract can never drift.

## Problem Statement / Motivation

Real usage produced both needs. Users hitting a genuine capability gap had
no shaped path to report it — the context died in chat screenshots — and
the maintainer workflow (hand a detailed issue to a coding agent to plan
the fix) works only when issues arrive with the goal, the attempt, and the
shortfall spelled out. Separately, the legacy connection-slug formula
(`<env>-<cluster resource name>`) assumed cluster resources are named
WITHOUT an env prefix, which collides with the fleet's own "always
env-prefix resource names" convention — chart authors following both rules
at once produced double-prefixed slugs and instantly-failing Kubernetes
workloads.

## Solution / What's New

- `.github/ISSUE_TEMPLATE/platform-gap.yml`: goal / attempt / shortfall /
  versions sections, labeled `platform-gap` + `needs-triage`. The section
  structure is deliberately machine-friendly — these issues are routinely
  handed to coding agents for fix planning.
- `_rules/charts/build-and-fix-planton-infra-charts.mdc`: the
  platform-behavior annotation conventions now lead with the explicit
  pattern — one `connection_name` values param used on the cluster's
  `planton.dev/connection-name` AND every workload's
  `planton.dev/connection` — with the legacy default formula documented
  for existing charts.

## Impact

- **Chart authors** get a naming pattern that composes with the fleet's
  env-prefix rule instead of contradicting it.
- **Users with unmet needs** get a first-class reporting path whose
  structure feeds the agent-driven fix pipeline directly.
