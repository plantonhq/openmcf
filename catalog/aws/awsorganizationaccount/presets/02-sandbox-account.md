# Sandbox Account

This preset creates a disposable sandbox account whose destroy
actually CLOSES it — for experiments that should not leave a live
orphan account behind.

## When to Use

- Per-developer or per-experiment sandboxes with a real end of life
- Hackathons, trainings, proof-of-concept accounts

## What You Get

- A member account in the sandbox OU (attach a spend-limiting SCP
  there)
- `closeOnDeletion: true` — destroy calls CloseAccount instead of
  removal, so nothing survives standalone

## Customize

- Know the closure mechanics: AWS suspends the account
  (~90 days PENDING_CLOSURE) before permanent deletion, holds its
  email through that window, and rate-limits closures (~10% of
  accounts per rolling 30 days) — churn sandboxes deliberately, not
  hourly
- Use per-sandbox plus-addressed emails so recreations never collide
- The default (removal) is the right contract for anything whose data
  must survive the resource — keep `closeOnDeletion` for true
  throwaways only
