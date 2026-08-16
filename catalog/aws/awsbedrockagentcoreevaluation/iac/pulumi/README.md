# AwsBedrockAgentCoreEvaluation — Pulumi module (Go)

Deploys an AgentCore Evaluations bundle: evaluators, harnesses, and
online evaluation configs (`bedrockagentcore` bridged types).

Module facts worth knowing before editing:

- **Every arm is optional; at least one is required** (spec CEL).
  `evaluators.go`, `harnesses.go`, and `online_configs.go` each
  iterate their collection independently.
- **Harness model vendors and evaluator rating scales are
  exactly-one.** The module renders the arm validation already
  admitted.
- **Evaluator and online-config names are identity seeds** (AWS
  derives `"<name>-<10 chars>"`) with no rename.
- **Harness inference floats are import-normalized** via
  `write_normalized_attributes` in the import catalog — do not round
  them in the module.
- **The spec's single `enabled` knob** on online configs fans out to
  the provider's two lifecycle fields (`enable_on_create` +
  `execution_status`). Keep the Terraform wiring in lockstep.

Outputs mirror the Terraform module key-for-key: name-keyed ID and ARN
maps for evaluators, harnesses, and online configs, plus the
server-assigned output log groups.
