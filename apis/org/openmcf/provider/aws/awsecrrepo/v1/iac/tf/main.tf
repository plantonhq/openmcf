resource "aws_ecr_repository" "this" {
  name                 = var.spec.repository_name
  image_tag_mutability = local.image_tag_mutability
  force_delete         = try(var.spec.force_delete, false)

  encryption_configuration {
    encryption_type = local.encryption_type
    kms_key         = local.kms_key_id
  }

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = local.tags
}

# Lifecycle policy to clean up untagged images older than 1 day
resource "aws_ecr_lifecycle_policy" "this" {
  repository = aws_ecr_repository.this.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 30 images, expire all others"
        selection = {
          tagStatus     = "untagged"
          countType     = "imageCountMoreThan"
          countNumber   = 30
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Remove untagged images older than 1 day"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 1
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}


