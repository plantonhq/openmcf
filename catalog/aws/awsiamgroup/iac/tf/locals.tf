locals {
  # inline_policies is free-form JSON (map<string, google.protobuf.Struct>), typed `any` in
  # variables.tf because its entries have heterogeneous shapes. Encode each policy document to a
  # JSON string here so the result is a homogeneous map(string): aws_iam_group_policy.for_each
  # accepts a map/set, and converting a heterogeneous object to a map would otherwise fail with
  # "all map elements must have the same type".
  inline_policies_json = {
    for k, v in var.spec.inline_policies : k => jsonencode(v)
  }

  # IAM groups (and their membership/policy satellites) are untaggable at
  # AWS, so this module carries no tag map - the one deliberate absence
  # against the catalog's tag convention (mirrored in the Pulumi module).
}
