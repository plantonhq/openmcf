# AWS VPC Deployment: From Console Clicks to Production-Grade Infrastructure

## Introduction

Here's a question that stumps many cloud newcomers: *"Should I use AWS's default VPC or create my own?"* The answer reveals a deeper truth about cloud networking—the path of least resistance rarely leads to production-ready infrastructure.

AWS gives every account a default VPC in each region. It's pre-configured with public subnets, an Internet Gateway, and auto-assigned public IPs. You can launch an EC2 instance and have it on the internet in minutes. Convenient? Absolutely. Production-ready? Not quite.

The default VPC uses the same IP range (172.31.0.0/16) in every AWS account. This means you can't peer it with other default VPCs without IP conflicts. All subnets are public by default, exposing every resource to the internet unless you explicitly lock it down. There's no network segmentation, no private tier for databases, and no consideration for multi-AZ resilience patterns.

This is where custom VPCs enter the picture—not as a theoretical "best practice," but as a practical necessity for any application that will carry production traffic. The question isn't whether to create a custom VPC, but *how* to deploy one reliably and consistently across environments.

This document explores the evolution of VPC deployment methods, from manual console work to production-grade Infrastructure as Code, and explains why OpenMCF embraces specific approaches for its AWS VPC provisioning.

## The Maturity Spectrum: How Teams Deploy VPCs

### Level 0: The Console Wizard (Learn, Don't Stay)

AWS's VPC creation wizard in the web console is where most cloud journeys begin. You can choose between "VPC Only" (manual mode where you specify just the CIDR block) or "VPC and More" (a guided wizard that creates subnets, route tables, and gateways in one shot).

The wizard is excellent for learning. It makes the relationships between components tangible—you see the Internet Gateway attach to the VPC, watch route tables associate with subnets, and understand why a NAT Gateway must live in a public subnet.

**Why it's not production-ready:**

The console creates *one-off configurations* that exist only in memory and clicks. There's no record of what you built, no way to replicate it exactly in another region or account, and no collaboration workflow. Common mistakes slip through: forgetting to enable DNS hostnames, placing a NAT Gateway in a private subnet (it won't work), or omitting the default route to the Internet Gateway.

AWS explicitly recommends against relying on the console for repeatable environments. As one AWS guide notes, rather than clicking through the same steps in dev, stage, and prod, teams should "create new VPCs using best practices—ideally via Infrastructure as Code for consistency."

**Verdict:** Use the console wizard to understand VPC components, then graduate immediately to code-based approaches for any environment that matters.

### Level 1: Scripts and CLIs (Automation Without Abstraction)

The next evolution is scripting VPC creation with the AWS CLI or SDKs. A bash script might call:

```bash
aws ec2 create-vpc --cidr-block 10.0.0.0/16
aws ec2 create-subnet --vpc-id $VPC_ID --cidr-block 10.0.1.0/24 --availability-zone us-west-2a
aws ec2 create-internet-gateway
aws ec2 attach-internet-gateway --vpc-id $VPC_ID --internet-gateway-id $IGW_ID
# ... dozens more lines
```

Similarly, Python's Boto3 or other language SDKs let you automate resource creation programmatically.

**The improvement:** These scripts are version-controlled and repeatable. You can commit them to a repository, review changes, and run them in CI/CD pipelines.

**The limitation:** You're managing dependencies manually. Did the VPC finish creating before you tried to create subnets? Are you handling errors and retries? Running the script twice might create duplicate resources unless you add explicit idempotency checks. You're essentially building your own state management layer on top of AWS APIs.

Configuration management tools like **Ansible** represent a step up here. Ansible's AWS modules (`amazon.aws.ec2_vpc_net`, `amazon.aws.ec2_vpc_subnet`) provide declarative resource definitions with built-in idempotency. Run the playbook multiple times, and Ansible ensures the desired state without creating duplicates.

**Verdict:** CLI scripts work for quick automation or learning, but they lack the structural rigor needed for managing complex, multi-environment VPC deployments. Ansible bridges the gap but is still more procedural than truly declarative Infrastructure as Code.

### Level 2: Declarative Infrastructure as Code (The Turning Point)

This is where VPC deployment becomes genuinely production-ready. Declarative IaC tools let you define *what* infrastructure you want, not *how* to create it. The tool handles dependencies, state management, and convergence to the desired configuration.

