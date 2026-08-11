# The option group managed for inline spec.options -- engine features
# like Oracle TDE/OEM or SQL Server native backup, activated as a named
# option list (glue, so it stays inside this module). Engine name and
# major version are derived from the spec (locals.tf). The provider's
# EC2-Classic db_security_group_memberships argument is deliberately
# unused -- security group access composes through
# vpc_security_group_memberships references.
resource "aws_db_option_group" "this" {
  count = local.manage_option_group ? 1 : 0

  name                     = local.instance_identifier
  engine_name              = var.spec.engine
  major_engine_version     = local.option_major_engine_version
  option_group_description = "Managed by Planton for ${local.instance_identifier}"

  dynamic "option" {
    for_each = var.spec.options
    content {
      option_name                    = option.value.option_name
      port                           = option.value.port != 0 ? option.value.port : null
      version                        = option.value.version != "" ? option.value.version : null
      vpc_security_group_memberships = length(option.value.vpc_security_group_memberships) > 0 ? option.value.vpc_security_group_memberships : null

      dynamic "option_settings" {
        for_each = option.value.option_settings
        content {
          name  = option_settings.value.name
          value = option_settings.value.value
        }
      }
    }
  }

  tags = local.aws_tags
}
