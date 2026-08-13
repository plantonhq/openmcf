# An IAM DENY policy — rules that BLOCK principals from using specific
# permissions regardless of any role grants they hold. Deny always outranks
# allow, which is what makes deny policies the guardrail layer rather than
# another access grant.
#
# No API-enablement resource here: deny policies attach to projects,
# folders, or organizations through the always-on IAM v2 surface, and
# creating them requires org-level permissions (roles/iam.denyAdmin) in any
# case.
resource "google_iam_deny_policy" "this" {
  name   = local.policy_name
  parent = local.parent

  display_name = var.spec.display_name != "" ? var.spec.display_name : null

  dynamic "rules" {
    for_each = var.spec.rules
    content {
      description = rules.value.description != "" ? rules.value.description : null

      deny_rule {
        denied_principals     = length(rules.value.deny_rule.denied_principals) > 0 ? rules.value.deny_rule.denied_principals : null
        exception_principals  = length(rules.value.deny_rule.exception_principals) > 0 ? rules.value.deny_rule.exception_principals : null
        denied_permissions    = length(rules.value.deny_rule.denied_permissions) > 0 ? rules.value.deny_rule.denied_permissions : null
        exception_permissions = length(rules.value.deny_rule.exception_permissions) > 0 ? rules.value.deny_rule.exception_permissions : null

        dynamic "denial_condition" {
          for_each = rules.value.deny_rule.denial_condition != null ? [rules.value.deny_rule.denial_condition] : []
          content {
            expression  = denial_condition.value.expression
            title       = denial_condition.value.title != "" ? denial_condition.value.title : null
            description = denial_condition.value.description != "" ? denial_condition.value.description : null
            location    = denial_condition.value.location != "" ? denial_condition.value.location : null
          }
        }
      }
    }
  }

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null
}
