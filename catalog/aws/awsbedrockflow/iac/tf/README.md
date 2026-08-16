# AwsBedrockFlow — Terraform/OpenTofu module

Deploys an Amazon Bedrock flow (`aws_bedrockagent_flow`) — a node graph
orchestrating prompts, agents, knowledge bases, Lambdas, and control-flow
logic.

Module facts worth knowing before editing:

- **Union members are derived from `type`.** Each configurable node class
  renders exactly its AWS configuration member (Agent → `agent`, Prompt →
  `prompt`, ...); the structural classes (Input, Output, Iterator,
  Collector) render an EMPTY member (`configuration { input {} }` — AWS
  requires the union member even when it carries nothing), and the Loop
  family renders no configuration at all (its member is not expressible
  at the pinned provider — an upstream gap, recorded in the parity
  manifest). Connection types (Data/Conditional) are derived the same
  way.
- **Graph topology is validated by AWS**, not the module: unreachable
  nodes, socket type mismatches, and missing connections fail at
  create/update server-side with named validation classes.
- **One-value vocabularies are module constants**: inline-code language
  `Python_3`, cache point type `default`, retrieval/storage service `S3`.
- **The inline prompt tree mirrors the AwsBedrockPrompt kind's module**
  arm-for-arm (upstream shares the same Go models between the two
  resources) — change them together.
- **The provider declares timeouts it never consumes** for this resource
  (create/update/delete run unwaited single calls at the pin) — do not
  add a timeouts block expecting behavior.

Outputs mirror the Pulumi module key-for-key: `flow_id`, `flow_arn`,
`draft_version`.
