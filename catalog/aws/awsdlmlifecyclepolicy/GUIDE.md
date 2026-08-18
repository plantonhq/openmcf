# AwsDlmLifecyclePolicy — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The policy never fails loudly

A DLM policy whose role lost permissions, or whose target tags match nothing, does NOT error your deploys — it silently stops producing snapshots (the policy state flips to ERROR at best). Alarm on snapshot age, not on the policy resource: the absence of new snapshots is the real signal.

## One default policy per type per region

AWS allows exactly one VOLUME and one INSTANCE default policy per region — the default arm is an account-global singleton in disguise. Custom tag-targeted policies coexist freely; prefer them anywhere two teams share an account.

## Overlapping target tags fail at apply, not plan

Two policies sharing the same target_tags is an AWS-side rejection the provider cannot see at plan time. Namespace your tags per policy (`backup:hourly`, `backup:daily`) instead of reusing one `backup:true` everywhere.

## Retention math compounds across rules

A schedule keeping 24 hourlies, cross-region-copying with 14-day retention, and archiving after that is three storage bills from one schedule. Read a policy as a cost graph: every rule that keeps something has its own meter.

## copy_tags is a schedule-replacing decision

Changing a schedule's `copy_tags` replaces the WHOLE schedule at the provider (ForceNew) — retention counting restarts from the new schedule's snapshots. Decide tag propagation before the first fire, not after a month of retained history.
