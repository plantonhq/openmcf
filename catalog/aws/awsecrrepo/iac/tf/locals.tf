locals {
  # Resource-identity tags follow the catalog convention. No Name tag: the
  # repository has a real name of its own (spec.repository_name, a
  # slash-namespaced registry path) — tagging Name with the graph node's
  # metadata.name would show a DIFFERENT value than the repository's actual
  # name in every console tag view. Kinds whose AWS resource carries its own
  # name omit the Name tag (the Global Accelerator / SageMaker convention).
  aws_tags = {
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
  # structurally; this rebuilds the exact document AWS expects. The "days"
  # count unit exists for all three "since..." count types and is
  # materialized here. Rules either expire images (the default) or
  # transition them to the archive storage class; the transition/target
  # coupling and the storageClass/countType pairing are CEL-enforced.
  # Every optional member is merged in only when set — both engines must
  # hand AWS the SAME document, so no member is ever emitted as null.
  lifecycle_rules = [
    for rule in var.spec.lifecycle_rules : merge(
      {
        rulePriority = rule.rule_priority
        selection = merge(
          {
            tagStatus   = rule.tag_status
            countType   = rule.count_type
            countNumber = rule.count_number
          },
          contains(["sinceImagePushed", "sinceImagePulled", "sinceImageTransitioned"], rule.count_type) ? { countUnit = "days" } : {},
          rule.storage_class != "" ? { storageClass = rule.storage_class } : {},
          length(rule.tag_prefixes) > 0 ? { tagPrefixList = rule.tag_prefixes } : {},
          length(rule.tag_patterns) > 0 ? { tagPatternList = rule.tag_patterns } : {},
        )
        action = merge(
          {
            type = rule.action_type != "" ? rule.action_type : "expire"
          },
          rule.target_storage_class != "" ? { targetStorageClass = rule.target_storage_class } : {},
        )
      },
      rule.description != "" ? { description = rule.description } : {},
    )
  ]

  # Repository policy — the Struct arrives from the tfvars layer as a nested
  # object; the provider wants the document as a JSON string.
  repository_policy = var.spec.repository_policy != null ? jsonencode(var.spec.repository_policy) : null
}
