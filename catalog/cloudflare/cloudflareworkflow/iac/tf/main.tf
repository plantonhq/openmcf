# Cloudflare Workflow registration: binds a class exported by a DEPLOYED
# Worker script to a named workflow in the account. Cloudflare's create IS a
# PUT (name-as-upsert): registering an existing name adopts and overwrites it
# rather than failing -- names must be chosen deliberately. account_id and
# workflow_name force replacement; class_name and script_name update in place
# (a full-body PUT to the same endpoint). Deletion is real, but the API keeps
# answering GET for deleted workflows with an is_deleted marker instead of a
# 404.
resource "cloudflare_workflow" "main" {
  account_id    = var.spec.account_id
  workflow_name = var.spec.workflow_name
  class_name    = var.spec.class_name
  script_name   = var.spec.script_name

  default_retention = local.default_retention
  limits            = local.limits
  schedules         = local.schedules
}
