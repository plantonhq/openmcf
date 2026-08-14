# AwsBedrockAgent — Terraform/OpenTofu module

Deploys an Amazon Bedrock agent with its folded satellites: action groups
(`aws_bedrockagent_agent_action_group`), aliases
(`aws_bedrockagent_agent_alias`), collaborators
(`aws_bedrockagent_agent_collaborator`), and knowledge-base associations
(`aws_bedrockagent_agent_knowledge_base_association`) — all keyed by their
stable spec entry names.

Module facts worth knowing before editing:

- **The DRAFT pivot.** Every satellite attaches to the agent's mutable
  DRAFT version, and creating an alias snapshots the draft into a new
  numbered version. The alias resource therefore carries `depends_on` on
  every other satellite — removing that edge makes alias snapshots race
  the rest of the draft assembly.
- **Prepare choreography is the provider's.** Each satellite CRUD
  re-prepares the agent and the provider retries the "agent is Preparing"
  and OptLock conflicts itself; the module adds no sleeps and no retries.
- **Apply-behavior knobs stay provider-default.** `prepare_agent` (true)
  and `skip_resource_in_use_check` (false) are deployment mechanics, not
  desired state — they are deliberately not spec surface (recorded in the
  parity manifest).
- **Prompt overrides are always OVERRIDDEN.** The provider strips
  non-overridden template configurations from state, so the module marks
  every authored entry `prompt_creation_mode = "OVERRIDDEN"` — rendering
  DEFAULT would drift forever.
- **Attribute-vs-block syntax.** `guardrail_configuration`,
  `memory_configuration`, `prompt_override_configuration` (agent) and
  `routing_configuration` (alias) are list ATTRIBUTES upstream —
  assignment syntax, not dynamic blocks.

Outputs mirror the Pulumi module key-for-key: `agent_id`, `agent_arn`,
`draft_version`, and the per-satellite maps keyed by entry name
(`alias_ids`, `alias_arns`, `action_group_ids`, `collaborator_ids`,
`associated_knowledge_base_ids`).
