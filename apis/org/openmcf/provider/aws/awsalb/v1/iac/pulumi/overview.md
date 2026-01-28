# Overview

The **AwsAlb** API resource simplifies the deployment of containerized applications on AWS ECS, whether running on
Fargate or EC2. By defining a single-container service with essential resource requirements, networking, and optional
IAM
roles, **AwsAlb** significantly reduces the complexity of managing individual ECS tasks, load balancers, and scaling
policies. It is part of OpenMCF’s broader multi-cloud deployment framework, supporting both Pulumi and Terraform
under the hood to fit seamlessly into your existing infrastructure-as-code workflows.

With a single YAML manifest conforming to the familiar Kubernetes-like structure (`apiVersion`, `kind`, `metadata`,
`spec`, `status`), teams can validate, provision, and maintain ECS services consistently across multiple environments.
Whether you’re running a simple background process or a web application behind a load balancer, **AwsAlb** handles
the heavy lifting for you—allocating CPU, memory, subnets, security groups, and optional IAM roles, all while adhering
to
best practices for high availability and security.

---

## Key Features

- **Fargate & EC2 Deployment**  
  Easily configure services to run on Fargate’s serverless compute or on self-managed EC2 instances, accommodating both
  cost and flexibility requirements.

- **Load Balancing & Networking**  
  Dynamically attach containers to an AWS ALB or NLB for inbound requests, including secure network isolation with
  custom
  VPC subnets and security groups.

- **Scalability & Monitoring**  
  Define your desired task count, or integrate with AWS Auto Scaling for dynamic scaling. Leverage built-in ECS metrics
  for performance monitoring and CloudWatch logging.

- **Consistent Resource Model**  
  Uses OpenMCF’s standard resource layout, ensuring familiarity and robust validations via Protobuf definitions.

- **Pulumi & Terraform Integration**  
  Provision the same ECS service specification using either Pulumi or Terraform, unified by the OpenMCF CLI’s
  straightforward commands and orchestration.

- **Environment & IAM Configuration**  
  Inject environment variables and link IAM roles for fine-grained AWS API access, simplifying secret management and
  compliance in enterprise settings.

---

## Next Steps

- Refer to the [README.md](./README.md) for detailed setup instructions, resource configuration fields, and best
  practices.
- Review the [examples.md](./examples.md) to explore common ECS service use cases, from simple web apps to production
  load-balanced services.
- Check out the wider OpenMCF documentation for deeper insights into multi-cloud deployments, advanced features,
  and the CLI usage patterns.
