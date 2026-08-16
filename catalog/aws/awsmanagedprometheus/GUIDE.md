# AwsManagedPrometheus — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Never set an alias you might want to remove

AWS offers no un-alias: once set, an alias can change but never clear — both engines REPLACE the workspace (new workspace ID, lost metrics) when the field empties. Renames in place are fine; removal is destruction.

## The configuration outlives your manifest

The workspace configuration (retention, label-set limits) is created via update and has NO delete API. Removing the block from the spec leaves the last-applied values on the workspace — set the values you want instead of removing the block and expecting defaults back.

## Alertmanager speaks SNS only

AMP's managed Alertmanager supports the SNS receiver (from which SNS fans out to email/Slack/PagerDuty). A pasted alertmanager.yml with webhook receivers validates at AWS but never fires — route through an AwsSnsTopic.

## Label-set limits are the noisy-neighbor control

`limits_per_label_set` caps active series per label population (e.g. per team label). A cap of 0 blocks ingestion for that set entirely — a kill switch for a misbehaving emitter, applied in place.

## Rule namespaces are whole files

Each namespace's `data` replaces the namespace's whole rules file per apply — treat one namespace as one owned file (per team, per service), never hand-merge rules across owners into one namespace.
