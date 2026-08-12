# Azure Monitor Autoscale Setting -- Operational Guide

Judgment calls that matter when you run autoscale in production.

## Scale out eagerly, scale in lazily

Flapping -- scale out, scale in, scale out again -- burns money and warm-up time. The standard defense is asymmetry: a scale-out rule with a short window and cooldown (react in minutes) paired with a scale-in rule with a LOWER threshold, a longer window, and a longer cooldown. Azure also applies its own "flapping" guard: it projects what the metric would look like after a scale-in and skips the action if it would immediately trigger a scale-out -- so a scale-in that "mysteriously never happens" usually means your in/out thresholds are too close together.

## The default capacity is your metrics-outage posture

`capacity.default` is what Azure applies when the metric is UNAVAILABLE -- and only if the current count is below it (metric loss never scales you in). Set it to the count that keeps the service healthy under normal peak, not the minimum: during a metrics outage you want to be over-provisioned, not at the floor.

## A schedule marks a START, not a window

A recurrence profile takes effect at its day/hour/minute and STAYS in effect until another profile's schedule begins -- there is no "end time". A business-hours profile therefore needs a partner: a second recurrence profile at 18:00 that returns to the overnight shape. A lone scheduled profile silently becomes the permanent shape after its first activation.

## Rules inside scheduled profiles still run

Exactly one profile is in effect at any moment (fixed-date beats recurrence beats default), and THAT profile's rules and envelope govern. Scheduled profiles without rules pin capacity to their default count -- the right shape for predictable loads. Scheduled profiles WITH rules give you a different elasticity envelope per time window -- the right shape when nights are quiet but not dead.

## Treat ExactCount changes as deliberate jumps

`ChangeCount` and `PercentChangeCount` nudge; `ExactCount` teleports. A misconfigured ExactCount rule (or a fixed-date profile with an aggressive default) is the classic accidental fleet-doubling. Review any ExactCount rule against the profile's maximum before shipping it.

## Predictive autoscale: forecast first, act later

Predictive autoscale (VM Scale Sets only, CPU-based) has a safe on-ramp: run `ForecastOnly` for a week and compare the forecast against reality in the portal before switching to `Enabled`. Use `look_ahead_time` to cover your instances' real warm-up (image pull, JIT, cache fill) -- capacity that arrives exactly on time is capacity that arrived late.

## One setting per target, so charts own the whole rule book

Azure allows a single autoscale setting per resource. Treat the setting as the one authoritative rule book for its target: when different teams want different rules on the same scale set, they edit this one object -- deploying a second setting fails, and hand-editing in the portal drifts the IaC state. Route every change through the manifest.

## When the count never moves, read Run history before touching rules

The portal's Run history records every evaluation with the observed metric value and the decision. The three usual culprits, in order: the rule's `metric_resource_id` points at the wrong resource (nothing emits the metric), the envelope leaves no room (min == max), or a paired scale-in threshold is so close to the scale-out threshold that the flapping guard vetoes every action.

## Email notifications name real inboxes, never "the administrator"

Azure retired classic subscription administrators in April 2024, and ARM now rejects any autoscale notification that asks to email "the subscription administrator" (`SendEmailsToAdminCoAdminsNotSupported`) -- which is why this kind has no such flags. Put your team's on-call or alias address in `notification.email.customEmails` (an email block needs at least one), or use a webhook to reach chat and paging systems.