The major players here are **Terraform**, **AWS CloudFormation**, **AWS CDK**, **Pulumi**, and (more recently) **OpenTofu** (an open-source Terraform fork).

#### Terraform

Terraform uses HashiCorp Configuration Language (HCL) to define resources. A basic VPC in Terraform looks like:

```hcl
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = { Name = "production-vpc" }
}

resource "aws_subnet" "public" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-west-2a"
}
```

Terraform maintains a **state file** that tracks what's deployed. When you run `terraform apply`, it compares the desired state (your code) with actual state (what's in AWS) and makes only the necessary changes. This is a game-changer for managing evolving infrastructure.

**Strengths:**
- Cloud-agnostic (can manage AWS, GCP, Azure with the same tool)
- Massive module ecosystem (the `terraform-aws-modules/vpc/aws` module is a community standard)
- Explicit state management with locking via S3 + DynamoDB
- Widely adopted with strong community support

**Considerations:**
- Requires managing state files (lost state = difficult recovery)
- HCL has limited abstraction compared to general-purpose languages
- Modularization requires discipline to avoid monolithic configurations

#### AWS CloudFormation

CloudFormation is AWS's native IaC service. You define resources in YAML or JSON templates, and AWS provisions and tracks them as a "stack."

```yaml
Resources:
  MyVPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: 10.0.0.0/16
      EnableDnsHostnames: true
```

**Strengths:**
- No external tools or state management—AWS handles everything
- Native AWS support and integration
- Stack drift detection catches manual changes

**Considerations:**
- Extremely verbose for complex VPCs (a three-tier, multi-AZ VPC can be 300+ lines)
- Limited abstraction (no native loops or functions, though Macros and Nested Stacks help)
- AWS-only (not suitable for multi-cloud)
- 200-500 resource limit per stack can force splitting large deployments

#### AWS CDK (Cloud Development Kit)

CDK lets you define infrastructure in TypeScript, Python, Java, or other languages, which it then synthesizes into CloudFormation templates.

A VPC in CDK can be as simple as:

```typescript
import * as ec2 from 'aws-cdk-lib/aws-ec2';

const vpc = new ec2.Vpc(this, 'ProductionVPC', {
  maxAzs: 3,
  natGateways: 3
});
```

This single construct creates a VPC with public and private subnets across three availability zones, complete with Internet Gateway, NAT Gateways, and route tables configured according to AWS best practices.

**Strengths:**
- High-level abstractions (one line can create a dozen resources)
- Use standard programming languages (loops, functions, testing frameworks)
- Leverages CloudFormation's state management (no separate state file)
- Strong AWS integration and rapid feature updates

**Considerations:**
- Tied to CloudFormation (subject to its limits)
- Relatively new compared to Terraform (smaller community)
- AWS-only

#### Pulumi

Pulumi is similar to CDK but uses real languages (TypeScript, Python, Go) and provisions resources directly via AWS SDKs or Terraform providers, not CloudFormation.

```typescript
import * as aws from "@pulumi/aws";

const vpc = new aws.ec2.Vpc("production-vpc", {
  cidrBlock: "10.0.0.0/16",
  enableDnsHostnames: true
});
```

Pulumi's AWSX library offers higher-level abstractions comparable to CDK.

**Strengths:**
- Multi-cloud with real programming languages
- Can use Terraform providers (broad resource coverage)
- Managed state backend or self-hosted options

**Considerations:**
- Default state backend is Pulumi's SaaS (some enterprises prefer self-hosted)
- Smaller ecosystem than Terraform

#### OpenTofu

OpenTofu is a community-driven fork of Terraform created after Terraform's license change to BSL (Business Source License). It maintains compatibility with Terraform's HCL syntax and provider ecosystem.

For teams committed to open-source tooling, OpenTofu provides the same functionality as Terraform under an MPL 2.0 license. From a VPC deployment perspective, OpenTofu and Terraform are interchangeable—modules, configurations, and workflows remain identical.

**Verdict on IaC Tools:** All these tools achieve the same end goal—a well-configured VPC. The choice depends on:
- **Multi-cloud needs?** → Terraform or Pulumi
- **Pure AWS with high-level abstractions?** → CDK
- **Organizational momentum and expertise?** → Likely Terraform (most widely adopted)
- **Open-source commitment?** → OpenTofu

