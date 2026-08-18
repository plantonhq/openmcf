# AwsEcrRegistrySettings — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Know each arm's destroy contract before you rely on it

Destroy is three different things here: the registry policy and the keyed collections (cache rules, templates, exclusions) genuinely delete; scanning and replication RESET to empty defaults (AWS has no delete — the modules put the defaults back); account settings PERSIST at their last-applied values. To "undo" an account setting, apply the value you want BEFORE destroying — removal changes nothing.

## Pair every cache rule with a creation template

A pull-through cache mints repositories on first pull — without a matching template they arrive with bare defaults (mutable tags, AES256, no lifecycle policy). The template with the same prefix is what makes cached repositories arrive governed: immutable tags, your KMS key, an expiry policy for stale cached images.

## Cache credentials do not un-set

Once a rule carries a credential or custom-role ARN, clearing the field back to empty is silently not propagated by the provider — the old credential stays attached. Replace the rule (a new prefix, or destroy-and-recreate) to genuinely drop credentials; rotating the secret's VALUE needs no rule change at all.

## Enhanced scanning is an Inspector billing decision

`scan_type: ENHANCED` hands scanning to Amazon Inspector: OS + language packages and continuous re-scanning as new CVEs publish — billed per image per month by Inspector, not free like basic scan-on-push. CONTINUOUS_SCAN rules only exist there; the manifest validation walls it off from BASIC.

## Exclude your CI roles from pull-time metrics early

Lifecycle policies that expire by days-since-last-pull are silently defeated by automation: replication, cache refreshes, and CI pulls all refresh the metric. Register those principals in `pull_time_update_exclusions` from day one — retrofitting after images have wrong last-pull stamps cannot rewrite history.
