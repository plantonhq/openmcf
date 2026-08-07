resource "aws_wafv2_regex_pattern_set" "this" {
  # The set's AWS name is the Planton resource name -- the stable identity
  # web ACL statements and operators see. Name and scope are create-time
  # immutable (ForceNew).
  name  = var.metadata.name
  scope = var.spec.scope

  # One block per expression. AWS validates the regex dialect server-side
  # (PCRE subset: no backreferences or lookaround), so an unsupported
  # pattern fails the apply with AWS's own message rather than a module
  # guess.
  dynamic "regular_expression" {
    for_each = var.spec.regular_expressions
    content {
      regex_string = regular_expression.value
    }
  }

  description = local.description

  tags = local.aws_tags
}
