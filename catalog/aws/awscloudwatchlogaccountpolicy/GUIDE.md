# AwsCloudwatchLogAccountPolicy — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Each type carries its own document schema

The `policy_document` is NOT one grammar: DATA_PROTECTION carries data-identifier statements, SUBSCRIPTION_FILTER a destination + filter pattern object, FIELD_INDEX a `Fields` list, TRANSFORMER a processor pipeline, METRIC_EXTRACTION a metric mapping. AWS validates server-side at Put — a wrong-shaped document fails the apply, never the plan.

## One policy per type is the practical quota

AWS bounds account policies per type (one for most). Treat the (name, type) pair as a singleton per capability: a second FIELD_INDEX_POLICY under a different name will be rejected server-side where the type is single-instance.

## Selection criteria exists only for subscription filters

AWS accepts `selection_criteria` on SUBSCRIPTION_FILTER_POLICY alone, and its one supported grammar is an exact-name exclusion list: `LogGroupName NOT IN ["noisy-group"]`. Every other type — data protection, field index, transformer, metric extraction — applies account-wide with no narrowing; PutAccountPolicy rejects any criteria string on them with "Invalid selection criteria provided" (live-verified against every documented grammar, including the prefix form older AWS blog posts show). If you need a subtree-scoped transformer or index, use the per-log-group resource instead of the account policy.

## Transformer account policies vs per-group transformers

A TRANSFORMER_POLICY here applies to every selected group; the per-log-group transformer (on AwsCloudwatchLogGroup) wins where both exist. Prefer per-group transformers for service-specific parsing, the account policy for org-wide normalization.