### Level 3: The Community Module Standard (Production Proven)

While you can write raw Terraform resources for every subnet, route table, and NAT Gateway, the community has already solved this problem.

The **Terraform AWS VPC Module** (`terraform-aws-modules/vpc/aws`) is a battle-tested, production-grade solution used by thousands of organizations. With approximately 10 lines of configuration:

```hcl
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  
  name = "production-vpc"
  cidr = "10.0.0.0/16"
  
  azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  
  enable_nat_gateway = true
  enable_dns_hostnames = true
  
  tags = {
    Environment = "production"
  }
}
```

This creates:
- A VPC with the specified CIDR block
- Three public subnets (one per AZ) with routes to an Internet Gateway
- Three private subnets (one per AZ) with routes to NAT Gateways
- One NAT Gateway per AZ for high availability
- Proper route table associations
- DNS settings configured correctly

**Why this matters:**

1. **Best practices by default:** Multi-AZ deployment, separate route tables per subnet type, one NAT Gateway per AZ (avoiding single points of failure)

2. **Flexible NAT strategies:** 
   - Default: One NAT per AZ (resilient, recommended for production)
   - `single_nat_gateway = true`: One NAT total (cost optimization for dev/test)
   - `one_nat_gateway_per_az = true`: Explicitly one per AZ regardless of subnet count

3. **Cost transparency:** The module makes trade-offs explicit. Want to save ~$64/month in dev? Use a single NAT Gateway. Need production resilience? Multi-AZ NAT is the default.

4. **Production features built-in:** VPC Flow Logs support, VPC Endpoints for S3/DynamoDB, VPN Gateway attachment, customizable NACLs, and more

The module embodies the **80/20 principle**: 20% of configuration options cover 80% of use cases. Common needs are simple toggles; complex scenarios are still achievable through additional parameters.

**Verdict:** Using proven community modules isn't "taking shortcuts"—it's standing on the shoulders of giants. Thousands of production deployments have battle-tested these patterns.

## What OpenMCF Supports (And Why)

OpenMCF's AWS VPC provisioning embraces the principles validated by the community module and production deployments worldwide.

### The API Design Philosophy

OpenMCF's VPC API exposes a minimal, focused set of inputs:

- **`vpc_cidr`**: The IP address range (required)
- **`availability_zones`**: Which AZs to span
- **`subnets_per_availability_zone`**: How many subnets per AZ
- **`subnet_size`**: The number of hosts per subnet
- **`is_nat_gateway_enabled`**: Enable NAT for private subnet internet access
- **`is_dns_hostnames_enabled`**: Enable public DNS names for instances
- **`is_dns_support_enabled`**: Enable AWS-provided DNS resolution

This isn't arbitrary minimalism—it reflects what matters for 80% of VPC deployments.

### What OpenMCF Handles For You

**Multi-AZ by Default:** OpenMCF creates subnets across all specified availability zones. This isn't optional because single-AZ deployments are an anti-pattern for any environment beyond quick experiments.

**NAT Gateway Best Practices:** When `is_nat_gateway_enabled = true`, OpenMCF deploys one NAT Gateway per availability zone. This follows AWS's explicit recommendation: "NAT Gateways within an AZ are redundant within that AZ, but if the entire AZ fails, instances in other AZs can't reach it. We recommend using at least one NAT Gateway per AZ."

This costs more than a single NAT (approximately $32/month per AZ), but it eliminates cross-AZ data transfer fees (which would partially offset the savings) and prevents a single AZ failure from cutting off internet access for the entire VPC.

**Route Table Organization:** Public subnets get a route table with a default route to the Internet Gateway. Private subnets get per-AZ route tables with routes to the local NAT Gateway. This ensures traffic from private subnets in AZ-A routes to NAT-in-AZ-A, not across availability zones.

