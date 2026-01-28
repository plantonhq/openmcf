# AWS ECS Service Examples

Below are several examples demonstrating how to define an AWS ECS Service component in
OpenMCF. After creating one of these YAML manifests, apply it with Terraform using the OpenMCF CLI:

```shell
openmcf tofu apply --manifest <yaml-path> --stack <stack-name>
```

---

## Basic ECS Service

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: basic-ecs-service
spec:
  clusterArn:
    value: "arn:aws:ecs:us-east-1:123456789012:cluster/my-existing-cluster"
  container:
    image:
      repo: "nginx"
      tag: "latest"
    port: 80
    replicas: 1
    cpu: 256
    memory: 512
  network:
    subnets:
      - value: "subnet-123abc"
      - value: "subnet-456def"
    securityGroups:
      - value: "sg-123abc"
```

This example creates a basic ECS service:
• Uses Fargate for serverless container execution.
• Single container replica with nginx image.
• Minimal CPU and memory allocation (256 CPU units, 512MB).
• Deployed across multiple subnets for high availability.
• Basic security group for network access.

---

## ECS Service with Environment Variables

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: env-ecs-service
spec:
  clusterArn:
    value: "arn:aws:ecs:us-east-1:123456789012:cluster/my-env-cluster"
  container:
    image:
      repo: "amazon/amazon-ecs-sample"
      tag: "latest"
    port: 8080
    replicas: 2
    cpu: 512
    memory: 1024
    env:
      variables:
        SPRING_PROFILES_ACTIVE: "prod"
        CUSTOM_VAR: "myvalue"
        LOG_LEVEL: "INFO"
  network:
    subnets:
      - value: "subnet-abc123"
      - value: "subnet-def456"
    securityGroups:
      - value: "sg-abc123"
```

This example includes environment configuration:
• Multiple container replicas for high availability.
• Environment variables for application configuration.
• Larger resource allocation for production workloads.
• Spring Boot application with production profile.

---

## ECS Service with ALB Path-Based Routing

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: path-based-ecs-service
spec:
  clusterArn:
    value: "arn:aws:ecs:us-east-1:123456789012:cluster/my-load-balanced-cluster"
  container:
    image:
      repo: "nginx"
      tag: "stable"
    port: 80
    replicas: 2
    cpu: 256
    memory: 512
  network:
    subnets:
      - value: "subnet-111111"
      - value: "subnet-222222"
    securityGroups:
      - value: "sg-111111"
      - value: "sg-222222"
  alb:
    enabled: true
    arn:
      value: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-application-lb/1234567890abcdef"
    routingType: "path"
    path: "/myapp"
    listenerPort: 80
    listenerPriority: 100
```

This example demonstrates ALB integration:
• Path-based routing for microservices architecture.
• Traffic routed to `/myapp` path on the ALB.
• Multiple replicas for load distribution.
• Health checks automatically configured.
• Listener priority for routing order.

---

## ECS Service with ALB Hostname-Based Routing

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: hostname-based-ecs-service
spec:
  clusterArn:
    value: "arn:aws:ecs:us-east-1:123456789012:cluster/my-hostname-cluster"
  container:
    image:
      repo: "myorg/myservice"
      tag: "v1.2.3"
    port: 8080
    replicas: 3
    cpu: 512
    memory: 1024
  network:
    subnets:
      - value: "subnet-111aaa"
      - value: "subnet-222bbb"
    securityGroups:
      - value: "sg-333ccc"
  alb:
    enabled: true
    arn:
      value: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-other-lb/abcdef1234567890"
    routingType: "hostname"
    hostname: "api.example.com"
    listenerPort: 443
    listenerPriority: 200
```

This example uses hostname-based routing:
• Dedicated hostname for the service.
• HTTPS listener for secure communication.
• Higher priority for routing precedence.
• Custom application image with version tag.
• Three replicas for production load handling.

---

## ECS Service with Custom IAM Roles

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: iam-ecs-service
spec:
  clusterArn:
    value: "arn:aws:ecs:us-east-1:123456789012:cluster/my-iam-cluster"
  container:
    image:
      repo: "amazon/amazon-ecs-sample"
      tag: "latest"
    port: 3000
    replicas: 1
    cpu: 512
    memory: 1024
  network:
    subnets:
      - value: "subnet-xyz123"
    securityGroups:
      - value: "sg-xyz123"
  iam:
    taskExecutionRoleArn:
      value: "arn:aws:iam::123456789012:role/my-custom-ecsTaskExecutionRole"
    taskRoleArn:
      value: "arn:aws:iam::123456789012:role/my-app-task-role"
```

This example includes custom IAM roles:
• Task execution role for pulling container images.
• Task role for application AWS service access.
• Custom permissions for specific use cases.
• Secure credential management.
• Application-specific IAM policies.

---

## Production ECS Service

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: production-ecs-service
spec:
  clusterArn:
    value: "arn:aws:ecs:us-east-1:123456789012:cluster/production-cluster"
  container:
    image:
      repo: "myorg/production-app"
      tag: "v2.1.0"
    port: 8080
    replicas: 3
    cpu: 1024
    memory: 2048
    env:
      variables:
        NODE_ENV: "production"
        LOG_LEVEL: "WARN"
        API_VERSION: "v2"
  network:
    subnets:
      - value: "subnet-private-1a"
      - value: "subnet-private-1b"
    securityGroups:
      - value: "sg-production-app"
      - value: "sg-monitoring"
  alb:
    enabled: true
    arn:
      value: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/prod-alb/abcdef1234567890"
    routingType: "hostname"
    hostname: "api.mycompany.com"
    listenerPort: 443
    listenerPriority: 100
  iam:
    taskExecutionRoleArn:
      value: "arn:aws:iam::123456789012:role/prod-ecs-task-execution-role"
    taskRoleArn:
      value: "arn:aws:iam::123456789012:role/prod-app-task-role"
```

This example is production-ready:
• High resource allocation for performance.
• Multiple replicas across availability zones.
• Production environment variables.
• Custom domain with HTTPS.
• Comprehensive IAM roles.
• Monitoring security group.
• Private subnets for security.

---

## Development ECS Service

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: development-ecs-service
spec:
  clusterArn:
    value: "arn:aws:ecs:us-east-1:123456789012:cluster/dev-cluster"
  container:
    image:
      repo: "myorg/dev-app"
      tag: "latest"
    port: 3000
    replicas: 1
    cpu: 256
    memory: 512
    env:
      variables:
        NODE_ENV: "development"
        LOG_LEVEL: "DEBUG"
        DEBUG: "true"
  network:
    subnets:
      - value: "subnet-public-1a"
    securityGroups:
      - value: "sg-dev-access"
```

This example is optimized for development:
• Single replica for cost efficiency.
• Debug logging enabled.
• Development environment variables.
• Public subnet for easy access.
• Minimal resource allocation.
• Latest tag for rapid iteration.

---

## After Deploying

Once you've applied your manifest with OpenMCF tofu, you can confirm that the ECS service is active via the AWS console or by
using the AWS CLI:

```shell
aws ecs list-services --cluster <your-cluster-name>
```

For detailed service information:

```shell
aws ecs describe-services --cluster <your-cluster-name> --services <your-service-name>
```

This will show service details including task count, load balancer configuration, and deployment status.
