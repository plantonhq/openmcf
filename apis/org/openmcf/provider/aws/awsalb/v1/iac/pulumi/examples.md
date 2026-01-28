```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: example
spec: {}
```

CLI:

```bash
openmcf pulumi preview \
  --manifest ../hack/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .

openmcf pulumi update \
  --manifest ../hack/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes
```

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
kind: AwsAlb
metadata:
  name: basic-aws-alb
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
kind: AwsAlb
metadata:
  name: ec2-aws-alb
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
kind: AwsAlb
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
kind: AwsAlb
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

**Next Steps**:

- Customize CPU/memory allocations, environment variables, or IAM roles based on application requirements.
- Refer to the [README.md](./README.md) for additional information on the ECS service resource fields and how to
  configure them for production workloads.
- Check out OpenMCF’s official documentation to explore advanced features like load balancer integration, auto
  scaling policies, and multi-environment workflows.
