# AwsBedrockAgentCoreEvaluation — Terraform/OpenTofu module

Deploys an AgentCore Evaluations bundle: evaluators
(`aws_bedrockagentcore_evaluator`), harnesses
(`aws_bedrockagentcore_harness`), and online evaluation configs
(`aws_bedrockagentcore_online_evaluation_config`).

Module facts worth knowing before editing:

- **Every arm is optional; at least one is required** (spec CEL).
  The module iterates each collection independently.
- **Harness model vendors and evaluator rating scales are
  exactly-one.** The module renders the arm validation already
  admitted.
- **Evaluator and online-config names are for_each keys** and AWS
  identifiers (letters, digits, underscores; no rename).
- **Harness inference floats are import-normalized** via
  `write_normalized_attributes` in the import catalog — do not round
  them in the module.
- **Online config `enable_on_create` is a REQUIRED create argument
  the API never echoes back** — it is config-only for IMPORT (the
  round-trip tolerates its absence), not inert. The spec's single
  `enabled` knob fans out to it plus the updatable
  `execution_status`.

Outputs mirror the Pulumi module key-for-key: name-keyed ID and ARN
maps for evaluators, harnesses, and online configs, plus the
server-assigned output log groups.
