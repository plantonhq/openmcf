resource "aws_ecr_repository" "this" {
  # The repository name is the immutable registry path (ForceNew). Changing
  # it replaces the repository; images are NOT migrated.
  name                 = var.spec.repository_name
  image_tag_mutability = local.image_tag_mutability
  force_delete         = var.spec.force_delete

  # Wildcard filters inverting the base mutability for matching tags — only
  # present in the *_WITH_EXCLUSION modes. WILDCARD is the only filter type
  # AWS supports today.
  dynamic "image_tag_mutability_exclusion_filter" {
    for_each = var.spec.image_tag_mutability_exclusion_filters
    content {
      filter      = image_tag_mutability_exclusion_filter.value
      filter_type = "WILDCARD"
    }
  }

  # The entire encryption configuration is create-time (ForceNew): changing
  # the type or the key replaces the repository.
  encryption_configuration {
    encryption_type = local.encryption_type
    kms_key         = local.kms_key_id
  }

  image_scanning_configuration {
    scan_on_push = local.scan_on_push
  }

  tags = local.aws_tags
}

# Lifecycle policy — a separate AWS API resource keyed 1:1 by the repository,
# folded into the spec as structured rules. Created only when rules exist;
# removing every rule removes the policy.
resource "aws_ecr_lifecycle_policy" "this" {
  count      = length(local.lifecycle_rules) > 0 ? 1 : 0
  repository = aws_ecr_repository.this.name

  policy = jsonencode({
    rules = local.lifecycle_rules
  })
}

# Repository policy — resource-based access control (cross-account pulls,
# service principals like Lambda pulling container images). Also a separate
# 1:1 API resource folded into the spec.
resource "aws_ecr_repository_policy" "this" {
  count      = local.repository_policy != null ? 1 : 0
  repository = aws_ecr_repository.this.name
  policy     = local.repository_policy
}
