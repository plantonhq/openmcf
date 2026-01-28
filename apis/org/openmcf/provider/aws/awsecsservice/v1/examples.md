# Create using CLI

Create a yaml manifest using one of the examples below. After the YAML is created, use the command below to apply with
OpenMCF:

```shell
openmcf pulumi up --manifest <yaml-path> --stack <stack-name>
```

Or, if using Terraform:

```shell
openmcf terraform apply --manifest <yaml-path> --stack <stack-name>
```

(You can also use a shorter form like `planton apply -f <yaml-path>` if your environment is configured accordingly.)

---

# Basic Fargate Example

This example deploys a simple ECS service onto an existing ECS cluster using AWS Fargate. It specifies the cluster ARN,
subnets, a single container replica, CPU, and memory.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: my-basic-ecs-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/my-existing-cluster"
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
      - "subnet-123abc"
      - "subnet-456def"
    securityGroups:
      - "sg-123abc"
```

---

# Example with Environment Variables and Secrets

This example demonstrates how to pass both plain environment variables and secrets to the container.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: ecs-service-with-env
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/my-env-cluster"
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
      secrets:
        DB_PASSWORD: "arn:aws:ssm:us-east-1:123456789012:parameter/db_password"
        API_KEY: "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-api-key"
  network:
    subnets:
      - "subnet-abc123"
      - "subnet-def456"
    securityGroups:
      - "sg-abc123"
```

---

# Example with ALB Path-Based Routing

In this example, an Application Load Balancer (ALB) is used to route traffic via a specific path (e.g., `/myapp`). Make
sure your ALB and target ECS cluster are in the same VPC.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: path-based-ecs-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/my-load-balanced-cluster"
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
      - "subnet-111111"
      - "subnet-222222"
    securityGroups:
      - "sg-111111"
      - "sg-222222"
  alb:
    enabled: true
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-application-lb/1234567890abcdef"
    routingType: 1  # path-based
    path: "/myapp"
```

---

# Example with ALB Hostname-Based Routing

Use this if you want to route traffic via a dedicated hostname (e.g., `api.example.com`). Ensure your DNS setup points
to
the ALB’s DNS name.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: hostname-based-ecs-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/my-hostname-cluster"
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
      - "subnet-111aaa"
      - "subnet-222bbb"
    securityGroups:
      - "sg-333ccc"
  alb:
    enabled: true
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-other-lb/abcdef1234567890"
    routingType: 2  # hostname-based
    hostname: "api.example.com"
```

---

# Example with Custom IAM Roles

Here, we show how to specify a custom task execution role for pulling images and a task role for your container to
access AWS services.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: ecs-service-with-iam
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/my-iam-cluster"
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
      - "subnet-xyz123"
    securityGroups:
      - "sg-xyz123"
  iam:
    taskExecutionRoleArn: "arn:aws:iam::123456789012:role/my-custom-ecsTaskExecutionRole"
    taskRoleArn: "arn:aws:iam::123456789012:role/my-app-task-role"
```

---

# Example with Autoscaling

This example demonstrates how to enable autoscaling for an ECS service using target tracking. The service will
automatically scale between 2 and 10 tasks based on CPU utilization, targeting 75% average CPU usage.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: autoscaling-ecs-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/my-production-cluster"
  container:
    image:
      repo: "myorg/production-app"
      tag: "v2.0.0"
    port: 8080
    replicas: 2  # Initial desired count
    cpu: 512
    memory: 1024
  network:
    subnets:
      - "subnet-prod-1"
      - "subnet-prod-2"
    securityGroups:
      - "sg-prod-app"
  alb:
    enabled: true
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/prod-lb/1234567890"
    routingType: "path"
    path: "/api/*"
    listenerPort: 443
  autoscaling:
    enabled: true
    minTasks: 2
    maxTasks: 10
    targetCpuPercent: 75
    targetMemoryPercent: 80  # Optional: also track memory utilization
```

---

# Example with Health Check Grace Period

This example shows how to configure a health check grace period, which is critical for applications that take time
to boot. This prevents ECS from prematurely killing tasks during startup when the ALB health checks initially fail.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: slow-start-ecs-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/my-app-cluster"
  container:
    image:
      repo: "myorg/springboot-app"  # Spring Boot apps often take 30-60s to start
      tag: "latest"
    port: 8080
    replicas: 2
    cpu: 1024
    memory: 2048
  network:
    subnets:
      - "subnet-app-1"
      - "subnet-app-2"
    securityGroups:
      - "sg-app"
  alb:
    enabled: true
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/main-lb/abc123"
    routingType: "hostname"
    hostname: "app.example.com"
    listenerPort: 443
    healthCheck:
      path: "/actuator/health"
      interval: 30
      timeout: 5
      healthyThreshold: 2
      unhealthyThreshold: 3
  # Grace period allows 120 seconds for the app to fully boot before considering health check failures
  healthCheckGracePeriodSeconds: 120
```

---

# Complete Production Example

This example combines all features: autoscaling, health check grace period, custom IAM roles, secrets, and ALB routing.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: production-complete-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/production-cluster"
  container:
    image:
      repo: "myorg/production-api"
      tag: "v3.1.0"
    port: 8080
    replicas: 3
    cpu: 1024
    memory: 2048
    env:
      variables:
        ENVIRONMENT: "production"
        LOG_LEVEL: "info"
      secrets:
        DATABASE_URL: "arn:aws:secretsmanager:us-east-1:123456789012:secret:prod/db-url"
        API_KEY: "arn:aws:secretsmanager:us-east-1:123456789012:secret:prod/api-key"
    logging:
      enabled: true  # Auto-creates CloudWatch Log Group
  network:
    subnets:
      - "subnet-private-1a"
      - "subnet-private-1b"
      - "subnet-private-1c"
    securityGroups:
      - "sg-production-app"
  iam:
    taskExecutionRoleArn: "arn:aws:iam::123456789012:role/ecsTaskExecutionRole"
    taskRoleArn: "arn:aws:iam::123456789012:role/production-app-task-role"
  alb:
    enabled: true
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/production-alb/xyz789"
    routingType: "hostname"
    hostname: "api.production.example.com"
    listenerPort: 443
    listenerPriority: 10
    healthCheck:
      protocol: "HTTP"
      path: "/health"
      port: "traffic-port"
      interval: 30
      timeout: 5
      healthyThreshold: 2
      unhealthyThreshold: 3
  healthCheckGracePeriodSeconds: 90
  autoscaling:
    enabled: true
    minTasks: 3
    maxTasks: 20
    targetCpuPercent: 70
    targetMemoryPercent: 75
```

---

After creating a YAML manifest for your ECS service, apply the configuration using either Pulumi or Terraform with the
OpenMCF CLI. The CLI will validate your manifest against the Protobuf schema, generate the required
infrastructure code, and provision the service on AWS.

For example:

```shell
openmcf pulumi up --manifest ecs-service.yaml --stack myorg/dev
```

Or:

```shell
openmcf terraform apply --manifest ecs-service.yaml --stack myorg/dev
```

You can verify the newly created ECS service in the AWS Console or with the AWS CLI:

```shell
aws ecs list-services --cluster my-hostname-cluster
```

This confirms that your ECS service has been created, is running the desired number of tasks, and is optionally load
balanced if ALB configuration was provided.

---

P.S If you encounter any issues, please ensure that the specified roles, subnets, and security groups already exist in
your AWS account and that your AWS credentials are configured properly in OpenMCF.
