# API Gateway account settings for one region -- a settings singleton:
# AWS keeps exactly one account object per account+region, and this
# module manages it. metadata.name never reaches AWS.
#
# Lifecycle facts the render below depends on:
#   - the ONE configurable lever is the CloudWatch Logs role; an empty
#     string resets it (the provider patches /cloudwatchRoleArn to nil
#     for "" and null alike), so passing the spec value through
#     unconditionally is faithful for both the set and clear postures;
#   - destroy RESETS the role to none -- the account object itself
#     always exists and cannot be deleted;
#   - AWS validates the role at apply (trust: apigateway.amazonaws.com
#     plus CloudWatch Logs write permissions) and the provider retries
#     through IAM propagation lag ("The role ARN does not have
#     required permissions").

resource "aws_api_gateway_account" "this" {
  cloudwatch_role_arn = var.spec.cloudwatch_role_arn
}
