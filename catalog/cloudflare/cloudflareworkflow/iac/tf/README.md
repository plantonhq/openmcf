# CloudflareWorkflow Terraform Module

Terraform IaC module for a Workflow registration: the binding of a durable-execution class exported by a deployed Worker script to a named workflow.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareWorkflowSpec (generated)
locals.tf     — Retention/limits/schedules shaping (absent trees stay null)
main.tf       — cloudflare_workflow
outputs.tf    — workflow_name, version_id
```

## Behavior

Create IS a PUT at the API (name-as-upsert): registering an existing name adopts and overwrites it. account_id and workflow_name force replacement; class_name and script_name update in place. Retention values pass through verbatim -- the API accepts integer milliseconds ("5000") or duration expressions ("5 minutes"). Deletion is real, but the API answers GET for deleted workflows with an is_deleted marker instead of a 404.

## Outputs

| Name | Description |
|------|-------------|
| `workflow_name` | The workflow's identity -- what Worker workflow bindings reference |
| `version_id` | The workflow version the registration produced |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
