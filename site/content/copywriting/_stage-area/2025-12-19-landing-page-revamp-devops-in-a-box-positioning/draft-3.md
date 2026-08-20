# Planton - New Landing Page Copy (Draft 2)

**Date:** December 19, 2025  
**Revision:** Draft 2 - Incorporating customer corrections and positioning updates

---

## Section 1: Hero Section

### Headline
**DevOps-in-a-Box—From Code to Cloud in Under an Hour**

### Subheadline
Multi-cloud automation platform that eliminates DevOps bottlenecks. Deploy infrastructure and backend services in minutes—no scripts, no waiting, no dedicated ops team required.

### Social Proof Badge
Trusted by 5+ IT consulting firms | 100% customer retention

### Primary CTA
**Start Free Trial**
(100 automation minutes, no credit card required)

### Secondary CTA
**Book Demo**
(See actual customer walkthrough)

### Trust Indicators
- **Multi-Cloud:** AWS • GCP • Azure • Digital Ocean • Cloudflare
- **Open Source Foundation** (OpenMCF)
- **Zero Vendor Lock-In**

### Visual Direction
Split-screen hero image showing:
- **Left:** Infra Hub dashboard with live Terraform deployment progress
- **Right:** Service Hub pipeline executing (Git commit → build → deploy)

---

## Section 2: The Problem/Solution

### Section Title
**The DevOps Bottleneck Is Killing Your Velocity**

### Section Subtitle
Whether you're a consulting firm, startup, or established team—infrastructure delays slow everyone down. Planton eliminates the wait.

### Three-Column Layout

#### Column 1: IT Consulting Firms

**Problem:**
New client projects delayed by infrastructure setup

**Impact:**
Non-billable days before developers can iterate on actual product features

**Solution with Planton:**
- ✅ <1 hour environment deployment
- ✅ Developers self-onboard services
- ✅ Multiple environments (dev, QA, staging, prod) in minutes

**Metric:** Deploy infrastructure 5x faster than manual cloud console work

---

#### Column 2: Startups & Small Product Companies

**Problem:**
Can't afford dedicated DevOps engineers

**Impact:**
$150K+/year overhead for basic operations—or bottlenecked on technical founders

**Solution with Planton:**
- ✅ $20/dev/month + usage pricing
- ✅ Platform replaces ops team for core workflows
- ✅ Self-service infrastructure and CI/CD

**Metric:** Full DevOps capabilities at 1/10th the cost of hiring

---

#### Column 3: Established Teams

**Problem:**
Snowflake infrastructure, compliance nightmares, visibility gaps

**Impact:**
Audit failures, security incidents, hours wasted troubleshooting "who changed what"

**Solution with Planton:**
- ✅ Standardized, repeatable deployments
- ✅ Full change tracking and audit logs
- ✅ Compliance-ready security controls

**Metric:** 100% visibility into every infrastructure change

---

## Section 3: Infra Hub

### Section Title
**Infra Hub: Infrastructure Deployment Without the Wait**

### Subtitle
Replaces Terraform Enterprise, Pulumi Cloud, and manual cloud console work

### Introduction Paragraph
Infra Hub turns infrastructure deployment from a multi-week manual process into a sub-hour sprint. Deploy individual components or complete environments with zero-configuration CI/CD, live progress visualization, and built-in best practices.

---

### Feature 1: Zero-Config CI/CD

**Headline:** Infrastructure Pipelines That Just Work

**Body Copy:**
No GitHub Actions to write. No Terraform state buckets to configure. No manual deployments through cloud consoles. Connect your cloud account and start deploying infrastructure through simple forms or CLI commands.

**Key Benefits:**
- **Live Terraform plan/apply visualization** - Watch your infrastructure come to life in real-time
- **5-10 minute deployments** vs weeks of manual cloud console work
- **Version history and rollback** - Every change tracked, easily reversible
- **Multi-cloud consistency** - Same workflow whether you're on AWS, GCP, or Azure

---

### Feature 2: Infra Charts: Complete Environments in One Form

**Headline:** Stop Assembling Lego Blocks—Deploy Complete Environments

**Body Copy:**
Deploying infrastructure piece-by-piece is tedious and error-prone. VPC first, then Route53, then ECR, then certificates, then load balancers... each waiting for the previous to complete. Infra Charts eliminate this pain.

**How It Works:**
Inspired by Kubernetes Helm Charts, Infra Charts bundle related infrastructure components into a single deployable unit. Fill out one form, click deploy, and the platform orchestrates everything in the correct order via DAG (Directed Acyclic Graph).

**Example: AWS ECS Environment**
One Infra Chart deploys:
- VPC with proper subnets and routing
- Route53 DNS zone for your domain
- ECR container registry for images
- SSL/TLS certificates via Certificate Manager
- Application Load Balancer with listeners
- Security groups with proper ingress/egress rules

**The Platform Knows:**
- VPC must exist before load balancer (load balancers require VPC ID)
- Certificates must be validated before HTTPS listeners
- Security groups must allow traffic to backend services

