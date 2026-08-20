locals {
  # Resource-identity tags match the Pulumi module key-for-key. NOTE:
  # neither aws_cloudwatch_event_connection nor
  # aws_cloudwatch_event_api_destination is taggable at AWS - the
  # deliberate tag-convention absence (the AwsCloudwatchDashboard
  # precedent). Kept for the day AWS adds tagging.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEventBridgeApiDestination"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
