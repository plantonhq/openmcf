locals {
  # Resource-identity tags follow the catalog convention. The repository's
  # cloud name comes from spec.repository_name (ECR names are slash-namespaced
  # registry paths, a different concept from the graph node's name).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEcrRepo"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Tag mutability — the spec default is MUTABLE (the AWS default). The
  # exclusion filters only apply to the *_WITH_EXCLUSION modes (CEL-enforced),
  # and AWS currently supports only wildcard-type filters, so filter_type is
  # materialized here rather than modeled in the spec.
  image_tag_mutability = coalesce(var.spec.image_tag_mutability, "MUTABLE")

  # Encryption — the whole configuration is create-time (ForceNew). A
  # customer-managed key only applies to the KMS types (CEL-enforced); with
  # AES256 the key must be null.
  encryption_type = coalesce(var.spec.encryption_type, "AES256")
  kms_key_id      = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # scan_on_push defaults to true (the spec's recommended security posture)
  # when the optional is pruned from the tfvars.
  scan_on_push = coalesce(var.spec.scan_on_push, true)

  # Lifecycle rules — the spec models the ECR lifecycle policy JSON
  # structurally; this rebuilds the exact document AWS expects. "expire" is
  # the only action ECR supports and the "days" count unit only exists for
  # sinceImagePushed, so both are materialized here. tagged rules carry
  # exactly one selector list (CEL-enforced); untagged/any rules carry none.
  lifecycle_rules = [
    for rule in var.spec.lifecycle_rules : {
      rulePriority = rule.rule_priority
      description  = rule.description != "" ? rule.description : null
      selection = merge(
        {
          tagStatus   = rule.tag_status
          countType   = rule.count_type
          countNumber = rule.count_number
        },
        rule.count_type == "sinceImagePushed" ? { countUnit = "days" } : {},
        length(rule.tag_prefixes) > 0 ? { tagPrefixList = rule.tag_prefixes } : {},
        length(rule.tag_patterns) > 0 ? { tagPatternList = rule.tag_patterns } : {},
      )
      action = {
        type = "expire"
      }
    }
  ]

  # Repository policy — the Struct arrives from the tfvars layer as a nested
  # object; the provider wants the document as a JSON string.
  repository_policy = var.spec.repository_policy != null ? jsonencode(var.spec.repository_policy) : null
}