**What This Means:**
You describe what you want. The platform figures out the how.

---

### Feature 3: Deployment Component Store

**Headline:** 20% of Cloud Services Used 80% of the Time—All Pre-Modeled

**Body Copy:**
Stop reinventing the wheel. Browse our carefully crafted deployment components for the most common cloud resources:

**AWS:** ECS, EKS, Lambda, RDS, S3, CloudFront, Route53, ALB, VPC, and more
**GCP:** Cloud Run, GKE, Cloud SQL, Cloud Storage, Cloud CDN, and more  
**Azure:** AKS, Azure Functions, Azure SQL, Storage Accounts, and more

**All Open Source:**
Every deployment component lives in the [OpenMCF repository](https://github.com/plantonhq/openmcf). Audit the exact Terraform and Pulumi code. Contribute improvements. No black boxes.

**Visual Direction:** Screenshot of Deployment Component Store browser showing categorized components by cloud provider

---

### Customer Proof Point

**Real Customer Result:**
"iorta TechNext deployed a complete AWS ECS environment for their SalesVerse product in under 1 hour—with 7 developers managing everything themselves. Previously, this would have required a dedicated DevOps engineer and taken at least a couple weeks of work."

---

## Section 4: Service Hub

### Section Title
**Service Hub: Vercel for Backend—In YOUR Cloud**

### Subtitle
Connect Git repo → Automatic pipelines → Multi-cloud deployment

### Introduction Paragraph
Vercel nailed the developer experience for frontend deployments. Service Hub brings that same magic to backend services—with one critical difference: everything runs in YOUR cloud account. You keep full infrastructure control and visibility.

---

### Feature 1: Git-to-Cloud in Minutes

**Headline:** From Repository to Running Service—Zero Configuration Required

**Body Copy:**
Connect your GitHub or GitLab repository. That's it. Every new commit automatically triggers a pipeline that builds your container image, pushes it to your registry, and deploys it to your target environment.

**No More:**
- ❌ Writing Dockerfiles (unless you want control—then it's optional)
- ❌ Configuring GitHub Actions or GitLab CI
- ❌ Managing container registries manually
- ❌ Orchestrating deployments across environments

**You Get:**
- ✅ Automatic pipeline creation on every commit
- ✅ Built-in Docker image building and registry pushes
- ✅ Live logs and deployment progress
- ✅ Multi-environment support (dev, QA, staging, prod)

---

### Feature 2: Multi-Cloud Consistency

**Headline:** Same Developer Experience—Any Cloud, Any Target

**Body Copy:**
Deploy to AWS ECS, AWS EKS, GCP Cloud Run, Azure Container Instances, or Cloudflare Workers. The Service Hub experience stays identical. No vendor lock-in. No relearning workflows when you switch clouds.

**Supported Targets:**
- **AWS ECS** - Container services without Kubernetes complexity
- **AWS EKS** - Full Kubernetes power when you need it
- **GCP Cloud Run** - Serverless containers that scale to zero
- **Azure Cloud Functions** - Event-driven serverless compute
- **Cloudflare Workers** - Edge compute for global performance

**Why This Matters:**
Start on AWS ECS for simplicity. Migrate to GCP Cloud Run for cost optimization. Deploy to Cloudflare Workers for edge use cases. Same code, same platform, zero refactoring.

---

### Feature 3: Pipeline Workflow

**Headline:** From Commit to Deployment—Fully Automated

**Body Copy:**
Every commit triggers an automated pipeline that handles the entire journey from code to running service.

**Pipeline Flow:**
- Every commit triggers a new pipeline run
- Live logs stream in real-time
- Deployment status visible at a glance
- Failed builds? Logs tell you exactly why
- Rollback to previous version with one click

**Traditional Alternative:**
Setting up this same workflow manually with Jenkins, GitHub Actions, or GitLab Pipelines typically takes several days to weeks of DevOps work.

---

### Customer Quote

**"Developers self-onboard their own services. That's a huge win."**
— Tynybay consulting, deploying the Odwen warehousing platform

**Context:**
3-person team (1 DevOps engineer, 2 developers). Infrastructure setup: <1 hour. Developers now onboard services independently without blocking on DevOps. Multiple environments (QA, staging, prod) managed through simple UI.

---

### Visual Direction
Animated flow showing:
1. Developer pushes code to GitHub
2. Service Hub automatically detects commit
3. Tekton pipeline builds container image
4. Image pushed to registry (ECR/GCR/etc)
5. Deployment to target (ECS/Cloud Run/etc)
6. Live service with URL

---

## Section 5: Built on Open Standards

### Section Title
**Built on Open Standards: Works With Your Existing Tooling**

### Subtitle
GitHub, Terraform, Pulumi, Tekton—Planton integrates with what you already use, no proprietary replacements required.

### Introduction Paragraph
Planton is built on open standards and industry-leading tools. Leverage your existing workflows, scripts, and expertise without being forced into a proprietary ecosystem.

---

### Standard 1: GitHub & GitLab Native

**Headline:** Works Seamlessly With Your Git Workflow

**Body Copy:**
Connect your GitHub or GitLab organization via OAuth. No need to migrate repositories or change your branching strategy. Planton integrates with your existing Git workflows—pull requests, branch protection, code reviews—everything stays the same.

**What's Supported:**
- GitHub repositories (public and private)
- GitLab repositories (self-hosted and GitLab.com)
- GitHub Actions (if you want to keep using them)
- GitLab CI (side-by-side with Planton pipelines)

**Why This Matters:**
Your code stays where it is. Your team keeps working the way they're used to. Planton adds automation without disrupting your established processes.

---

### Standard 2: Terraform & Pulumi First-Class

**Headline:** Your Infrastructure-as-Code, Your Choice

**Body Copy:**
Planton supports both Terraform and Pulumi as first-class infrastructure-as-code frameworks. All deployment components are available in both formats, and you can mix and match based on your team's preferences.

**What's Open Source:**
- Every Terraform module in the [OpenMCF repository](https://github.com/plantonhq/openmcf)
- Every Pulumi component (TypeScript, Python, Go)
- CLI tools for local development and testing
- Schema definitions and validation rules

**Integration Points:**
- Use Planton UI for simple deployments
- Use Terraform/Pulumi CLI directly for complex orchestration
- Import existing Terraform state
- Export Planton configurations as standard Terraform/Pulumi code

---

### Standard 3: Tekton-Powered CI/CD

**Headline:** Open-Source Kubernetes-Native Pipelines

**Body Copy:**
Service Hub is powered by Tekton, the CNCF graduated project for Kubernetes-native CI/CD. Planton provides transparent access to customize and extend your pipelines—no proprietary abstractions hiding the underlying technology.

**Why Tekton:**
- **Cloud-native** - Runs on Kubernetes, scales automatically
- **Portable** - Your pipelines work anywhere Kubernetes runs
- **Extensible** - Rich catalog of community tasks and integrations
- **Battle-tested** - Used by Google Cloud Build, Red Hat OpenShift Pipelines

**What You Can Do:**
- Use Planton's zero-config pipelines for common workflows
- Write custom Tekton Tasks for specialized build steps
- Integrate with Tekton Catalog for pre-built components
- Export and run pipelines outside Planton platform

**Transparency:**
Planton publishes all Tekton pipeline definitions openly. Audit, fork, and customize them. No black boxes, no vendor lock-in.

---

### Why Open Standards Matter

**w/ Vendor Lock-In Risk:**
- Proprietary pipeline formats
- Custom IaC languages  
- Migration costs in engineer-months
- Closed-source implementations

**w/ Planton:**
- Standard Terraform/Pulumi/Tekton workflows
- All code public and auditable
- Migration costs measured in hours
- Works with your existing tools (kubectl, terraform, git)

**The Promise:**
Use Planton because it accelerates your team—not because you're locked in. Every abstraction Planton provides has a standard alternative underneath.

---

### Visual Direction
- Diagram showing Planton integrations: GitHub ↔ Planton ↔ Terraform/Pulumi ↔ Tekton
- Side-by-side: "w/ Vendor Lock-In Risk" vs "w/ Planton" comparison table
- Screenshot of Tekton pipeline definition (showing transparency)

---

## Section 6: Open Source Foundation

### Section Title
**Open Source Foundation: Your Infrastructure, Your Code, Your Freedom**

### Subtitle
OpenMCF - All deployment components publicly available. No vendor lock-in, ever.

### Introduction Paragraph
Most DevOps platforms are black boxes. You trust them with your infrastructure, but you can't audit how they work. You're locked in—switching costs are prohibitive. Planton is different. Everything is transparent. Everything is portable. Your infrastructure, your rules.

---

### Pillar 1: Transparent Infrastructure

**Headline:** Audit Every Line of Infrastructure Code

**Body Copy:**
Every Terraform module and Pulumi component used by Planton lives in the public [OpenMCF repository](https://github.com/plantonhq/openmcf). Want to know exactly what IAM permissions the AWS ALB deployment requires? Read the code. Concerned about security groups? Audit the Terraform files yourself.

**What's Open Source:**
- All Terraform modules for every deployment component
- All Pulumi implementations (alternative IaC framework)
- CLI tools for local infrastructure management
- Deployment component schemas and validation rules

**Why This Matters:**
- **Security audits** - Your team or third-party auditors can verify everything
- **Compliance requirements** - Prove to regulators exactly what's deployed and how
- **Learning and contribution** - Understand best practices, contribute improvements
- **No surprises** - No hidden infrastructure, no mystery permissions

---

### Pillar 2: Built-In Exit Strategy

**Headline:** Outgrow Planton? Take Your Infrastructure With You

**Body Copy:**
Unlike proprietary platforms that trap you in their ecosystem, Planton is designed for portability. If you outgrow the platform, get acquired, or decide to hire an in-house DevOps team—your infrastructure investment comes with you.

**How It Works:**

**1. Export Your Configurations**
One API call exports all your infrastructure configurations as YAML manifests:
- All deployment component specs
- All environment configurations
- All service definitions
- Complete infrastructure-as-code

**2. Use the Open Source CLI**
The OpenMCF CLI is independent of the proprietary platform. Install it anywhere:
```bash
brew install openmcf/tap/openmcf
```

**3. Manage Independently**
Continue managing your infrastructure using the same YAML manifests:
- Run `openmcf pulumi up` directly from your CI/CD
- Integrate with GitHub Actions, GitLab CI, or Jenkins
- No platform dependency, no subscription required

**4. Gradual or Immediate Transition**
- **Overlap period:** Run Planton platform AND independent CLI simultaneously
- **Zero downtime:** Transition infrastructure management without production disruption
- **Support available:** Planton helps you migrate successfully (Build, Operate, Transfer model)

---

### Pillar 3: Community & Contributions

**Headline:** Built By DevOps Engineers, For DevOps Engineers

**Body Copy:**
OpenMCF isn't just open source—it's community-driven. Found a bug in a deployment component? Submit a PR. Need support for a new AWS service? Contribute the Terraform module. Want a better way to handle multi-region deployments? Open an issue, propose the design.

**Contribution Areas:**
- New deployment components for emerging cloud services
- Improved Terraform/Pulumi implementations
- Multi-cloud abstractions and compatibility
- Security and compliance enhancements
- Documentation and examples

---

### Why This Matters: The Vendor Lock-In Problem

**w/ Vendor Lock-In Risk:**
- Proprietary infrastructure code
- Custom workflows
- Migration costs in engineer-months
- Closed-source security models

**w/ Planton:**
- All code public and auditable
- Standard Terraform/Pulumi workflows
- Migration costs measured in hours
- Open security model

**The Promise:**
Use Planton because it's the best platform for your needs—not because switching is too expensive.

---

### Visual Direction
- GitHub repository screenshot showing OpenMCF repo with stars, forks, recent commits
- Side-by-side comparison: "w/ Vendor Lock-In Risk" vs "w/ Planton"
- Flow diagram: Export configs → Independent CLI → Continue managing infrastructure

---

## Section 7: Agent Fleet

### Section Title
**Agent Fleet: AI That Actually Works for DevOps**

### Subtitle
Specialized agents, not generic copilots. Deterministic outputs, not 80% accuracy.

### Introduction: The Copilot Experiment That Failed

**Honest Context:**
Planton spent 8 months building a generic AI copilot for DevOps. It hallucinated constantly. Customers were hesitant to use it for production because accuracy wasn't good enough. DevOps requires deterministic outputs—if you're creating or destroying a database, 80% accuracy isn't acceptable.

**The Pivot:**
Learning from Cursor's success, Planton rebuilt the AI layer from scratch as "Agent Fleet"—specialized agents for specific DevOps tasks combined with real tooling, deep context, and narrow scope. Not a one-size-fits-all chatbot.

---

### What Agent Fleet Is

**Not a Generic Copilot:**
Agent Fleet doesn't try to do everything. Each agent is specialized for a specific domain—pipeline troubleshooting, ECS debugging, Kubernetes diagnostics. Narrow scope = higher accuracy = production-ready.

**Agentic Architecture:**
These aren't chatbots with fancy prompts. They're agents with real capabilities:
- **Terminal access** - Run commands, inspect logs, analyze configurations
- **Repository operations** - Clone repos, read code, understand context
- **Cloud API integration** - Query AWS, GCP, Azure for real-time resource state
- **Platform awareness** - Full context of your Planton infrastructure and services

**Technology Under the Hood:**
- **MCP (Model Context Protocol) servers** - AWS MCP, GitHub MCP, Planton MCP
- **Specialized tooling** - Each agent has custom tools for its domain
- **Context-aware LLMs** - Agents understand your specific infrastructure, not generic examples

---

### Current Agents (Beta)

#### Agent 1: Pipeline Manager

**What It Does:**
Troubleshoots CI/CD pipeline failures—whether they're due to Docker build errors, deployment misconfigurations, or code issues.

**Capabilities:**
- Analyzes build logs to identify root cause
- Understands Tekton, GitHub Actions, and GitLab CI configurations
- Can clone your repository to inspect Dockerfiles and code
- Suggests fixes based on actual errors, not generic advice

**Example Use Case:**
Your pipeline fails with "ERROR: failed to solve: process '/bin/sh -c npm install' did not complete successfully"

Pipeline Manager:
1. Reads the full build log
2. Identifies the specific npm dependency causing the issue
3. Checks your `package.json` for compatibility problems
4. Suggests the fix: update Node.js version or pin dependency version

---

#### Agent 2: ECS Troubleshooter

**What It Does:**
Debugs AWS ECS service issues—task failures, networking problems, resource constraints.

**Capabilities:**
- Integrates with AWS MCP server to query ECS task status
- Reads CloudWatch logs for your failed tasks
- Analyzes service event history
- Understands Planton service configurations and deployment history

**Example Use Case:**
Your ECS service shows "STOPPED (Essential container exited)" but you don't know why.

ECS Troubleshooter:
1. Queries ECS for task details via AWS API
2. Pulls CloudWatch logs for the stopped container
3. Identifies exit code and error message
4. Cross-references with your service configuration
5. Suggests fix: increase memory limit or fix application crash

---

### Why This Approach Works

**Deterministic Outputs:**
Specialized agents with real tooling produce consistent, reliable results. They don't hallucinate because they work with actual data—your logs, your configs, your infrastructure state.

**Production-Ready:**
Customers trust these agents for real troubleshooting because they behave like experienced DevOps engineers—methodical investigation, root cause analysis, specific fixes.

**Continuously Improving:**
As Planton adds more specialized agents (Kubernetes Diagnostics, Database Performance, Security Scanner), the fleet becomes more capable without sacrificing accuracy.

---

### Roadmap (Coming Soon)

**Kubernetes Diagnostics Agent:**
- Pod crash loops and restarts
- Resource quota and limit issues
- Network policy debugging

**Database Performance Agent:**
- Slow query analysis
- Connection pool optimization
- Index recommendations

**Security Scanner Agent:**
- Infrastructure misconfigurations
- Compliance violations (HIPAA, SOC 2)
- Overly permissive IAM roles

---

### CTA
**Request Early Access to Agent Fleet Beta**

Currently available for select customers. Join the waitlist to get priority access when we open beta to all users.

---

### Visual Direction
- Agent interface screenshot showing Pipeline Manager analyzing a failed build
- Terminal output showing agent running commands and analyzing logs
- Before/After: Generic copilot response vs Agent Fleet specific fix

---

## Section 8: Customer Stories

### Section Title
**Real Customers, Real Results**

### Subtitle
From consulting firms to fintech platforms—see how teams use Planton to accelerate delivery and eliminate ops bottlenecks.

---

### Customer Story 1: Tynybay → Odwen Project

**Customer:** Tynybay (IT Consulting Firm)  
**Project:** Odwen - Online warehousing platform  
**Cloud:** Google Cloud Platform  
**Team Size:** 3 people (1 DevOps engineer, 2 developers)

**The Challenge:**
Tight timeline to deliver a working warehousing platform for their client. Infrastructure setup was blocking developers from making progress on product features.

**Before Planton:**
- DevOps engineer spent at least a couple weeks setting up GCP Cloud Run environment
- Developers waited for infrastructure before writing code
- Each new environment (QA, staging, prod) required manual replication
- Service onboarding required DevOps engineer involvement

**With Planton:**
- **Infrastructure deployment:** <1 hour for complete GCP Cloud Run environment
- **Developer autonomy:** Developers self-onboard services without blocking on DevOps
- **Multi-environment:** QA, staging, and prod environments created in minutes
- **Ongoing velocity:** Developers manage their own deployments through Service Hub

**The Result:**
> "Developers self-onboard their own services. That's a huge win."

The DevOps engineer now focuses on architecture and optimization—not repetitive environment setup and service deployments.

---

### Customer Story 2: iorta TechNext

**Customer:** iorta TechNext (via Tynybay consulting)  
**Product:** SalesVerse - Sales cycle management platform  
**Industry:** Financial Services (BFSI)  
**Team Size:** 7 developers  
**Cloud:** AWS (ECS)  
**Monthly Spend:** ~$450 (platform + usage)

**The Challenge:**
Small team building a sales management platform without dedicated DevOps resources. Needed secure, compliant infrastructure without the overhead of hiring ops specialists.

**Before Planton:**
- Technical founder managing infrastructure manually through AWS Console
- Inconsistent deployments across environments
- Limited visibility into what changed and when
- Slow iteration due to infrastructure being a bottleneck

**With Planton:**
- **Complete self-service:** 7 developers manage all infrastructure and deployments
- **Highest platform usage** among all Planton customers (most active account)
- **No dedicated DevOps hire needed** - $150K+/year savings
- **Full audit trail** - Every infrastructure change tracked for compliance

**The Economics:**
- Team size: 7 developers
- Subscription: $140/month (7 × $20/dev)
- Usage: ~$310/month (automation minutes)
- Total: ~$450/month vs $150K+/year for DevOps engineer

**The Result:**
iorta TechNext gets enterprise-grade DevOps capabilities at 1/30th the cost of hiring. Team focuses on product, not infrastructure.

---

### Metrics Bar

**Infrastructure Deployment Time:**
<1 hour (vs couple weeks manually, days with basic automation)

**Customer Retention:**
100% (despite platform being less stable at launch—customer loyalty through rapid iteration)

**Velocity Improvement:**
5x faster than manual cloud console work

**Cost to Start:**
$20/developer/month + usage

**Support:**
24/7 availability for all customers

---

### Visual Direction
- Customer logo grid (if available) or testimonial cards
- Before/After timeline showing deployment time reduction
- Metrics visualization (bar charts showing time/cost savings)

---

## Section 9: Pricing

### Section Title
**Simple, Usage-Based Pricing**

### Subtitle
Pay for what you use. Start free. Scale as you grow. No hidden costs.

---

### Pricing Model Overview

**Two Components:**

**1. Seat-Based Subscription**
$20 per developer per month

**2. Automation Minutes**
$0.006 per minute (0.6¢/min)

Metered separately for infrastructure runners and service CI/CD runners, but charged at the same rate.

**For detailed pricing tiers and calculator, visit [planton.ai/pricing](https://planton.ai/pricing)**

---

### Real Customer Example: iorta TechNext

**Team Size:** 7 developers  
**Monthly Spend:** ~$450 (platform + usage)

**What They Get:**
- Complete infrastructure management (AWS ECS environment for SalesVerse)
- Automated CI/CD for all services
- Multi-environment deployments (dev, QA, staging, prod)
- 24/7 support access
- Full audit trail for compliance

**Alternative Cost:**
Hiring one DevOps engineer: $150K+/year = $12,500/month

**Savings:** $12,050/month | 96% cost reduction

---

### CTA
**Calculate Your Pricing**

See exactly what you'll pay based on team size and deployment frequency.

[Visit Pricing Page →](https://planton.ai/pricing)

---

## Section 10: Built by DevOps Engineers

### Section Title
**Built by DevOps Engineers, For Real-World Use**

### Subtitle
Not a theory. Not a prototype. A production platform born from 10+ years of solving infrastructure problems at scale.

---

### Dogfooding: Planton Runs on Planton

**We Don't Ask You to Trust Something We Don't Use:**

Planton's production infrastructure is 100% managed through Planton itself. Zero exceptions.

**Our Stack:**
- **10+ database types:** ClickHouse, Redis, Postgres, MongoDB, and more
- **Deployment target:** Google Cloud Platform (multi-region)
- **Orchestration:** All databases run on Kubernetes, deployed via Planton
- **Service Hub:** All backend services deployed through Planton pipelines (Tekton)
- **Zero manual deployments:** Everything through the platform or CLI

**What This Means:**
Every feature we ship, we use ourselves. Every bug we fix, we've likely hit in production. Every performance optimization, we've measured on our own infrastructure.

---

### Support & Operations

**24/7 Support Capability:**
Distributed team across time zones ensures customer issues get addressed around the clock—not just business hours.

**Customer Success Focus:**
100% retention rate didn't happen by accident. We iterate rapidly on customer feedback and treat early adopters as partners, not just users.

**Production SLAs:**
- Plus tier: 24-hour email response
- Pro tier: 12-hour email response
- Enterprise: Custom SLA (Slack/phone support available)

---

### Engineering Philosophy

**Grounded in Reality:**
Planton solves real problems experienced by DevOps engineers—no hype, no theoretical solutions.

**Customer-Driven Roadmap:**
Features are prioritized based on actual customer requests from paying users—not hypothetical use cases or marketing trends.

**Iterate Fast, Break Nothing:**
Platform stability is non-negotiable. New features ship to beta customers first, then roll out gradually with monitoring at every stage.

**Open by Default:**
All deployment components open source. Transparent about how things work. No black boxes.

---

### Visual Direction
- Dogfooding visualization: "Planton managing Planton" infrastructure diagram
- Team photo or distributed team map showing 24/7 support coverage
- Customer success metrics dashboard

---

## Section 11: Security & Compliance

### Section Title
**Secure by Design**

### Subtitle
Multiple security models and scoped permissions—because trust is earned, not assumed.

---

### Deployment Runner Security Models

**You Choose Your Security Posture:**

Planton supports three deployment runner topologies. Start simple, upgrade to more secure models as your requirements grow.

---

#### Model 1: Basic - AWS Access Keys

**How It Works:**
Provide AWS access key and secret key. Planton-hosted runners use these credentials to deploy infrastructure in your account.

**Security:**
- Credentials encrypted at rest
- Stored in secure secrets management (Google Secret Manager)
- Never logged or exposed in UI
- Rotatable anytime via platform settings

**Best For:**
- Development environments
- Initial platform evaluation
- Non-production workloads

**Tradeoffs:**
- Long-lived credentials (rotation recommended)
- Highest trust in Planton (we hold the keys)

---

#### Model 2: Trust Relationship - AWS Assume Role

**How It Works:**
Configure cross-account trust policy in your AWS account. Planton platform assumes a role in your account with time-limited credentials.

**Security:**
- No long-lived credentials stored
- Temporary session tokens (auto-expire)
- Trust policy controls which Planton account can assume role
- Revocable instantly (delete trust relationship)

**Best For:**
- Production environments
- Compliance requirements (HIPAA, SOC 2)
- Organizations with centralized AWS governance

**Setup:**
1. Create IAM role in your AWS account
2. Add trust policy allowing Planton's AWS account
3. Attach policies with scoped permissions (see below)
4. Provide role ARN to Planton platform

**Tradeoffs:**
- Requires AWS IAM configuration (5-10 minutes)
- Still requires trust in Planton platform

---

#### Model 3: Customer-Hosted Runner

**How It Works:**
Planton provides Docker image (`planton-deployer`). You deploy it as an ECS service in YOUR AWS account. Runner communicates back to Planton platform for job orchestration, but all infrastructure operations run within your boundary.

**Security:**
- Credentials never leave your AWS account
- You control IAM role attached to runner
- You control network egress/ingress
- Runner code auditable (Docker image layers inspectable)

**Best For:**
- Maximum security requirements
- Zero-trust environments
- Highly regulated industries (healthcare, finance)
- Air-gapped or restricted networks

**Setup:**
1. Pull `planton-deployer` Docker image from Planton registry
2. Deploy as ECS service in your AWS account
3. Attach IAM role with scoped permissions
4. Configure runner to communicate with Planton control plane
5. Platform detects runner and routes jobs to it

**Tradeoffs:**
- You manage runner infrastructure (minimal overhead, ~$10-20/month for ECS task)
- Requires networking configuration (if using VPC endpoints)

---

### Scoped Permissions: Least Privilege by Design

**No Blanket Admin Required:**

Planton never asks for `AdministratorAccess` or wildcard permissions. Every deployment component specifies exact IAM permissions needed.

**How It Works:**

**1. Audit the Terraform Code:**
All deployment components are open source in [OpenMCF repository](https://github.com/plantonhq/openmcf). Read the exact resources being created.

**2. Generate Scoped IAM Policy:**
Each component includes documented IAM permissions. Example for AWS ALB:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "elasticloadbalancing:CreateLoadBalancer",
        "elasticloadbalancing:DescribeLoadBalancers",
        "elasticloadbalancing:ModifyLoadBalancerAttributes",
        "elasticloadbalancing:DeleteLoadBalancer",
        "elasticloadbalancing:CreateListener",
        "elasticloadbalancing:CreateTargetGroup"
      ],
      "Resource": "*"
    }
  ]
}
```

**3. Combine Permissions for Your Stack:**
Deploying Infra Chart with VPC + ALB + ECR? Combine IAM policies for each component.

**4. Attach to Runner:**
Whether using assume role or customer-hosted runner, attach the scoped policy.

**Why This Matters:**
- Security teams can audit exact permissions before granting access
- Principle of least privilege enforced at infrastructure level
- No risk of Planton platform over-reaching permissions

---

### Gradual Trust Model for Regulated Industries

**Start Non-Production:**
Deploy dev, QA, and staging environments first. Build trust through rigorous testing:
- Performance benchmarks
- Vulnerability scanning (static + dynamic)
- Security audits
- Compliance spot-checks

**Earn Production Access:**
Only after validation cycle, promote to production. Planton supports overlap period—run legacy and Planton-managed infrastructure in parallel during transition.

---

### Visual Direction
- Diagram showing three deployment runner topologies (basic → trust → customer-hosted)
- Scoped permissions example (IAM policy JSON snippet)

---

## Section 12: Final CTA

### Section Title
**Ready to Eliminate Your DevOps Bottleneck?**

### Subtitle
Join IT consulting firms and startups deploying infrastructure in minutes, not weeks.

---

### Primary CTA

**Start Free Trial**

- 100 automation minutes included
- No credit card required
- Full platform access (Infra Hub + Service Hub)
- Deploy your first environment in <10 minutes

**Button Text:** Start Free Trial →

---

### Secondary CTA

**Book Live Demo**

30-minute walkthrough with actual customer setup. See how Planton deploys infrastructure and services in real-time.

What you'll see:
- Infrastructure deployment (complete environment in <1 hour)
- Service onboarding (Git to cloud in minutes)
- Multi-cloud capabilities (AWS, GCP, Azure)
- Agent Fleet (if interested in beta access)

**Button Text:** Book Demo →

---

### Tertiary CTA

**View Pricing Calculator**

See exactly what you'll pay based on team size and deployment frequency. No surprises.

**Button Text:** Calculate Pricing →

---

### Trust Bar

Display as badges/icons with text:

✅ **5 Paying Customers**  
✅ **100% Retention Rate**  
✅ **Open Source Foundation**  
✅ **Multi-Cloud Support** (AWS • GCP • Azure)  
✅ **24/7 Support Available**

---

### Customer Proof

**"Developers self-onboard their own services. That's a huge win."**
— Tynybay, deploying Odwen warehousing platform

**"Complete AWS ECS environment deployed in under 1 hour with 7 developers managing everything themselves."**
— iorta TechNext, SalesVerse platform

---

### Footer Navigation

**Product**
- Infra Hub
- Service Hub
- Agent Fleet (Beta)
- Deployment Components
- Pricing

**Resources**
- Documentation
- API Reference
- OpenMCF (GitHub)
- Customer Stories
- Blog

**Company**
- About
- Careers
- Contact
- Support

**Legal**
- Privacy Policy
- Terms of Service
- Security
- Compliance

**Community**
- Discord
- GitHub Discussions
- Twitter
- LinkedIn

---

### Bottom Footer

**Planton**  
DevOps-in-a-Box: From Code to Cloud in Under an Hour

**Open Source:**  
[github.com/plantonhq/openmcf](https://github.com/plantonhq/openmcf)

**Contact:**  
support@planton.ai

**© 2025 Planton. All rights reserved.**

---

### Visual Direction
- Clean, spacious CTA section with gradient background
- Three CTA buttons with clear visual hierarchy (primary = filled, secondary = outlined, tertiary = text link)
- Trust bar with icon badges
- Customer quote cards with attribution
- Standard footer layout with organized link columns

---

# End of Landing Page Copy

---

## Implementation Notes for Web Developer

### Content Structure
This document provides complete copy for 12 landing page sections. Each section includes:
- Headlines and subheadlines
- Body copy
- Bullet points and feature lists
- CTAs (calls to action)
- Customer quotes (where applicable)
- Visual direction notes

### Design System Alignment
Maintain existing Planton design system:
- Typography: Use existing TypoH1, TypoH2, TypoP1 components
- Colors: Gradient backgrounds, border colors, text colors from current site
- Spacing: Maintain current section padding and gaps
- Responsive: Mobile-first, ensure all copy works on small screens

### Component Mapping

**Replace:**
- `conversational.tsx` → Section 1 (Hero)
- `experience-the-future.tsx` → Section 12 (Final CTA)

**Evolve:**
- `home-features.tsx` → Section 3 (Infra Hub)
- `home-integrated-iac.tsx` → Can be repurposed or replaced with new content

**New:**
- Section 2 (Problem/Solution)
- Section 4 (Service Hub) - Critical, currently missing
- Section 5 (Built on Open Standards) - NEW in Draft 2
- Section 6 (Open Source)
- Section 7 (Agent Fleet)
- Section 8 (Customer Stories)
- Section 9 (Pricing - simplified)
- Section 10 (Built by DevOps Engineers - simplified, no founder story/architecture)
- Section 11 (Security & Compliance)

### Key Changes in Draft 2

**Customer Name Corrections:**
- "IOTA Technics" → "iorta TechNext" (iorta.in subsidiary)
- "TynyBay" → "Tynybay" (consulting firm)
- "Oddwin" → "Odwen" (odwen.co.in)
- Tynybay is the consulting company that brought us iorta TechNext and Odwen as clients
- TechNext deploys SalesVerse product to AWS ECS
- Odwen deploys to GCP Cloud Run

**Content Updates:**
- Removed Founder Story subsection
- Removed Platform Architecture subsection
- Added new "Built on Open Standards" section highlighting GitHub, Terraform, Pulumi, Tekton transparency
- Changed "one day" deployment time to "couple weeks" for manual work
- Changed vendor comparisons from "Most Platforms" vs "Planton" to "w/ Vendor Lock-In Risk" vs "w/ Planton"
- Changed all "planton.cloud" to "planton.ai"
- Simplified Pricing section (Section 9) - removed detailed tiers, kept only overview and customer example
- Added transparency about Tekton being the CI/CD engine behind Service Hub

### Copy Guidelines

**Tone:**
- Technical but accessible
- Honest (no marketing fluff)
- Grounded in real customer usage
- Confident but not arrogant
- Transparent about technology choices

**Emphasis:**
Bold and italics already applied in markdown. Translate to appropriate HTML/CSS styling.

**Links:**
- All GitHub links point to `github.com/plantonhq/openmcf`
- All planton links use `planton.ai` domain
- CTAs should link to:
  - Trial signup flow
  - Demo booking (Calendly or equivalent)
  - Pricing page at planton.ai/pricing

### Visual Assets Needed

**Screenshots/Diagrams:**
- Infra Hub deployment flow (DAG visualization)
- Service Hub pipeline flow (Git → Tekton → deploy)
- Deployment Component Store browser
- Agent Fleet interface
- Three deployment runner topologies diagram
- Open Standards integration diagram (GitHub ↔ Planton ↔ Terraform/Pulumi ↔ Tekton)
- Tekton pipeline definition screenshot

**Customer Logos:**
Request from customers (Tynybay, iorta TechNext, Odwen) if available for social proof.

### Metrics to Track

Once live, track:
- CTA click rates (trial signup vs demo booking)
- Section scroll depth (which sections get read)
- Time on page
- Bounce rate from hero (clarity test)

---

**Ready for Review**

This is Draft 2 incorporating all customer corrections and positioning updates. Review, provide feedback, and iterate as needed before final implementation.
