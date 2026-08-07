locals {
  # SVMs carry TWO names: the ONTAP-internal spec.name (the SVM identity in
  # junction paths, SnapMirror, and DNS — underscore-only charset) and the
  # cloud resource's metadata.name, which becomes the Name tag so the AWS
  # console shows the same identity both engines pin.
  resource_name = var.metadata.name

  # Resource-identity tags follow the catalog convention; user labels merge in
  # without being able to override the identity keys.
  aws_tags = merge(try(var.metadata.labels, {}), {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsFsxOntapStorageVirtualMachine"
    "planton.ai/resource-id"   = var.metadata.id
  })

  # Empty string means "no vsadmin access" — omit the argument entirely so
  # AWS never sees an empty password.
  svm_admin_password = var.spec.svm_admin_password != "" ? var.spec.svm_admin_password : null
}
