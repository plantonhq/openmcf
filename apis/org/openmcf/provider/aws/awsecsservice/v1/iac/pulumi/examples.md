# Create using CLI

Create a YAML file using one of the examples below. After the YAML is created, use the following command to apply with
OpenMCF (under the hood, you can choose Pulumi or Terraform):

```shell
openmcf pulumi up --manifest <yaml-path> --stack <stack-name>

# or, if you prefer Terraform:

openmcf terraform apply --manifest <yaml-path> --stack <stack-name>
```

If your environment is set up for shorthand, you might also use:

```shell
planton apply -f <yaml-path>
```

---

# Basic Web Service

A simple ECS service running on AWS Fargate, listening on a container port.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: basic-aws-ecs-service
  version:
    message: "First ECS service"
spec:
  cluster_name: my-ecs-cluster
  service_name: basic-web-service
  image: "amazonlinux:2"
  container_port: 80
  desired_count: 1
  cpu: 256
  memory: 512
  subnets:
    - subnet-0abc123
    - subnet-1def456
  security_groups:
    - sg-09876abc
  assign_public_ip: false
```

**Key Points**:

- **Fargate**: Specifies `cpu` and `memory` suitable for a small workload.
- **Private Subnets**: Using subnets typically not exposed to the internet.
- **Security**: Attaches a custom security group to the tasks.

---

# Example with EC2 Launch Type

If your ECS cluster is configured with an EC2 capacity provider, specify an ECS service that runs on EC2.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: ec2-aws-ecs-service
  version:
    message: "Running on EC2"
spec:
  cluster_name: my-ec2-ecs-cluster
  service_name: ec2-app-service
  image: "123456789012.dkr.ecr.us-east-1.amazonaws.com/my-ec2-app:latest"
  container_port: 8080
  desired_count: 2
  cpu: 512
  memory: 1024
  subnets:
    - subnet-0abc123
    - subnet-1def456
  security_groups:
    - sg-01234567
  assign_public_ip: false
  task_execution_role_arn: arn:aws:iam::123456789012:role/ecsTaskExecutionRole
  task_role_arn: arn:aws:iam::123456789012:role/myAppTaskRole
```

**Key Points**:

- **EC2 Launch**: The ECS cluster must already be set up with an EC2 capacity provider.
- **IAM Roles**: Using custom roles for task execution and AWS API access within the container.

---

# Example with Environment Variables

Inject environment variables into your container for configuration or secrets.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: service-with-env
  version:
    message: "Using environment variables"
spec:
  cluster_name: my-ecs-cluster
  service_name: env-service
  image: "123456789012.dkr.ecr.us-east-1.amazonaws.com/myapp:latest"
  container_port: 3000
  desired_count: 2
  cpu: 512
  memory: 1024
  subnets:
    - subnet-11111111
    - subnet-22222222
  security_groups:
    - sg-33333333
  environment:
    - name: "LOG_LEVEL"
      value: "DEBUG"
    - name: "API_KEY"
      value: "some-api-key"
```

**Key Points**:

- **Environment Vars**: Pass sensitive data or config parameters directly to the container.
- **Scaling**: Increase `desired_count` as needed for high availability.

---

# Minimal Example

This minimal spec relies on default values for `desired_count` and `assign_public_ip`. Great for quick POCs or internal
services.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: minimal-service
  version:
    message: "Minimal ECS service example"
spec:
  cluster_name: my-simple-cluster
  service_name: minimal
  image: "amazonlinux:latest"
  cpu: 256
  memory: 512
  subnets:
    - subnet-12345abc
```

**Key Points**:

- **Defaults**: `desired_count` defaults to 1, `assign_public_ip` defaults to false.
- **Private Deployments**: Without a security group or container port, you can run purely internal workloads.

---

# Example with Autoscaling

This example enables target tracking auto scaling, automatically adjusting task count based on CPU utilization.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: autoscaling-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/production-cluster"
  container:
    image:
      repo: "123456789012.dkr.ecr.us-east-1.amazonaws.com/app"
      tag: "v2.0.0"
    port: 8080
    replicas: 2
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
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/prod-lb/xyz"
    routingType: "path"
    path: "/api/*"
    listenerPort: 443
  autoscaling:
    enabled: true
    minTasks: 2
    maxTasks: 10
    targetCpuPercent: 75
```

**Key Points**:

- **Target Tracking**: Automatically scales to maintain 75% average CPU utilization.
- **Bounded Scaling**: Maintains between 2-10 tasks regardless of load.
- **Production-Ready**: Uses AWS Application Auto Scaling with proper cooldown periods.

---

# Example with Health Check Grace Period

This example demonstrates configuring a health check grace period for applications with slow startup times.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: springboot-service
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/app-cluster"
  container:
    image:
      repo: "123456789012.dkr.ecr.us-east-1.amazonaws.com/springboot-app"
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
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/main-lb/abc"
    routingType: "hostname"
    hostname: "app.example.com"
    listenerPort: 443
    healthCheck:
      path: "/actuator/health"
      interval: 30
  healthCheckGracePeriodSeconds: 120
```

**Key Points**:

- **Grace Period**: Allows 120 seconds for application to fully boot.
- **Prevents Failed Deployments**: ECS ignores ALB health check failures during startup.
- **Essential for JVM Apps**: Critical for Spring Boot, Quarkus, and other slow-starting applications.

---

# Complete Production Example

This comprehensive example combines all production features: autoscaling, health check grace period, secrets, and ALB integration.

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsEcsService
metadata:
  name: production-api
spec:
  clusterArn: "arn:aws:ecs:us-east-1:123456789012:cluster/production"
  container:
    image:
      repo: "123456789012.dkr.ecr.us-east-1.amazonaws.com/production-api"
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
      enabled: true
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
    arn: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/production-alb/xyz"
    routingType: "hostname"
    hostname: "api.production.example.com"
    listenerPort: 443
    listenerPriority: 10
    healthCheck:
      protocol: "HTTP"
      path: "/health"
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

**Key Points**:

- **Full Production Stack**: Combines all features for production-ready deployments.
- **Dual Scaling Metrics**: Scales on both CPU (70%) and memory (75%) utilization.
- **Secure Secrets**: Uses AWS Secrets Manager for sensitive configuration.
- **High Availability**: Spans multiple AZs with 3-20 task range.

---

**Next Steps**:

- Customize CPU/memory allocations, environment variables, or IAM roles based on application requirements.
- Refer to the [README.md](./README.md) for additional information on the ECS service resource fields and how to
  configure them for production workloads.
- Check out OpenMCF’s official documentation to explore advanced features like load balancer integration, auto
  scaling policies, and multi-environment workflows.
