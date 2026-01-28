# Overview

The Pulumi module for **CronJobKubernetes** simplifies and standardizes how developers create and manage scheduled tasks
on Kubernetes. By consuming a `CronJobKubernetesStackInput`, which encapsulates everything from cluster credentials and
Docker image references to scheduling rules and environment variable mappings, this module programmatically generates
the necessary Kubernetes resources. It harnesses Pulumi's Kubernetes provider under the hood, ensuring an
infrastructure-as-code approach that is both powerful and straightforward.

Key highlights of this module include:

- **Automated Resource Creation**  
  Generates and manages all required resources (e.g., Namespace, CronJob, Secrets) purely based on the declarative
  specification found in the input. This eliminates repetitive configuration and keeps your workflow simple.

- **Support for Secrets and Environment Variables**  
  Securely injects secrets into your CronJob pods via Kubernetes Secrets, while also supporting environment variables
  for non-sensitive parameters. This separation follows best practices by preventing sensitive data from being stored in
  plain text within code or container images.

- **Flexible Scheduling**  
  With fields such as `schedule`, `concurrencyPolicy`, `suspend`, and `startingDeadlineSeconds`, developers can tailor
  when and how often tasks should run. Whether you need daily backups, periodic data processing, or event-driven tasks,
  you can configure them all through a standard Cron expression.

- **Resource Requests and Limits**  
  Manage CPU and memory usage with built-in resource requests and limits. This ensures that CronJob containers have the
  correct allocations, preventing resource contention or starvation in the Kubernetes cluster.

- **Optional Image Pull Credentials**  
  If your container images are hosted in private registries, you can easily provide Docker credentials (
  `docker_config_json`), which the module uses to create an `image-pull-secret`. This keeps authentication secure and
  streamlined.

- **Consistent Naming and Labeling**  
  A standardized labeling strategy (e.g., `resource`, `resource_kind`) and alignment with the `metadata` fields ensure
  clear resource identification and traceability. This unifies deployments across various environments and teams.

- **Flexible Namespace Management**  
  Control whether the module creates a new namespace or references an existing one via the `create_namespace` flag. When
  set to `true`, the module creates and manages the namespace with appropriate labels. When set to `false`, it references
  an existing namespace, enabling multi-tenant scenarios where multiple CronJobs share a namespace or GitOps workflows
  where namespace lifecycle is managed separately. This flexibility supports both isolated and shared deployment patterns.

Overall, the **CronJobKubernetes** Pulumi module provides a robust, opinionated foundation for running scheduled jobs on
Kubernetes. By codifying best practices and reducing boilerplate, it empowers you to focus on application logic rather
than low-level infrastructure details. With minimal configuration, you can confidently provision CronJobs, integrate
secrets, handle concurrency concerns, and manage your cluster resources in a clean, controlled fashion.
