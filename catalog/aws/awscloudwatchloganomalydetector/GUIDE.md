# AwsCloudwatchLogAnomalyDetector — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The list that takes one entry

AWS's PutLogAnomalyDetector models `logGroupArnList` as a list but currently rejects more than one entry. The spec stays list-shaped (AWS's own forward-compatible contract); give it exactly one log group until AWS lifts the cap — a second entry fails at apply, not at validation.

## Pause, don't delete

`enabled: false` stops evaluation but keeps the trained baseline; deleting the detector discards it and retraining takes up to a day of live traffic. Pause through incident noise, delete only on decommission.

## AccessDenied looks like deletion

The provider drops the detector from state when reads return AccessDenied — so an IAM regression shows up as "detector vanished" in plans. Check permissions before concluding someone deleted it.

## Expect a quiet first day

Anomalies only surface after the model builds a baseline (hours to ~24h). A freshly deployed detector reporting nothing is training, not broken.
