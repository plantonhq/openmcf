locals {
  # The client's cloud name is metadata.name -- the same basis the Pulumi
  # module uses, so both engines create identically-named app clients.
  #
  # No aws_tags map here: the aws_cognito_user_pool_client resource is not
  # taggable (identity tagging lives on the pool).
  resource_name = var.metadata.name
}
