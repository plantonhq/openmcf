locals {
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  password  = var.spec.password != "" ? var.spec.password : null
  user_type = var.spec.user_type
}