**DNS Configuration:** Most VPCs need DNS hostnames enabled (so instances get friendly DNS names), but OpenMCF makes it a toggle rather than assuming. DNS support is virtually always enabled (it's required for basic AWS service resolution), but the API exposes it for completeness.

### What's Intentionally Left Out

**Secondary CIDR Blocks:** While AWS allows adding additional IP ranges to a VPC, this is relatively rare (typically only when the original range is exhausted). OpenMCF focuses on getting the initial CIDR sizing right.

**Custom DHCP Option Sets:** Most deployments use AWS's default DNS. Custom DHCP options (for on-premises DNS integration or custom domain names) are advanced scenarios that can be handled separately.

**VPN and Direct Connect:** While important for hybrid connectivity, these are separate concerns. A VPC is the foundation; hybrid networking is layered on top.

**Custom Network ACLs:** Security Groups handle most access control needs. NACLs are available for subnet-level filtering but are often left at default (allow-all) with security enforced via Security Groups.

### The Underlying Implementation

OpenMCF provides both **Terraform** and **Pulumi** implementations, ensuring teams can use their preferred IaC tool.

The Terraform implementation aligns with the patterns from `terraform-aws-modules/vpc/aws`, providing consistency with community standards. The Pulumi implementation offers equivalent functionality using Pulumi's AWS library.

Both implementations:
- Calculate subnet CIDRs automatically based on the VPC CIDR and subnet size requirements
- Create route tables with correct associations
- Handle NAT Gateway Elastic IP allocation
- Tag all resources consistently
- Output subnet IDs, NAT Gateway IPs, and other values for downstream use

## Cost Considerations: The NAT Gateway Question

The most significant cost variable in VPC design is NAT Gateway placement.

### The Numbers

A NAT Gateway costs approximately:
- **$0.045/hour** (~$32/month) per gateway
- **$0.045/GB** for data processed

For a three-AZ VPC with NAT per AZ:
- **Base cost:** ~$96/month (3 × $32)
- **Data processing:** $0.045/GB of outbound traffic

### The Alternatives

**Single NAT Gateway (Dev/Test Pattern):**
- Cost: ~$32/month + data processing
- Trade-off: AZ failure takes down internet access for the entire VPC
- Additional cost: Cross-AZ data transfer (~$0.01/GB) for traffic from other AZs to the NAT

**NAT Instance (Budget Optimization):**
- Cost: EC2 instance hours (e.g., t3.small at ~$15/month) + standard data transfer
- Trade-off: You manage the instance, patching, and failover. No auto-scaling. Limited bandwidth.
- When it makes sense: Very low-traffic dev environments where you can tolerate manual failover

**VPC Endpoints (Cost Reduction):**
- **S3 and DynamoDB Gateway Endpoints:** FREE (no hourly charge, no data transfer cost)
- **PrivateLink Interface Endpoints:** ~$7/month per endpoint per AZ + $0.01/GB

The strategic play: Enable S3 and DynamoDB Gateway Endpoints in every VPC. They cost nothing and bypass NAT entirely for those services. Evaluate Interface Endpoints for heavily-used services (CloudWatch, Secrets Manager, etc.) where the ~$7/month/AZ pays for itself by avoiding NAT data processing fees.

### OpenMCF's Approach

OpenMCF makes NAT a toggle (`is_nat_gateway_enabled`) with clear defaults:
- **Production:** Enable NAT, deploy one per AZ (resilience over cost)
- **Development:** Enable NAT, consider single-NAT mode if cost is critical and you accept the availability trade-off

The API doesn't expose "how many NAT Gateways" because the answer is almost always "one per AZ if you care about availability, one total if you don't." Hiding that complexity makes the right choice (multi-AZ) the easy choice.

## The Path Forward

AWS VPC deployment has evolved from manual console work through scripted automation to declarative Infrastructure as Code, culminating in proven community modules and patterns.

OpenMCF doesn't reinvent this evolution—it distills the lessons into an opinionated API that makes best practices the default path:

- **Multi-AZ** is assumed, not optional
- **NAT per AZ** is the default for production resilience
- **80% of configuration** is handled automatically
- **The 20% you must specify** (CIDR, AZs, subnet sizing) remains explicit and clear

This isn't about removing flexibility—the underlying Terraform and Pulumi implementations remain accessible and extensible. It's about recognizing that most VPCs should follow similar patterns, and automating those patterns shouldn't require deep AWS networking expertise.

For teams deploying their first VPC or their hundredth, OpenMCF provides a path from zero to production-grade networking that embodies years of community knowledge and AWS best practices—without requiring you to relive every mistake that led to those practices.

The console wizard taught us how VPCs work. Infrastructure as Code taught us how to deploy them consistently. Community modules taught us what "production-ready" actually means. OpenMCF brings all of that together into an API that just works.

