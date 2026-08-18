# Platform Metrics Workspace

The EKS-fleet landing zone: 90-day retention, a recording+alerting rules namespace, and Alertmanager routing to SNS (wire your AwsSnsTopic's ARN). Pair it with an AwsManagedPrometheusScraper per cluster, or point remote-write agents at `prometheus_endpoint`.
