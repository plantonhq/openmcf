# Azure Backup Policy (VM) -- Operational Guide

Judgment that saves real time when designing VM backup policies. The field reference lives in the API Explorer; this is the operational layer above it.

## Start V2 unless you have a reason not to

`policyType` defaults to V1 (the provider's own default, kept for parity), but V2 -- the "enhanced" policy -- is what new designs should say explicitly: it unlocks hourly schedules, zonally resilient instant restore, and instant-restore retention beyond 5 days, at no extra cost class. The catch: `policyType` is ForceNew, and a policy replacement re-binds every protected VM under it -- so decide the generation BEFORE VMs bind, not after.

## The retention ladder is your backup bill

Backup storage cost is retention math: `retentionDaily.count` dailies + weeklies + monthlies + yearlies, each a recovery point at (incremental) storage prices. The classic shape -- 30 daily, 12 weekly, 12 monthly, 7 yearly -- suits most workloads; resist the instinct to keep everything daily for years. For multi-year retention, add `tieringPolicy` (TierRecommended is the safe default) so aged points move to archive-tier prices.

## Azure's daily floor is 1 or 7+

The service rejects 2-6 days of daily retention at CREATE time -- a rule that surfaces nowhere in the portal until it fails. The spec front-loads it: `count: 1` (a single rolling daily) or `count: 7+`.

## The time field is one dial for everything

`backup.time` sets the backup start AND the retention times of every layer (the provider wires them together). It must land on the hour or half past. Pick a low-traffic window in the policy's `timezone` -- and remember VM snapshots briefly elevate IO.

## Hourly is a window, not a metronome

An hourly V2 schedule runs inside a WINDOW: `time` starts it, `hourDuration` bounds it, `hourInterval` spaces backups within it. The duration must be a multiple of the interval (4/12 gives three backups a day). RPO-driven workloads pick interval 4; the window keeps backup churn out of business-peak hours.

## Retention forms: week-of-month vs month-days

Monthly and yearly layers pick WHICH backup to keep, in exactly one of two grammars: `weeks` + `weekdays` ("First Sunday") or `days` / `includeLastDays` ("the 1st and the last day"). Mixing them is rejected at manifest time (the provider's own contract). `includeLastDays` exists because months differ in length -- it is the only correct way to say "the last day".

## Changing a policy changes every VM under it

A policy update applies to all protected VMs at the next scheduled backup -- there is no per-VM override. Shortening retention DELETES existing recovery points beyond the new horizon on the service side; treat retention reductions as a deliberate, announced operation (and consider vault immutability, which blocks them outright).
