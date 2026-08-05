# Build Contract

The compile loop's machine interface. Agents and CI must treat this as a
published contract — exit codes and JSON shape are stable.

## Command

```
planton chart build <chart-dir> -o json
```

Read-only: nothing is written to disk unless you pass `--output-dir`. JSON mode
puts the **entire report on stdout** and keeps stderr for human messages only
(connection errors, deprecation warnings).

Build a variant without editing the chart with `--set name=value` (repeatable):

```
planton chart build . -o json --set dns_enabled=true --set aws_region=eu-west-1
```

The value is parsed as YAML, exactly as if values.yaml carried it: `true` is a
bool, `3` is a number, `"1.29"` stays a string, `[a, b]` is a list. The chart
on disk is never touched — this is the way to prove a bool toggle in both
positions (build once per position, compare the `resources[]` arrays).

Pin the control plane when needed (local instance, CI, desktop studio):

```
PLANTON_API_ENDPOINT_INFRA_HUB=127.0.0.1:23802 planton chart build . -o json
```

## Exit codes

| Code | Meaning | Agent action |
|------|---------|--------------|
| **0** | Chart valid. Warnings may be present; they never fail the build. | Phase 4 self-check, then stop or offer publish. |
| **1** | Chart has fixable errors. JSON report on stdout lists every issue. | Parse `issues[]`, fix templates/values, rebuild. |
| **2** | Check could not run. **Stdout is empty.** | Do not edit the chart. Fix environment (not a chart dir, bad flags, an unknown `--set` param name, control plane down) or report and stop. |

Exit 2 is the critical guardrail: an unreachable control plane has never been
fixed by changing a template.

## JSON report shape

```json
{
  "result": "passed",
  "errors": 0,
  "warnings": 0,
  "issues": [],
  "resources": [
    {
      "kind": "AwsVpc",
      "name": "dev-vpc",
      "sourceFile": "network.yaml"
    }
  ]
}
```

| Field | Meaning |
|-------|---------|
| `result` | `"passed"` or `"failed"` |
| `errors` / `warnings` | Counts (only errors affect exit code) |
| `issues[]` | Every problem. Each issue has `severity`, `message`, and usually `file`, `resourceKind`, `resourceName`, `fieldPath` when attributable. |
| `resources[]` | **Rendered truth** after Jinjava — every resource the chart produces post-conditionals. Use this to verify edits landed. |
| `overrides[]` | The applied `--set` values (name + parsed value). **Present only when `--set` was used** — its absence means the report reflects values.yaml as written. |

Resource names render with the literal placeholder `<env>` where
`{{ values.env }}` appears (JSON-escaped as `\u003cenv\u003e`) — the real
environment is only bound at deployment. Seeing `<env>-vpc` in the report
means substitution worked, not that it failed.

Local parse failures (malformed `Chart.yaml` or `values.yaml`) and
server-side rejection of the chart manifest itself (e.g. a missing
`spec.selector`) all arrive through the same `issues[]` channel with
exit 1 — one channel for everything fixable.

## Human mode (no `-o json`)

Exit codes are identical. The formatted report goes to **stderr**; stdout stays
empty unless `--show` renders template YAML. Agents should always use `-o json`.

## Flag interactions

- `--show` + `-o json` → exit 2 (conflicting outputs).
- `--set unknown_name=x` → exit 2, stderr lists the params values.yaml declares. Fix the flag, never the chart.
- `--no-browser`, `--copy` → deprecated no-ops; fleet Makefiles may still pass them.
- `--output-dir <dir>` → writes rendered `template.yaml`; optional for agents.

## Offline gates (no control plane)

| Command | Purpose |
|---------|---------|
| `planton explain <Kind>` | Field names, types, docs, outputs (drill with dotted paths) |
| `planton validate manifest.yaml` | Validate one draft manifest offline |

These reduce errors before the compile loop but do **not** replace
`chart build` — only the control plane runs full template rendering,
valueFrom validation, and cross-resource checks.

## CI usage

```bash
planton chart build ./charts/my-chart -o json > report.json
code=$?
if [ "$code" -eq 1 ]; then
  jq '.issues[] | select(.severity=="error")' report.json
  exit 1
fi
if [ "$code" -eq 2 ]; then
  echo "build could not run" >&2
  exit 2
fi
```

Treat exit 2 as infrastructure failure, not chart failure.
